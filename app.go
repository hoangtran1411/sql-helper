package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hoangtran1411/sql-helper/internal/excel"
	"github.com/hoangtran1411/sql-helper/internal/sql"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// App struct
type App struct{}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// ExcelResult represents the result of opening an Excel file
type ExcelResult struct {
	FilePath   string   `json:"filePath"`
	SheetNames []string `json:"sheetNames"`
}

// OpenExcelFile opens a file dialog and parses the selected Excel file
func (a *App) OpenExcelFile() (*ExcelResult, error) {
	dialog := application.Get().Dialog.OpenFileWithOptions(&application.OpenFileDialogOptions{
		Title: "Select Excel File",
		Filters: []application.FileFilter{
			{
				DisplayName: "Excel Files (*.xlsx, *.xls)",
				Pattern:     "*.xlsx;*.xls",
			},
		},
	})

	filePath, err := dialog.PromptForSingleSelection()
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "cancel") {
			return nil, nil // User cancelled
		}
		return nil, err
	}

	if filePath == "" {
		return nil, nil // User cancelled
	}

	sheetNames, err := excel.ParseExcelFile(filePath)
	if err != nil {
		return nil, err
	}

	return &ExcelResult{
		FilePath:   filePath,
		SheetNames: sheetNames,
	}, nil
}

// ProcessSheet reads data from a specific sheet
func (a *App) ProcessSheet(filePath, sheetName string) (*excel.SheetData, error) {
	// Limit preview to 100 rows for performance
	return excel.GetPreview(filePath, sheetName, 100)
}

// GenerateSQL generates SQL INSERT values from the data
func (a *App) GenerateSQL(headers []string, dataRows [][]interface{}, numberColumns []string) (string, error) {
	result := sql.GenerateSQLValues(headers, dataRows, numberColumns)
	return result, nil
}

// FindAndReplace replaces values in the data rows
func (a *App) FindAndReplace(dataRows [][]interface{}, findValue, replaceWith string) [][]interface{} {
	return sql.FindAndReplace(dataRows, findValue, replaceWith)
}

// CopyToClipboard copies text to the clipboard
func (a *App) CopyToClipboard(text string) error {
	if !application.Get().Clipboard.SetText(text) {
		return fmt.Errorf("failed to copy text to clipboard")
	}
	return nil
}

// Replacement defines a find-and-replace rule
type Replacement struct {
	Find    string `json:"find"`
	Replace string `json:"replace"`
}

// GenerateAndSaveSQL streams data from Excel directly to a SQL file
// This is memory efficient (O(1)) and can handle massive files
func (a *App) GenerateAndSaveSQL(filePath, sheetName string, headers []string, numberColumns []string, replacements []Replacement) (bool, error) {
	dialog := application.Get().Dialog.SaveFileWithOptions(&application.SaveFileDialogOptions{
		Title:    "Save SQL File",
		Filename: "output.sql",
		Filters: []application.FileFilter{
			{
				DisplayName: "SQL Files (*.sql)",
				Pattern:     "*.sql",
			},
		},
	})

	savePath, err := dialog.PromptForSingleSelection()
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "cancel") {
			return false, nil // User cancelled
		}
		return false, err
	}
	if savePath == "" {
		return false, nil // User cancelled
	}

	// Create output file with buffered writer for performance
	outFile, err := os.Create(savePath)
	if err != nil {
		return false, fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	// Prepare column set for fast lookup
	numColSet := make(map[string]bool)
	for _, col := range numberColumns {
		numColSet[col] = true
	}

	firstRow := true

	// Use the streaming iterator
	err = excel.IterateSheet(filePath, sheetName, func(row []string) error {
		// Apply replacements to raw strings first
		for i := range row {
			for _, r := range replacements {
				if row[i] == r.Find {
					row[i] = r.Replace
				}
			}
		}

		// Convert string row to interface row for the formatter
		interfaceRow := make([]interface{}, len(headers))
		for i := 0; i < len(headers); i++ {
			if i < len(row) && row[i] != "" {
				interfaceRow[i] = row[i]
			} else {
				interfaceRow[i] = nil
			}
		}

		// Add comma separator for subsequent rows
		if !firstRow {
			if _, err := writer.WriteString(",\n"); err != nil {
				return err
			}
		}

		// Format and write the row
		valStr := sql.FormatRowSQL(interfaceRow, headers, numColSet)
		if _, err := writer.WriteString(valStr); err != nil {
			return err
		}

		firstRow = false
		return nil
	})

	if err != nil {
		return false, fmt.Errorf("streaming failed: %w", err)
	}

	return true, nil
}
