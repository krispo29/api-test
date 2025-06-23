package compare

import (
	"context"
	"fmt"
	"log"
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
	dbValuesMap := make(map[string]DBDetails)
	hsCodeMap := make(map[string][]DBDetails)
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
			for _, existingRow := range hsCodeMap[row.HSCode] {
				if existingRow.GoodsEN != row.GoodsEN || existingRow.GoodsTH != row.GoodsTH {
					log.Printf("Warning: Potential inconsistent data for HSCode %s. Record 1: GoodsEN='%s', GoodsTH='%s'. Record 2: GoodsEN='%s', GoodsTH='%s'", row.HSCode, existingRow.GoodsEN, existingRow.GoodsTH, row.GoodsEN, row.GoodsTH)
				}
			}
			hsCodeMap[row.HSCode] = append(hsCodeMap[row.HSCode], row)
		}
	}

	matchedRows := 0
	excelItems := make([]ExcelItem, 0, len(excelValues))

	for _, excelVal := range excelValues {
		item := ExcelItem{Value: excelVal.Value}
		if dbRow, exists := dbValuesMap[excelVal.Value]; exists {
			item.IsMatch = true
			item.MatchedBy = "column"
			item.DBDetails = &dbRow
			matchedRows++
		} else if (columnName == "goods_en" || columnName == "goods_th") && excelVal.HSCode != "" {
			if dbRecords, ok := hsCodeMap[excelVal.HSCode]; ok && len(dbRecords) > 0 {
				matchType := "hs_code_fallback"
				for _, record := range dbRecords {
					if (columnName == "goods_en" && record.GoodsEN == excelVal.Value) ||
						(columnName == "goods_th" && record.GoodsTH == excelVal.Value) {
						item.DBDetails = &record
						item.IsMatch = true
						if columnName == "goods_en" {
							matchType = "hs_code_specific_en"
						} else {
							matchType = "hs_code_specific_th"
						}
						matchedRows++
						break
					}
				}
				if !item.IsMatch {
					item.DBDetails = &dbRecords[0]
					item.IsMatch = true
					matchedRows++
				}
				item.MatchedBy = matchType
			}
		}
		excelItems = append(excelItems, item)
	}

	return &CompareResponse{
		TotalExcelRows: len(excelValues),
		TotalDBRows:    len(dbValuesMap),
		MatchedRows:    matchedRows,
		ExcelItems:     excelItems,
	}, nil
}
