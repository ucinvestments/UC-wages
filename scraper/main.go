package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const scrapeURL = "https://ucannualwage.ucop.edu/wage/search"

var locations = []string{
	"ASUCLA",
	"Berkeley",
	"Davis",
	"UC SF Law",
	"Irvine",
	"Los Angeles",
	"Merced",
	"Riverside",
	"San Diego",
	"San Francisco",
	"Santa Barbara",
	"Santa Cruz",
	"UCOP",
}

var years = []int{
	2024, 2023, 2022, 2021, 2020,
	2019, 2018, 2017, 2016, 2015,
	2014, 2013, 2012, 2011, 2010,
}

// SearchPayload represents the API request payload
type SearchPayload struct {
	Op       string `json:"op"`
	Page     int    `json:"page"`
	Rows     int    `json:"rows"`
	Sidx     string `json:"sidx"`
	Sord     string `json:"sord"`
	Count    int    `json:"count"`
	Year     string `json:"year"`
	Firstname string `json:"firstname"`
	Location string `json:"location"`
	Lastname string `json:"lastname"`
	Title    string `json:"title"`
	StartSal string `json:"startSal"`
	EndSal   string `json:"endSal"`
}

// APIResponse represents the API response structure
type APIResponse struct {
	Records int                      `json:"records"`
	Page    int                      `json:"page"`
	Total   int                      `json:"total"`
	Rows    []map[string]interface{} `json:"rows"`
}

// OutputData represents the structure of saved JSON files
type OutputData struct {
	Location     string                   `json:"location"`
	Year         int                      `json:"year"`
	ScrapedAt    string                   `json:"scraped_at"`
	TotalRecords int                      `json:"total_records"`
	Records      []map[string]interface{} `json:"records"`
}

// Task represents a scraping task
type Task struct {
	Location string
	Year     int
}

// Progress represents the scraping progress
type Progress struct {
	CompletedTasks map[string]bool `json:"completed_tasks"`
	StartTime      string          `json:"start_time"`
	LastUpdated    string          `json:"last_updated"`
	TotalTasks     int             `json:"total_tasks"`
	CompletedCount int             `json:"completed_count"`
}

// Scraper handles the scraping operations
type Scraper struct {
	client       *http.Client
	delay        time.Duration
	workers      int
	successful   int32
	failed       int32
	fileMutex    sync.Mutex
	dataDir      string
	progress     *Progress
	progressFile string
	progressMutex sync.Mutex
}

// NewScraper creates a new scraper instance
func NewScraper(workers int, delaySeconds float64, dataDir string) *Scraper {
	progressFile := filepath.Join(dataDir, "scrape_progress.json")
	scraper := &Scraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		delay:        time.Duration(delaySeconds * float64(time.Second)),
		workers:      workers,
		dataDir:      dataDir,
		progressFile: progressFile,
	}

	// Load existing progress or create new
	scraper.loadProgress()
	return scraper
}

// resetProgress clears all progress and starts fresh
func (s *Scraper) resetProgress() {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()

	s.progress = &Progress{
		CompletedTasks: make(map[string]bool),
		StartTime:      time.Now().Format(time.RFC3339),
		LastUpdated:    time.Now().Format(time.RFC3339),
		TotalTasks:     0,
		CompletedCount: 0,
	}

	// Remove progress file
	os.Remove(s.progressFile)
	fmt.Println("Progress reset successfully")
}

// showProgressStatus displays current progress information
func (s *Scraper) showProgressStatus() {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()

	fmt.Printf("\n📊 Current Progress Status\n")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Total tasks: %d\n", s.progress.TotalTasks)
	fmt.Printf("Completed: %d\n", s.progress.CompletedCount)
	fmt.Printf("Remaining: %d\n", s.progress.TotalTasks-s.progress.CompletedCount)
	if s.progress.TotalTasks > 0 {
		percentage := float64(s.progress.CompletedCount) / float64(s.progress.TotalTasks) * 100
		fmt.Printf("Progress: %.1f%%\n", percentage)
	}
	fmt.Printf("Started: %s\n", s.progress.StartTime)
	fmt.Printf("Last updated: %s\n", s.progress.LastUpdated)

	if s.progress.CompletedCount > 0 {
		fmt.Printf("\n📝 Recently completed tasks:\n")

		// Convert map to sorted slice for display
		var completedTasks []string
		for task := range s.progress.CompletedTasks {
			completedTasks = append(completedTasks, task)
		}
		sort.Strings(completedTasks)

		// Show last 10 completed tasks
		start := len(completedTasks) - 10
		if start < 0 {
			start = 0
		}
		for _, task := range completedTasks[start:] {
			fmt.Printf("  ✓ %s\n", task)
		}
	}
}

// loadProgress loads existing progress or creates new progress tracking
func (s *Scraper) loadProgress() {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()

	if _, err := os.Stat(s.progressFile); os.IsNotExist(err) {
		// Create new progress
		s.progress = &Progress{
			CompletedTasks: make(map[string]bool),
			StartTime:      time.Now().Format(time.RFC3339),
			LastUpdated:    time.Now().Format(time.RFC3339),
			TotalTasks:     0,
			CompletedCount: 0,
		}
		return
	}

	// Load existing progress
	data, err := os.ReadFile(s.progressFile)
	if err != nil {
		fmt.Printf("Warning: Could not read progress file: %v\n", err)
		s.progress = &Progress{
			CompletedTasks: make(map[string]bool),
			StartTime:      time.Now().Format(time.RFC3339),
			LastUpdated:    time.Now().Format(time.RFC3339),
		}
		return
	}

	var progress Progress
	if err := json.Unmarshal(data, &progress); err != nil {
		fmt.Printf("Warning: Could not parse progress file: %v\n", err)
		s.progress = &Progress{
			CompletedTasks: make(map[string]bool),
			StartTime:      time.Now().Format(time.RFC3339),
			LastUpdated:    time.Now().Format(time.RFC3339),
		}
		return
	}

	s.progress = &progress
	fmt.Printf("Loaded progress: %d/%d tasks completed\n", s.progress.CompletedCount, s.progress.TotalTasks)
}

// saveProgress saves the current progress to file
func (s *Scraper) saveProgress() error {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()

	s.progress.LastUpdated = time.Now().Format(time.RFC3339)

	data, err := json.MarshalIndent(s.progress, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling progress: %v", err)
	}

	return os.WriteFile(s.progressFile, data, 0644)
}

// getTaskKey generates a unique key for a location-year combination
func getTaskKey(location string, year int) string {
	return fmt.Sprintf("%s-%d", location, year)
}

// isTaskCompleted checks if a task has already been completed
func (s *Scraper) isTaskCompleted(location string, year int) bool {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()

	key := getTaskKey(location, year)
	return s.progress.CompletedTasks[key]
}

// markTaskCompleted marks a task as completed and saves progress
func (s *Scraper) markTaskCompleted(location string, year int) {
	s.progressMutex.Lock()
	key := getTaskKey(location, year)
	if !s.progress.CompletedTasks[key] {
		s.progress.CompletedTasks[key] = true
		s.progress.CompletedCount++
	}
	s.progressMutex.Unlock()

	// Save progress after each completed task
	if err := s.saveProgress(); err != nil {
		fmt.Printf("Warning: Could not save progress: %v\n", err)
	}
}

// doesDataExist checks if data file already exists for a location-year combination
func (s *Scraper) doesDataExist(location string, year int) bool {
	locationDir := strings.ReplaceAll(strings.ReplaceAll(location, " ", "_"), "/", "_")
	fullDirPath := filepath.Join(s.dataDir, locationDir)
	filename := filepath.Join(fullDirPath, fmt.Sprintf("wages_%d.json", year))

	if _, err := os.Stat(filename); err == nil {
		// File exists, check if it has valid data
		data, err := os.ReadFile(filename)
		if err != nil {
			return false
		}

		var outputData OutputData
		if err := json.Unmarshal(data, &outputData); err != nil {
			return false
		}

		// Consider it exists if it has records
		return outputData.TotalRecords > 0
	}

	return false
}

// createPayload creates the search payload
func (s *Scraper) createPayload(location string, year int, page int) SearchPayload {
	return SearchPayload{
		Op:       "search",
		Page:     page,
		Rows:     100,
		Sidx:     "lastname",
		Sord:     "asc",
		Count:    0,
		Year:     fmt.Sprintf("%d", year),
		Firstname: "",
		Location: location,
		Lastname: "",
		Title:    "",
		StartSal: "",
		EndSal:   "",
	}
}

// fetchPage fetches a single page of data
func (s *Scraper) fetchPage(location string, year int, page int) (*APIResponse, error) {
	payload := s.createPayload(location, year, page)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", scrapeURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	return &apiResp, nil
}

// fetchAllPages fetches all pages for a location and year
func (s *Scraper) fetchAllPages(location string, year int) ([]map[string]interface{}, error) {
	var allRecords []map[string]interface{}
	page := 1

	for {
		fmt.Printf("  [%s-%d] Fetching page %d...\n", location, year, page)

		apiResp, err := s.fetchPage(location, year, page)
		if err != nil {
			return nil, err
		}

		if apiResp.Rows == nil || len(apiResp.Rows) == 0 {
			break
		}

		allRecords = append(allRecords, apiResp.Rows...)

		// Check if we've fetched all records
		if len(allRecords) >= apiResp.Records || len(apiResp.Rows) == 0 {
			break
		}

		page++
		time.Sleep(s.delay)
	}

	return allRecords, nil
}

// saveData saves the scraped data to a JSON file
func (s *Scraper) saveData(location string, year int, records []map[string]interface{}) error {
	s.fileMutex.Lock()
	defer s.fileMutex.Unlock()

	// Create directory for the location
	dirName := strings.ReplaceAll(strings.ReplaceAll(location, " ", "_"), "/", "_")
	fullDirPath := filepath.Join(s.dataDir, dirName)
	if err := os.MkdirAll(fullDirPath, 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Create the output data structure
	outputData := OutputData{
		Location:     location,
		Year:         year,
		ScrapedAt:    time.Now().Format(time.RFC3339),
		TotalRecords: len(records),
		Records:      records,
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling output data: %v", err)
	}

	// Write to file
	filename := filepath.Join(fullDirPath, fmt.Sprintf("wages_%d.json", year))
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Printf("  ✓ Saved %d records to %s\n", len(records), filename)
	return nil
}

// processTask processes a single scraping task
func (s *Scraper) processTask(task Task) {
	// Check if task is already completed
	if s.isTaskCompleted(task.Location, task.Year) {
		fmt.Printf("[Worker] Skipping %s %d (already completed)\n", task.Location, task.Year)
		atomic.AddInt32(&s.successful, 1)
		return
	}

	// Check if data already exists
	if s.doesDataExist(task.Location, task.Year) {
		fmt.Printf("[Worker] Skipping %s %d (data already exists)\n", task.Location, task.Year)
		s.markTaskCompleted(task.Location, task.Year)
		atomic.AddInt32(&s.successful, 1)
		return
	}

	fmt.Printf("[Worker] Starting: %s for year %d...\n", task.Location, task.Year)

	records, err := s.fetchAllPages(task.Location, task.Year)
	if err != nil {
		fmt.Printf("  ✗ Error scraping %s %d: %v\n", task.Location, task.Year, err)
		atomic.AddInt32(&s.failed, 1)
		return
	}

	if len(records) == 0 {
		fmt.Printf("  ✗ No data found for %s %d\n", task.Location, task.Year)
		atomic.AddInt32(&s.failed, 1)
		return
	}

	if err := s.saveData(task.Location, task.Year, records); err != nil {
		fmt.Printf("  ✗ Error saving data for %s %d: %v\n", task.Location, task.Year, err)
		atomic.AddInt32(&s.failed, 1)
		return
	}

	fmt.Printf("  ✓ Successfully scraped %d records for %s %d\n", len(records), task.Location, task.Year)

	// Mark task as completed
	s.markTaskCompleted(task.Location, task.Year)
	atomic.AddInt32(&s.successful, 1)
}

// scrapeAll runs the scraper with worker pool
func (s *Scraper) scrapeAll(locationsToScrape []string, yearsToScrape []int) {
	// Create all possible tasks
	var allTasks []Task
	for _, location := range locationsToScrape {
		for _, year := range yearsToScrape {
			allTasks = append(allTasks, Task{Location: location, Year: year})
		}
	}

	totalTasks := len(allTasks)

	// Initialize progress if this is the first run
	if s.progress.TotalTasks == 0 {
		s.progress.TotalTasks = totalTasks
		s.saveProgress()
	}

	// Filter out already completed tasks for efficient processing
	var pendingTasks []Task
	skippedCount := 0
	for _, task := range allTasks {
		if s.isTaskCompleted(task.Location, task.Year) || s.doesDataExist(task.Location, task.Year) {
			skippedCount++
			continue
		}
		pendingTasks = append(pendingTasks, task)
	}

	// If we found existing data files that aren't marked as completed, mark them now
	if skippedCount > s.progress.CompletedCount {
		fmt.Printf("Found %d existing data files, updating progress...\n", skippedCount-s.progress.CompletedCount)
		for _, task := range allTasks {
			if s.doesDataExist(task.Location, task.Year) && !s.isTaskCompleted(task.Location, task.Year) {
				s.markTaskCompleted(task.Location, task.Year)
			}
		}
	}

	fmt.Printf("\nUC Wage Data Scraper - Go Concurrent Version\n")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total tasks: %d\n", totalTasks)
	fmt.Printf("Already completed: %d\n", skippedCount)
	fmt.Printf("Remaining tasks: %d\n", len(pendingTasks))
	fmt.Printf("Using %d concurrent workers\n", s.workers)
	fmt.Println(strings.Repeat("=", 60))

	if len(pendingTasks) == 0 {
		fmt.Println("\n🎉 All tasks already completed! No work to do.")
		return
	}

	startTime := time.Now()

	// Create task channel and wait group
	taskChan := make(chan Task, len(pendingTasks))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range taskChan {
				s.processTask(task)
				time.Sleep(s.delay) // Rate limiting between tasks
			}
		}(i)
	}

	// Send pending tasks to channel
	for _, task := range pendingTasks {
		taskChan <- task
	}
	close(taskChan)

	// Wait for all workers to complete
	wg.Wait()

	elapsed := time.Since(startTime)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("Scraping complete in %.2f seconds!\n", elapsed.Seconds())
	fmt.Printf("  ✓ Successful: %d\n", atomic.LoadInt32(&s.successful))
	fmt.Printf("  ✗ Failed: %d\n", atomic.LoadInt32(&s.failed))
	fmt.Printf("  Total: %d\n", totalTasks)
	fmt.Printf("  Average time per item: %.2fs\n", elapsed.Seconds()/float64(totalTasks))
}

func main() {
	// Parse command-line arguments
	workers := flag.Int("workers", 5, "Number of concurrent workers")
	delay := flag.Float64("delay", 1.0, "Delay between requests in seconds")
	dataDir := flag.String("data-dir", "./data", "Directory to save scraped data")
	specificLocations := flag.String("locations", "", "Comma-separated list of specific locations to scrape")
	specificYears := flag.String("years", "", "Comma-separated list of specific years to scrape")
	resetProgress := flag.Bool("reset", false, "Reset progress and start from beginning")
	showProgress := flag.Bool("status", false, "Show current progress and exit")
	flag.Parse()

	// Parse specific locations if provided
	locationsToScrape := locations
	if *specificLocations != "" {
		locationsToScrape = strings.Split(*specificLocations, ",")
		for i := range locationsToScrape {
			locationsToScrape[i] = strings.TrimSpace(locationsToScrape[i])
		}
	}

	// Parse specific years if provided
	yearsToScrape := years
	if *specificYears != "" {
		yearStrs := strings.Split(*specificYears, ",")
		yearsToScrape = []int{}
		for _, yearStr := range yearStrs {
			var year int
			if _, err := fmt.Sscanf(strings.TrimSpace(yearStr), "%d", &year); err == nil {
				if year >= 2010 && year <= 2024 {
					yearsToScrape = append(yearsToScrape, year)
				}
			}
		}
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	// Create scraper
	scraper := NewScraper(*workers, *delay, *dataDir)

	// Handle reset flag
	if *resetProgress {
		scraper.resetProgress()
		if *showProgress {
			scraper.showProgressStatus()
		}
		return
	}

	// Handle status flag
	if *showProgress {
		scraper.showProgressStatus()
		return
	}

	fmt.Println("UC Wage Data Scraper")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Workers: %d\n", *workers)
	fmt.Printf("  Delay: %.1fs\n", *delay)
	fmt.Printf("  Data Directory: %s\n", *dataDir)
	fmt.Printf("  Locations: %v\n", locationsToScrape)
	fmt.Printf("  Years: %v\n", yearsToScrape)
	fmt.Printf("  Progress File: %s\n", filepath.Join(*dataDir, "scrape_progress.json"))
	fmt.Println(strings.Repeat("=", 60))

	// Run scraper
	scraper.scrapeAll(locationsToScrape, yearsToScrape)
}