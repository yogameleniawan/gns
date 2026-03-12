package models

import (
	"fmt"
	"time"
)

// ParseTradeDate parses date string from format "YYYYMMDD" to time.Time
// Example: "20260115" -> 2026-01-15
func ParseTradeDate(dateStr string) (time.Time, error) {
	if len(dateStr) != 8 {
		return time.Time{}, fmt.Errorf("invalid date format, expected YYYYMMDD, got: %s", dateStr)
	}

	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date %s: %w", dateStr, err)
	}

	return t, nil
}

// FormatTradeDate formats time.Time to "YYYYMMDD" string
func FormatTradeDate(t time.Time) string {
	return t.Format("20060102")
}
