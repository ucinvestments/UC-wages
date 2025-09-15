.PHONY: build run docker-build docker-run clean test

# Build the Go application
build:
	go build -o scraper main.go

# Run the scraper locally
run:
	go run main.go -data-dir=./data

# Run with custom parameters
run-fast:
	go run main.go -workers=10 -delay=0.5 -data-dir=./data

# Run for specific locations
run-berkeley:
	go run main.go -locations="Berkeley" -years="2023,2024" -data-dir=./data

# Show current progress
status:
	go run main.go -data-dir=./data -status

# Reset progress and start fresh
reset:
	go run main.go -data-dir=./data -reset

# Resume scraping (same as run, but will skip completed tasks)
resume:
	go run main.go -data-dir=./data

# Build Docker image
docker-build:
	docker build -t uc-wage-scraper .

# Run Docker container
docker-run:
	docker run -v $$(pwd)/data:/app/data uc-wage-scraper -data-dir=/app/data

# Run Docker container to check status
docker-status:
	docker run -v $$(pwd)/data:/app/data uc-wage-scraper -data-dir=/app/data -status

# Reset progress in Docker
docker-reset:
	docker run -v $$(pwd)/data:/app/data uc-wage-scraper -data-dir=/app/data -reset

# Run with docker-compose
compose-up:
	docker-compose up

# Clean up generated files
clean:
	rm -f scraper
	rm -rf data/

# Clean only scraped data
clean-data:
	rm -rf data/

# Test the scraper with a small subset
test:
	go run main.go -locations="Berkeley" -years="2024" -workers=2 -data-dir=./data