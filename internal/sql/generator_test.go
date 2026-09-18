package sql

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

func TestGenerateSQLValues(t *testing.T) {
	tests := []struct {
		name          string
		headers       []string
		dataRows      [][]any
		numberColumns []string
		expected      string
	}{
		{
			name:          "empty rows",
			headers:       []string{"Name", "Age"},
			dataRows:      [][]any{},
			numberColumns: []string{},
			expected:      "",
		},
		{
			name:          "single row no numeric",
			headers:       []string{"Name", "Email"},
			dataRows:      [][]any{{"John", "john@example.com"}},
			numberColumns: []string{},
			expected:      "('John', 'john@example.com')",
		},
		{
			name:          "single row with numeric",
			headers:       []string{"Name", "Age"},
			dataRows:      [][]any{{"John", 30}},
			numberColumns: []string{"Age"},
			expected:      "('John', 30)",
		},
		{
			name:    "multiple rows",
			headers: []string{"Name", "Age", "Active"},
			dataRows: [][]any{
				{"John", 30, "TRUE"},
				{"Jane", 25, "FALSE"},
			},
			numberColumns: []string{"Age"},
			expected:      "('John', 30, 'TRUE'),\n('Jane', 25, 'FALSE')",
		},
		{
			name:    "null values",
			headers: []string{"Name", "Age", "Email"},
			dataRows: [][]any{
				{"John", nil, nil},
			},
			numberColumns: []string{"Age"},
			expected:      "('John', NULL, '')",
		},
		{
			name:    "mixed types",
			headers: []string{"ID", "Name", "Score", "Date"},
			dataRows: [][]any{
				{1, "Alice", 95.5, "2024-01-15"},
				{2, "Bob", nil, "2024-02-20"},
			},
			numberColumns: []string{"ID", "Score"},
			expected:      "(1, 'Alice', 95.5, '2024-01-15'),\n(2, 'Bob', NULL, '2024-02-20')",
		},
		{
			name:          "row shorter than headers",
			headers:       []string{"A", "B", "C"},
			dataRows:      [][]any{{"X", "Y"}},
			numberColumns: []string{},
			expected:      "('X', 'Y', '')",
		},
		{
			name:          "quote escaping",
			headers:       []string{"Description"},
			dataRows:      [][]any{{"It's O'Brien's book"}},
			numberColumns: []string{},
			expected:      "('It''s O''Brien''s book')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSQLValues(tt.headers, tt.dataRows, tt.numberColumns)
			if result != tt.expected {
				t.Errorf("GenerateSQLValues() = %q; want %q", result, tt.expected)
			}
		})
	}
}

func TestGenerateSQLValues_LargeDataset(t *testing.T) {
	// Test with larger dataset to verify performance
	headers := []string{"ID", "Name", "Value"}
	dataRows := make([][]any, 100)
	for i := range 100 {
		dataRows[i] = []any{i, "Name" + string(rune('A'+i%26)), float64(i) * 1.5}
	}

	result := GenerateSQLValues(headers, dataRows, []string{"ID", "Value"})

	// Verify it starts and ends correctly
	if !strings.HasPrefix(result, "(0, 'NameA', 0)") {
		t.Errorf("Result should start with first row, got: %s", result[:50])
	}

	lines := strings.Split(result, "\n")
	if len(lines) != 100 {
		t.Errorf("Expected 100 lines, got %d", len(lines))
	}
}

func TestFindAndReplace(t *testing.T) {
	tests := []struct {
		name        string
		dataRows    [][]any
		findValue   string
		replaceWith string
		expected    [][]any
	}{
		{
			name: "simple replace",
			dataRows: [][]any{
				{"Hello", "World"},
				{"Hello", "There"},
			},
			findValue:   "Hello",
			replaceWith: "Hi",
			expected: [][]any{
				{"Hi", "World"},
				{"Hi", "There"},
			},
		},
		{
			name: "no match",
			dataRows: [][]any{
				{"A", "B"},
			},
			findValue:   "X",
			replaceWith: "Y",
			expected: [][]any{
				{"A", "B"},
			},
		},
		{
			name: "replace with empty",
			dataRows: [][]any{
				{"Remove", "Keep"},
			},
			findValue:   "Remove",
			replaceWith: "",
			expected: [][]any{
				{"", "Keep"},
			},
		},
		{
			name: "replace number as string",
			dataRows: [][]any{
				{123, "text"},
			},
			findValue:   "123",
			replaceWith: "456",
			expected: [][]any{
				{"456", "text"},
			},
		},
		{
			name: "nil handling",
			dataRows: [][]any{
				{nil, "value"},
			},
			findValue:   "<nil>",
			replaceWith: "NULL",
			expected: [][]any{
				{"NULL", "value"},
			},
		},
		{
			name:        "empty rows",
			dataRows:    [][]any{},
			findValue:   "X",
			replaceWith: "Y",
			expected:    [][]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindAndReplace(tt.dataRows, tt.findValue, tt.replaceWith)

			if len(result) != len(tt.expected) {
				t.Errorf("FindAndReplace() returned %d rows; want %d", len(result), len(tt.expected))
				return
			}

			for i, row := range result {
				for j, cell := range row {
					expectedCell := tt.expected[i][j]
					if cell != expectedCell {
						t.Errorf("FindAndReplace()[%d][%d] = %v; want %v", i, j, cell, expectedCell)
					}
				}
			}
		})
	}
}

func TestFormatIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"col", "col"},
		{"productioncode", "productioncode"},
		{"Order ID", "[Order ID]"},
		{"first name", "[first name]"},
		{"[pre_bracketed]", "[pre_bracketed]"},
		{"", ""},
		{"   ", ""},
	}

	for _, tt := range tests {
		result := FormatIdentifier(tt.input)
		if result != tt.expected {
			t.Errorf("FormatIdentifier(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestBuildInsertPrefix(t *testing.T) {
	tests := []struct {
		tableName string
		columns   []string
		expected  string
	}{
		{
			tableName: "users",
			columns:   []string{"id", "name"},
			expected:  "INSERT INTO users (id, name) VALUES\n",
		},
		{
			tableName: "",
			columns:   []string{"col1"},
			expected:  "INSERT INTO my_table (col1) VALUES\n",
		},
		{
			tableName: "Order Details",
			columns:   []string{"order id", "item code"},
			expected:  "INSERT INTO [Order Details] ([order id], [item code]) VALUES\n",
		},
	}

	for _, tt := range tests {
		result := BuildInsertPrefix(tt.tableName, tt.columns)
		if result != tt.expected {
			t.Errorf("BuildInsertPrefix(%q, %v) = %q; want %q", tt.tableName, tt.columns, result, tt.expected)
		}
	}
}

func TestGenerateBatchSQL(t *testing.T) {
	headers := []string{"id", "name", "age", "status"}
	dataRows := [][]any{
		{1, "Alice", 25, "active"},
		{2, "Bob", 30, "pending"},
		{3, "Charlie", 35, "active"},
	}

	t.Run("default batch all columns single statement", func(t *testing.T) {
		opts := GenerateOptions{
			TableName:     "users",
			NumberColumns: []string{"id", "age"},
			BatchSize:     1000,
		}
		expected := "INSERT INTO users (id, name, age, status) VALUES\n" +
			"(1, 'Alice', 25, 'active'),\n" +
			"(2, 'Bob', 30, 'pending'),\n" +
			"(3, 'Charlie', 35, 'active');"

		result := GenerateBatchSQL(headers, dataRows, opts)
		if result != expected {
			t.Errorf("GenerateBatchSQL() =\n%s\nwant:\n%s", result, expected)
		}
	})

	t.Run("chunked batches with batchSize 2", func(t *testing.T) {
		opts := GenerateOptions{
			TableName:     "users",
			NumberColumns: []string{"id", "age"},
			BatchSize:     2,
		}
		expected := "INSERT INTO users (id, name, age, status) VALUES\n" +
			"(1, 'Alice', 25, 'active'),\n" +
			"(2, 'Bob', 30, 'pending');\n\n" +
			"INSERT INTO users (id, name, age, status) VALUES\n" +
			"(3, 'Charlie', 35, 'active');"

		result := GenerateBatchSQL(headers, dataRows, opts)
		if result != expected {
			t.Errorf("GenerateBatchSQL() =\n%s\nwant:\n%s", result, expected)
		}
	})

	t.Run("column selection subset", func(t *testing.T) {
		opts := GenerateOptions{
			TableName:       "production",
			SelectedColumns: []string{"name", "status"},
			BatchSize:       1000,
		}
		expected := "INSERT INTO production (name, status) VALUES\n" +
			"('Alice', 'active'),\n" +
			"('Bob', 'pending'),\n" +
			"('Charlie', 'active');"

		result := GenerateBatchSQL(headers, dataRows, opts)
		if result != expected {
			t.Errorf("GenerateBatchSQL() =\n%s\nwant:\n%s", result, expected)
		}
	})

	t.Run("values only mode", func(t *testing.T) {
		opts := GenerateOptions{
			SelectedColumns: []string{"id", "name"},
			NumberColumns:   []string{"id"},
			ValuesOnly:      true,
		}
		expected := "(1, 'Alice'),\n(2, 'Bob'),\n(3, 'Charlie')"

		result := GenerateBatchSQL(headers, dataRows, opts)
		if result != expected {
			t.Errorf("GenerateBatchSQL() =\n%s\nwant:\n%s", result, expected)
		}
	})

	t.Run("empty rows", func(t *testing.T) {
		opts := GenerateOptions{
			TableName: "users",
		}
		result := GenerateBatchSQL(headers, nil, opts)
		if result != "" {
			t.Errorf("Expected empty string for nil rows, got %q", result)
		}
	})

	t.Run("invalid or non-matching selected columns", func(t *testing.T) {
		opts := GenerateOptions{
			TableName:       "users",
			SelectedColumns: []string{"nonexistent"},
		}
		result := GenerateBatchSQL(headers, dataRows, opts)
		if result != "" {
			t.Errorf("Expected empty string for non-matching columns, got %q", result)
		}
	})
}

func TestFormatRowSQL(t *testing.T) {
	headers := []string{"id", "name", "price"}
	row := []any{1, "Widget", 19.99}
	numCols := map[string]bool{"id": true, "price": true}

	result := FormatRowSQL(row, headers, numCols)
	expected := "(1, 'Widget', 19.99)"
	if result != expected {
		t.Errorf("FormatRowSQL() = %q, want %q", result, expected)
	}
}

func TestResolveColumns(t *testing.T) {
	headers := []string{"id", "name", "age", "status"}

	t.Run("default to all headers when selected is empty", func(t *testing.T) {
		colMap, err := ResolveColumns(headers, nil, []string{"id", "age"})
		if err != nil {
			t.Fatalf("ResolveColumns() returned error: %v", err)
		}
		if len(colMap.ColIndices) != 4 {
			t.Errorf("expected 4 column indices, got %d", len(colMap.ColIndices))
		}
		if !colMap.NumColSet["id"] || !colMap.NumColSet["age"] || colMap.NumColSet["name"] {
			t.Errorf("unexpected NumColSet: %v", colMap.NumColSet)
		}
	})

	t.Run("subset selection", func(t *testing.T) {
		colMap, err := ResolveColumns(headers, []string{"name", "status"}, nil)
		if err != nil {
			t.Fatalf("ResolveColumns() returned error: %v", err)
		}
		if len(colMap.ColIndices) != 2 || colMap.ColIndices[0] != 1 || colMap.ColIndices[1] != 3 {
			t.Errorf("unexpected colIndices: %v", colMap.ColIndices)
		}
	})

	t.Run("empty headers and selected columns error", func(t *testing.T) {
		_, err := ResolveColumns(nil, nil, nil)
		if err == nil {
			t.Error("expected error for empty headers and selected columns, got nil")
		}
	})

	t.Run("no matching columns error", func(t *testing.T) {
		_, err := ResolveColumns(headers, []string{"nonexistent"}, nil)
		if err == nil {
			t.Error("expected error for non-matching selected columns, got nil")
		}
	})
}

func TestNewBatchWriter(t *testing.T) {
	headers := []string{"id", "name"}

	tests := []struct {
		name      string
		w         io.Writer
		headers   []string
		opts      GenerateOptions
		expectErr bool
	}{
		{
			name:      "nil writer",
			w:         nil,
			headers:   headers,
			opts:      GenerateOptions{},
			expectErr: true,
		},
		{
			name:      "empty headers",
			w:         &bytes.Buffer{},
			headers:   nil,
			opts:      GenerateOptions{},
			expectErr: true,
		},
		{
			name:    "invalid selected columns",
			w:       &bytes.Buffer{},
			headers: headers,
			opts: GenerateOptions{
				SelectedColumns: []string{"unknown_col"},
			},
			expectErr: true,
		},
		{
			name:    "valid writer and headers",
			w:       &bytes.Buffer{},
			headers: headers,
			opts: GenerateOptions{
				TableName: "users",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bw, err := NewBatchWriter(tt.w, tt.headers, tt.opts)
			if (err != nil) != tt.expectErr {
				t.Fatalf("NewBatchWriter() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !tt.expectErr && bw == nil {
				t.Fatal("expected non-nil BatchWriter")
			}
		})
	}
}

func TestBatchWriter_TableDriven(t *testing.T) {
	headers := []string{"id", "name", "price"}

	tests := []struct {
		name     string
		opts     GenerateOptions
		dataRows [][]any
		expected string
	}{
		{
			name: "single row default batch",
			opts: GenerateOptions{
				TableName:     "items",
				NumberColumns: []string{"id", "price"},
				BatchSize:     10,
			},
			dataRows: [][]any{
				{1, "Apple", 1.25},
			},
			expected: "INSERT INTO items (id, name, price) VALUES\n" +
				"(1, 'Apple', 1.25);",
		},
		{
			name: "multiple rows single batch",
			opts: GenerateOptions{
				TableName:     "items",
				NumberColumns: []string{"id", "price"},
				BatchSize:     5,
			},
			dataRows: [][]any{
				{1, "Apple", 1.25},
				{2, "Banana", 0.75},
			},
			expected: "INSERT INTO items (id, name, price) VALUES\n" +
				"(1, 'Apple', 1.25),\n" +
				"(2, 'Banana', 0.75);",
		},
		{
			name: "exact multiple of batch size",
			opts: GenerateOptions{
				TableName:     "items",
				NumberColumns: []string{"id", "price"},
				BatchSize:     2,
			},
			dataRows: [][]any{
				{1, "Apple", 1.25},
				{2, "Banana", 0.75},
			},
			expected: "INSERT INTO items (id, name, price) VALUES\n" +
				"(1, 'Apple', 1.25),\n" +
				"(2, 'Banana', 0.75);",
		},
		{
			name: "chunked batches across boundary",
			opts: GenerateOptions{
				TableName:     "items",
				NumberColumns: []string{"id", "price"},
				BatchSize:     2,
			},
			dataRows: [][]any{
				{1, "Apple", 1.25},
				{2, "Banana", 0.75},
				{3, "Cherry", 2.5},
			},
			expected: "INSERT INTO items (id, name, price) VALUES\n" +
				"(1, 'Apple', 1.25),\n" +
				"(2, 'Banana', 0.75);\n\n" +
				"INSERT INTO items (id, name, price) VALUES\n" +
				"(3, 'Cherry', 2.5);",
		},
		{
			name: "batch size 1",
			opts: GenerateOptions{
				TableName:     "items",
				NumberColumns: []string{"id", "price"},
				BatchSize:     1,
			},
			dataRows: [][]any{
				{1, "Apple", 1.25},
				{2, "Banana", 0.75},
			},
			expected: "INSERT INTO items (id, name, price) VALUES\n" +
				"(1, 'Apple', 1.25);\n\n" +
				"INSERT INTO items (id, name, price) VALUES\n" +
				"(2, 'Banana', 0.75);",
		},
		{
			name: "batch size 0 unlimited",
			opts: GenerateOptions{
				TableName:     "items",
				NumberColumns: []string{"id", "price"},
				BatchSize:     0,
			},
			dataRows: [][]any{
				{1, "Apple", 1.25},
				{2, "Banana", 0.75},
			},
			expected: "INSERT INTO items (id, name, price) VALUES\n" +
				"(1, 'Apple', 1.25),\n" +
				"(2, 'Banana', 0.75);",
		},
		{
			name: "values only mode",
			opts: GenerateOptions{
				ValuesOnly:    true,
				NumberColumns: []string{"id", "price"},
			},
			dataRows: [][]any{
				{1, "Apple", 1.25},
				{2, "Banana", 0.75},
			},
			expected: "(1, 'Apple', 1.25),\n" +
				"(2, 'Banana', 0.75)",
		},
		{
			name: "column selection subset",
			opts: GenerateOptions{
				TableName:       "items",
				SelectedColumns: []string{"name"},
				BatchSize:       10,
			},
			dataRows: [][]any{
				{1, "Apple", 1.25},
			},
			expected: "INSERT INTO items (name) VALUES\n" +
				"('Apple');",
		},
		{
			name: "zero rows written",
			opts: GenerateOptions{
				TableName: "items",
			},
			dataRows: nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bw, err := NewBatchWriter(&buf, headers, tt.opts)
			if err != nil {
				t.Fatalf("NewBatchWriter() error = %v", err)
			}

			for _, row := range tt.dataRows {
				if err := bw.WriteRow(row); err != nil {
					t.Errorf("WriteRow() error = %v", err)
					break
				}
			}

			if err := bw.Close(); err != nil {
				t.Fatalf("Close() error = %v", err)
			}

			result := buf.String()
			if result != tt.expected {
				t.Errorf("output =\n%q\nwant:\n%q", result, tt.expected)
			}
		})
	}
}

func TestBatchWriter_TotalRows(t *testing.T) {
	headers := []string{"id", "name"}
	var buf bytes.Buffer
	bw, err := NewBatchWriter(&buf, headers, GenerateOptions{TableName: "users"})
	if err != nil {
		t.Fatalf("NewBatchWriter() error = %v", err)
	}

	if bw.TotalRows() != 0 {
		t.Errorf("expected 0 total rows, got %d", bw.TotalRows())
	}

	_ = bw.WriteRow([]any{1, "Alice"}) //nolint:errcheck // total rows counter test; write error tested separately
	if bw.TotalRows() != 1 {
		t.Errorf("expected 1 total rows, got %d", bw.TotalRows())
	}

	_ = bw.WriteRow([]any{2, "Bob"}) //nolint:errcheck // total rows counter test; write error tested separately
	if bw.TotalRows() != 2 {
		t.Errorf("expected 2 total rows, got %d", bw.TotalRows())
	}

	_ = bw.Close() //nolint:errcheck // total rows persistence test; close error tested separately
	if bw.TotalRows() != 2 {
		t.Errorf("expected 2 total rows after close, got %d", bw.TotalRows())
	}
}

func TestBatchWriter_BufioFlush(t *testing.T) {
	headers := []string{"id", "name"}
	var buf bytes.Buffer
	bufWriter := bufio.NewWriter(&buf)

	bw, err := NewBatchWriter(bufWriter, headers, GenerateOptions{
		TableName:     "users",
		NumberColumns: []string{"id"},
	})
	if err != nil {
		t.Fatalf("NewBatchWriter() error = %v", err)
	}

	if err := bw.WriteRow([]any{1, "Alice"}); err != nil {
		t.Fatalf("WriteRow() error = %v", err)
	}

	// Close finalizes SQL without flushing the buffer
	if err := bw.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	// Buffer flushing is caller's responsibility
	if err := bufWriter.Flush(); err != nil {
		t.Fatalf("bufWriter.Flush() error = %v", err)
	}

	expected := "INSERT INTO users (id, name) VALUES\n(1, 'Alice');"
	if buf.String() != expected {
		t.Errorf("buf.String() = %q, want %q", buf.String(), expected)
	}
}

type flushErrorWriter struct {
	err error
}

func (f *flushErrorWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (f *flushErrorWriter) Flush() error {
	return f.err
}

func TestBatchWriter_FlushError(t *testing.T) {
	headers := []string{"id", "name"}
	expectedErr := errors.New("simulated flush failure")
	fw := &flushErrorWriter{err: expectedErr}

	bw, err := NewBatchWriter(fw, headers, GenerateOptions{TableName: "users"})
	if err != nil {
		t.Fatalf("NewBatchWriter() error = %v", err)
	}

	if err := bw.WriteRow([]any{1, "Alice"}); err != nil {
		t.Fatalf("WriteRow() error = %v", err)
	}

	// Close should not flush; caller is responsible for flushing buffer.
	if err := bw.Close(); err != nil {
		t.Errorf("Close() error = %v, expected nil as Close does not flush", err)
	}

	// Calling Flush on writer directly returns the simulated error.
	if err := fw.Flush(); !errors.Is(err, expectedErr) {
		t.Errorf("fw.Flush() error = %v, want %v", err, expectedErr)
	}
}

type errorWriter struct {
	err error
}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, e.err
}

type errorAfterNWriter struct {
	limit int
	count int
	err   error
}

func (e *errorAfterNWriter) Write(p []byte) (n int, err error) {
	if e.count >= e.limit {
		return 0, e.err
	}
	e.count++
	return len(p), nil
}

func TestBatchWriter_WriteErrors(t *testing.T) {
	headers := []string{"id", "name"}
	expectedErr := errors.New("write failure")

	t.Run("WriteRow error on initial write", func(t *testing.T) {
		ew := &errorWriter{err: expectedErr}
		bw, err := NewBatchWriter(ew, headers, GenerateOptions{TableName: "users"})
		if err != nil {
			t.Fatalf("NewBatchWriter() error = %v", err)
		}

		err = bw.WriteRow([]any{1, "Alice"})
		if err == nil {
			t.Error("expected error on WriteRow, got nil")
		}
	})

	t.Run("WriteRow error in valuesOnly mode", func(t *testing.T) {
		ew := &errorWriter{err: expectedErr}
		bw, err := NewBatchWriter(ew, headers, GenerateOptions{ValuesOnly: true})
		if err != nil {
			t.Fatalf("NewBatchWriter() error = %v", err)
		}

		err = bw.WriteRow([]any{1, "Alice"})
		if err == nil {
			t.Error("expected error on WriteRow, got nil")
		}
	})

	t.Run("Close error writing semicolon", func(t *testing.T) {
		// Allows initial writes, then fails when writing semicolon in Close
		ew := &errorAfterNWriter{limit: 2, err: expectedErr}
		bw, err := NewBatchWriter(ew, headers, GenerateOptions{TableName: "users"})
		if err != nil {
			t.Fatalf("NewBatchWriter() error = %v", err)
		}

		_ = bw.WriteRow([]any{1, "Alice"}) // error is verified on Close

		err = bw.Close()
		if err == nil {
			t.Error("expected error on Close writing semicolon, got nil")
		}
	})
}

func BenchmarkGenerateBatchSQL(b *testing.B) {
	headers := []string{"id", "name", "age", "status"}
	dataRows := make([][]any, 1000)
	for i := range 1000 {
		dataRows[i] = []any{i, "Alice", 30, "active"}
	}
	opts := GenerateOptions{
		TableName:     "users",
		NumberColumns: []string{"id", "age"},
		BatchSize:     500,
	}

	b.ResetTimer()
	for b.Loop() {
		_ = GenerateBatchSQL(headers, dataRows, opts) //nolint:errcheck // benchmark execution result intentionally discarded
	}
}
