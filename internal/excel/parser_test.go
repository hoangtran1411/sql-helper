package excel

import (
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestSheetDataStructure(t *testing.T) {
	// Test SheetData struct initialization
	data := &SheetData{
		Headers:  []string{"Name", "Age", "Email"},
		DataRows: [][]interface{}{{"John", 30, "john@example.com"}},
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
		DataRows: [][]interface{}{},
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
		DataRows: [][]interface{}{
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
		DataRows: [][]interface{}{
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
	found := false
	for _, s := range sheets {
		if s == "Sheet1" {
			found = true
			break
		}
	}
	if !found {
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
