package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/ucinvestments/uc-wages-analysis/pkg/models"
)

type Config struct {
	DatabaseURL string
	OutputDir   string
	DryRun      bool
	BatchSize   int
}

func main() {
	var config Config
	flag.StringVar(&config.DatabaseURL, "db", "", "PostgreSQL database URL")
	flag.StringVar(&config.OutputDir, "output", "./output", "Analysis output directory")
	flag.BoolVar(&config.DryRun, "dry-run", false, "Show what would be uploaded without actually doing it")
	flag.IntVar(&config.BatchSize, "batch", 50, "Batch size for uploads")
	flag.Parse()

	// Use environment variable if no flag provided
	if config.DatabaseURL == "" {
		config.DatabaseURL = os.Getenv("DATABASE_URL")
	}

	if config.DatabaseURL == "" {
		log.Fatal("Database URL is required. Use -db flag or set DATABASE_URL environment variable")
	}

	fmt.Println("🚀 UC Wages Analysis Uploader")
	fmt.Println("==============================")

	// Connect to database
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("✅ Connected to database")

	// Create tables if they don't exist (for development)
	if err := createTables(db); err != nil {
		log.Printf("Warning: Could not create tables: %v", err)
	}

	// Upload summaries
	if err := uploadSummaries(db, config); err != nil {
		log.Printf("Error uploading summaries: %v", err)
	}

	// Upload pyramids
	if err := uploadPyramids(db, config); err != nil {
		log.Printf("Error uploading pyramids: %v", err)
	}

	// Upload title analyses
	if err := uploadTitleAnalyses(db, config); err != nil {
		log.Printf("Error uploading title analyses: %v", err)
	}

	fmt.Println("\n✨ Upload completed successfully!")
}

func uploadSummaries(db *sql.DB, config Config) error {
	fmt.Println("\n▶ Uploading wage summaries...")

	summaryFiles, err := filepath.Glob(filepath.Join(config.OutputDir, "sums", "*.json"))
	if err != nil {
		return err
	}

	fmt.Printf("Found %d summary files\n", len(summaryFiles))

	if config.DryRun {
		fmt.Println("DRY RUN: Would upload summaries")
		return nil
	}

	// Prepare upsert statement
	stmt, err := db.Prepare(`
		INSERT INTO wage_summaries (
			location, year, employee_count, total_gross_pay, avg_gross_pay,
			median_pay, std_dev, min_pay, max_pay, percentiles,
			total_base, total_overtime, total_adjustments,
			avg_base, avg_overtime, avg_adjustments, generated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		) ON CONFLICT (location, year) DO UPDATE SET
			employee_count = EXCLUDED.employee_count,
			total_gross_pay = EXCLUDED.total_gross_pay,
			avg_gross_pay = EXCLUDED.avg_gross_pay,
			median_pay = EXCLUDED.median_pay,
			std_dev = EXCLUDED.std_dev,
			min_pay = EXCLUDED.min_pay,
			max_pay = EXCLUDED.max_pay,
			percentiles = EXCLUDED.percentiles,
			total_base = EXCLUDED.total_base,
			total_overtime = EXCLUDED.total_overtime,
			total_adjustments = EXCLUDED.total_adjustments,
			avg_base = EXCLUDED.avg_base,
			avg_overtime = EXCLUDED.avg_overtime,
			avg_adjustments = EXCLUDED.avg_adjustments,
			generated_at = EXCLUDED.generated_at,
			uploaded_at = NOW()
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, file := range summaryFiles {
		summary, err := loadSummary(file)
		if err != nil {
			log.Printf("Error loading %s: %v", file, err)
			continue
		}

		// Convert percentiles to JSON
		percentilesJSON, err := json.Marshal(summary.Percentiles)
		if err != nil {
			log.Printf("Error marshaling percentiles for %s: %v", file, err)
			continue
		}

		_, err = stmt.Exec(
			summary.Location, summary.Year, summary.EmployeeCount,
			summary.TotalGrossPay, summary.AvgGrossPay, summary.MedianPay,
			summary.StdDev, summary.MinPay, summary.MaxPay, string(percentilesJSON),
			summary.PayComponents.TotalBase, summary.PayComponents.TotalOvertime,
			summary.PayComponents.TotalAdjustments, summary.PayComponents.AvgBase,
			summary.PayComponents.AvgOvertime, summary.PayComponents.AvgAdjustments,
			summary.GeneratedAt,
		)

		if err != nil {
			log.Printf("Error inserting summary for %s %d: %v", summary.Location, summary.Year, err)
			continue
		}

		fmt.Printf("✓ Uploaded summary: %s %d\n", summary.Location, summary.Year)
	}

	return nil
}

func uploadPyramids(db *sql.DB, config Config) error {
	fmt.Println("\n▶ Uploading wage pyramids...")

	pyramidFiles, err := filepath.Glob(filepath.Join(config.OutputDir, "pyramid", "*.json"))
	if err != nil {
		return err
	}

	fmt.Printf("Found %d pyramid files\n", len(pyramidFiles))

	if config.DryRun {
		fmt.Println("DRY RUN: Would upload pyramids")
		return nil
	}

	// Prepare upsert statement
	stmt, err := db.Prepare(`
		INSERT INTO wage_pyramids (
			location, year, total_employees, total_pay, brackets, generated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) ON CONFLICT (location, year) DO UPDATE SET
			total_employees = EXCLUDED.total_employees,
			total_pay = EXCLUDED.total_pay,
			brackets = EXCLUDED.brackets,
			generated_at = EXCLUDED.generated_at,
			uploaded_at = NOW()
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, file := range pyramidFiles {
		pyramid, err := loadPyramid(file)
		if err != nil {
			log.Printf("Error loading %s: %v", file, err)
			continue
		}

		// Convert brackets to JSON
		bracketsJSON, err := json.Marshal(pyramid.Brackets)
		if err != nil {
			log.Printf("Error marshaling brackets for %s: %v", file, err)
			continue
		}

		_, err = stmt.Exec(
			pyramid.Location, pyramid.Year, pyramid.TotalEmployees,
			pyramid.TotalPay, string(bracketsJSON), pyramid.GeneratedAt,
		)

		if err != nil {
			log.Printf("Error inserting pyramid for %s %d: %v", pyramid.Location, pyramid.Year, err)
			continue
		}

		fmt.Printf("✓ Uploaded pyramid: %s %d\n", pyramid.Location, pyramid.Year)
	}

	return nil
}

func uploadTitleAnalyses(db *sql.DB, config Config) error {
	fmt.Println("\n▶ Uploading title analyses...")

	titleFiles, err := filepath.Glob(filepath.Join(config.OutputDir, "titles", "*.json"))
	if err != nil {
		return err
	}

	fmt.Printf("Found %d title analysis files\n", len(titleFiles))

	if config.DryRun {
		fmt.Println("DRY RUN: Would upload title analyses")
		return nil
	}

	// Prepare upsert statement
	stmt, err := db.Prepare(`
		INSERT INTO title_analysis (
			location, year, unique_titles, top_titles, generated_at
		) VALUES (
			$1, $2, $3, $4, $5
		) ON CONFLICT (location, year) DO UPDATE SET
			unique_titles = EXCLUDED.unique_titles,
			top_titles = EXCLUDED.top_titles,
			generated_at = EXCLUDED.generated_at,
			uploaded_at = NOW()
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, file := range titleFiles {
		analysis, err := loadTitleAnalysis(file)
		if err != nil {
			log.Printf("Error loading %s: %v", file, err)
			continue
		}

		// Convert top titles to JSON
		titlesJSON, err := json.Marshal(analysis.TopTitles)
		if err != nil {
			log.Printf("Error marshaling titles for %s: %v", file, err)
			continue
		}

		_, err = stmt.Exec(
			analysis.Location, analysis.Year, analysis.UniqueTitles,
			string(titlesJSON), analysis.GeneratedAt,
		)

		if err != nil {
			log.Printf("Error inserting title analysis for %s %d: %v", analysis.Location, analysis.Year, err)
			continue
		}

		fmt.Printf("✓ Uploaded title analysis: %s %d\n", analysis.Location, analysis.Year)
	}

	return nil
}

func loadSummary(filepath string) (*models.Summary, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var summary models.Summary
	if err := json.Unmarshal(data, &summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

func loadPyramid(filepath string) (*models.Pyramid, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var pyramid models.Pyramid
	if err := json.Unmarshal(data, &pyramid); err != nil {
		return nil, err
	}

	return &pyramid, nil
}

func loadTitleAnalysis(filepath string) (*models.TitleAnalysis, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var analysis models.TitleAnalysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		return nil, err
	}

	return &analysis, nil
}

func createTables(db *sql.DB) error {
	// This is a simplified version for development
	// In production, use proper migrations
	queries := []string{
		`CREATE TABLE IF NOT EXISTS wage_summaries (
			id SERIAL PRIMARY KEY,
			location VARCHAR(50) NOT NULL,
			year INTEGER NOT NULL,
			employee_count INTEGER NOT NULL,
			total_gross_pay DECIMAL(15,2) NOT NULL,
			avg_gross_pay DECIMAL(12,2) NOT NULL,
			median_pay DECIMAL(12,2) NOT NULL,
			std_dev DECIMAL(12,2) NOT NULL,
			min_pay DECIMAL(12,2) NOT NULL,
			max_pay DECIMAL(12,2) NOT NULL,
			percentiles TEXT NOT NULL,
			total_base DECIMAL(15,2) NOT NULL,
			total_overtime DECIMAL(15,2) NOT NULL,
			total_adjustments DECIMAL(15,2) NOT NULL,
			avg_base DECIMAL(12,2) NOT NULL,
			avg_overtime DECIMAL(12,2) NOT NULL,
			avg_adjustments DECIMAL(12,2) NOT NULL,
			generated_at TIMESTAMP NOT NULL,
			uploaded_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(location, year)
		)`,
		`CREATE TABLE IF NOT EXISTS wage_pyramids (
			id SERIAL PRIMARY KEY,
			location VARCHAR(50) NOT NULL,
			year INTEGER NOT NULL,
			total_employees INTEGER NOT NULL,
			total_pay DECIMAL(15,2) NOT NULL,
			brackets TEXT NOT NULL,
			generated_at TIMESTAMP NOT NULL,
			uploaded_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(location, year)
		)`,
		`CREATE TABLE IF NOT EXISTS title_analysis (
			id SERIAL PRIMARY KEY,
			location VARCHAR(50) NOT NULL,
			year INTEGER NOT NULL,
			unique_titles INTEGER NOT NULL,
			top_titles TEXT NOT NULL,
			generated_at TIMESTAMP NOT NULL,
			uploaded_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(location, year)
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}