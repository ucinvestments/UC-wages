# UC Wages Data Analysis and Preprocessing Plan (Go Implementation)

## Problem Statement
- Loading 1.3GB of raw wage data (173 JSON files) for real-time calculations is inefficient
- Need to precalculate common metrics and aggregations for faster web app performance
- Implement in Go for consistency with existing scraper and high performance

## Data Structure Overview
- 14 UC locations + ASUCLA
- 15 years of data (2010-2024)
- ~37,000 records per location per year (Berkeley 2024 example)
- Each record contains: basepay, overtimepay, adjustpay, grosspay, title, name (often redacted)

## Proposed Precalculated Metrics

### 1. Core Summary Statistics (sums/)
For each location-year combination:
- **Employee Count**: Total number of employees
- **Average Wage**: Mean gross pay
- **Median Wage**: Median gross pay (more representative than mean)
- **Total Wage**: Sum of all gross pay
- **Pay Components**:
  - Total base pay
  - Total overtime pay
  - Total adjustment pay
- **Percentiles**: 25th, 75th, 90th, 95th, 99th percentiles
- **Standard Deviation**: Wage distribution spread
- **Min/Max**: Lowest and highest gross pay

### 2. Wage Pyramid Analysis (pyramid/)
Group employees by earning ranges with logarithmic scaling:
- **Brackets**:
  - < $25k
  - $25k - $50k
  - $50k - $75k
  - $75k - $100k
  - $100k - $150k
  - $150k - $200k
  - $200k - $300k
  - $300k - $500k
  - $500k - $1M
  - > $1M
- **Per Bracket**: Count, percentage, top 10 job titles (with frequency)

### 3. Advanced Analytics

#### Time Series Aggregations (trends/)
- **Year-over-Year Growth**: By location
- **Inflation-Adjusted Wages**: Using CPI data
- **Employee Count Trends**: Growth/decline patterns
- **Pay Component Trends**: How overtime/adjustments change

#### Job Title Analysis (titles/)
- **Top 100 Job Titles**: By count and average pay
- **Title Categories**: Academic, Administrative, Medical, Executive
- **Pay Ranges by Title**: Min, max, average, median
- **Title Frequency**: Count per title

#### Distribution Analysis (distributions/)
- **Gini Coefficient**: Income inequality measure
- **Lorenz Curve Data Points**: For inequality visualization
- **Skewness & Kurtosis**: Distribution shape metrics
- **Decile Analysis**: Income by deciles

#### Comparative Analysis (comparisons/)
- **Campus Rankings**: By various metrics
- **System-wide Aggregates**: UC-wide statistics
- **Year-over-Year Comparisons**: Growth rates

## Go Implementation Structure

### Data Structures
```go
type WageRecord struct {
    ID          int     `json:"id"`
    FirstName   string  `json:"firstname"`
    LastName    string  `json:"lastname"`
    Title       string  `json:"title"`
    Location    string  `json:"location"`
    Year        string  `json:"year"`
    BasePay     float64 `json:"basepay"`
    OvertimePay float64 `json:"overtimepay"`
    AdjustPay   float64 `json:"adjustpay"`
    GrossPay    float64 `json:"grosspay"`
}

type Summary struct {
    Location       string            `json:"location"`
    Year          int               `json:"year"`
    GeneratedAt   time.Time         `json:"generated_at"`
    EmployeeCount int               `json:"employee_count"`
    TotalGrossPay float64           `json:"total_gross_pay"`
    AvgGrossPay   float64           `json:"avg_gross_pay"`
    MedianPay     float64           `json:"median_gross_pay"`
    StdDev        float64           `json:"std_dev"`
    MinPay        float64           `json:"min_pay"`
    MaxPay        float64           `json:"max_pay"`
    Percentiles   map[int]float64   `json:"percentiles"`
    PayComponents PayComponents     `json:"pay_components"`
}

type WageBracket struct {
    Range      string       `json:"range"`
    MinValue   float64      `json:"min_value"`
    MaxValue   float64      `json:"max_value"`
    Count      int          `json:"count"`
    Percentage float64      `json:"percentage"`
    AvgPay     float64      `json:"avg_pay"`
    TopTitles  []TitleCount `json:"top_titles"`
}

type Pyramid struct {
    Location    string        `json:"location"`
    Year        int           `json:"year"`
    GeneratedAt time.Time     `json:"generated_at"`
    Brackets    []WageBracket `json:"brackets"`
}
```

### Directory Structure
```
anal/
├── cmd/
│   ├── calculate_sums/
│   │   └── main.go
│   ├── generate_pyramid/
│   │   └── main.go
│   ├── analyze_titles/
│   │   └── main.go
│   └── run_all/
│       └── main.go
├── pkg/
│   ├── models/
│   │   └── types.go
│   ├── parser/
│   │   └── json.go
│   ├── calculator/
│   │   ├── statistics.go
│   │   ├── pyramid.go
│   │   └── titles.go
│   └── utils/
│       └── helpers.go
├── output/
│   ├── sums/
│   ├── pyramid/
│   ├── titles/
│   └── trends/
└── go.mod
```

## Implementation Plan

### Phase 1: Core Analysis Tools
1. **calculate_sums** - Calculate basic statistics for each location-year
   - Parse JSON wage files
   - Calculate mean, median, percentiles
   - Generate summary JSON files

2. **generate_pyramid** - Create wage distribution pyramids
   - Group employees by wage brackets
   - Extract top job titles per bracket
   - Output pyramid JSON files

3. **analyze_titles** - Aggregate job title statistics
   - Count frequency of each title
   - Calculate average pay per title
   - Identify top earners by title

### Phase 2: Advanced Analytics
4. **calculate_trends** - Time series analysis
   - Year-over-year growth calculations
   - Multi-year trend analysis
   - System-wide aggregations

5. **analyze_distributions** - Statistical metrics
   - Gini coefficient calculation
   - Distribution shape analysis
   - Inequality metrics

### Phase 3: Integration
6. **run_all** - Orchestrator to run all analyses
   - Concurrent processing with goroutines
   - Progress tracking
   - Error handling and retry logic

7. **Database Upload** - Modify existing uploader
   - Create new tables for aggregated data
   - Batch upload precalculated metrics
   - API endpoints for serving data

## Sample Output Formats

### Summary JSON (sums/Berkeley_2024.json)
```json
{
  "location": "Berkeley",
  "year": 2024,
  "generated_at": "2024-01-18T10:30:00Z",
  "employee_count": 37078,
  "total_gross_pay": 3500000000.50,
  "avg_gross_pay": 94385.50,
  "median_gross_pay": 72450.00,
  "std_dev": 85234.12,
  "min_pay": 100.00,
  "max_pay": 1250000.00,
  "percentiles": {
    "25": 45000.00,
    "50": 72450.00,
    "75": 110000.00,
    "90": 175000.00,
    "95": 235000.00,
    "99": 450000.00
  },
  "pay_components": {
    "total_base": 3200000000.00,
    "total_overtime": 150000000.00,
    "total_adjustments": 150000000.50,
    "avg_base": 86298.45,
    "avg_overtime": 4045.23,
    "avg_adjustments": 4041.82
  }
}
```

### Pyramid JSON (pyramid/Berkeley_2024.json)
```json
{
  "location": "Berkeley",
  "year": 2024,
  "generated_at": "2024-01-18T10:30:00Z",
  "brackets": [
    {
      "range": "0-25k",
      "min_value": 0,
      "max_value": 25000,
      "count": 5234,
      "percentage": 14.1,
      "avg_pay": 18500.50,
      "top_titles": [
        {"title": "STUDENT AST", "count": 1200, "avg_pay": 15000},
        {"title": "READER", "count": 800, "avg_pay": 12000}
      ]
    },
    {
      "range": "25k-50k",
      "min_value": 25000,
      "max_value": 50000,
      "count": 8456,
      "percentage": 22.8,
      "avg_pay": 38500.75,
      "top_titles": [
        {"title": "ADMIN AST 2", "count": 1500, "avg_pay": 42000}
      ]
    }
  ]
}
```

## Performance Optimizations

### Concurrent Processing
- Use goroutines to process multiple files simultaneously
- Worker pool pattern for controlled concurrency
- Channel-based coordination

### Memory Management
- Stream large JSON files instead of loading entirely
- Use buffered channels for data pipelines
- Clear processed data from memory promptly

### Caching Strategy
- Cache parsed wage data during processing
- Reuse calculations across different analyses
- Incremental updates for new data

## Benefits
- **Performance**: Reduce data transfer from ~1.3GB to ~10MB
- **Speed**: Eliminate real-time calculation overhead
- **Scalability**: Support more complex analytics
- **Consistency**: Single source of truth for metrics
- **Maintenance**: Easier to update and extend analyses

## Next Steps
1. ✅ Review and approve this plan
2. Implement Phase 1 Go programs
3. Test with sample data
4. Run full analysis on all data
5. Update database schema for aggregated tables
6. Modify web app to use precalculated data