# UC Wages Analysis Pipeline

High-performance Go-based analysis tools for preprocessing University of California wage data. This pipeline processes 1.3GB of raw JSON data to generate precalculated metrics, reducing database query times from seconds to milliseconds.

## Features

- **Statistical Summaries**: Employee counts, averages, medians, percentiles, and standard deviations
- **Wage Pyramids**: Income distribution across 10 logarithmic brackets with job title analysis
- **Title Analysis**: Top job titles with pay statistics and frequency counts
- **Concurrent Processing**: Utilizes Go routines for efficient parallel processing
- **Database Integration**: Direct upload to PostgreSQL with conflict resolution
- **Docker Support**: Containerized for consistent deployment

## Quick Start

### Using Docker (Recommended)

```bash
# Build the analysis pipeline
docker build -t uc-wages-analysis .

# Run complete analysis
docker run --rm \
  -v /path/to/your/data:/data:ro \
  -v $(pwd)/output:/app/output \
  uc-wages-analysis

# Run specific analysis
docker run --rm \
  -v /path/to/your/data:/data:ro \
  -v $(pwd)/output:/app/output \
  uc-wages-analysis /bin/calculate_sums -data /data -output /app/output/sums -workers 8
```

### Local Development

```bash
# Install dependencies
go mod download

# Build all tools
go build -o bin/calculate_sums ./cmd/calculate_sums/
go build -o bin/generate_pyramid ./cmd/generate_pyramid/
go build -o bin/analyze_titles ./cmd/analyze_titles/
go build -o bin/run_all ./cmd/run_all/
go build -o bin/upload_analysis ./cmd/upload_analysis/

# Run analysis pipeline
./bin/run_all -data ../data -output ./output -workers 8
```

## Analysis Tools

### 1. Summary Statistics (`calculate_sums`)

Generates comprehensive statistical summaries for each location-year combination:

- **Employee Count**: Total number of employees
- **Wage Statistics**: Mean, median, min, max, standard deviation
- **Percentiles**: 25th, 50th, 75th, 90th, 95th, 99th
- **Pay Components**: Base, overtime, and adjustment totals/averages

**Output**: `output/sums/[Location]_[Year].json`

```bash
docker run --rm \
  -v /path/to/data:/data:ro \
  -v $(pwd)/output:/app/output \
  uc-wages-analysis /bin/calculate_sums -data /data -output /app/output/sums -workers 8
```

### 2. Wage Pyramids (`generate_pyramid`)

Creates income distribution analysis with job title breakdowns:

**Wage Brackets**:
- < $25k, $25k-50k, $50k-75k, $75k-100k
- $100k-150k, $150k-200k, $200k-300k
- $300k-500k, $500k-1M, $1M+

**Per Bracket**: Count, percentage, average pay, top 10 job titles

**Output**: `output/pyramid/[Location]_[Year].json`

```bash
docker run --rm \
  -v /path/to/data:/data:ro \
  -v $(pwd)/output:/app/output \
  uc-wages-analysis /bin/generate_pyramid -data /data -output /app/output/pyramid -workers 8
```

### 3. Title Analysis (`analyze_titles`)

Analyzes job title patterns and compensation:

- **Top N Titles**: Most common job titles (default: 100)
- **Per Title Stats**: Count, average, median, min, max, standard deviation
- **Total Pay**: Sum of all compensation for each title

**Output**: `output/titles/[Location]_[Year].json`

```bash
docker run --rm \
  -v /path/to/data:/data:ro \
  -v $(pwd)/output:/app/output \
  uc-wages-analysis /bin/analyze_titles -data /data -output /app/output/titles -top 100 -workers 8
```

### 4. Database Upload (`upload_analysis`)

Uploads precalculated analysis to PostgreSQL database:

- **Conflict Resolution**: Uses UPSERT for data consistency
- **Batch Processing**: Efficient bulk uploads
- **JSON Storage**: Complex data structures stored as JSON
- **Auto-Schema**: Creates tables if they don't exist

```bash
docker run --rm \
  -v $(pwd)/output:/app/output \
  uc-wages-analysis /bin/upload_analysis \
  -db "postgresql://user:pass@host:5432/dbname" \
  -output /app/output
```

## Performance Benefits

### Before (Raw Data Queries)
- **Data Transfer**: ~1.3GB per complex query
- **Processing Time**: 10-30 seconds for aggregations
- **Memory Usage**: High due to large dataset processing
- **Scalability**: Limited by real-time calculation overhead

### After (Precalculated Analytics)
- **Data Transfer**: ~10MB for comprehensive metrics
- **Processing Time**: <100ms for complex visualizations
- **Memory Usage**: Minimal - pre-aggregated data
- **Scalability**: Supports real-time dashboards and APIs

## Output Format Examples

### Summary Statistics
```json
{
  "location": "Berkeley",
  "year": 2024,
  "employee_count": 36673,
  "total_gross_pay": 1803384251,
  "avg_gross_pay": 49174.71,
  "median_gross_pay": 17898,
  "percentiles": {
    "p25": 3850,
    "p50": 17896.5,
    "p75": 72786,
    "p90": 132141,
    "p95": 180289.5,
    "p99": 339073.5
  },
  "pay_components": {
    "total_base": 1650836310,
    "avg_overtime": 361.14,
    "avg_adjustments": 3798.54
  }
}
```

### Wage Pyramid
```json
{
  "location": "Berkeley",
  "year": 2024,
  "total_employees": 36673,
  "brackets": [
    {
      "range": "0-25k",
      "count": 20223,
      "percentage": 55.14,
      "avg_pay": 6926.18,
      "top_titles": [
        {"title": "STDT 2", "count": 7092, "avg_pay": 4045.20},
        {"title": "TEACHG ASST-1/10-GSHIP", "count": 1762, "avg_pay": 14341.29}
      ]
    }
  ]
}
```

## Database Schema

The pipeline creates three main tables:

### `wage_summaries`
- Statistical summaries per location-year
- Percentiles stored as JSON
- Pay component breakdowns

### `wage_pyramids`
- Income distribution brackets
- Job title analysis per bracket
- Stored as JSON for flexibility

### `title_analysis`
- Job title frequency and statistics
- Top N titles per location-year
- Compensation ranges per title

## Architecture

```
UC Wages Analysis Pipeline
├── Data Input (173 JSON files, 1.3GB)
├── Concurrent Processing (Go routines)
├── Statistical Analysis (percentiles, distributions)
├── Output Generation (JSON)
└── Database Upload (PostgreSQL)
```

### Key Components

- **Parser Package**: JSON data loading and currency parsing
- **Calculator Package**: Statistical computations and aggregations
- **Models Package**: Type definitions and data structures
- **Command Tools**: Individual analysis executables
- **Docker Container**: Production deployment

## Configuration

### Command Line Flags

All tools support these common flags:
- `-data`: Input data directory path
- `-output`: Output directory path
- `-workers`: Number of concurrent workers (default: 4)

### Environment Variables

- `DATABASE_URL`: PostgreSQL connection string for uploads
- `WORKERS`: Default worker count

## Monitoring & Logging

The pipeline provides detailed progress reporting:
- ✓ File processing confirmations
- Error handling with detailed messages
- Performance metrics and timing
- Final summary statistics

## Development

### Project Structure
```
anal/
├── cmd/                    # Command-line tools
│   ├── calculate_sums/
│   ├── generate_pyramid/
│   ├── analyze_titles/
│   ├── run_all/
│   └── upload_analysis/
├── pkg/                    # Shared packages
│   ├── models/            # Data structures
│   ├── parser/            # JSON processing
│   └── calculator/        # Analysis algorithms
├── output/                # Generated analysis files
├── Dockerfile             # Container definition
└── docker-compose.yml     # Service orchestration
```

### Contributing

1. **Add New Analysis**: Create new calculator functions in `pkg/calculator/`
2. **Extend Models**: Update type definitions in `pkg/models/types.go`
3. **Create Tools**: Add new commands in `cmd/` directory
4. **Database Integration**: Update schema and upload logic

### Testing

```bash
# Test with single campus
docker run --rm \
  -v /path/to/data/Berkeley:/data:ro \
  -v $(pwd)/output:/app/output \
  uc-wages-analysis

# Dry run database upload
docker run --rm \
  -v $(pwd)/output:/app/output \
  uc-wages-analysis /bin/upload_analysis -dry-run -db "connection_string"
```

## Deployment

### Production Deployment
1. Build Docker image
2. Configure database connection
3. Set up data volume mounts
4. Run analysis pipeline
5. Upload to production database

### Automated Pipeline
The tools can be integrated into CI/CD for:
- Scheduled analysis updates
- Data validation workflows
- Performance monitoring
- Automated reporting

---

**Generated Analysis Coverage**
- 14 UC Locations + ASUCLA
- 15 Years (2010-2024)
- 173 Location-Year Combinations
- 2M+ Individual Records Processed
- ~10MB Precalculated Metrics Generated