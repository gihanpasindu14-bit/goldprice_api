package utils

import (
	"encoding/csv"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"goldprice-api/models"
)

// ValidateAndParseCSV validates CSV format and extracts data
// Returns: parsed data, error
func ValidateAndParseCSV(content string) ([]models.GoldPrice, error) {
	reader := csv.NewReader(strings.NewReader(content))

	// Read all lines
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	if len(lines) < 3 {
		return nil, fmt.Errorf("CSV file must have at least 3 lines (header + data)")
	}

	// Parse data (skip first 2 header rows)
	var data []models.GoldPrice

	for i := 1; i < len(lines); i++ {
		row := lines[i]

		// Skip empty rows
		if len(row) < 3 {
			continue
		}

		dateStr := strings.TrimSpace(row[1])
		priceStr := strings.TrimSpace(row[2])

		// Skip if empty
		if dateStr == "" || priceStr == "" {
			continue
		}

		// Parse price
		pricePerOunce, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			continue
		}

		// Parse and format date
		dateObj, err := parseDate(dateStr)
		if err != nil {
			continue
		}

		// Convert to standard format: DD-MMM-YY
		formattedDate := dateObj.Format("02-Jan-06")
		// Also create sortable format: YYYY-MM-DD
		sortableDate := dateObj.Format("2006-01-02")

		// Convert from per ounce to per gram (1 oz = 31.1035 grams)
		pricePerGram := pricePerOunce / 31.1035

		// Calculate prices for different carats
		record := models.GoldPrice{
			Date:          formattedDate,
			DateSortable:  sortableDate,
			PricePerOunce: math.Round(pricePerOunce*100) / 100,
			Price24K:      int(math.Round(pricePerGram)),
			Price22K:      int(math.Round(pricePerGram * 0.9167)),
			Price18K:      int(math.Round(pricePerGram * 0.75)),
		}

		data = append(data, record)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no valid data rows found in CSV")
	}

	return data, nil
}

// parseDate tries multiple date formats
func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",  // YYYY-MM-DD
		"01/02/2006",  // MM/DD/YYYY (Excel format)
		"1/2/2006",    // M/D/YYYY (Excel short format)
		"02-Jan-06",   // DD-Mon-YY
		"02-Jan-2006", // DD-Mon-YYYY
		"2006/01/02",  // YYYY/MM/DD
		"02/01/2006",  // DD/MM/YYYY
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
