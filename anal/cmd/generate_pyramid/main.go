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
	outputDir := flag.String("output", "./output/pyramid", "Output directory for pyramids")
	workers := flag.Int("workers", 4, "Number of concurrent workers")
	flag.Parse()

	// Get all JSON files
	files, err := findWageFiles(*dataDir)
	if err != nil {
		log.Fatal("Error finding wage files:", err)
	}

	fmt.Printf("Found %d wage files to process\n", len(files))

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

			if err := processFile(filepath, *outputDir); err != nil {
				errorsChan <- fmt.Errorf("error processing %s: %w", filepath, err)
			} else {
				fmt.Printf("✓ Generated pyramid for %s\n", filepath)
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
		fmt.Println("\n✅ All pyramids generated successfully!")
		printBracketInfo()
	}
}

func processFile(filepath, outputDir string) error {
	// Load wage data
	data, err := parser.LoadWageData(filepath)
	if err != nil {
		return err
	}

	// Generate pyramid
	pyramid, err := calculator.CalculatePyramid(data)
	if err != nil {
		return err
	}

	if pyramid == nil {
		return fmt.Errorf("no valid wage data found")
	}

	// Generate output filename
	filename := fmt.Sprintf("%s_%d.json",
		strings.ReplaceAll(data.Location, " ", "_"),
		data.Year)
	outputPath := fmt.Sprintf("%s/%s", outputDir, filename)

	// Save pyramid
	if err := parser.SaveJSON(outputPath, pyramid); err != nil {
		return err
	}

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

func printBracketInfo() {
	fmt.Println("\nWage Brackets Used:")
	brackets := calculator.GetWageBrackets()
	for _, bracket := range brackets {
		fmt.Printf("  • %s: $%.0f - $%.0f\n", bracket.Range, bracket.MinValue, bracket.MaxValue)
	}
}