package excel

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestSheetDataStructure(t *testing.T) {
	// Test SheetData struct initialization
	data := &SheetData{
		Headers:  []string{"Name", "Age", "Email"},
		DataRows: [][]any{{"John", 30, "john@example.com"}},
	}

	if len(data.Headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(data.Headers))
	}

	if len(data.DataRows) != 1 {
		t.Errorf("Expected 1 data row, got %d", len(data.DataRows))
	}

	if data.Headers[0] != "Name" {
		t.Errorf("Expected first header to be 'Name', got '%s'", data.Headers[0])
	}
}

func TestSheetDataEmpty(t *testing.T) {
	data := &SheetData{
		Headers:  []string{},
		DataRows: [][]any{},
	}

	if len(data.Headers) != 0 {
		t.Errorf("Expected 0 headers, got %d", len(data.Headers))
	}

	if len(data.DataRows) != 0 {
		t.Errorf("Expected 0 data rows, got %d", len(data.DataRows))
	}
}

func TestSheetDataWithNilValues(t *testing.T) {
	data := &SheetData{
		Headers: []string{"A", "B", "C"},
		DataRows: [][]any{
			{nil, "value", nil},
			{"x", nil, "z"},
		},
	}

	if len(data.DataRows) != 2 {
		t.Errorf("Expected 2 data rows, got %d", len(data.DataRows))
	}

	// Check nil handling
	if data.DataRows[0][0] != nil {
		t.Errorf("Expected nil at [0][0], got %v", data.DataRows[0][0])
	}

	if data.DataRows[0][1] != "value" {
		t.Errorf("Expected 'value' at [0][1], got %v", data.DataRows[0][1])
	}
}

func TestSheetDataMultipleTypes(t *testing.T) {
	data := &SheetData{
		Headers: []string{"String", "Int", "Float", "Bool"},
		DataRows: [][]any{
			{"text", 42, 3.14, true},
			{"more", 100, 2.5, false},
		},
	}

	// Verify types are preserved
	row := data.DataRows[0]

	if _, ok := row[0].(string); !ok {
		t.Errorf("Expected string at index 0, got %T", row[0])
	}

	if _, ok := row[1].(int); !ok {
		t.Errorf("Expected int at index 1, got %T", row[1])
	}

	if _, ok := row[2].(float64); !ok {
		t.Errorf("Expected float64 at index 2, got %T", row[2])
	}

	if _, ok := row[3].(bool); !ok {
		t.Errorf("Expected bool at index 3, got %T", row[3])
	}
}

func TestParseExcelFile_InvalidPath(t *testing.T) {
	_, err := ParseExcelFile("nonexistent_file.xlsx")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestGetPreview_InvalidPath(t *testing.T) {
	_, err := GetPreview("nonexistent_file.xlsx", "Sheet1", 100)
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

// createTestExcelFile creates a temporary Excel file for testing
func createTestExcelFile(t *testing.T, dir string, data [][]string) string {
	t.Helper()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			t.Logf("Warning: failed to close file: %v", err)
		}
	}()

	// Write data to Sheet1
	for i, row := range data {
		for j, cell := range row {
			cellName, _ := excelize.CoordinatesToCellName(j+1, i+1)
			if err := f.SetCellValue("Sheet1", cellName, cell); err != nil {
				t.Fatalf("Failed to set cell value: %v", err)
			}
		}
	}

	filePath := filepath.Join(dir, "test.xlsx")
	if err := f.SaveAs(filePath); err != nil {
		t.Fatalf("Failed to save test Excel file: %v", err)
	}
	return filePath
}

func TestParseExcelFile_ValidFile(t *testing.T) {
	tempDir := t.TempDir()
	testData := [][]string{
		{"Name", "Age", "Email"},
		{"John", "30", "john@example.com"},
		{"Jane", "25", "jane@example.com"},
	}
	filePath := createTestExcelFile(t, tempDir, testData)

	sheets, err := ParseExcelFile(filePath)
	if err != nil {
		t.Fatalf("ParseExcelFile failed: %v", err)
	}

	if len(sheets) == 0 {
		t.Error("Expected at least one sheet")
	}

	// Default sheet name should be "Sheet1"
	if !slices.Contains(sheets, "Sheet1") {
		t.Errorf("Expected 'Sheet1' in sheets, got %v", sheets)
	}
}

func TestGetPreview_ValidFile(t *testing.T) {
	tempDir := t.TempDir()
	testData := [][]string{
		{"ID", "Name", "Value"},
		{"1", "Alice", "100"},
		{"2", "Bob", "200"},
		{"3", "Charlie", "300"},
	}
	filePath := createTestExcelFile(t, tempDir, testData)

	data, err := GetPreview(filePath, "Sheet1", 100)
	if err != nil {
		t.Fatalf("GetPreview failed: %v", err)
	}

	// Check headers
	if len(data.Headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(data.Headers))
	}
	if data.Headers[0] != "ID" || data.Headers[1] != "Name" || data.Headers[2] != "Value" {
		t.Errorf("Unexpected headers: %v", data.Headers)
	}

	// Check data rows
	if len(data.DataRows) != 3 {
		t.Errorf("Expected 3 data rows, got %d", len(data.DataRows))
	}

	// Check first row values
	if data.DataRows[0][0] != "1" {
		t.Errorf("Expected '1' at [0][0], got %v", data.DataRows[0][0])
	}
	if data.DataRows[0][1] != "Alice" {
		t.Errorf("Expected 'Alice' at [0][1], got %v", data.DataRows[0][1])
	}
}

func TestGetPreview_EmptySheet(t *testing.T) {
	tempDir := t.TempDir()
	testData := [][]string{} // Empty
	filePath := createTestExcelFile(t, tempDir, testData)

	data, err := GetPreview(filePath, "Sheet1", 100)
	if err != nil {
		t.Fatalf("GetPreview failed: %v", err)
	}

	if len(data.Headers) != 0 {
		t.Errorf("Expected 0 headers for empty sheet, got %d", len(data.Headers))
	}
	if len(data.DataRows) != 0 {
		t.Errorf("Expected 0 data rows for empty sheet, got %d", len(data.DataRows))
	}
}

func TestGetPreview_HeaderOnly(t *testing.T) {
	tempDir := t.TempDir()
	testData := [][]string{
		{"A", "B", "C"}, // Only headers, no data
	}
	filePath := createTestExcelFile(t, tempDir, testData)

	data, err := GetPreview(filePath, "Sheet1", 100)
	if err != nil {
		t.Fatalf("GetPreview failed: %v", err)
	}

	if len(data.Headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(data.Headers))
	}
	if len(data.DataRows) != 0 {
		t.Errorf("Expected 0 data rows (header only), got %d", len(data.DataRows))
	}
}

func TestGetPreview_InvalidSheetName(t *testing.T) {
	tempDir := t.TempDir()
	testData := [][]string{{"A", "B"}}
	filePath := createTestExcelFile(t, tempDir, testData)

	_, err := GetPreview(filePath, "NonExistentSheet", 100)
	if err == nil {
		t.Error("Expected error for non-existent sheet, got nil")
	}
}

func TestGetPreview_SparseData(t *testing.T) {
	tempDir := t.TempDir()
	// Sparse data: some cells are empty
	testData := [][]string{
		{"A", "B", "C"},
		{"1", "", "3"},  // B is empty
		{"", "2", ""},   // A and C are empty
		{"x", "y", "z"}, // All filled
	}
	filePath := createTestExcelFile(t, tempDir, testData)

	data, err := GetPreview(filePath, "Sheet1", 100)
	if err != nil {
		t.Fatalf("GetPreview failed: %v", err)
	}

	// Check that empty cells are nil
	if data.DataRows[0][1] != nil {
		t.Errorf("Expected nil at [0][1], got %v", data.DataRows[0][1])
	}
	if data.DataRows[1][0] != nil {
		t.Errorf("Expected nil at [1][0], got %v", data.DataRows[1][0])
	}

	// Non-empty cells should have values
	if data.DataRows[0][0] != "1" {
		t.Errorf("Expected '1' at [0][0], got %v", data.DataRows[0][0])
	}
}

func TestIterateSheet(t *testing.T) {
	tempDir := t.TempDir()
	testData := [][]string{
		{"ID", "Name", "City"},
		{"1", "Alice", "Hanoi"},
		{"2", "Bob", "Danang"},
		{"3", "Charlie", "Saigon"},
	}
	filePath := createTestExcelFile(t, tempDir, testData)

	t.Run("successful streaming", func(t *testing.T) {
		var collected [][]string
		err := IterateSheet(filePath, "Sheet1", func(row []string) error {
			collected = append(collected, row)
			return nil
		})
		if err != nil {
			t.Fatalf("IterateSheet failed: %v", err)
		}
		if len(collected) != 3 {
			t.Fatalf("Expected 3 rows, got %d", len(collected))
		}
		if collected[0][1] != "Alice" || collected[2][2] != "Saigon" {
			t.Errorf("Unexpected row content: %v", collected)
		}
	})

	t.Run("callback error terminates iteration", func(t *testing.T) {
		count := 0
		customErr := excelize.ErrSheetNotExist{SheetName: "test"}
		err := IterateSheet(filePath, "Sheet1", func(row []string) error {
			count++
			return customErr
		})
		if err == nil {
			t.Fatal("Expected error from callback, got nil")
		}
		if count != 1 {
			t.Errorf("Expected exactly 1 iteration before error, got %d", count)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		err := IterateSheet("nonexistent.xlsx", "Sheet1", func(row []string) error {
			return nil
		})
		if err == nil {
			t.Error("Expected error for nonexistent file, got nil")
		}
	})

	t.Run("nonexistent sheet", func(t *testing.T) {
		err := IterateSheet(filePath, "NonExistentSheet", func(row []string) error {
			return nil
		})
		if err == nil {
			t.Error("Expected error for nonexistent sheet, got nil")
		}
	})
}

func TestNormalizeRow(t *testing.T) {
	tests := []struct {
		name         string
		rawRow       []string
		numHeaders   int
		replacements []Replacement
		expected     []any
	}{
		{
			name:         "exact length without replacements",
			rawRow:       []string{"Alice", "30", "Paris"},
			numHeaders:   3,
			replacements: nil,
			expected:     []any{"Alice", "30", "Paris"},
		},
		{
			name:         "empty strings converted to nil",
			rawRow:       []string{"Alice", "", "Paris"},
			numHeaders:   3,
			replacements: nil,
			expected:     []any{"Alice", nil, "Paris"},
		},
		{
			name:         "shorter raw row padded with nil",
			rawRow:       []string{"Alice"},
			numHeaders:   3,
			replacements: nil,
			expected:     []any{"Alice", nil, nil},
		},
		{
			name:         "longer raw row truncated to numHeaders",
			rawRow:       []string{"Alice", "30", "Paris", "extra"},
			numHeaders:   3,
			replacements: nil,
			expected:     []any{"Alice", "30", "Paris"},
		},
		{
			name:       "apply replacements",
			rawRow:     []string{"N/A", "active", "NULL"},
			numHeaders: 3,
			replacements: []Replacement{
				{Find: "N/A", Replace: "Unknown"},
				{Find: "NULL", Replace: ""},
			},
			expected: []any{"Unknown", "active", nil},
		},
		{
			name:       "chained replacements",
			rawRow:     []string{"foo"},
			numHeaders: 1,
			replacements: []Replacement{
				{Find: "foo", Replace: "bar"},
				{Find: "bar", Replace: "baz"},
			},
			expected: []any{"baz"},
		},
		{
			name:         "empty raw row with numHeaders",
			rawRow:       []string{},
			numHeaders:   2,
			replacements: nil,
			expected:     []any{nil, nil},
		},
		{
			name:         "zero numHeaders",
			rawRow:       []string{"foo", "bar"},
			numHeaders:   0,
			replacements: nil,
			expected:     []any{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeRow(tc.rawRow, tc.numHeaders, tc.replacements)
			if len(got) != len(tc.expected) {
				t.Fatalf("expected length %d, got %d", len(tc.expected), len(got))
			}
			for i := range got {
				if got[i] != tc.expected[i] {
					t.Errorf("at index %d: expected %v (%T), got %v (%T)", i, tc.expected[i], tc.expected[i], got[i], got[i])
				}
			}
		})
	}
}

func TestFindAndReplace(t *testing.T) {
	t.Run("nil input returns empty slice", func(t *testing.T) {
		res := FindAndReplace(nil, "foo", "bar")
		if res == nil || len(res) != 0 {
			t.Errorf("expected empty non-nil slice, got %v", res)
		}
	})

	t.Run("empty input returns empty slice", func(t *testing.T) {
		input := [][]any{}
		res := FindAndReplace(input, "foo", "bar")
		if len(res) != 0 {
			t.Errorf("expected 0 rows, got %d", len(res))
		}
	})

	t.Run("replaces matching string and preserves non-matching", func(t *testing.T) {
		input := [][]any{
			{"hello", "world", 123},
			{nil, "hello", true},
		}
		res := FindAndReplace(input, "hello", "hi")
		if len(res) != 2 {
			t.Fatalf("expected 2 rows, got %d", len(res))
		}
		if res[0][0] != "hi" {
			t.Errorf("expected res[0][0] == 'hi', got %v", res[0][0])
		}
		if res[0][1] != "world" {
			t.Errorf("expected res[0][1] == 'world', got %v", res[0][1])
		}
		if res[0][2] != 123 {
			t.Errorf("expected res[0][2] == 123, got %v", res[0][2])
		}
		if res[1][0] != nil {
			t.Errorf("expected res[1][0] == nil, got %v", res[1][0])
		}
		if res[1][1] != "hi" {
			t.Errorf("expected res[1][1] == 'hi', got %v", res[1][1])
		}
		if res[1][2] != true {
			t.Errorf("expected res[1][2] == true, got %v", res[1][2])
		}
	})

	t.Run("replaces matching formatted number", func(t *testing.T) {
		input := [][]any{
			{100, "200"},
		}
		res := FindAndReplace(input, "100", "replaced")
		if res[0][0] != "replaced" {
			t.Errorf("expected 'replaced', got %v", res[0][0])
		}
		if res[0][1] != "200" {
			t.Errorf("expected '200', got %v", res[0][1])
		}
	})
}

func TestStreamRows(t *testing.T) {
	tempDir := t.TempDir()
	testData := [][]string{
		{"ID", "Name", "Status"},
		{"1", "Alice", "pending"},
		{"2", "Bob", ""},
		{"3", "Charlie", "active"},
	}
	filePath := createTestExcelFile(t, tempDir, testData)

	t.Run("successful streaming with normalization and replacements", func(t *testing.T) {
		headers := []string{"ID", "Name", "Status"}
		replacements := []Replacement{
			{Find: "pending", Replace: "queued"},
		}

		var streamedRows [][]any
		err := StreamRows(filePath, "Sheet1", headers, replacements, func(row []any) error {
			streamedRows = append(streamedRows, row)
			return nil
		})
		if err != nil {
			t.Fatalf("StreamRows failed: %v", err)
		}

		if len(streamedRows) != 3 {
			t.Fatalf("expected 3 rows, got %d", len(streamedRows))
		}

		// Row 1: "1", "Alice", "queued" (replacement applied)
		if streamedRows[0][0] != "1" || streamedRows[0][1] != "Alice" || streamedRows[0][2] != "queued" {
			t.Errorf("unexpected row 0: %v", streamedRows[0])
		}

		// Row 2: "2", "Bob", nil (empty string normalized to nil)
		if streamedRows[1][0] != "2" || streamedRows[1][1] != "Bob" || streamedRows[1][2] != nil {
			t.Errorf("unexpected row 1: %v", streamedRows[1])
		}

		// Row 3: "3", "Charlie", "active"
		if streamedRows[2][0] != "3" || streamedRows[2][1] != "Charlie" || streamedRows[2][2] != "active" {
			t.Errorf("unexpected row 2: %v", streamedRows[2])
		}
	})

	t.Run("streaming with row padding when raw row is shorter than headers", func(t *testing.T) {
		headers := []string{"ID", "Name", "Status", "Extra"}
		var streamedRows [][]any
		err := StreamRows(filePath, "Sheet1", headers, nil, func(row []any) error {
			streamedRows = append(streamedRows, row)
			return nil
		})
		if err != nil {
			t.Fatalf("StreamRows failed: %v", err)
		}

		for i, row := range streamedRows {
			if len(row) != 4 {
				t.Errorf("row %d expected length 4, got %d", i, len(row))
			}
			if row[3] != nil {
				t.Errorf("row %d expected nil at index 3, got %v", i, row[3])
			}
		}
	})

	t.Run("callback error terminates streaming", func(t *testing.T) {
		headers := []string{"ID", "Name", "Status"}
		count := 0
		customErr := excelize.ErrSheetNotExist{SheetName: "test"}
		err := StreamRows(filePath, "Sheet1", headers, nil, func(row []any) error {
			count++
			return customErr
		})
		if err == nil {
			t.Fatal("expected error from callback, got nil")
		}
		if count != 1 {
			t.Errorf("expected exactly 1 iteration before error, got %d", count)
		}
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		headers := []string{"ID", "Name"}
		err := StreamRows("nonexistent.xlsx", "Sheet1", headers, nil, func(row []any) error {
			return nil
		})
		if err == nil {
			t.Error("expected error for nonexistent file, got nil")
		}
	})

	t.Run("nonexistent sheet returns error", func(t *testing.T) {
		headers := []string{"ID", "Name"}
		err := StreamRows(filePath, "NonExistentSheet", headers, nil, func(row []any) error {
			return nil
		})
		if err == nil {
			t.Error("expected error for nonexistent sheet, got nil")
		}
	})
}
