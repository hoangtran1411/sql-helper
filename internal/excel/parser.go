package excel

import (
	"github.com/xuri/excelize/v2"
)

// SheetData represents data from a processed sheet
type SheetData struct {
	Headers  []string `json:"headers"`
	DataRows [][]any  `json:"dataRows"`
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

		// Create interface slice for the row
		interfaceRow := make([]any, len(headers))
		for j := range headers {
			if j < len(row) && row[j] != "" {
				interfaceRow[j] = row[j]
			} else {
				interfaceRow[j] = nil
			}
		}
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
