package calculator

import (
	"sort"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/ucinvestments/uc-wages-analysis/pkg/models"
	"github.com/ucinvestments/uc-wages-analysis/pkg/parser"
)

// BracketDefinition defines a wage bracket
type BracketDefinition struct {
	Range    string
	MinValue float64
	MaxValue float64
}

// GetWageBrackets returns the standard wage brackets
func GetWageBrackets() []BracketDefinition {
	return []BracketDefinition{
		{"0-25k", 0, 25000},
		{"25k-50k", 25000, 50000},
		{"50k-75k", 50000, 75000},
		{"75k-100k", 75000, 100000},
		{"100k-150k", 100000, 150000},
		{"150k-200k", 150000, 200000},
		{"200k-300k", 200000, 300000},
		{"300k-500k", 300000, 500000},
		{"500k-1M", 500000, 1000000},
		{"1M+", 1000000, 999999999},
	}
}

// CalculatePyramid generates wage distribution pyramid
func CalculatePyramid(data *models.WageData) (*models.Pyramid, error) {
	brackets := GetWageBrackets()
	bracketData := make(map[string]*bracketInfo)

	// Initialize bracket data
	for _, bracket := range brackets {
		bracketData[bracket.Range] = &bracketInfo{
			definition: bracket,
			wages:      []float64{},
			titles:     make(map[string][]float64),
		}
	}

	totalEmployees := 0
	totalPay := 0.0

	// Process each record
	for _, record := range data.Records {
		_, _, _, gross := parser.ConvertRecordToFloat(record)

		if gross <= 0 {
			continue
		}

		totalEmployees++
		totalPay += gross

		// Find appropriate bracket
		for _, bracket := range brackets {
			if gross >= bracket.MinValue && gross < bracket.MaxValue {
				bd := bracketData[bracket.Range]
				bd.wages = append(bd.wages, gross)

				// Track title data
				if record.Title != "" && record.Title != "*****" {
					bd.titles[record.Title] = append(bd.titles[record.Title], gross)
				}
				break
			}
		}
	}

	// Build pyramid structure
	pyramid := &models.Pyramid{
		Location:       data.Location,
		Year:           data.Year,
		GeneratedAt:    time.Now(),
		TotalEmployees: totalEmployees,
		TotalPay:       totalPay,
		Brackets:       []models.WageBracket{},
	}

	// Calculate statistics for each bracket
	for _, bracket := range brackets {
		bd := bracketData[bracket.Range]

		if len(bd.wages) == 0 {
			continue
		}

		wageBracket := models.WageBracket{
			Range:      bracket.Range,
			MinValue:   bracket.MinValue,
			MaxValue:   bracket.MaxValue,
			Count:      len(bd.wages),
			Percentage: float64(len(bd.wages)) / float64(totalEmployees) * 100,
		}

		// Calculate bracket statistics
		wageBracket.AvgPay, _ = stats.Mean(bd.wages)
		wageBracket.MedianPay, _ = stats.Median(bd.wages)
		wageBracket.TotalPay = sumFloat64(bd.wages)

		// Get top titles
		wageBracket.TopTitles = getTopTitles(bd.titles, 10)

		pyramid.Brackets = append(pyramid.Brackets, wageBracket)
	}

	return pyramid, nil
}

// bracketInfo holds temporary data for bracket calculations
type bracketInfo struct {
	definition BracketDefinition
	wages      []float64
	titles     map[string][]float64
}

// getTopTitles returns the top N titles by frequency
func getTopTitles(titles map[string][]float64, limit int) []models.TitleCount {
	// Create slice for sorting
	titleStats := []models.TitleCount{}

	for title, wages := range titles {
		avgPay, _ := stats.Mean(wages)
		titleStats = append(titleStats, models.TitleCount{
			Title:  title,
			Count:  len(wages),
			AvgPay: avgPay,
		})
	}

	// Sort by count descending
	sort.Slice(titleStats, func(i, j int) bool {
		return titleStats[i].Count > titleStats[j].Count
	})

	// Return top N
	if len(titleStats) > limit {
		return titleStats[:limit]
	}
	return titleStats
}