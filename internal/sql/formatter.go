package sql

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FormatCellValue formats a cell value for SQL
func FormatCellValue(value interface{}, isNumeric bool) string {
	// Handle nil
	if value == nil {
		if isNumeric {
			return "NULL"
		}
		return "''"
	}

	// Handle numeric columns
	if isNumeric {
		return FormatNumericValue(value)
	}

	// Handle text columns
	return FormatTextValue(value)
}

// FormatNumericValue formats a value for a numeric column
func FormatNumericValue(value interface{}) string {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return "NULL"
		}
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			return s
		}
		return "NULL"
	case bool:
		if v {
			return "1"
		}
		return "0"
	default:
		return "NULL"
	}
}

// FormatTextValue formats a value for a text column
func FormatTextValue(value interface{}) string {
	switch v := value.(type) {
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	case string:

		// Optimization: heuristic check for date-like string (YYYY-MM-DD...)
		// This avoids expensive time.Parse calls for normal text
		if len(v) >= 10 && v[4] == '-' && v[7] == '-' {
			// Check if it's an ISO date string
			if t, err := time.Parse(time.RFC3339, v); err == nil {
				return fmt.Sprintf("'%s'", t.Format("2006-01-02 15:04:05"))
			}
			// Try other date formats
			if t, err := time.Parse("2006-01-02T15:04:05", v); err == nil {
				return fmt.Sprintf("'%s'", t.Format("2006-01-02 15:04:05"))
			}
		}
		// Escape single quotes and wrap
		escaped := strings.ReplaceAll(v, "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	case bool:
		if v {
			return "'TRUE'"
		}
		return "'FALSE'"
	default:
		// Convert to string and escape
		s := fmt.Sprintf("%v", v)
		escaped := strings.ReplaceAll(s, "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	}
}
