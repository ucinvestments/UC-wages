# UC Wages Scraper

A high-performance, concurrent Go application for scraping University of California employee wage data from the UC Annual Wage website.

## Data Storage

**All scraped data is saved to the `./data` directory** by default. The structure is:

```
data/
├── ASUCLA/
│   ├── wages_2024.json
│   ├── wages_2023.json
│   └── ...
├── Berkeley/
│   ├── wages_2024.json
│   ├── wages_2023.json
│   └── ...
├── Davis/
│   └── ...
└── ...
```

Each JSON file contains:
- Location name
- Year
- Scrape timestamp
- Total records count
- Array of employee wage records

## Quick Start

### Run with Go
```bash
# Scrape all locations and years (saves to ./data)
go run main.go

# Custom data directory
go run main.go -data-dir=/path/to/save/data

# Specific locations and years
go run main.go -locations="Berkeley,Los Angeles" -years="2023,2024"

# Faster scraping with more workers
go run main.go -workers=10 -delay=0.5
```

### Run with Make
```bash
make run          # Run with defaults (saves to ./data)
make run-fast     # Run with 10 workers
make test         # Test with Berkeley 2024
make clean-data   # Remove all scraped data
```

### Run with Docker
```bash
# Build and run
make docker-build
make docker-run

# Or with docker-compose
docker-compose up
```

Data will be saved to `./data` directory in all cases.

## Progress Tracking & Resume

The scraper automatically tracks progress and can resume from where it left off:

- **Progress file**: `data/scrape_progress.json` tracks completed tasks
- **Automatic resume**: Simply re-run the scraper to continue from where it stopped
- **Skip existing data**: Won't re-scrape locations/years that already have data
- **Progress status**: Use `-status` to see current progress

```bash
# Check progress
make status

# Resume scraping after interruption
make resume

# Start fresh (reset all progress)
make reset
```

## Command-line Options

- `-data-dir`: Directory to save scraped data (default: `./data`)
- `-workers`: Number of concurrent workers (default: 5)
- `-delay`: Delay between requests in seconds (default: 1.0)
- `-locations`: Comma-separated list of locations to scrape
- `-years`: Comma-separated list of years to scrape
- `-status`: Show current progress and exit
- `-reset`: Reset progress and start from beginning

## Available Locations

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

## Years Available

2010-2024