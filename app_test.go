package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/hoangtran1411/sql-helper/internal/excel"
	"github.com/xuri/excelize/v2"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Error("NewApp() returned nil")
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
	data := excel.SheetData{
		Headers:  []string{"Name", "Age"},
		DataRows: [][]any{{"John", 30}, {"Jane", 25}},
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
	dataRows := [][]any{
		{1, "John", 30},
		{2, "Jane", 25},
	}
	opts := SQLOptions{
		TableName:       "users",
		SelectedColumns: []string{"id", "name"},
		NumberColumns:   []string{"id"},
		BatchSize:       1000,
	}

	result, err := app.GenerateSQL(headers, dataRows, opts)
	if err != nil {
		t.Errorf("GenerateSQL() returned error: %v", err)
	}
	if result == "" {
		t.Error("GenerateSQL() returned empty string")
	}

	// Should contain INSERT INTO and VALUES keyword
	if !strings.Contains(result, "INSERT INTO users (id, name) VALUES") {
		t.Errorf("GenerateSQL() result does not contain expected prefix: %s", result)
	}
}

func TestAppGenerateSQLEmptyData(t *testing.T) {
	app := NewApp()

	headers := []string{}
	dataRows := [][]any{}
	opts := SQLOptions{}

	result, err := app.GenerateSQL(headers, dataRows, opts)
	if err != nil {
		t.Errorf("GenerateSQL() returned error: %v", err)
	}

	// Empty data should return empty string
	if result != "" {
		t.Errorf("GenerateSQL() with empty data should return empty string, got %q", result)
	}
}

func TestAppFindAndReplace_ExcelDelegation(t *testing.T) {
	app := NewApp()

	input := [][]any{
		{"apple", 100, nil},
		{"banana", "apple", true},
		{"", "cherry", false},
	}

	// Prepare identical input to verify app.FindAndReplace delegates cleanly to excel.FindAndReplace
	expectedInput := [][]any{
		{"apple", 100, nil},
		{"banana", "apple", true},
		{"", "cherry", false},
	}

	expected := excel.FindAndReplace(expectedInput, "apple", "orange")
	result := app.FindAndReplace(input, "apple", "orange")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("app.FindAndReplace delegation mismatch: got %v, want %v", result, expected)
	}

	if result[0][0] != "orange" || result[1][1] != "orange" {
		t.Errorf("expected 'orange' at [0][0] and [1][1], got %v and %v", result[0][0], result[1][1])
	}
	if result[0][1] != 100 || result[0][2] != nil || result[1][0] != "banana" {
		t.Errorf("non-target cells modified unexpectedly: %v", result)
	}
}

func TestAppFindAndReplace(t *testing.T) {
	app := NewApp()

	dataRows := [][]any{
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

	dataRows := [][]any{
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

	dataRows := [][]any{
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

func TestSQLOptionsStruct(t *testing.T) {
	opts := SQLOptions{
		TableName:       "orders",
		SelectedColumns: []string{"id", "total"},
		NumberColumns:   []string{"id", "total"},
		BatchSize:       500,
		ValuesOnly:      false,
	}

	if opts.TableName != "orders" {
		t.Errorf("TableName = %q, want 'orders'", opts.TableName)
	}
	if len(opts.SelectedColumns) != 2 {
		t.Errorf("len(SelectedColumns) = %d, want 2", len(opts.SelectedColumns))
	}
	if opts.BatchSize != 500 {
		t.Errorf("BatchSize = %d, want 500", opts.BatchSize)
	}
	if opts.ValuesOnly {
		t.Errorf("ValuesOnly = true, want false")
	}
}

// createTestExcelFile creates a temporary Excel file for testing in the specified directory.
func createTestExcelFile(t *testing.T, dir string, data [][]string) string {
	t.Helper()

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			t.Logf("Warning: failed to close file: %v", err)
		}
	}()

	for i, row := range data {
		for j, cell := range row {
			cellName, err := excelize.CoordinatesToCellName(j+1, i+1)
			if err != nil {
				t.Errorf("CoordinatesToCellName(%d, %d) failed: %v", j+1, i+1, err)
				return ""
			}
			if err := f.SetCellValue("Sheet1", cellName, cell); err != nil {
				t.Errorf("SetCellValue failed: %v", err)
				return ""
			}
		}
	}

	filePath := filepath.Join(dir, "test.xlsx")
	if err := f.SaveAs(filePath); err != nil {
		t.Fatalf("Failed to save test Excel file: %v", err)
	}
	return filePath
}

func TestExportSQLStream_Basic(t *testing.T) {
	tempDir := t.TempDir()
	data := [][]string{
		{"id", "name", "age", "status"},
		{"1", "Alice", "25", "active"},
		{"2", "Bob", "30", "inactive"},
	}
	excelPath := createTestExcelFile(t, tempDir, data)

	var buf bytes.Buffer
	headers := []string{"id", "name", "age", "status"}
	opts := SQLOptions{
		TableName:     "users",
		BatchSize:     1000,
		NumberColumns: []string{"id", "age"},
	}

	err := ExportSQLStream(&buf, excelPath, "Sheet1", headers, opts, nil)
	if err != nil {
		t.Fatalf("ExportSQLStream() failed: %v", err)
	}

	got := buf.String()
	expected := "INSERT INTO users (id, name, age, status) VALUES\n" +
		"(1, 'Alice', 25, 'active'),\n" +
		"(2, 'Bob', 30, 'inactive');"

	if got != expected {
		t.Errorf("ExportSQLStream() output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
	}

	if !strings.HasPrefix(got, "INSERT INTO users (id, name, age, status) VALUES\n") {
		t.Errorf("expected INSERT INTO prefix, got %q", got)
	}
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (prefix + 2 rows), got %d", len(lines))
	}
	if !strings.HasSuffix(got, ";") {
		t.Errorf("expected output to end with semicolon, got %q", got)
	}
}

func TestExportSQLStream_BatchSplitting(t *testing.T) {
	tempDir := t.TempDir()
	data := [][]string{
		{"id", "name"},
		{"1", "Alice"},
		{"2", "Bob"},
		{"3", "Charlie"},
	}
	excelPath := createTestExcelFile(t, tempDir, data)

	t.Run("split with remainder", func(t *testing.T) {
		var buf bytes.Buffer
		headers := []string{"id", "name"}
		opts := SQLOptions{
			TableName:     "users",
			BatchSize:     2,
			NumberColumns: []string{"id"},
		}

		err := ExportSQLStream(&buf, excelPath, "Sheet1", headers, opts, nil)
		if err != nil {
			t.Fatalf("ExportSQLStream() failed: %v", err)
		}

		got := buf.String()
		expected := "INSERT INTO users (id, name) VALUES\n" +
			"(1, 'Alice'),\n" +
			"(2, 'Bob');\n\n" +
			"INSERT INTO users (id, name) VALUES\n" +
			"(3, 'Charlie');"

		if got != expected {
			t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
		}

		insertCount := strings.Count(got, "INSERT INTO users (id, name) VALUES")
		if insertCount != 2 {
			t.Errorf("expected 2 INSERT INTO statements, got %d", insertCount)
		}
		semicolonCount := strings.Count(got, ";")
		if semicolonCount != 2 {
			t.Errorf("expected 2 semicolons, got %d", semicolonCount)
		}
	})

	t.Run("split exact multiple", func(t *testing.T) {
		var buf bytes.Buffer
		headers := []string{"id", "name"}
		opts := SQLOptions{
			TableName:     "users",
			BatchSize:     1,
			NumberColumns: []string{"id"},
		}

		err := ExportSQLStream(&buf, excelPath, "Sheet1", headers, opts, nil)
		if err != nil {
			t.Fatalf("ExportSQLStream() failed: %v", err)
		}

		got := buf.String()
		expected := "INSERT INTO users (id, name) VALUES\n" +
			"(1, 'Alice');\n\n" +
			"INSERT INTO users (id, name) VALUES\n" +
			"(2, 'Bob');\n\n" +
			"INSERT INTO users (id, name) VALUES\n" +
			"(3, 'Charlie');"

		if got != expected {
			t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
		}
	})
}

func TestExportSQLStream_Replacements(t *testing.T) {
	tempDir := t.TempDir()
	data := [][]string{
		{"id", "name", "status"},
		{"1", "Alice", "Active"},
		{"2", "Bob", "Pending"},
		{"3", "Charlie", "Active"},
	}
	excelPath := createTestExcelFile(t, tempDir, data)

	t.Run("replace string values", func(t *testing.T) {
		var buf bytes.Buffer
		headers := []string{"id", "name", "status"}
		opts := SQLOptions{
			TableName:     "users",
			BatchSize:     1000,
			NumberColumns: []string{"id"},
		}
		replacements := []Replacement{
			{Find: "Active", Replace: "ENABLED"},
			{Find: "Pending", Replace: "PENDING_VERIFY"},
		}

		err := ExportSQLStream(&buf, excelPath, "Sheet1", headers, opts, replacements)
		if err != nil {
			t.Fatalf("ExportSQLStream() failed: %v", err)
		}

		got := buf.String()
		expected := "INSERT INTO users (id, name, status) VALUES\n" +
			"(1, 'Alice', 'ENABLED'),\n" +
			"(2, 'Bob', 'PENDING_VERIFY'),\n" +
			"(3, 'Charlie', 'ENABLED');"

		if got != expected {
			t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
		}
		if strings.Contains(got, "Active") || strings.Contains(got, "Pending") {
			t.Errorf("output contains unreplaced values: %s", got)
		}
	})

	t.Run("replace to numeric value", func(t *testing.T) {
		var buf bytes.Buffer
		headers := []string{"id", "name", "status"}
		opts := SQLOptions{
			TableName:     "users",
			BatchSize:     1000,
			NumberColumns: []string{"id", "status"},
		}
		replacements := []Replacement{
			{Find: "Active", Replace: "1"},
			{Find: "Pending", Replace: "0"},
		}

		err := ExportSQLStream(&buf, excelPath, "Sheet1", headers, opts, replacements)
		if err != nil {
			t.Fatalf("ExportSQLStream() failed: %v", err)
		}

		got := buf.String()
		expected := "INSERT INTO users (id, name, status) VALUES\n" +
			"(1, 'Alice', 1),\n" +
			"(2, 'Bob', 0),\n" +
			"(3, 'Charlie', 1);"

		if got != expected {
			t.Errorf("output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
		}
	})
}

func TestExportSQLStream_ValuesOnly(t *testing.T) {
	tempDir := t.TempDir()
	data := [][]string{
		{"id", "name", "score"},
		{"1", "Alice", "95.5"},
		{"2", "Bob", "88.0"},
	}
	excelPath := createTestExcelFile(t, tempDir, data)

	var buf bytes.Buffer
	headers := []string{"id", "name", "score"}
	opts := SQLOptions{
		NumberColumns: []string{"id", "score"},
		ValuesOnly:    true,
	}

	err := ExportSQLStream(&buf, excelPath, "Sheet1", headers, opts, nil)
	if err != nil {
		t.Fatalf("ExportSQLStream() failed: %v", err)
	}

	got := buf.String()
	expected := "(1, 'Alice', 95.5),\n(2, 'Bob', 88.0)"

	if got != expected {
		t.Errorf("ExportSQLStream() output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
	}

	if strings.Contains(got, "INSERT INTO") {
		t.Errorf("ValuesOnly mode should not contain 'INSERT INTO', got %q", got)
	}
	if strings.HasSuffix(strings.TrimSpace(got), ";") {
		t.Errorf("ValuesOnly mode should not end with semicolon, got %q", got)
	}
}

func TestExportSQLStream_ColumnSelection(t *testing.T) {
	tempDir := t.TempDir()
	data := [][]string{
		{"id", "first_name", "last_name", "salary", "department"},
		{"101", "John", "Doe", "50000", "Engineering"},
		{"102", "Jane", "Smith", "60000", "Design"},
	}
	excelPath := createTestExcelFile(t, tempDir, data)

	var buf bytes.Buffer
	headers := []string{"id", "first_name", "last_name", "salary", "department"}
	opts := SQLOptions{
		TableName:       "employees",
		SelectedColumns: []string{"first_name", "salary"},
		NumberColumns:   []string{"salary"},
		BatchSize:       1000,
	}

	err := ExportSQLStream(&buf, excelPath, "Sheet1", headers, opts, nil)
	if err != nil {
		t.Fatalf("ExportSQLStream() failed: %v", err)
	}

	got := buf.String()
	expected := "INSERT INTO employees (first_name, salary) VALUES\n" +
		"('John', 50000),\n" +
		"('Jane', 60000);"

	if got != expected {
		t.Errorf("ExportSQLStream() output mismatch:\ngot:\n%s\nwant:\n%s", got, expected)
	}

	for _, excluded := range []string{"last_name", "department", "Doe", "Smith", "Engineering", "Design"} {
		if strings.Contains(got, excluded) {
			t.Errorf("output should not contain excluded column/value %q", excluded)
		}
	}
}

func TestExportSQLStream_Errors(t *testing.T) {
	tempDir := t.TempDir()
	data := [][]string{
		{"id", "name"},
		{"1", "Alice"},
	}
	validExcelPath := createTestExcelFile(t, tempDir, data)
	headers := []string{"id", "name"}

	t.Run("non-existent excel file", func(t *testing.T) {
		var buf bytes.Buffer
		opts := SQLOptions{TableName: "users"}
		err := ExportSQLStream(&buf, filepath.Join(tempDir, "nonexistent.xlsx"), "Sheet1", headers, opts, nil)
		if err == nil {
			t.Error("expected error for non-existent Excel file, got nil")
		}
	})

	t.Run("invalid sheet name", func(t *testing.T) {
		var buf bytes.Buffer
		opts := SQLOptions{TableName: "users"}
		err := ExportSQLStream(&buf, validExcelPath, "NonExistentSheet", headers, opts, nil)
		if err == nil {
			t.Error("expected error for non-existent sheet, got nil")
		}
	})

	t.Run("invalid non-matching selected columns", func(t *testing.T) {
		var buf bytes.Buffer
		opts := SQLOptions{
			TableName:       "users",
			SelectedColumns: []string{"unknown_col_1", "unknown_col_2"},
		}
		err := ExportSQLStream(&buf, validExcelPath, "Sheet1", headers, opts, nil)
		if err == nil {
			t.Error("expected error for non-matching selected columns, got nil")
		}
	})

	t.Run("nil writer", func(t *testing.T) {
		opts := SQLOptions{TableName: "users"}
		err := ExportSQLStream(nil, validExcelPath, "Sheet1", headers, opts, nil)
		if err == nil {
			t.Error("expected error for nil writer, got nil")
		}
	})

	t.Run("empty headers", func(t *testing.T) {
		var buf bytes.Buffer
		opts := SQLOptions{TableName: "users"}
		err := ExportSQLStream(&buf, validExcelPath, "Sheet1", nil, opts, nil)
		if err == nil {
			t.Error("expected error for empty headers, got nil")
		}
	})
}

func TestExportSQLStream_ToFile(t *testing.T) {
	tempDir := t.TempDir()
	data := [][]string{
		{"id", "title"},
		{"10", "First Post"},
		{"20", "Second Post"},
	}
	excelPath := createTestExcelFile(t, tempDir, data)

	outPath := filepath.Join(tempDir, "export.sql")
	outFile, err := os.Create(outPath)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}

	headers := []string{"id", "title"}
	opts := SQLOptions{
		TableName:     "posts",
		NumberColumns: []string{"id"},
		BatchSize:     100,
	}

	writer := bufio.NewWriter(outFile)
	if err := ExportSQLStream(writer, excelPath, "Sheet1", headers, opts, nil); err != nil {
		_ = outFile.Close() // best-effort file cleanup on failure
		t.Fatalf("ExportSQLStream() failed: %v", err)
	}
	if err := writer.Flush(); err != nil {
		_ = outFile.Close() // best-effort file cleanup on failure
		t.Fatalf("writer.Flush() failed: %v", err)
	}
	if err := outFile.Close(); err != nil {
		t.Fatalf("Failed to close output file: %v", err)
	}

	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	expected := "INSERT INTO posts (id, title) VALUES\n" +
		"(10, 'First Post'),\n" +
		"(20, 'Second Post');"

	if string(content) != expected {
		t.Errorf("File content mismatch:\ngot:\n%s\nwant:\n%s", string(content), expected)
	}
}
