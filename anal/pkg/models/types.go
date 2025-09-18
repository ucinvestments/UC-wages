package models

import "time"

// WageRecord represents a single employee wage record
type WageRecord struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Title       string `json:"title"`
	Location    string `json:"location"`
	Year        string `json:"year"`
	BasePay     string `json:"basepay"`
	OvertimePay string `json:"overtimepay"`
	AdjustPay   string `json:"adjustpay"`
	GrossPay    string `json:"grosspay"`
}

// WageData represents the complete JSON file structure
type WageData struct {
	Location      string       `json:"location"`
	Year          int          `json:"year"`
	ScrapedAt     string       `json:"scraped_at"`
	TotalRecords  int          `json:"total_records"`
	Records       []WageRecord `json:"records"`
}

// PayComponents breakdown of pay types
type PayComponents struct {
	TotalBase        float64 `json:"total_base"`
	TotalOvertime    float64 `json:"total_overtime"`
	TotalAdjustments float64 `json:"total_adjustments"`
	AvgBase          float64 `json:"avg_base"`
	AvgOvertime      float64 `json:"avg_overtime"`
	AvgAdjustments   float64 `json:"avg_adjustments"`
}

// Summary contains statistical summary for a location-year
type Summary struct {
	Location       string            `json:"location"`
	Year           int               `json:"year"`
	GeneratedAt    time.Time         `json:"generated_at"`
	EmployeeCount  int               `json:"employee_count"`
	TotalGrossPay  float64           `json:"total_gross_pay"`
	AvgGrossPay    float64           `json:"avg_gross_pay"`
	MedianPay      float64           `json:"median_gross_pay"`
	StdDev         float64           `json:"std_dev"`
	MinPay         float64           `json:"min_pay"`
	MaxPay         float64           `json:"max_pay"`
	Percentiles    map[string]float64 `json:"percentiles"`
	PayComponents  PayComponents     `json:"pay_components"`
}

// TitleCount represents job title frequency
type TitleCount struct {
	Title   string  `json:"title"`
	Count   int     `json:"count"`
	AvgPay  float64 `json:"avg_pay"`
}

// WageBracket represents a wage range bracket
type WageBracket struct {
	Range       string       `json:"range"`
	MinValue    float64      `json:"min_value"`
	MaxValue    float64      `json:"max_value"`
	Count       int          `json:"count"`
	Percentage  float64      `json:"percentage"`
	AvgPay      float64      `json:"avg_pay"`
	MedianPay   float64      `json:"median_pay"`
	TotalPay    float64      `json:"total_pay"`
	TopTitles   []TitleCount `json:"top_titles"`
}

// Pyramid contains wage distribution data
type Pyramid struct {
	Location       string        `json:"location"`
	Year           int           `json:"year"`
	GeneratedAt    time.Time     `json:"generated_at"`
	TotalEmployees int           `json:"total_employees"`
	TotalPay       float64       `json:"total_pay"`
	Brackets       []WageBracket `json:"brackets"`
}

// TitleAnalysis contains job title statistics
type TitleAnalysis struct {
	Location      string       `json:"location"`
	Year          int          `json:"year"`
	GeneratedAt   time.Time    `json:"generated_at"`
	UniqueTitles  int          `json:"unique_titles"`
	TopTitles     []TitleStats `json:"top_titles"`
}

// TitleStats contains statistics for a specific job title
type TitleStats struct {
	Title      string  `json:"title"`
	Count      int     `json:"count"`
	AvgPay     float64 `json:"avg_pay"`
	MedianPay  float64 `json:"median_pay"`
	MinPay     float64 `json:"min_pay"`
	MaxPay     float64 `json:"max_pay"`
	StdDev     float64 `json:"std_dev"`
	TotalPay   float64 `json:"total_pay"`
}