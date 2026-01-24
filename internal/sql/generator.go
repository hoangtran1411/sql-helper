package sql

import (
	"fmt"
	"strings"
)

// GenerateSQLValues generates SQL INSERT values from the data
func GenerateSQLValues(headers []string, dataRows [][]interface{}, numberColumns []string) string {
	if len(dataRows) == 0 {
		return ""
	}

	// Convert numberColumns to map for O(1) lookup
	numColSet := make(map[string]bool)
	for _, col := range numberColumns {
		numColSet[col] = true
	}

	columnCount := len(headers)

	// Pre-allocate builder with estimated size
	estimatedSize := len(dataRows) * columnCount * 20
	var builder strings.Builder
	builder.Grow(estimatedSize)

	for rowIdx, row := range dataRows {
		if rowIdx > 0 {
			builder.WriteString(",\n")
		}
		builder.WriteByte('(')

		for i := 0; i < columnCount; i++ {
			if i > 0 {
				builder.WriteString(", ")
			}

			var cellValue interface{}
			if i < len(row) {
				cellValue = row[i]
			}

			header := ""
			if i < len(headers) {
				header = headers[i]
			}

			formatted := FormatCellValue(cellValue, numColSet[header])
			builder.WriteString(formatted)
		}

		builder.WriteByte(')')
	}

	return builder.String()
}

// FindAndReplace replaces values in the data rows
func FindAndReplace(dataRows [][]interface{}, findValue, replaceWith string) [][]interface{} {
	if dataRows == nil {
		return make([][]interface{}, 0)
	}

	for _, row := range dataRows {
		for j, cell := range row {
			// Convert to string to check value - simple robust check
			// Optimization: could be type-specific but generic is safer for now
			if fmt.Sprintf("%v", cell) == findValue {
				row[j] = replaceWith
			}
		}
	}

	return dataRows
}
