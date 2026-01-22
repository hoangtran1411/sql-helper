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

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &SheetData{
			Headers:  []string{},
			DataRows: [][]interface{}{},
		}, nil
	}

	// First row is headers
	headers := rows[0]

	// Convert remaining rows to []interface{}
	dataRows := make([][]interface{}, 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
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
