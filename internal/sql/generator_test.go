package sql

import (
	"strings"
	"testing"
)

func TestGenerateSQLValues(t *testing.T) {
	tests := []struct {
		name          string
		headers       []string
		dataRows      [][]interface{}
		numberColumns []string
		expected      string
	}{
		{
			name:          "empty rows",
			headers:       []string{"Name", "Age"},
			dataRows:      [][]interface{}{},
			numberColumns: []string{},
			expected:      "",
		},
		{
			name:          "single row no numeric",
			headers:       []string{"Name", "Email"},
			dataRows:      [][]interface{}{{"John", "john@example.com"}},
			numberColumns: []string{},
			expected:      "('John', 'john@example.com')",
		},
		{
			name:          "single row with numeric",
			headers:       []string{"Name", "Age"},
			dataRows:      [][]interface{}{{"John", 30}},
			numberColumns: []string{"Age"},
			expected:      "('John', 30)",
		},
		{
			name:    "multiple rows",
			headers: []string{"Name", "Age", "Active"},
			dataRows: [][]interface{}{
				{"John", 30, "TRUE"},
				{"Jane", 25, "FALSE"},
			},
			numberColumns: []string{"Age"},
			expected:      "('John', 30, 'TRUE'),\n('Jane', 25, 'FALSE')",
		},
		{
			name:    "null values",
			headers: []string{"Name", "Age", "Email"},
			dataRows: [][]interface{}{
				{"John", nil, nil},
			},
			numberColumns: []string{"Age"},
			expected:      "('John', NULL, '')",
		},
		{
			name:    "mixed types",
			headers: []string{"ID", "Name", "Score", "Date"},
			dataRows: [][]interface{}{
				{1, "Alice", 95.5, "2024-01-15"},
				{2, "Bob", nil, "2024-02-20"},
			},
			numberColumns: []string{"ID", "Score"},
			expected:      "(1, 'Alice', 95.5, '2024-01-15'),\n(2, 'Bob', NULL, '2024-02-20')",
		},
		{
			name:          "row shorter than headers",
			headers:       []string{"A", "B", "C"},
			dataRows:      [][]interface{}{{"X", "Y"}},
			numberColumns: []string{},
			expected:      "('X', 'Y', '')",
		},
		{
			name:          "quote escaping",
			headers:       []string{"Description"},
			dataRows:      [][]interface{}{{"It's O'Brien's book"}},
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
	dataRows := make([][]interface{}, 100)
	for i := 0; i < 100; i++ {
		dataRows[i] = []interface{}{i, "Name" + string(rune('A'+i%26)), float64(i) * 1.5}
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
		dataRows    [][]interface{}
		findValue   string
		replaceWith string
		expected    [][]interface{}
	}{
		{
			name: "simple replace",
			dataRows: [][]interface{}{
				{"Hello", "World"},
				{"Hello", "There"},
			},
			findValue:   "Hello",
			replaceWith: "Hi",
			expected: [][]interface{}{
				{"Hi", "World"},
				{"Hi", "There"},
			},
		},
		{
			name: "no match",
			dataRows: [][]interface{}{
				{"A", "B"},
			},
			findValue:   "X",
			replaceWith: "Y",
			expected: [][]interface{}{
				{"A", "B"},
			},
		},
		{
			name: "replace with empty",
			dataRows: [][]interface{}{
				{"Remove", "Keep"},
			},
			findValue:   "Remove",
			replaceWith: "",
			expected: [][]interface{}{
				{"", "Keep"},
			},
		},
		{
			name: "replace number as string",
			dataRows: [][]interface{}{
				{123, "text"},
			},
			findValue:   "123",
			replaceWith: "456",
			expected: [][]interface{}{
				{"456", "text"},
			},
		},
		{
			name: "nil handling",
			dataRows: [][]interface{}{
				{nil, "value"},
			},
			findValue:   "<nil>",
			replaceWith: "NULL",
			expected: [][]interface{}{
				{"NULL", "value"},
			},
		},
		{
			name:        "empty rows",
			dataRows:    [][]interface{}{},
			findValue:   "X",
			replaceWith: "Y",
			expected:    [][]interface{}{},
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
	dataRows := [][]interface{}{
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
