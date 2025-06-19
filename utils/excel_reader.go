package utils

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ReadExcelColumn(fileBytes []byte, columnName string) (map[string]struct{}, error) {
	values := make(map[string]struct{})
	file, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer file.Close()

	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read rows: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows found in sheet")
	}

	// หา index ของคอลัมน์
	columnIndex := -1
	for i, header := range rows[0] {
		if header == columnName {
			columnIndex = i
			break
		}
	}
	if columnIndex == -1 {
		return nil, fmt.Errorf("column '%s' not found in Excel file", columnName)
	}

	// อ่านข้อมูลจากคอลัมน์
	for _, row := range rows[1:] { // ข้าม header
		if columnIndex < len(row) && row[columnIndex] != "" {
			values[row[columnIndex]] = struct{}{}
		}
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("no valid values found in column '%s'", columnName)
	}

	return values, nil
}
