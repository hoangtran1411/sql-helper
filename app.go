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

// SQLOptions defines configuration for generating SQL statements
type SQLOptions struct {
	TableName       string   `json:"tableName"`
	SelectedColumns []string `json:"selectedColumns"`
	NumberColumns   []string `json:"numberColumns"`
	BatchSize       int      `json:"batchSize"`
	ValuesOnly      bool     `json:"valuesOnly"`
}

// GenerateSQL generates SQL INSERT statements (or raw values) from the data
func (a *App) GenerateSQL(headers []string, dataRows [][]interface{}, options SQLOptions) (string, error) {
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

	// Determine selected columns and map indices
	selectedCols := options.SelectedColumns
	if len(selectedCols) == 0 {
		selectedCols = headers
	}
	if len(selectedCols) == 0 {
		return false, fmt.Errorf("no columns selected")
	}

	headerIndexMap := make(map[string]int, len(headers))
	for i, h := range headers {
		if _, exists := headerIndexMap[h]; !exists {
			headerIndexMap[h] = i
		}
	}

	colIndices := make([]int, 0, len(selectedCols))
	validSelectedCols := make([]string, 0, len(selectedCols))
	for _, col := range selectedCols {
		if idx, ok := headerIndexMap[col]; ok {
			colIndices = append(colIndices, idx)
			validSelectedCols = append(validSelectedCols, col)
		}
	}
	if len(validSelectedCols) == 0 {
		return false, fmt.Errorf("no valid columns found to export")
	}

	// Prepare column set for fast lookup
	numColSet := make(map[string]bool, len(options.NumberColumns))
	for _, col := range options.NumberColumns {
		numColSet[col] = true
	}

	insertPrefix := sql.BuildInsertPrefix(options.TableName, validSelectedCols)
	batchSize := options.BatchSize
	if batchSize < 0 {
		batchSize = 0
	}

	inBatchCount := 0
	totalRows := 0

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

		if options.ValuesOnly {
			if totalRows > 0 {
				if _, err := writer.WriteString(",\n"); err != nil {
					return err
				}
			}
			valStr := sql.FormatRowSQLSelected(interfaceRow, colIndices, headers, numColSet)
			if _, err := writer.WriteString(valStr); err != nil {
				return err
			}
		} else {
			if inBatchCount == 0 {
				if totalRows > 0 {
					if _, err := writer.WriteString("\n\n"); err != nil {
						return err
					}
				}
				if _, err := writer.WriteString(insertPrefix); err != nil {
					return err
				}
			} else {
				if _, err := writer.WriteString(",\n"); err != nil {
					return err
				}
			}

			valStr := sql.FormatRowSQLSelected(interfaceRow, colIndices, headers, numColSet)
			if _, err := writer.WriteString(valStr); err != nil {
				return err
			}
			inBatchCount++

			if batchSize > 0 && inBatchCount == batchSize {
				if _, err := writer.WriteString(";"); err != nil {
					return err
				}
				inBatchCount = 0
			}
		}

		totalRows++
		return nil
	})

	if err != nil {
		return false, fmt.Errorf("streaming failed: %w", err)
	}

	if !options.ValuesOnly && inBatchCount > 0 {
		if _, err := writer.WriteString(";"); err != nil {
			return false, err
		}
	}

	if err := writer.Flush(); err != nil {
		return false, fmt.Errorf("failed to flush buffer: %w", err)
	}

	return true, nil
}
