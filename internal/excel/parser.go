package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// SheetData represents data from a processed sheet
type SheetData struct {
	Headers  []string `json:"headers"`
	DataRows [][]any  `json:"dataRows"`
}

// Replacement defines a find-and-replace rule
type Replacement struct {
	Find    string `json:"find"`
	Replace string `json:"replace"`
}

// ParseExcelFile opens an Excel file and returns the list of sheet names
func ParseExcelFile(filePath string) ([]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return f.GetSheetList(), nil
}

// GetPreview reads only the first N rows for UI display
// limit <= 0 means read all (not recommended for large files)
func GetPreview(filePath, sheetName string, limit int) (*SheetData, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.Rows(sheetName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var headers []string
	var dataRows [][]any

	isFirstRow := true
	rowCount := 0

	for rows.Next() {
		// Stop if we reached the limit
		if limit > 0 && rowCount >= limit {
			break
		}

		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		if isFirstRow {
			headers = row
			isFirstRow = false
			continue
		}

		interfaceRow := NormalizeRow(row, len(headers), nil)
		dataRows = append(dataRows, interfaceRow)
		rowCount++
	}

	return &SheetData{
		Headers:  headers,
		DataRows: dataRows,
	}, nil
}

// IterateSheet streams through the excel file and executes a callback for each row
// This allows processing massive files with O(1) memory
func IterateSheet(filePath, sheetName string, onRow func(row []string) error) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	rows, err := f.Rows(sheetName)
	if err != nil {
		return err
	}
	defer rows.Close()

	isFirstRow := true
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return err
		}

		if isFirstRow {
			isFirstRow = false
			continue // Skip processing header in callback, consumer can handle header separately if needed logic
		}

		if err := onRow(row); err != nil {
			return err
		}
	}
	return nil
}

// NormalizeRow normalizes a raw string slice into an interface row of length numHeaders,
// applying any replacements and converting empty strings to nil.
func NormalizeRow(rawRow []string, numHeaders int, replacements []Replacement) []any {
	row := make([]any, numHeaders)
	for j := range numHeaders {
		if j < len(rawRow) {
			val := rawRow[j]
			for _, r := range replacements {
				if val == r.Find {
					val = r.Replace
				}
			}
			if val != "" {
				row[j] = val
				continue
			}
		}
		row[j] = nil
	}
	return row
}

// FindAndReplace replaces matching string cell values in 2D interface rows.
func FindAndReplace(dataRows [][]any, findValue, replaceWith string) [][]any {
	if dataRows == nil {
		return make([][]any, 0)
	}

	for _, row := range dataRows {
		for j, cell := range row {
			if fmt.Sprintf("%v", cell) == findValue {
				row[j] = replaceWith
			}
		}
	}

	return dataRows
}

// StreamRows streams normalized rows from an Excel sheet to a callback.
func StreamRows(filePath, sheetName string, headers []string, replacements []Replacement, onRow func(row []any) error) error {
	return IterateSheet(filePath, sheetName, func(rawRow []string) error {
		normRow := NormalizeRow(rawRow, len(headers), replacements)
		return onRow(normRow)
	})
}
