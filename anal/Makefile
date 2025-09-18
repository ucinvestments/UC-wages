# UC Wages Analysis Makefile

.PHONY: all build clean run-sums run-pyramid run-titles run-all test deps

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_SUMS=calculate_sums
BINARY_PYRAMID=generate_pyramid
BINARY_TITLES=analyze_titles
BINARY_ALL=run_all

# Directories
CMD_DIR=cmd
OUTPUT_DIR=output
DATA_DIR=../data

all: deps build

deps:
	$(GOMOD) download
	$(GOGET) github.com/montanaflynn/stats

build: build-sums build-pyramid build-titles build-all

build-sums:
	cd $(CMD_DIR)/calculate_sums && $(GOBUILD) -o $(BINARY_SUMS) -v

build-pyramid:
	cd $(CMD_DIR)/generate_pyramid && $(GOBUILD) -o $(BINARY_PYRAMID) -v

build-titles:
	cd $(CMD_DIR)/analyze_titles && $(GOBUILD) -o $(BINARY_TITLES) -v

build-all:
	cd $(CMD_DIR)/run_all && $(GOBUILD) -o $(BINARY_ALL) -v

clean:
	$(GOCLEAN)
	rm -f $(CMD_DIR)/calculate_sums/$(BINARY_SUMS)
	rm -f $(CMD_DIR)/generate_pyramid/$(BINARY_PYRAMID)
	rm -f $(CMD_DIR)/analyze_titles/$(BINARY_TITLES)
	rm -f $(CMD_DIR)/run_all/$(BINARY_ALL)
	rm -rf $(OUTPUT_DIR)

run-sums: build-sums
	cd $(CMD_DIR)/calculate_sums && ./$(BINARY_SUMS) -data $(DATA_DIR) -output ../../$(OUTPUT_DIR)/sums -workers 8

run-pyramid: build-pyramid
	cd $(CMD_DIR)/generate_pyramid && ./$(BINARY_PYRAMID) -data $(DATA_DIR) -output ../../$(OUTPUT_DIR)/pyramid -workers 8

run-titles: build-titles
	cd $(CMD_DIR)/analyze_titles && ./$(BINARY_TITLES) -data $(DATA_DIR) -output ../../$(OUTPUT_DIR)/titles -top 100 -workers 8

run-all: build
	cd $(CMD_DIR)/run_all && ./$(BINARY_ALL) -data $(DATA_DIR) -output ../../$(OUTPUT_DIR) -workers 8

test:
	$(GOTEST) -v ./...

# Test with single file
test-berkeley:
	@mkdir -p $(OUTPUT_DIR)/test
	cd $(CMD_DIR)/calculate_sums && $(GOBUILD) -o $(BINARY_SUMS) && \
		./$(BINARY_SUMS) -data ../../data/Berkeley -output ../../$(OUTPUT_DIR)/test/sums -workers 1
	cd $(CMD_DIR)/generate_pyramid && $(GOBUILD) -o $(BINARY_PYRAMID) && \
		./$(BINARY_PYRAMID) -data ../../data/Berkeley -output ../../$(OUTPUT_DIR)/test/pyramid -workers 1
	cd $(CMD_DIR)/analyze_titles && $(GOBUILD) -o $(BINARY_TITLES) && \
		./$(BINARY_TITLES) -data ../../data/Berkeley -output ../../$(OUTPUT_DIR)/test/titles -workers 1

# View output samples
view-sums:
	@echo "Sample Summary Output:"
	@head -50 $(OUTPUT_DIR)/sums/Berkeley_2024.json 2>/dev/null || echo "Run 'make run-sums' first"

view-pyramid:
	@echo "Sample Pyramid Output:"
	@head -100 $(OUTPUT_DIR)/pyramid/Berkeley_2024.json 2>/dev/null || echo "Run 'make run-pyramid' first"

view-titles:
	@echo "Sample Titles Output:"
	@head -50 $(OUTPUT_DIR)/titles/Berkeley_2024.json 2>/dev/null || echo "Run 'make run-titles' first"

help:
	@echo "UC Wages Analysis - Makefile commands:"
	@echo ""
	@echo "  make all          - Download dependencies and build all binaries"
	@echo "  make build        - Build all analysis binaries"
	@echo "  make clean        - Clean build files and output"
	@echo ""
	@echo "  make run-sums     - Calculate summary statistics"
	@echo "  make run-pyramid  - Generate wage pyramids"
	@echo "  make run-titles   - Analyze job titles"
	@echo "  make run-all      - Run complete analysis pipeline"
	@echo ""
	@echo "  make test         - Run tests"
	@echo "  make test-berkeley - Test with Berkeley data only"
	@echo ""
	@echo "  make view-sums    - View sample summary output"
	@echo "  make view-pyramid - View sample pyramid output"
	@echo "  make view-titles  - View sample titles output"