package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/hoangtran1411/sql-helper/internal/excel"
	"github.com/hoangtran1411/sql-helper/internal/sql"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ExcelResult represents the result of opening an Excel file
type ExcelResult struct {
	FilePath   string   `json:"filePath"`
	SheetNames []string `json:"sheetNames"`
}

// SheetData represents data from a processed sheet
type SheetData struct {
	Headers  []string        `json:"headers"`
	DataRows [][]interface{} `json:"dataRows"`
}

// OpenExcelFile opens a file dialog and parses the selected Excel file
func (a *App) OpenExcelFile() (*ExcelResult, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Excel File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Excel Files (*.xlsx, *.xls)",
				Pattern:     "*.xlsx;*.xls",
			},
		},
	})
	if err != nil {
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
func (a *App) ProcessSheet(filePath, sheetName string) (*SheetData, error) {
	// Limit preview to 100 rows for performance
	data, err := excel.GetPreview(filePath, sheetName, 100)
	if err != nil {
		return nil, err
	}

	return &SheetData{
		Headers:  data.Headers,
		DataRows: data.DataRows,
	}, nil
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
	return runtime.ClipboardSetText(a.ctx, text)
}

// ExportToFile saves content to a file
func (a *App) ExportToFile(content string) error {
	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "result.txt",
		Filters: []runtime.FileFilter{
			{DisplayName: "Text Files", Pattern: "*.txt"},
		},
	})
	if err != nil {
		return err
	}

	if savePath == "" {
		return nil // User cancelled
	}

	return os.WriteFile(savePath, []byte(content), 0644)
}

// Replacement defines a find-and-replace rule
type Replacement struct {
	Find    string `json:"find"`
	Replace string `json:"replace"`
}

// GenerateAndSaveSQL streams data from Excel directly to a SQL file
// This is memory efficient (O(1)) and can handle massive files
func (a *App) GenerateAndSaveSQL(filePath, sheetName string, headers []string, numberColumns []string, replacements []Replacement) error {
	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: "output.sql",
		Title:           "Save SQL File",
		Filters: []runtime.FileFilter{
			{DisplayName: "SQL Files", Pattern: "*.sql"},
		},
	})
	if err != nil {
		return err
	}
	if savePath == "" {
		return nil // User cancelled
	}

	// Create output file with buffered writer for performance
	outFile, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
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
		// This matches the behavior of FindAndReplace on the interface{} slice
		// since we are dealing with strings from Excel
		for i := range row {
			for _, r := range replacements {
				if row[i] == r.Find {
					row[i] = r.Replace
				}
			}
		}

		// Convert string row to interface row for the formatter
		// Note: We use the headers from the frontend configuration to ensure alignment
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
		return fmt.Errorf("streaming failed: %w", err)
	}

	return nil
}
