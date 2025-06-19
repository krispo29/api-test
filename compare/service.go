package compare

import (
	"context"
	"fmt"
)

type ExcelServiceInterface interface {
	CompareExcelWithDB(ctx context.Context, excelData map[string]ExcelValue, columnName string) (*CompareResponse, error)
}

type ExcelValue struct {
	Value  string // ค่าจากคอลัมน์ที่เลือก
	HSCode string // hs_code จาก Excel (ถ้ามี)
}

type excelService struct {
	repo ExcelRepositoryInterface
}

func NewExcelService(repo ExcelRepositoryInterface) ExcelServiceInterface {
	return &excelService{repo: repo}
}

func (s *excelService) CompareExcelWithDB(ctx context.Context, excelValues map[string]ExcelValue, columnName string) (*CompareResponse, error) {
	// ดึงข้อมูลจากฐานข้อมูล
	dbValuesSlice, err := s.repo.GetValuesFromDB(ctx, columnName)
	if err != nil {
		return nil, fmt.Errorf("failed to get values from DB: %w", err)
	}

	// สร้าง map สำหรับเปรียบเทียบตาม columnName และ hs_code
	dbValuesMap := make(map[string]DBDetails) // สำหรับ columnName
	hsCodeMap := make(map[string][]DBDetails) // สำหรับ hs_code
	for _, row := range dbValuesSlice {
		var val string
		switch columnName {
		case "goods_en":
			val = row.GoodsEN
		case "goods_th":
			val = row.GoodsTH
		case "hs_code":
			val = row.HSCode
		case "tariff":
			val = fmt.Sprintf("%d", row.Tariff)
		case "unit_code":
			val = row.UnitCode
		case "duty_rate":
			val = fmt.Sprintf("%f", row.DutyRate)
		}
		if val != "" {
			dbValuesMap[val] = row
		}
		if row.HSCode != "" {
			hsCodeMap[row.HSCode] = append(hsCodeMap[row.HSCode], row)
		}
	}

	// เปรียบเทียบข้อมูล
	matchedRows := 0
	excelItems := []ExcelItem{}

	for _, excelVal := range excelValues {
		item := ExcelItem{
			Value:   excelVal.Value,
			IsMatch: false,
		}

		// ตรวจสอบการจับคู่โดย columnName
		if dbRow, exists := dbValuesMap[excelVal.Value]; exists {
			matchedRows++
			item.IsMatch = true
			item.MatchedBy = "column"
			item.DBDetails = &dbRow
		} else if (columnName == "goods_en" || columnName == "goods_th") && excelVal.HSCode != "" {
			// ตรวจสอบ hs_code สำหรับ goods_en หรือ goods_th
			if rows, exists := hsCodeMap[excelVal.HSCode]; exists {
				matchedRows++
				item.IsMatch = true
				item.MatchedBy = "hs_code"
				item.DBDetails = &rows[0] // เลือกแถวแรกหากมีหลายแถว
			}
		}

		excelItems = append(excelItems, item)
	}

	// สร้าง response
	response := &CompareResponse{
		TotalExcelRows: len(excelValues),
		TotalDBRows:    len(dbValuesMap),
		MatchedRows:    matchedRows,
		ExcelItems:     excelItems,
	}

	return response, nil
}
