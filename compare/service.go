package compare

import (
	"context"
	"fmt"
)

type ExcelServiceInterface interface {
	CompareExcelWithDB(ctx context.Context, excelData map[string]struct{}, columnName string) (*CompareResponse, error)
}

type excelService struct {
	repo ExcelRepositoryInterface
}

func NewExcelService(repo ExcelRepositoryInterface) ExcelServiceInterface {
	return &excelService{repo: repo}
}

func (s *excelService) CompareExcelWithDB(ctx context.Context, excelValues map[string]struct{}, columnName string) (*CompareResponse, error) {
	// ดึงข้อมูลจากฐานข้อมูล
	dbValuesSlice, err := s.repo.GetValuesFromDB(ctx, columnName)
	if err != nil {
		return nil, fmt.Errorf("failed to get values from DB: %w", err)
	}

	// สร้าง map สำหรับข้อมูลใน DB
	dbValuesMap := make(map[string]struct{})
	for _, val := range dbValuesSlice {
		dbValuesMap[val] = struct{}{}
	}

	// เปรียบเทียบข้อมูล
	matchedRows := 0
	excelItems := []ExcelItem{}

	for excelVal := range excelValues {
		isMatch := false
		if _, exists := dbValuesMap[excelVal]; exists {
			matchedRows++
			isMatch = true
		}
		excelItems = append(excelItems, ExcelItem{
			Value:   excelVal,
			IsMatch: isMatch,
		})
	}

	// สร้าง response
	response := &CompareResponse{
		TotalExcelRows: len(excelValues),
		TotalDBRows:    len(dbValuesMap),
		MatchedRows:    matchedRows,
		MismatchedRows: len(excelValues) - matchedRows,
		ExcelItems:     excelItems,
	}

	return response, nil
}
