package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ucinvestments/uc-wages-analysis/pkg/models"
)

// LoadWageData loads and parses a wage JSON file
func LoadWageData(filepath string) (*models.WageData, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filepath, err)
	}
	defer file.Close()

	var data models.WageData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding JSON from %s: %w", filepath, err)
	}

	return &data, nil
}

// ParseCurrency converts currency string to float64
func ParseCurrency(amount string) float64 {
	// Remove commas and dollar signs
	cleaned := strings.ReplaceAll(amount, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.TrimSpace(cleaned)

	// Handle empty or invalid values
	if cleaned == "" || cleaned == "N/A" || cleaned == "*****" {
		return 0.0
	}

	value, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0.0
	}

	return value
}

// ConvertRecordToFloat converts string wage values to floats
func ConvertRecordToFloat(record models.WageRecord) (basePay, overtimePay, adjustPay, grossPay float64) {
	basePay = ParseCurrency(record.BasePay)
	overtimePay = ParseCurrency(record.OvertimePay)
	adjustPay = ParseCurrency(record.AdjustPay)
	grossPay = ParseCurrency(record.GrossPay)

	// If gross pay is 0 but components exist, calculate it
	if grossPay == 0 && (basePay > 0 || overtimePay > 0 || adjustPay > 0) {
		grossPay = basePay + overtimePay + adjustPay
	}

	return
}

// SaveJSON saves any struct to a JSON file
func SaveJSON(filepath string, data interface{}) error {
	// Create directory if it doesn't exist
	dir := strings.Split(filepath, "/")
	if len(dir) > 1 {
		dirPath := strings.Join(dir[:len(dir)-1], "/")
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", filepath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding JSON to %s: %w", filepath, err)
	}

	return nil
}