package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ucinvestments/uc-wages-analysis/pkg/calculator"
	"github.com/ucinvestments/uc-wages-analysis/pkg/parser"
)

func main() {
	// Command line flags
	dataDir := flag.String("data", "../../data", "Path to data directory")
	outputDir := flag.String("output", "./output/titles", "Output directory for title analysis")
	topN := flag.Int("top", 100, "Number of top titles to include")
	workers := flag.Int("workers", 4, "Number of concurrent workers")
	flag.Parse()

	// Get all JSON files
	files, err := findWageFiles(*dataDir)
	if err != nil {
		log.Fatal("Error finding wage files:", err)
	}

	fmt.Printf("Found %d wage files to process\n", len(files))
	fmt.Printf("Will extract top %d titles per file\n", *topN)

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatal("Error creating output directory:", err)
	}

	// Process files concurrently
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, *workers)
	errorsChan := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		go func(filepath string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := processFile(filepath, *outputDir, *topN); err != nil {
				errorsChan <- fmt.Errorf("error processing %s: %w", filepath, err)
			} else {
				fmt.Printf("✓ Analyzed titles for %s\n", filepath)
			}
		}(file)
	}

	wg.Wait()
	close(errorsChan)

	// Report errors
	var hasErrors bool
	for err := range errorsChan {
		log.Println(err)
		hasErrors = true
	}

	if !hasErrors {
		fmt.Println("\n✅ All title analyses completed successfully!")
	}
}

func processFile(filepath, outputDir string, topN int) error {
	// Load wage data
	data, err := parser.LoadWageData(filepath)
	if err != nil {
		return err
	}

	// Analyze titles
	analysis, err := calculator.AnalyzeTitles(data, topN)
	if err != nil {
		return err
	}

	if analysis == nil {
		return fmt.Errorf("no valid title data found")
	}

	// Generate output filename
	filename := fmt.Sprintf("%s_%d.json",
		strings.ReplaceAll(data.Location, " ", "_"),
		data.Year)
	outputPath := fmt.Sprintf("%s/%s", outputDir, filename)

	// Save analysis
	if err := parser.SaveJSON(outputPath, analysis); err != nil {
		return err
	}

	// Print summary
	fmt.Printf("  - %s %d: %d unique titles, %d in top list\n",
		data.Location, data.Year,
		analysis.UniqueTitles,
		len(analysis.TopTitles))

	return nil
}

func findWageFiles(dataDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-wage JSON files
		if strings.HasSuffix(path, ".json") &&
		   strings.Contains(path, "wages_") &&
		   !strings.Contains(path, "scrape_progress") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}