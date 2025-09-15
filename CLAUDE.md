# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

UC Wages Data Repository - A Go-based high-performance scraper for collecting University of California employee wage data from the UC Annual Wage website (https://ucannualwage.ucop.edu/wage/).

## Architecture

### Core Components
- **scraper/main.go**: Concurrent Go scraper with worker pool pattern
  - Multi-threaded data collection using goroutines
  - Progress tracking with resume capability
  - Rate limiting and retry logic
  - JSON data persistence

### Data Structure
```
data/
├── [Campus_Name]/        # One directory per UC location
│   ├── wages_YYYY.json  # Annual wage data files
│   └── ...
└── scrape_progress.json # Progress tracking for resume capability
```

Each wage file contains:
- Location, year, and timestamp
- Total records count
- Array of employee wage records

## Commands

### Running the Scraper
```bash
# From scraper directory
make run                 # Run with default settings
make run-fast           # Run with 10 workers
make test               # Test with Berkeley 2024 data
make status             # Check scraping progress
make resume             # Resume interrupted scraping
make reset              # Reset progress and start fresh
make clean-data         # Remove all scraped data

# Direct Go commands
go run scraper/main.go -workers=5 -delay=1.0 -data-dir=./data
go run scraper/main.go -locations="Berkeley,Irvine" -years="2023,2024"
```

### Docker Support
```bash
make docker-build       # Build Docker image
make docker-run         # Run scraper in Docker
docker-compose up       # Run with docker-compose
```

## UC Locations Available
- ASUCLA
- Berkeley
- Davis
- UC SF Law
- Irvine
- Los Angeles
- Merced
- Riverside
- San Diego
- San Francisco
- Santa Barbara
- Santa Cruz
- UCOP

## Years Covered
2010-2024 (15 years of data per location)

## Progress Tracking
The scraper maintains progress in `data/scrape_progress.json`:
- Automatically resumes interrupted scraping
- Skips already completed location-year combinations
- Tracks completion status and timestamps