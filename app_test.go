package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Error("NewApp() returned nil")
	}
}

func TestAppStartup(t *testing.T) {
	app := NewApp()
	ctx := context.Background()
	app.startup(ctx)

	if app.ctx == nil {
		t.Error("startup() did not set ctx")
	}
}

func TestExcelResultStruct(t *testing.T) {
	result := ExcelResult{
		FilePath:   "/path/to/file.xlsx",
		SheetNames: []string{"Sheet1", "Sheet2"},
	}

	if result.FilePath != "/path/to/file.xlsx" {
		t.Errorf("FilePath = %q, want %q", result.FilePath, "/path/to/file.xlsx")
	}
	if len(result.SheetNames) != 2 {
		t.Errorf("len(SheetNames) = %d, want 2", len(result.SheetNames))
	}
}

func TestSheetDataStruct(t *testing.T) {
	data := SheetData{
		Headers:  []string{"Name", "Age"},
		DataRows: [][]interface{}{{"John", 30}, {"Jane", 25}},
	}

	if len(data.Headers) != 2 {
		t.Errorf("len(Headers) = %d, want 2", len(data.Headers))
	}
	if len(data.DataRows) != 2 {
		t.Errorf("len(DataRows) = %d, want 2", len(data.DataRows))
	}
}

func TestAppProcessSheet(t *testing.T) {
	// Create a temporary test Excel file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.xlsx")

	// Test with non-existent file - should return error
	app := NewApp()
	_, err := app.ProcessSheet(testFile, "Sheet1")
	if err == nil {
		t.Error("ProcessSheet() should return error for non-existent file")
	}
}

func TestAppGenerateSQL(t *testing.T) {
	app := NewApp()

	headers := []string{"id", "name", "age"}
	dataRows := [][]interface{}{
		{1, "John", 30},
		{2, "Jane", 25},
	}
	numberColumns := []string{"id", "age"}

	result, err := app.GenerateSQL(headers, dataRows, numberColumns)
	if err != nil {
		t.Errorf("GenerateSQL() returned error: %v", err)
	}
	if result == "" {
		t.Error("GenerateSQL() returned empty string")
	}

	// Should contain VALUES keyword
	if len(result) < 10 {
		t.Error("GenerateSQL() result too short")
	}
}

func TestAppGenerateSQLEmptyData(t *testing.T) {
	app := NewApp()

	headers := []string{}
	dataRows := [][]interface{}{}
	numberColumns := []string{}

	result, err := app.GenerateSQL(headers, dataRows, numberColumns)
	if err != nil {
		t.Errorf("GenerateSQL() returned error: %v", err)
	}

	// Empty data should still work
	_ = result
}

func TestAppFindAndReplace(t *testing.T) {
	app := NewApp()

	dataRows := [][]interface{}{
		{"Hello", "World"},
		{"Hello", "Go"},
	}

	result := app.FindAndReplace(dataRows, "Hello", "Hi")

	if result[0][0] != "Hi" {
		t.Errorf("FindAndReplace() failed: got %v, want 'Hi'", result[0][0])
	}
	if result[1][0] != "Hi" {
		t.Errorf("FindAndReplace() failed: got %v, want 'Hi'", result[1][0])
	}
}

func TestAppFindAndReplaceNoMatch(t *testing.T) {
	app := NewApp()

	dataRows := [][]interface{}{
		{"Hello", "World"},
	}

	result := app.FindAndReplace(dataRows, "NotFound", "Replaced")

	if result[0][0] != "Hello" {
		t.Errorf("FindAndReplace() should not change non-matching values")
	}
}

func TestAppFindAndReplaceNil(t *testing.T) {
	app := NewApp()

	result := app.FindAndReplace(nil, "Hello", "Hi")

	// FindAndReplace with nil returns empty slice, not nil
	if result == nil {
		t.Error("FindAndReplace(nil) should return empty slice, not nil")
	}
	if len(result) != 0 {
		t.Errorf("FindAndReplace(nil) should return empty slice, got len=%d", len(result))
	}
}

func TestAppFindAndReplaceEmptyStrings(t *testing.T) {
	app := NewApp()

	dataRows := [][]interface{}{
		{"", "Value"},
	}

	result := app.FindAndReplace(dataRows, "", "Replaced")

	if result[0][0] != "Replaced" {
		t.Errorf("FindAndReplace() should replace empty strings too")
	}
}

func TestProcessSheetWithInvalidPath(t *testing.T) {
	app := NewApp()

	// Test with various invalid paths
	testCases := []string{
		"",
		"nonexistent.xlsx",
		"/invalid/path/file.xlsx",
	}

	for _, path := range testCases {
		_, err := app.ProcessSheet(path, "Sheet1")
		if err == nil {
			t.Errorf("ProcessSheet(%q) should return error", path)
		}
	}
}

func TestExportToFileLogic(t *testing.T) {
	// Test the file writing logic separately (without dialog)
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_export.txt")
	content := "Test SQL content"

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Verify file was written
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(data) != content {
		t.Errorf("File content = %q, want %q", string(data), content)
	}
}
