.PHONY: build run test clean help

# Default target
help:
	@echo "UC Wages Scraper - Available commands:"
	@echo "  make run       - Run scraper with default settings"
	@echo "  make fast      - Run with 10 workers for faster scraping"
	@echo "  make test      - Test with Berkeley 2024 data"
	@echo "  make build     - Build Go binary"
	@echo "  make clean     - Remove binary and clean up"

# Build binary
build:
	go build -o scraper main.go

# Run with defaults
run:
	go run main.go

# Fast mode with more workers
fast:
	go run main.go -workers=10 -delay=0.5

# Test with single location/year
test:
	go run main.go -locations="Berkeley" -years="2024" -workers=2

# Clean up
clean:
	rm -f scraper
	go clean