package calculator

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/ucinvestments/uc-wages-analysis/pkg/models"
	"github.com/ucinvestments/uc-wages-analysis/pkg/parser"
)

// CalculateSummary computes statistical summary for wage data
func CalculateSummary(data *models.WageData) (*models.Summary, error) {
	var grossPays []float64
	var basePays []float64
	var overtimePays []float64
	var adjustPays []float64

	// Extract all pay values
	for _, record := range data.Records {
		base, overtime, adjust, gross := parser.ConvertRecordToFloat(record)

		// Only include records with valid gross pay
		if gross > 0 {
			grossPays = append(grossPays, gross)
			basePays = append(basePays, base)
			overtimePays = append(overtimePays, overtime)
			adjustPays = append(adjustPays, adjust)
		}
	}

	if len(grossPays) == 0 {
		return nil, nil
	}

	// Sort for percentile calculations
	sort.Float64s(grossPays)

	// Calculate basic statistics
	mean, _ := stats.Mean(grossPays)
	median, _ := stats.Median(grossPays)
	stdDev, _ := stats.StandardDeviation(grossPays)
	min, _ := stats.Min(grossPays)
	max, _ := stats.Max(grossPays)

	// Calculate percentiles
	percentiles := make(map[string]float64)
	percentileValues := []float64{25, 50, 75, 90, 95, 99}
	for _, p := range percentileValues {
		value, _ := stats.Percentile(grossPays, p)
		percentiles[formatPercentileKey(p)] = value
	}

	// Calculate pay components
	totalBase := sumFloat64(basePays)
	totalOvertime := sumFloat64(overtimePays)
	totalAdjust := sumFloat64(adjustPays)
	totalGross := sumFloat64(grossPays)

	employeeCount := len(grossPays)

	summary := &models.Summary{
		Location:      data.Location,
		Year:          data.Year,
		GeneratedAt:   time.Now(),
		EmployeeCount: employeeCount,
		TotalGrossPay: totalGross,
		AvgGrossPay:   mean,
		MedianPay:     median,
		StdDev:        stdDev,
		MinPay:        min,
		MaxPay:        max,
		Percentiles:   percentiles,
		PayComponents: models.PayComponents{
			TotalBase:        totalBase,
			TotalOvertime:    totalOvertime,
			TotalAdjustments: totalAdjust,
			AvgBase:          totalBase / float64(employeeCount),
			AvgOvertime:      totalOvertime / float64(employeeCount),
			AvgAdjustments:   totalAdjust / float64(employeeCount),
		},
	}

	return summary, nil
}

// CalculateGiniCoefficient calculates income inequality
func CalculateGiniCoefficient(wages []float64) float64 {
	if len(wages) == 0 {
		return 0
	}

	// Sort wages
	sorted := make([]float64, len(wages))
	copy(sorted, wages)
	sort.Float64s(sorted)

	n := float64(len(sorted))
	index := 0.0
	gini := 0.0

	for _, wage := range sorted {
		index++
		gini += (2*index - n - 1) * wage
	}

	total := sumFloat64(sorted)
	if total == 0 {
		return 0
	}

	return gini / (n * total)
}

// CalculateSkewness calculates the skewness of distribution
func CalculateSkewness(wages []float64) float64 {
	if len(wages) < 3 {
		return 0
	}

	mean, _ := stats.Mean(wages)
	stdDev, _ := stats.StandardDeviation(wages)

	if stdDev == 0 {
		return 0
	}

	n := float64(len(wages))
	sum := 0.0

	for _, wage := range wages {
		sum += math.Pow((wage - mean) / stdDev, 3)
	}

	return (n / ((n - 1) * (n - 2))) * sum
}

// CalculateKurtosis calculates the kurtosis of distribution
func CalculateKurtosis(wages []float64) float64 {
	if len(wages) < 4 {
		return 0
	}

	mean, _ := stats.Mean(wages)
	stdDev, _ := stats.StandardDeviation(wages)

	if stdDev == 0 {
		return 0
	}

	n := float64(len(wages))
	sum := 0.0

	for _, wage := range wages {
		sum += math.Pow((wage - mean) / stdDev, 4)
	}

	return ((n*(n+1))/((n-1)*(n-2)*(n-3)))*sum - (3*(n-1)*(n-1))/((n-2)*(n-3))
}

// Helper functions
func sumFloat64(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum
}

func formatPercentileKey(p float64) string {
	return fmt.Sprintf("p%d", int(p))
}