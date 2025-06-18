package utils

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ReadExcelColumn อ่านค่าจาก Column ที่ระบุในไฟล์ Excel
// และส่งคืนเป็น map[string]struct{} เพื่อความรวดเร็วในการค้นหาและไม่เก็บค่าซ้ำ
func ReadExcelColumn(fileBytes []byte, columnName string) (map[string]struct{}, error) {
	f, err := excelize.OpenReader(strings.NewReader(string(fileBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing Excel file: %v\n", err)
		}
	}()

	sheetName := f.GetSheetName(0) // ได้ชื่อ Sheet แรก
	if sheetName == "" {
		return nil, fmt.Errorf("no sheet found in Excel file")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet: %w", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("excel file is empty or has no data rows")
	}

	// ค้นหา Index ของ Column ที่ Frontend ส่งมา
	columnIndex := -1
	headerRow := rows[0] // แถวแรกคือ Header
	for i, cellValue := range headerRow {
		if strings.EqualFold(strings.TrimSpace(cellValue), strings.TrimSpace(columnName)) { // เปรียบเทียบแบบไม่สนใจ Case และตัดช่องว่าง
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		return nil, fmt.Errorf("column '%s' not found in Excel file", columnName)
	}

	excelValues := make(map[string]struct{})
	for i, row := range rows {
		if i == 0 { // ข้าม Header Row
			continue
		}
		if columnIndex < len(row) {
			val := strings.TrimSpace(row[columnIndex])
			if val != "" { // ไม่เก็บค่าว่าง
				excelValues[val] = struct{}{}
			}
		}
	}

	return excelValues, nil
}
