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
	dbValuesSlice, err := s.repo.GetValuesFromDB(ctx, columnName)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get data from database for column %s: %w", columnName, err)
	}

	dbValuesMap := make(map[string]struct{})
	for _, val := range dbValuesSlice {
		dbValuesMap[val] = struct{}{}
	}

	matchedRows := 0
	var excelMissingInDB []string
	var dbMissingInExcel []string

	for excelVal := range excelValues {
		if _, exists := dbValuesMap[excelVal]; exists {
			matchedRows++
		} else {
			excelMissingInDB = append(excelMissingInDB, excelVal)
		}
	}

	for _, dbVal := range dbValuesSlice {
		if _, exists := excelValues[dbVal]; !exists {
			dbMissingInExcel = append(dbMissingInExcel, dbVal)
		}
	}

	// 3. สร้างผลลัพธ์
	response := &CompareResponse{
		TotalExcelRows: len(excelValues),
		TotalDBRows:    len(dbValuesMap),
		MatchedRows:    matchedRows,
		MismatchedRows: len(excelMissingInDB),
		MissingInExcel: dbMissingInExcel,
		MissingInDB:    excelMissingInDB,
	}

	return response, nil
}
