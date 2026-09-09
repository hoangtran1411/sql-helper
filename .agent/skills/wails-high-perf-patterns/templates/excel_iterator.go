package patterns

import (
	"bufio"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

// IterateSheet streams through an Excel sheet row-by-row with O(1) memory usage.
// onRow callback is invoked for every data row (skipping header).
func IterateSheet(filePath, sheetName string, onRow func(row []string) error) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.Rows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get sheet rows: %w", err)
	}
	defer rows.Close()

	isFirstRow := true
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to read row: %w", err)
		}

		if isFirstRow {
			isFirstRow = false
			continue // Skip header row
		}

		if err := onRow(row); err != nil {
			return err
		}
	}
	return nil
}

// StreamToFile streams processed Excel data directly to an output file using buffered I/O.
func StreamToFile(filePath, sheetName, outputPath string, transform func(row []string) (string, error)) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)

	err = IterateSheet(filePath, sheetName, func(row []string) error {
		line, err := transform(row)
		if err != nil {
			return err
		}
		if _, err := writer.WriteString(line); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	return nil
}
