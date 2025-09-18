package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	// Command line flags
	dataDir := flag.String("data", "../../data", "Path to data directory")
	outputDir := flag.String("output", "./output", "Base output directory")
	workers := flag.Int("workers", 4, "Number of concurrent workers")
	flag.Parse()

	fmt.Println("🚀 UC Wages Analysis Pipeline")
	fmt.Println("================================")
	start := time.Now()

	// Create output directories
	dirs := []string{
		*outputDir,
		fmt.Sprintf("%s/sums", *outputDir),
		fmt.Sprintf("%s/pyramid", *outputDir),
		fmt.Sprintf("%s/titles", *outputDir),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal("Error creating directory:", err)
		}
	}

	// Run each analysis
	analyses := []struct {
		name    string
		command string
		args    []string
	}{
		{
			name:    "Summary Statistics",
			command: "calculate_sums",
			args:    []string{"-data", *dataDir, "-output", fmt.Sprintf("%s/sums", *outputDir), "-workers", fmt.Sprintf("%d", *workers)},
		},
		{
			name:    "Wage Pyramids",
			command: "generate_pyramid",
			args:    []string{"-data", *dataDir, "-output", fmt.Sprintf("%s/pyramid", *outputDir), "-workers", fmt.Sprintf("%d", *workers)},
		},
		{
			name:    "Title Analysis",
			command: "analyze_titles",
			args:    []string{"-data", *dataDir, "-output", fmt.Sprintf("%s/titles", *outputDir), "-workers", fmt.Sprintf("%d", *workers), "-top", "100"},
		},
	}

	for _, analysis := range analyses {
		fmt.Printf("\n▶ Running %s...\n", analysis.name)
		cmd := exec.Command(analysis.command, analysis.args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Error running %s: %v", analysis.name, err)
		} else {
			fmt.Printf("✅ %s completed\n", analysis.name)
		}
	}

	// Generate summary report
	generateReport(*outputDir)

	elapsed := time.Since(start)
	fmt.Printf("\n✨ Analysis pipeline completed in %s\n", elapsed.Round(time.Second))
}

func generateReport(outputDir string) {
	fmt.Println("\n📊 Analysis Report")
	fmt.Println("==================")

	// Count output files
	analyses := []string{"sums", "pyramid", "titles"}
	for _, dir := range analyses {
		path := fmt.Sprintf("%s/%s", outputDir, dir)
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		count := 0
		for _, file := range files {
			if !file.IsDir() {
				count++
			}
		}
		fmt.Printf("• %s: %d files generated\n", dir, count)
	}

	fmt.Println("\nOutput structure:")
	fmt.Printf("%s/\n", outputDir)
	fmt.Println("├── sums/       # Statistical summaries per location-year")
	fmt.Println("├── pyramid/    # Wage distribution pyramids")
	fmt.Println("└── titles/     # Job title analysis")
}