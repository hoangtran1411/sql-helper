package sql

import (
	"fmt"
	"strings"
)

// GenerateOptions holds configuration for SQL generation
type GenerateOptions struct {
	TableName       string
	SelectedColumns []string
	NumberColumns   []string
	BatchSize       int
	ValuesOnly      bool
}

// FormatIdentifier formats a table or column name.
// If it contains spaces or brackets, it cleans and wraps it in [name].
func FormatIdentifier(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	if strings.ContainsAny(name, " \t\n\r[]") {
		clean := strings.ReplaceAll(name, "[", "")
		clean = strings.ReplaceAll(clean, "]", "")
		return "[" + clean + "]"
	}
	return name
}

// BuildInsertPrefix builds the "INSERT INTO table (col1, col2) VALUES\n" prefix.
func BuildInsertPrefix(tableName string, columns []string) string {
	tableName = strings.TrimSpace(tableName)
	if tableName == "" {
		tableName = "my_table"
	}

	var b strings.Builder
	b.WriteString("INSERT INTO ")
	b.WriteString(FormatIdentifier(tableName))
	b.WriteString(" (")
	for i, col := range columns {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(FormatIdentifier(col))
	}
	b.WriteString(") VALUES\n")
	return b.String()
}

// GenerateBatchSQL generates SQL INSERT statements (or raw VALUES) according to opts.
func GenerateBatchSQL(headers []string, dataRows [][]interface{}, opts GenerateOptions) string {
	if len(dataRows) == 0 {
		return ""
	}

	// Default to all headers if SelectedColumns is empty or nil
	selectedCols := opts.SelectedColumns
	if len(selectedCols) == 0 {
		selectedCols = headers
	}
	if len(selectedCols) == 0 {
		return ""
	}

	// Map headers to original indices
	headerIndexMap := make(map[string]int, len(headers))
	for i, h := range headers {
		if _, exists := headerIndexMap[h]; !exists {
			headerIndexMap[h] = i
		}
	}

	colIndices := make([]int, 0, len(selectedCols))
	validSelectedCols := make([]string, 0, len(selectedCols))
	for _, col := range selectedCols {
		if idx, ok := headerIndexMap[col]; ok {
			colIndices = append(colIndices, idx)
			validSelectedCols = append(validSelectedCols, col)
		}
	}
	if len(validSelectedCols) == 0 {
		return ""
	}

	// Convert numberColumns to map for O(1) lookup
	numColSet := make(map[string]bool, len(opts.NumberColumns))
	for _, col := range opts.NumberColumns {
		numColSet[col] = true
	}

	// Estimated allocation
	estimatedSize := len(dataRows) * len(validSelectedCols) * 25
	var builder strings.Builder
	builder.Grow(estimatedSize)

	// If ValuesOnly mode, generate comma-separated tuples
	if opts.ValuesOnly {
		for rowIdx, row := range dataRows {
			if rowIdx > 0 {
				builder.WriteString(",\n")
			}
			builder.WriteString(FormatRowSQLSelected(row, colIndices, headers, numColSet))
		}
		return builder.String()
	}

	// Batch INSERT mode
	prefix := BuildInsertPrefix(opts.TableName, validSelectedCols)
	batchSize := opts.BatchSize
	if batchSize < 0 {
		batchSize = 0
	}

	inBatchCount := 0
	for rowIdx, row := range dataRows {
		if inBatchCount == 0 {
			if rowIdx > 0 {
				builder.WriteString("\n\n")
			}
			builder.WriteString(prefix)
		} else {
			builder.WriteString(",\n")
		}

		builder.WriteString(FormatRowSQLSelected(row, colIndices, headers, numColSet))
		inBatchCount++

		if batchSize > 0 && inBatchCount == batchSize {
			builder.WriteString(";")
			inBatchCount = 0
		}
	}

	if inBatchCount > 0 {
		builder.WriteString(";")
	}

	return builder.String()
}

// GenerateSQLValues generates SQL INSERT values from the data (legacy wrapper)
func GenerateSQLValues(headers []string, dataRows [][]interface{}, numberColumns []string) string {
	return GenerateBatchSQL(headers, dataRows, GenerateOptions{
		SelectedColumns: headers,
		NumberColumns:   numberColumns,
		ValuesOnly:      true,
	})
}

// FormatRowSQL formats a single row into SQL values format (val1, val2, ...) using all headers
func FormatRowSQL(row []interface{}, headers []string, numColSet map[string]bool) string {
	indices := make([]int, len(headers))
	for i := range headers {
		indices[i] = i
	}
	return FormatRowSQLSelected(row, indices, headers, numColSet)
}

// FormatRowSQLSelected formats selected columns of a row into SQL values format (val1, val2, ...)
func FormatRowSQLSelected(row []interface{}, colIndices []int, headers []string, numColSet map[string]bool) string {
	var builder strings.Builder
	builder.WriteByte('(')

	for i, colIdx := range colIndices {
		if i > 0 {
			builder.WriteString(", ")
		}

		var cellValue interface{}
		if colIdx >= 0 && colIdx < len(row) {
			cellValue = row[colIdx]
		}

		header := ""
		if colIdx >= 0 && colIdx < len(headers) {
			header = headers[colIdx]
		}

		formatted := FormatCellValue(cellValue, numColSet[header])
		builder.WriteString(formatted)
	}

	builder.WriteByte(')')
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
