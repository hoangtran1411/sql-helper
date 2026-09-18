package main

import (
	"bufio"
	"fmt"
	"io"
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

// SQLOptions defines configuration for generating SQL statements
type SQLOptions struct {
	TableName       string   `json:"tableName"`
	SelectedColumns []string `json:"selectedColumns"`
	NumberColumns   []string `json:"numberColumns"`
	BatchSize       int      `json:"batchSize"`
	ValuesOnly      bool     `json:"valuesOnly"`
}

// GenerateSQL generates SQL INSERT statements (or raw values) from the data
func (a *App) GenerateSQL(headers []string, dataRows [][]any, options SQLOptions) (string, error) {
	result := sql.GenerateBatchSQL(headers, dataRows, sql.GenerateOptions{
		TableName:       options.TableName,
		SelectedColumns: options.SelectedColumns,
		NumberColumns:   options.NumberColumns,
		BatchSize:       options.BatchSize,
		ValuesOnly:      options.ValuesOnly,
	})
	return result, nil
}

// FindAndReplace replaces values in the data rows
func (a *App) FindAndReplace(dataRows [][]any, findValue, replaceWith string) [][]any {
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

// ExportSQLStream streams data from an Excel sheet to an io.Writer using sql.BatchWriter.
// This decouples SQL streaming from UI dialogs, enabling O(1) memory exports and direct testing.
func ExportSQLStream(w io.Writer, filePath, sheetName string, headers []string, options SQLOptions, replacements []Replacement) error {
	bw, err := sql.NewBatchWriter(w, headers, sql.GenerateOptions{
		TableName:       options.TableName,
		SelectedColumns: options.SelectedColumns,
		NumberColumns:   options.NumberColumns,
		BatchSize:       options.BatchSize,
		ValuesOnly:      options.ValuesOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize batch writer: %w", err)
	}

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
		interfaceRow := make([]any, len(headers))
		for i := range headers {
			if i < len(row) && row[i] != "" {
				interfaceRow[i] = row[i]
			} else {
				interfaceRow[i] = nil
			}
		}

		return bw.WriteRow(interfaceRow)
	})
	if err != nil {
		return fmt.Errorf("streaming failed: %w", err)
	}

	if err := bw.Close(); err != nil {
		return fmt.Errorf("failed to finalize SQL output: %w", err)
	}

	return nil
}

// GenerateAndSaveSQL streams data from Excel directly to a SQL file
// This is memory efficient (O(1)) and can handle massive files
func (a *App) GenerateAndSaveSQL(filePath, sheetName string, headers []string, options SQLOptions, replacements []Replacement) (bool, error) {
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
	defer func() {
		_ = outFile.Close()
	}()
 
	writer := bufio.NewWriter(outFile)
	if err := ExportSQLStream(writer, filePath, sheetName, headers, options, replacements); err != nil {
		return false, err
	}
	return true, nil
}
