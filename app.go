package main

import (
	"context"
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
	data, err := excel.ProcessSheet(filePath, sheetName)
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
