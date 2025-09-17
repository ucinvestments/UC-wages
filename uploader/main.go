package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/schollz/progressbar/v3"
)

type Config struct {
	Workers     int
	DataDir     string
	DatabaseURL string
	BatchSize   int
}

type WageData struct {
	Location     string                   `json:"location"`
	Year         int                      `json:"year"`
	ScrapedAt    string                   `json:"scraped_at"`
	TotalRecords int                      `json:"total_records"`
	Records      []map[string]interface{} `json:"records"`
}

type WageRecord struct {
	Location   string
	Year       int
	EmployeeID int
	Firstname  string
	Lastname   string
	Title      string
	Basepay    float64
	Overtimepay float64
	Adjustpay  float64
	Grosspay   float64
	ScrapedAt  time.Time
}

type Task struct {
	FilePath string
	Location string
	Year     int
}

type Uploader struct {
	db           *sql.DB
	config       Config
	success      int64
	failed       int64
	total        int64
	processed    int64
	mu           sync.Mutex
	progressBar  *progressbar.ProgressBar
	totalFiles   int64
}

func NewUploader(cfg Config) (*Uploader, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Uploader{
		db:     db,
		config: cfg,
	}, nil
}

func (u *Uploader) Close() error {
	return u.db.Close()
}

func (u *Uploader) initSchema() error {
	schemaSQL := `
	-- Create table for wage records
	CREATE TABLE IF NOT EXISTS uc_wages (
		id SERIAL PRIMARY KEY,
		location VARCHAR(50) NOT NULL,
		year INTEGER NOT NULL,
		employee_id INTEGER,
		firstname VARCHAR(100),
		lastname VARCHAR(100),
		title VARCHAR(200),
		basepay DECIMAL(12,2) DEFAULT 0.00,
		overtimepay DECIMAL(12,2) DEFAULT 0.00,
		adjustpay DECIMAL(12,2) DEFAULT 0.00,
		grosspay DECIMAL(12,2) DEFAULT 0.00,
		scraped_at TIMESTAMP,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(location, year, employee_id)
	);

	-- Create indexes for faster queries
	CREATE INDEX IF NOT EXISTS idx_uc_wages_location ON uc_wages(location);
	CREATE INDEX IF NOT EXISTS idx_uc_wages_year ON uc_wages(year);
	CREATE INDEX IF NOT EXISTS idx_uc_wages_location_year ON uc_wages(location, year);
	CREATE INDEX IF NOT EXISTS idx_uc_wages_grosspay ON uc_wages(grosspay);
	CREATE INDEX IF NOT EXISTS idx_uc_wages_title ON uc_wages(title);

	-- Create table for upload tracking
	CREATE TABLE IF NOT EXISTS upload_progress (
		id SERIAL PRIMARY KEY,
		location VARCHAR(50) NOT NULL,
		year INTEGER NOT NULL,
		total_records INTEGER,
		uploaded_records INTEGER DEFAULT 0,
		status VARCHAR(20) DEFAULT 'pending',
		started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		completed_at TIMESTAMP,
		error_message TEXT,
		UNIQUE(location, year)
	);`

	_, err := u.db.Exec(schemaSQL)
	return err
}

func (u *Uploader) parseAmount(amount interface{}) float64 {
	if amount == nil {
		return 0.0
	}

	amountStr := fmt.Sprintf("%v", amount)
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amountStr = strings.TrimSpace(amountStr)

	if amountStr == "" {
		return 0.0
	}

	value, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Printf("Warning: failed to parse amount '%s': %v", amountStr, err)
		return 0.0
	}

	return value
}

func (u *Uploader) parseScrapedAt(scrapedAt string) time.Time {
	if scrapedAt == "" {
		return time.Now()
	}

	t, err := time.Parse(time.RFC3339, scrapedAt)
	if err != nil {
		log.Printf("Warning: failed to parse scraped_at '%s': %v", scrapedAt, err)
		return time.Now()
	}

	return t
}

func (u *Uploader) convertRecord(record map[string]interface{}, location string, year int, scrapedAt time.Time) WageRecord {
	employeeID := 0
	if id, ok := record["id"]; ok {
		if idFloat, ok := id.(float64); ok {
			employeeID = int(idFloat)
		} else if idStr := fmt.Sprintf("%v", id); idStr != "" {
			if parsed, err := strconv.Atoi(idStr); err == nil {
				employeeID = parsed
			}
		}
	}

	firstname := ""
	if fn, ok := record["firstname"]; ok {
		firstname = fmt.Sprintf("%v", fn)
	}

	lastname := ""
	if ln, ok := record["lastname"]; ok {
		lastname = fmt.Sprintf("%v", ln)
	}

	title := ""
	if t, ok := record["title"]; ok {
		title = fmt.Sprintf("%v", t)
	}

	return WageRecord{
		Location:    location,
		Year:        year,
		EmployeeID:  employeeID,
		Firstname:   firstname,
		Lastname:    lastname,
		Title:       title,
		Basepay:     u.parseAmount(record["basepay"]),
		Overtimepay: u.parseAmount(record["overtimepay"]),
		Adjustpay:   u.parseAmount(record["adjustpay"]),
		Grosspay:    u.parseAmount(record["grosspay"]),
		ScrapedAt:   scrapedAt,
	}
}

func (u *Uploader) updateProgress(location string, year int, status string, totalRecords, uploadedRecords int, errorMsg string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	var completedAt interface{}
	if status == "completed" || status == "failed" {
		completedAt = time.Now()
	} else {
		completedAt = nil
	}

	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO upload_progress (location, year, total_records, uploaded_records, status, error_message, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (location, year)
		DO UPDATE SET
			uploaded_records = EXCLUDED.uploaded_records,
			status = EXCLUDED.status,
			error_message = EXCLUDED.error_message,
			completed_at = EXCLUDED.completed_at
	`

	_, err = tx.Exec(query, location, year, totalRecords, uploadedRecords, status, errorMsg, completedAt)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *Uploader) uploadBatch(records []WageRecord) error {
	if len(records) == 0 {
		return nil
	}

	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO uc_wages (location, year, employee_id, firstname, lastname, title, basepay, overtimepay, adjustpay, grosspay, scraped_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (location, year, employee_id) DO UPDATE SET
			firstname = EXCLUDED.firstname,
			lastname = EXCLUDED.lastname,
			title = EXCLUDED.title,
			basepay = EXCLUDED.basepay,
			overtimepay = EXCLUDED.overtimepay,
			adjustpay = EXCLUDED.adjustpay,
			grosspay = EXCLUDED.grosspay,
			scraped_at = EXCLUDED.scraped_at,
			uploaded_at = CURRENT_TIMESTAMP
	`

	for _, record := range records {
		_, err := tx.Exec(query,
			record.Location,
			record.Year,
			record.EmployeeID,
			record.Firstname,
			record.Lastname,
			record.Title,
			record.Basepay,
			record.Overtimepay,
			record.Adjustpay,
			record.Grosspay,
			record.ScrapedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert record %d: %w", record.EmployeeID, err)
		}
	}

	return tx.Commit()
}

func (u *Uploader) processFile(task Task) {
	startTime := time.Now()

	if err := u.updateProgress(task.Location, task.Year, "processing", 0, 0, ""); err != nil {
		log.Printf("❌ Error updating progress for %s %d: %v", task.Location, task.Year, err)
	}

	log.Printf("📂 Reading file: %s", task.FilePath)
	data, err := os.ReadFile(task.FilePath)
	if err != nil {
		log.Printf("❌ Error reading file %s: %v", task.FilePath, err)
		u.updateProgress(task.Location, task.Year, "failed", 0, 0, err.Error())
		atomic.AddInt64(&u.failed, 1)
		u.progressBar.Add(1)
		return
	}

	log.Printf("🔍 Parsing JSON data for %s %d (size: %.2f MB)", task.Location, task.Year, float64(len(data))/(1024*1024))
	var wageData WageData
	if err := json.Unmarshal(data, &wageData); err != nil {
		log.Printf("❌ Error parsing JSON from %s: %v", task.FilePath, err)
		u.updateProgress(task.Location, task.Year, "failed", 0, 0, err.Error())
		atomic.AddInt64(&u.failed, 1)
		u.progressBar.Add(1)
		return
	}

	scrapedAt := u.parseScrapedAt(wageData.ScrapedAt)
	totalRecords := len(wageData.Records)
	uploadedRecords := 0

	if totalRecords == 0 {
		log.Printf("⚠️  No records found in %s", task.FilePath)
		u.updateProgress(task.Location, task.Year, "completed", 0, 0, "")
		atomic.AddInt64(&u.success, 1)
		u.progressBar.Add(1)
		return
	}

	atomic.AddInt64(&u.total, int64(totalRecords))
	log.Printf("📊 Starting upload for %s %d: %d records to process", task.Location, task.Year, totalRecords)

	var batch []WageRecord
	batchCount := 0

	for i, record := range wageData.Records {
		wageRecord := u.convertRecord(record, wageData.Location, wageData.Year, scrapedAt)
		batch = append(batch, wageRecord)

		if len(batch) >= u.config.BatchSize || i == len(wageData.Records)-1 {
			batchStartTime := time.Now()
			if err := u.uploadBatch(batch); err != nil {
				log.Printf("❌ Error uploading batch %d for %s: %v", batchCount, task.FilePath, err)
				u.updateProgress(task.Location, task.Year, "failed", totalRecords, uploadedRecords, err.Error())
				atomic.AddInt64(&u.failed, 1)
				u.progressBar.Add(1)
				return
			}

			uploadedRecords += len(batch)
			batchCount++
			atomic.AddInt64(&u.processed, int64(len(batch)))

			batchDuration := time.Since(batchStartTime)
			recordsPerSec := float64(len(batch)) / batchDuration.Seconds()

			log.Printf("✅ Batch %d complete for %s %d: %d records (%.1f rec/sec)",
				batchCount, task.Location, task.Year, len(batch), recordsPerSec)

			batch = nil

			if batchCount%5 == 0 {
				progress := float64(uploadedRecords) / float64(totalRecords) * 100
				totalProcessed := atomic.LoadInt64(&u.processed)
				log.Printf("📈 Progress: %s %d - %d/%d (%.1f%%) | Overall: %d records processed",
					task.Location, task.Year, uploadedRecords, totalRecords, progress, totalProcessed)
			}
		}
	}

	if err := u.updateProgress(task.Location, task.Year, "completed", totalRecords, uploadedRecords, ""); err != nil {
		log.Printf("❌ Error updating final progress for %s %d: %v", task.Location, task.Year, err)
	}

	duration := time.Since(startTime)
	recordsPerSec := float64(uploadedRecords) / duration.Seconds()

	log.Printf("🎉 Completed: %s %d - %d records in %v (%.1f rec/sec)",
		task.Location, task.Year, uploadedRecords, duration, recordsPerSec)

	atomic.AddInt64(&u.success, 1)
	u.progressBar.Add(1)
}

func (u *Uploader) run() error {
	if err := u.initSchema(); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	var tasks []Task
	err := filepath.WalkDir(u.config.DataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasPrefix(d.Name(), "wages_") && strings.HasSuffix(d.Name(), ".json") {
			parts := strings.Split(path, string(filepath.Separator))
			if len(parts) >= 2 {
				location := parts[len(parts)-2]
				yearStr := strings.TrimSuffix(strings.TrimPrefix(d.Name(), "wages_"), ".json")
				if year, err := strconv.Atoi(yearStr); err == nil {
					tasks = append(tasks, Task{
						FilePath: path,
						Location: location,
						Year:     year,
					})
				}
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan data directory: %w", err)
	}

	if len(tasks) == 0 {
		log.Println("📭 No wage files found to upload")
		return nil
	}

	u.totalFiles = int64(len(tasks))
	log.Printf("🚀 Found %d files to upload", len(tasks))

	// Create progress bar
	u.progressBar = progressbar.NewOptions(len(tasks),
		progressbar.OptionSetDescription("📊 Uploading files"),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("files"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetElapsedTime(true),
	)

	ch := make(chan Task, len(tasks))
	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < u.config.Workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("👷 Worker %d starting", workerID)
			for task := range ch {
				u.processFile(task)
			}
			log.Printf("✅ Worker %d finished", workerID)
		}(i)
	}

	for _, task := range tasks {
		ch <- task
	}
	close(ch)

	wg.Wait()
	u.progressBar.Finish()

	duration := time.Since(startTime)
	successCount := atomic.LoadInt64(&u.success)
	failedCount := atomic.LoadInt64(&u.failed)
	totalRecords := atomic.LoadInt64(&u.total)
	recordsPerSec := float64(totalRecords) / duration.Seconds()

	fmt.Printf("\n")
	log.Printf("🎉 Upload complete!")
	log.Printf("📈 Summary:")
	log.Printf("   • Files processed: %d", len(tasks))
	log.Printf("   • Files succeeded: %d", successCount)
	log.Printf("   • Files failed: %d", failedCount)
	log.Printf("   • Total records: %d", totalRecords)
	log.Printf("   • Duration: %v", duration)
	log.Printf("   • Average rate: %.1f records/sec", recordsPerSec)

	if failedCount > 0 {
		log.Printf("⚠️  Some files failed to upload. Check the logs above for details.")
	}

	return nil
}

func main() {
	var (
		workers  = flag.Int("workers", 5, "number of concurrent workers")
		dataDir  = flag.String("data", "../data", "data directory containing wage files")
		batchSize = flag.Int("batch", 1000, "batch size for database inserts")
		envFile  = flag.String("env", "", "path to .env file (optional)")
	)
	flag.Parse()

	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil {
			log.Printf("Warning: failed to load .env file: %v", err)
		}
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	config := Config{
		Workers:     *workers,
		DataDir:     *dataDir,
		DatabaseURL: databaseURL,
		BatchSize:   *batchSize,
	}

	log.Printf("Starting upload: %d workers, batch size %d, data dir: %s",
		config.Workers, config.BatchSize, config.DataDir)

	uploader, err := NewUploader(config)
	if err != nil {
		log.Fatal(err)
	}
	defer uploader.Close()

	if err := uploader.run(); err != nil {
		log.Fatal(err)
	}
}