package sql

import (
	"cmp"
	"fmt"
	"io"
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
	tableName = cmp.Or(strings.TrimSpace(tableName), "my_table")

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

// ColumnMapping holds resolved column indices, valid column names, and numeric lookup set
type ColumnMapping struct {
	ColIndices        []int
	ValidSelectedCols []string
	NumColSet         map[string]bool
}

// ResolveColumns maps selected column names to indices and resolves numeric columns.
// Returns an error if no columns are selected or no selected columns exist in headers.
func ResolveColumns(headers, selectedColumns, numberColumns []string) (*ColumnMapping, error) {
	selectedCols := selectedColumns
	if len(selectedCols) == 0 {
		selectedCols = headers
	}
	if len(selectedCols) == 0 {
		return nil, fmt.Errorf("no columns selected")
	}

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
		return nil, fmt.Errorf("no valid columns found to export")
	}

	numColSet := make(map[string]bool, len(numberColumns))
	for _, col := range numberColumns {
		numColSet[col] = true
	}

	return &ColumnMapping{
		ColIndices:        colIndices,
		ValidSelectedCols: validSelectedCols,
		NumColSet:         numColSet,
	}, nil
}

// BatchWriter streams SQL statements or values into an io.Writer.
type BatchWriter struct {
	w            io.Writer
	headers      []string
	colMap       *ColumnMapping
	prefix       string
	batchSize    int
	valuesOnly   bool
	inBatchCount int
	totalRows    int
}

// NewBatchWriter creates a new BatchWriter for streaming SQL statements.
func NewBatchWriter(w io.Writer, headers []string, opts GenerateOptions) (*BatchWriter, error) {
	if w == nil {
		return nil, fmt.Errorf("writer cannot be nil")
	}

	colMap, err := ResolveColumns(headers, opts.SelectedColumns, opts.NumberColumns)
	if err != nil {
		return nil, fmt.Errorf("resolve columns: %w", err)
	}

	prefix := BuildInsertPrefix(opts.TableName, colMap.ValidSelectedCols)
	batchSize := max(0, opts.BatchSize)

	return &BatchWriter{
		w:          w,
		headers:    headers,
		colMap:     colMap,
		prefix:     prefix,
		batchSize:  batchSize,
		valuesOnly: opts.ValuesOnly,
	}, nil
}

func (bw *BatchWriter) writeString(s string) error {
	if _, err := io.WriteString(bw.w, s); err != nil {
		return fmt.Errorf("write sql: %w", err)
	}
	return nil
}

// WriteRow writes a single row of data to the SQL batch stream.
func (bw *BatchWriter) WriteRow(row []any) error {
	if bw.valuesOnly || bw.inBatchCount > 0 {
		if bw.totalRows > 0 {
			if err := bw.writeString(",\n"); err != nil {
				return err
			}
		}
	} else {
		if bw.totalRows > 0 {
			if err := bw.writeString("\n\n"); err != nil {
				return err
			}
		}
		if err := bw.writeString(bw.prefix); err != nil {
			return err
		}
	}

	valStr := FormatRowSQLSelected(row, bw.colMap.ColIndices, bw.headers, bw.colMap.NumColSet)
	if err := bw.writeString(valStr); err != nil {
		return err
	}
	bw.totalRows++

	if !bw.valuesOnly {
		bw.inBatchCount++
		if bw.batchSize > 0 && bw.inBatchCount == bw.batchSize {
			if err := bw.writeString(";"); err != nil {
				return err
			}
			bw.inBatchCount = 0
		}
	}

	return nil
}

// Close finalizes the SQL batch output, writing any trailing semicolon.
func (bw *BatchWriter) Close() error {
	if !bw.valuesOnly && bw.inBatchCount > 0 {
		if err := bw.writeString(";"); err != nil {
			return err
		}
		bw.inBatchCount = 0
	}

	return nil
}

// TotalRows returns the total number of rows written by this BatchWriter.
func (bw *BatchWriter) TotalRows() int {
	return bw.totalRows
}

// GenerateBatchSQL generates SQL INSERT statements (or raw VALUES) according to opts.
func GenerateBatchSQL(headers []string, dataRows [][]any, opts GenerateOptions) string {
	if len(dataRows) == 0 {
		return ""
	}

	var builder strings.Builder
	bw, err := NewBatchWriter(&builder, headers, opts)
	if err != nil {
		return ""
	}
	builder.Grow(len(dataRows) * len(bw.colMap.ValidSelectedCols) * 25)

	for _, row := range dataRows {
		if err := bw.WriteRow(row); err != nil {
			return ""
		}
	}

	if err := bw.Close(); err != nil {
		return ""
	}

	return builder.String()
}

// GenerateSQLValues generates SQL INSERT values from the data (legacy wrapper)
func GenerateSQLValues(headers []string, dataRows [][]any, numberColumns []string) string {
	return GenerateBatchSQL(headers, dataRows, GenerateOptions{
		SelectedColumns: headers,
		NumberColumns:   numberColumns,
		ValuesOnly:      true,
	})
}

// FormatRowSQL formats a single row into SQL values format (val1, val2, ...) using all headers
func FormatRowSQL(row []any, headers []string, numColSet map[string]bool) string {
	indices := make([]int, len(headers))
	for i := range headers {
		indices[i] = i
	}
	return FormatRowSQLSelected(row, indices, headers, numColSet)
}

// FormatRowSQLSelected formats selected columns of a row into SQL values format (val1, val2, ...)
func FormatRowSQLSelected(row []any, colIndices []int, headers []string, numColSet map[string]bool) string {
	var builder strings.Builder
	builder.WriteByte('(')

	for i, colIdx := range colIndices {
		if i > 0 {
			builder.WriteString(", ")
		}

		var cellValue any
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
