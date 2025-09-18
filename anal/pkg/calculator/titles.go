package calculator

import (
	"sort"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/ucinvestments/uc-wages-analysis/pkg/models"
	"github.com/ucinvestments/uc-wages-analysis/pkg/parser"
)

// AnalyzeTitles generates comprehensive title statistics
func AnalyzeTitles(data *models.WageData, topN int) (*models.TitleAnalysis, error) {
	titleMap := make(map[string]*titleData)

	// Aggregate data by title
	for _, record := range data.Records {
		if record.Title == "" || record.Title == "*****" {
			continue
		}

		_, _, _, gross := parser.ConvertRecordToFloat(record)
		if gross <= 0 {
			continue
		}

		if _, exists := titleMap[record.Title]; !exists {
			titleMap[record.Title] = &titleData{
				title: record.Title,
				wages: []float64{},
			}
		}

		titleMap[record.Title].wages = append(titleMap[record.Title].wages, gross)
	}

	// Calculate statistics for each title
	var titleStats []models.TitleStats

	for title, data := range titleMap {
		if len(data.wages) == 0 {
			continue
		}

		// Sort wages for median calculation
		sort.Float64s(data.wages)

		avgPay, _ := stats.Mean(data.wages)
		medianPay, _ := stats.Median(data.wages)
		minPay, _ := stats.Min(data.wages)
		maxPay, _ := stats.Max(data.wages)
		stdDev, _ := stats.StandardDeviation(data.wages)
		totalPay := sumFloat64(data.wages)

		titleStats = append(titleStats, models.TitleStats{
			Title:     title,
			Count:     len(data.wages),
			AvgPay:    avgPay,
			MedianPay: medianPay,
			MinPay:    minPay,
			MaxPay:    maxPay,
			StdDev:    stdDev,
			TotalPay:  totalPay,
		})
	}

	// Sort by count (most common titles first)
	sort.Slice(titleStats, func(i, j int) bool {
		// Primary sort by count, secondary by average pay
		if titleStats[i].Count == titleStats[j].Count {
			return titleStats[i].AvgPay > titleStats[j].AvgPay
		}
		return titleStats[i].Count > titleStats[j].Count
	})

	// Limit to top N titles
	topTitles := titleStats
	if len(titleStats) > topN {
		topTitles = titleStats[:topN]
	}

	analysis := &models.TitleAnalysis{
		Location:     data.Location,
		Year:         data.Year,
		GeneratedAt:  time.Now(),
		UniqueTitles: len(titleMap),
		TopTitles:    topTitles,
	}

	return analysis, nil
}

// CategorizeTitle attempts to categorize a job title
func CategorizeTitle(title string) string {
	// Convert to uppercase for comparison
	upperTitle := title

	// Academic titles
	academicKeywords := []string{"PROF", "LECTURER", "INSTRUCTOR", "TEACHER", "DEAN", "CHAIR", "RESEARCHER", "POST DOC", "STUDENT"}
	for _, keyword := range academicKeywords {
		if contains(upperTitle, keyword) {
			return "Academic"
		}
	}

	// Medical titles
	medicalKeywords := []string{"PHYSICIAN", "NURSE", "DOCTOR", "SURGEON", "MEDICAL", "CLINICAL", "THERAPIST", "PHARMACY", "HEALTH"}
	for _, keyword := range medicalKeywords {
		if contains(upperTitle, keyword) {
			return "Medical"
		}
	}

	// Executive titles
	execKeywords := []string{"PRESIDENT", "VICE PRESIDENT", "VP ", "CHIEF", "CEO", "CFO", "CTO", "DIRECTOR", "EXECUTIVE"}
	for _, keyword := range execKeywords {
		if contains(upperTitle, keyword) {
			return "Executive"
		}
	}

	// IT titles
	itKeywords := []string{"PROGRAMMER", "DEVELOPER", "ENGINEER", "ANALYST", "DATA", "IT ", "SOFTWARE", "SYSTEM", "NETWORK", "DATABASE"}
	for _, keyword := range itKeywords {
		if contains(upperTitle, keyword) {
			return "IT/Technical"
		}
	}

	// Administrative titles
	adminKeywords := []string{"ADMIN", "ASSISTANT", "COORDINATOR", "MANAGER", "CLERK", "SECRETARY", "RECEPTIONIST", "OFFICE"}
	for _, keyword := range adminKeywords {
		if contains(upperTitle, keyword) {
			return "Administrative"
		}
	}

	// Facilities/Operations
	facilityKeywords := []string{"CUSTODIAN", "MAINTENANCE", "GROUNDS", "FACILITIES", "SECURITY", "POLICE", "PARKING", "UTILITY"}
	for _, keyword := range facilityKeywords {
		if contains(upperTitle, keyword) {
			return "Facilities"
		}
	}

	return "Other"
}

// titleData holds temporary data for title calculations
type titleData struct {
	title string
	wages []float64
}

// contains checks if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr ||
		   len(str) > len(substr) && containsHelper(str[1:], substr)
}

func containsHelper(str, substr string) bool {
	if len(str) < len(substr) {
		return false
	}
	if str[:len(substr)] == substr {
		return true
	}
	return containsHelper(str[1:], substr)
}