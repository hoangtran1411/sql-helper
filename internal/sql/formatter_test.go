package sql

import (
	"testing"
	"time"
)

func TestFormatCellValue(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		isNumeric bool
		expected  string
	}{
		// Numeric column tests
		{"nil numeric", nil, true, "NULL"},
		{"int numeric", 42, true, "42"},
		{"int64 numeric", int64(100), true, "100"},
		{"float64 numeric", 3.14, true, "3.14"},
		{"float32 numeric", float32(2.5), true, "2.5"},
		{"empty string numeric", "", true, "NULL"},
		{"whitespace string numeric", "   ", true, "NULL"},
		{"invalid string numeric", "abc", true, "NULL"},
		{"valid string number", "123", true, "123"},
		{"valid string float", "45.67", true, "45.67"},
		{"negative number", -10, true, "-10"},
		{"zero value", 0, true, "0"},
		{"bool true numeric", true, true, "1"},
		{"bool false numeric", false, true, "0"},

		// Text column tests
		{"nil text", nil, false, "''"},
		{"simple string", "Hello", false, "'Hello'"},
		{"string with quote", "It's OK", false, "'It''s OK'"},
		{"string with multiple quotes", "O'Brien's", false, "'O''Brien''s'"},
		{"number as text", 123, false, "'123'"},
		{"float as text", 45.67, false, "'45.67'"},
		{"bool true text", true, false, "'TRUE'"},
		{"bool false text", false, false, "'FALSE'"},
		{"unicode string", "Việt Nam", false, "'Việt Nam'"},
		{"empty string text", "", false, "''"},
		{"html entities", "<div>", false, "'<div>'"},
		{"sql injection attempt", "'; DROP TABLE--", false, "'''; DROP TABLE--'"},

		// Date tests
		{"date value", time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC), false, "'2024-01-15 10:30:00'"},
		{"iso date string", "2024-01-15T10:30:00Z", false, "'2024-01-15 10:30:00'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCellValue(tt.value, tt.isNumeric)
			if result != tt.expected {
				t.Errorf("FormatCellValue(%v, %v) = %s; want %s",
					tt.value, tt.isNumeric, result, tt.expected)
			}
		})
	}
}

func TestFormatNumericValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"int", 42, "42"},
		{"int8", int8(10), "10"},
		{"int16", int16(100), "100"},
		{"int32", int32(1000), "1000"},
		{"int64", int64(10000), "10000"},
		{"uint", uint(5), "5"},
		{"uint8", uint8(15), "15"},
		{"uint16", uint16(150), "150"},
		{"uint32", uint32(1500), "1500"},
		{"uint64", uint64(15000), "15000"},
		{"float32", float32(3.14), "3.14"},
		{"float64", 3.14159, "3.14159"},
		{"negative float", -99.99, "-99.99"},
		{"valid string", "123.45", "123.45"},
		{"empty string", "", "NULL"},
		{"invalid string", "not a number", "NULL"},
		{"bool true", true, "1"},
		{"bool false", false, "0"},
		{"nil", nil, "NULL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatNumericValue(tt.value)
			if result != tt.expected {
				t.Errorf("FormatNumericValue(%v) = %s; want %s",
					tt.value, result, tt.expected)
			}
		})
	}
}

func TestFormatTextValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"simple string", "Hello World", "'Hello World'"},
		{"single quote escape", "It's a test", "'It''s a test'"},
		{"multiple quotes", "He said 'Hello'", "'He said ''Hello'''"},
		{"number", 42, "'42'"},
		{"float", 3.14, "'3.14'"},
		{"bool true", true, "'TRUE'"},
		{"bool false", false, "'FALSE'"},
		{"time value", time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC), "'2024-06-15 14:30:45'"},
		{"rfc3339 string", "2024-12-25T08:00:00Z", "'2024-12-25 08:00:00'"},
		{"iso datetime string", "2024-03-15T09:30:00", "'2024-03-15 09:30:00'"},
		{"unicode", "日本語", "'日本語'"},
		{"special chars", "Line1\nLine2\tTab", "'Line1\nLine2\tTab'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTextValue(tt.value)
			if result != tt.expected {
				t.Errorf("FormatTextValue(%v) = %s; want %s",
					tt.value, result, tt.expected)
			}
		})
	}
}
