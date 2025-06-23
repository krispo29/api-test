package utils

import (
	"bytes"
	"fmt"
	"hpc-express-service/compare"

	"github.com/xuri/excelize/v2"
)

func ReadExcelColumn(fileBytes []byte, columnName string) (map[string]compare.ExcelValue, error) {
	values := make(map[string]compare.ExcelValue)
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
	hsCodeIndex := -1
	for i, header := range rows[0] {
		if header == columnName {
			columnIndex = i
		}
		if header == "hs_code" {
			hsCodeIndex = i
		}
	}
	if columnIndex == -1 {
		return nil, fmt.Errorf("column '%s' not found in Excel file", columnName)
	}
	if (columnName == "goods_en" || columnName == "goods_th") && hsCodeIndex == -1 {
		return nil, fmt.Errorf("column 'hs_code' is required when comparing '%s'", columnName)
	}

	for _, row := range rows[1:] { // ข้าม header
		if columnIndex < len(row) && row[columnIndex] != "" {
			hsCode := ""
			if hsCodeIndex != -1 && hsCodeIndex < len(row) {
				hsCode = row[hsCodeIndex]
			}
			values[row[columnIndex]] = compare.ExcelValue{
				Value:  row[columnIndex],
				HSCode: hsCode,
			}
		}
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("no valid values found in column '%s'", columnName)
	}

	return values, nil
}
