package compare

import (
	"fmt"

	"hpc-express-service/compare"
)

// ExcelServiceInterface กำหนด Contract สำหรับ Service
type ExcelServiceInterface interface {
	CompareExcelWithDB(excelData map[string]struct{}, columnName string) (*compare.CompareResponse, error)
}

// excelService implements ExcelServiceInterface
type excelService struct {
	repo compare.ExcelRepositoryInterface
}

// NewExcelService สร้าง instance ของ ExcelService
func NewExcelService(repo compare.ExcelRepositoryInterface) ExcelServiceInterface {
	return &excelService{repo: repo}
}

// CompareExcelWithDB ดำเนินการเปรียบเทียบข้อมูล Excel กับข้อมูลในฐานข้อมูล
func (s *excelService) CompareExcelWithDB(excelValues map[string]struct{}, columnName string) (*compare.CompareResponse, error) {
	// 1. ดึงข้อมูลจากฐานข้อมูลผ่าน Repository
	dbValuesSlice, err := s.repo.GetValuesFromDB(columnName) // Pass columnName here
	if err != nil {
		return nil, fmt.Errorf("service: failed to get data from database for column %s: %w", columnName, err)
	}

	// แปลง slice เป็น map เพื่อการค้นหาที่รวดเร็ว
	dbValuesMap := make(map[string]struct{})
	for _, val := range dbValuesSlice {
		dbValuesMap[val] = struct{}{}
	}

	// 2. เปรียบเทียบข้อมูล (Business Logic) - This part remains largely the same
	matchedRows := 0
	// mismatchedRows := 0 // This was calculated incorrectly before, should be len(excelMissingInDB)
	var excelMissingInDB []string // ข้อมูลที่มีใน Excel แต่ไม่มีใน DB
	var dbMissingInExcel []string // ข้อมูลที่มีใน DB แต่ไม่มีใน Excel

	// ตรวจสอบข้อมูลใน Excel เทียบกับ DB
	for excelVal := range excelValues {
		if _, exists := dbValuesMap[excelVal]; exists {
			matchedRows++
		} else {
			excelMissingInDB = append(excelMissingInDB, excelVal)
		}
	}

	// ตรวจสอบข้อมูลใน DB เทียบกับ Excel (เพื่อหา MissingInExcel)
	for _, dbVal := range dbValuesSlice {
		if _, exists := excelValues[dbVal]; !exists {
			dbMissingInExcel = append(dbMissingInExcel, dbVal)
		}
	}

	// 3. สร้างผลลัพธ์
	response := &compare.CompareResponse{
		TotalExcelRows: len(excelValues),
		TotalDBRows:    len(dbValuesMap),
		MatchedRows:    matchedRows,
		MismatchedRows: len(excelMissingInDB), // Corrected calculation
		MissingInExcel: dbMissingInExcel,
		MissingInDB:    excelMissingInDB,
	}

	return response, nil
}
