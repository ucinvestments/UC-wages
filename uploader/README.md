# UC Wages Database Uploader

A multithreaded Go application that uploads UC wage data from JSON files to a Neon PostgreSQL database.

## Features

- **Multithreaded Processing**: Configurable number of worker goroutines for concurrent uploads
- **Batch Processing**: Efficient database inserts with configurable batch sizes
- **Progress Tracking**: Database table tracks upload progress for each location/year
- **Error Handling**: Robust error handling with detailed logging
- **Resume Capability**: Skips already uploaded data using UPSERT operations
- **Environment Configuration**: Secure database credentials via .env file

## Prerequisites

- Go 1.21 or later
- PostgreSQL database (Neon)
- Scraped UC wage data in JSON format

## Setup

1. **Install Dependencies**:
   ```bash
   make deps
   ```

2. **Configure Environment**:
   Ensure your `.env` file contains:
   ```bash
   DATABASE_URL='postgresql://username:password@host/database?sslmode=require'
   ```

3. **Build the Application**:
   ```bash
   make build
   ```

## Usage

### Quick Start
```bash
# Upload all wage data with default settings (5 workers, 1000 batch size)
make run

# Fast upload with more workers
make upload-fast

# Test with Berkeley data only
make test
```

### Command Line Options
```bash
./uploader [options]

Options:
  -workers int     Number of concurrent workers (default: 5)
  -batch int       Batch size for database inserts (default: 1000)
  -data string     Data directory containing wage files (default: "../data")
  -env string      Path to .env file (default: "../.env")
```

### Examples
```bash
# Upload with 10 workers and larger batches
./uploader -workers=10 -batch=2000

# Upload specific location data
./uploader -data="../data/Berkeley"

# Use custom .env file location
./uploader -env="/path/to/.env"
```

## Database Schema

The uploader creates two tables:

### `uc_wages`
Stores individual wage records with indexes for efficient querying:
- `location`, `year`, `employee_id` (unique constraint)
- `firstname`, `lastname`, `title`
- `basepay`, `overtimepay`, `adjustpay`, `grosspay`
- `scraped_at`, `uploaded_at`

### `upload_progress`
Tracks upload status for each location/year combination:
- `location`, `year` (unique constraint)
- `total_records`, `uploaded_records`
- `status` (pending, processing, completed, failed)
- `started_at`, `completed_at`, `error_message`

## Performance

- **Concurrent Processing**: Multiple workers process files simultaneously
- **Batch Inserts**: Database writes are batched for optimal performance
- **Connection Pooling**: PostgreSQL driver handles connection pooling
- **Memory Efficient**: Processes files individually, not loading entire dataset

### Recommended Settings
- **Small datasets** (< 1GB): 5 workers, 1000 batch size
- **Large datasets** (> 1GB): 10+ workers, 2000+ batch size
- **Network limitations**: Reduce workers, increase batch size

## Error Handling

- **File Errors**: Logged and skipped, other files continue processing
- **Database Errors**: Transaction rollback, detailed error logging
- **Network Issues**: Automatic retry via PostgreSQL driver
- **Progress Tracking**: Failed uploads marked in progress table

## Data Processing

The uploader handles various data formats:
- **Amounts**: Parses comma-separated numbers, handles empty values
- **Dates**: RFC3339 timestamp parsing with fallback to current time
- **IDs**: Flexible parsing of employee IDs from various formats
- **Text Fields**: Safe handling of empty/null values

## Monitoring

Check upload progress:
```sql
SELECT location, year, status, uploaded_records, total_records,
       ROUND(uploaded_records::float / total_records * 100, 2) as progress_pct
FROM upload_progress
ORDER BY location, year;
```

## Troubleshooting

1. **Connection Issues**: Verify DATABASE_URL and network connectivity
2. **Memory Issues**: Reduce number of workers or batch size
3. **Slow Performance**: Increase workers for CPU-bound, reduce for I/O-bound
4. **Duplicate Data**: Uses UPSERT operations, safe to re-run

## File Structure
```
uploader/
├── main.go           # Main application
├── go.mod           # Go module definition
├── Makefile         # Build and run commands
├── schema.sql       # Database schema (reference)
├── test_small.sh    # Small dataset test script
└── README.md        # This file
```