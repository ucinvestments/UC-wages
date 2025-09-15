package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	apiURL      = "https://ucannualwage.ucop.edu/wage/search"
	defaultRows = 100
	timeout     = 30 * time.Second
)

var (
	locations = []string{
		"ASUCLA", "Berkeley", "Davis", "UC SF Law", "Irvine",
		"Los Angeles", "Merced", "Riverside", "San Diego",
		"San Francisco", "Santa Barbara", "Santa Cruz", "UCOP",
	}
	years = []int{
		2024, 2023, 2022, 2021, 2020, 2019, 2018, 2017,
		2016, 2015, 2014, 2013, 2012, 2011, 2010,
	}
)

type Config struct {
	Workers   int
	Delay     time.Duration
	DataDir   string
	Locations []string
	Years     []int
}

type Request struct {
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

type Response struct {
	Records int                      `json:"records"`
	Page    int                      `json:"page"`
	Total   int                      `json:"total"`
	Rows    []map[string]interface{} `json:"rows"`
}

type Output struct {
	Location     string                   `json:"location"`
	Year         int                      `json:"year"`
	ScrapedAt    string                   `json:"scraped_at"`
	TotalRecords int                      `json:"total_records"`
	Records      []map[string]interface{} `json:"records"`
}

type Task struct {
	Location string
	Year     int
}

type Scraper struct {
	client  *http.Client
	config  Config
	success int32
	failed  int32
	mu      sync.Mutex
}

func NewScraper(cfg Config) *Scraper {
	return &Scraper{
		client:  &http.Client{Timeout: timeout},
		config:  cfg,
	}
}

func (s *Scraper) fetchPage(location string, year int, page int) (*Response, error) {
	req := Request{
		Op:       "search",
		Page:     page,
		Rows:     defaultRows,
		Sidx:     "lastname",
		Sord:     "asc",
		Count:    0,
		Year:     fmt.Sprintf("%d", year),
		Location: location,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Post(apiURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response
	return &result, json.Unmarshal(body, &result)
}

func (s *Scraper) fetchAll(location string, year int) ([]map[string]interface{}, error) {
	var records []map[string]interface{}
	page := 1

	for {
		resp, err := s.fetchPage(location, year, page)
		if err != nil {
			return nil, err
		}

		if len(resp.Rows) == 0 {
			break
		}

		records = append(records, resp.Rows...)

		if len(records) >= resp.Records {
			break
		}

		page++
		time.Sleep(s.config.Delay)
	}

	return records, nil
}

func (s *Scraper) save(location string, year int, records []map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := filepath.Join(s.config.DataDir, sanitize(location))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	output := Output{
		Location:     location,
		Year:         year,
		ScrapedAt:    time.Now().Format(time.RFC3339),
		TotalRecords: len(records),
		Records:      records,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(dir, fmt.Sprintf("wages_%d.json", year))
	return os.WriteFile(path, data, 0644)
}

func (s *Scraper) process(task Task) {
	path := filepath.Join(s.config.DataDir, sanitize(task.Location), fmt.Sprintf("wages_%d.json", task.Year))
	if _, err := os.Stat(path); err == nil {
		log.Printf("Skip: %s %d (exists)", task.Location, task.Year)
		atomic.AddInt32(&s.success, 1)
		return
	}

	log.Printf("Fetch: %s %d", task.Location, task.Year)

	records, err := s.fetchAll(task.Location, task.Year)
	if err != nil {
		log.Printf("Error: %s %d - %v", task.Location, task.Year, err)
		atomic.AddInt32(&s.failed, 1)
		return
	}

	if len(records) == 0 {
		log.Printf("No data: %s %d", task.Location, task.Year)
		atomic.AddInt32(&s.failed, 1)
		return
	}

	if err := s.save(task.Location, task.Year, records); err != nil {
		log.Printf("Save error: %s %d - %v", task.Location, task.Year, err)
		atomic.AddInt32(&s.failed, 1)
		return
	}

	log.Printf("Done: %s %d (%d records)", task.Location, task.Year, len(records))
	atomic.AddInt32(&s.success, 1)
}

func (s *Scraper) run() {
	var tasks []Task
	for _, loc := range s.config.Locations {
		for _, year := range s.config.Years {
			tasks = append(tasks, Task{Location: loc, Year: year})
		}
	}

	ch := make(chan Task, len(tasks))
	var wg sync.WaitGroup

	for i := 0; i < s.config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range ch {
				s.process(task)
				time.Sleep(s.config.Delay)
			}
		}()
	}

	for _, task := range tasks {
		ch <- task
	}
	close(ch)

	wg.Wait()

	log.Printf("Complete: %d success, %d failed",
		atomic.LoadInt32(&s.success), atomic.LoadInt32(&s.failed))
}

func sanitize(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, " ", "_"), "/", "_")
}

func parseList(s string, parser func(string) interface{}) []interface{} {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]interface{}, 0, len(parts))
	for _, p := range parts {
		if v := parser(strings.TrimSpace(p)); v != nil {
			result = append(result, v)
		}
	}
	return result
}

func main() {
	var (
		workers   = flag.Int("workers", 5, "concurrent workers")
		delay     = flag.Float64("delay", 1.0, "delay between requests (seconds)")
		dataDir   = flag.String("data", "../data", "data directory")
		locList   = flag.String("locations", "", "locations (comma-separated)")
		yearList  = flag.String("years", "", "years (comma-separated)")
	)
	flag.Parse()

	cfg := Config{
		Workers:   *workers,
		Delay:     time.Duration(*delay * float64(time.Second)),
		DataDir:   *dataDir,
		Locations: locations,
		Years:     years,
	}

	if *locList != "" {
		locs := parseList(*locList, func(s string) interface{} { return s })
		cfg.Locations = make([]string, len(locs))
		for i, l := range locs {
			cfg.Locations[i] = l.(string)
		}
	}

	if *yearList != "" {
		yrs := parseList(*yearList, func(s string) interface{} {
			var y int
			if _, err := fmt.Sscanf(s, "%d", &y); err == nil && y >= 2010 && y <= 2024 {
				return y
			}
			return nil
		})
		cfg.Years = make([]int, len(yrs))
		for i, y := range yrs {
			cfg.Years[i] = y.(int)
		}
	}

	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting: %d workers, %.1fs delay, %d locations, %d years",
		cfg.Workers, cfg.Delay.Seconds(), len(cfg.Locations), len(cfg.Years))

	scraper := NewScraper(cfg)
	scraper.run()
}