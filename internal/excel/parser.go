package excel

import (
	"github.com/xuri/excelize/v2"
)

// SheetData represents data from a processed sheet
type SheetData struct {
	Headers  []string
	DataRows [][]interface{}
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

// ProcessSheet reads data from a specific sheet
func ProcessSheet(filePath, sheetName string) (*SheetData, error) {
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
	var dataRows [][]interface{}

	isFirstRow := true
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		if isFirstRow {
			headers = row
			isFirstRow = false
			continue
		}

		// Create interface slice for the row, ensuring it matches header length
		interfaceRow := make([]interface{}, len(headers))
		for j := 0; j < len(headers); j++ {
			if j < len(row) && row[j] != "" {
				interfaceRow[j] = row[j]
			} else {
				interfaceRow[j] = nil
			}
		}
		dataRows = append(dataRows, interfaceRow)
	}

	return &SheetData{
		Headers:  headers,
		DataRows: dataRows,
	}, nil
}
