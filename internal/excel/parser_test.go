package excel

import (
	"testing"
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

// Note: Integration tests for ParseExcelFile and ProcessSheet
// require actual Excel files and are better suited for E2E tests.
// These tests focus on the data structures and logic that can be
// unit tested without file I/O.

func TestParseExcelFile_InvalidPath(t *testing.T) {
	_, err := ParseExcelFile("nonexistent_file.xlsx")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestProcessSheet_InvalidPath(t *testing.T) {
	_, err := ProcessSheet("nonexistent_file.xlsx", "Sheet1")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}
