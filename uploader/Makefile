.PHONY: deps run test clean upload-all upload-fast upload-status

# Download dependencies
deps:
	go mod tidy
	go mod download

# Build the uploader
build:
	go build -o uploader main.go

# Run with default settings (5 workers, batch size 1000)
run: build
	./uploader

# Run with more workers for faster uploads
upload-fast: build
	./uploader -workers=10 -batch=2000

# Upload specific data directory
upload-dir: build
	./uploader -data=$(DIR)

# Check upload status in database
upload-status:
	@echo "Checking upload progress..."
	@go run -tags tools tools/status.go

# Test with a small subset of data
test: build
	./uploader -data="../data/Berkeley" -workers=2

# Clean build artifacts
clean:
	rm -f uploader

# Initialize database schema only
init-schema:
	@echo "Initializing database schema..."
	@go run -tags tools tools/init_schema.go

# Show help
help:
	@echo "Available commands:"
	@echo "  deps        - Download Go dependencies"
	@echo "  build       - Build the uploader binary"
	@echo "  run         - Run uploader with default settings"
	@echo "  upload-fast - Run with 10 workers and larger batches"
	@echo "  test        - Test upload with Berkeley data only"
	@echo "  clean       - Remove build artifacts"
	@echo "  help        - Show this help message"
	@echo ""
	@echo "Usage examples:"
	@echo "  make run                    # Upload all data with 5 workers"
	@echo "  make upload-fast           # Upload with 10 workers"
	@echo "  make upload-dir DIR=../data/Davis  # Upload specific location"