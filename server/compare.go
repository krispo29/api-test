package server

import (
	"encoding/json"
	"fmt"
	"hpc-express-service/compare"
	"hpc-express-service/utils"
	"io"
	"net/http"
)

// ExcelHandlerInterface กำหนด Contract สำหรับ Handler
type ExcelHandlerInterface interface {
	CompareExcel(w http.ResponseWriter, r *http.Request)
}

// excelHandler implements ExcelHandlerInterface
type excelHandler struct {
	service compare.ExcelServiceInterface
}

// NewExcelHandler สร้าง instance ของ ExcelHandler
func NewExcelHandler(svc compare.ExcelServiceInterface) ExcelHandlerInterface {
	return &excelHandler{service: svc}
}

// CompareExcel จัดการ HTTP Request สำหรับการเปรียบเทียบไฟล์ Excel
func (h *excelHandler) CompareExcel(w http.ResponseWriter, r *http.Request) {
	// 1. รับไฟล์ Excel และข้อมูลอื่นๆ จาก HTTP Request
	// ตั้งค่าขนาดไฟล์สูงสุดที่อนุญาต (เช่น 10MB)
	r.ParseMultipartForm(10 << 20) // 10 MB

	file, _, err := r.FormFile("excelFile") // "excelFile" คือชื่อ field ใน form data ที่ Frontend ส่งมา
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving the file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	columnName := r.FormValue("columnName") // "columnName" คือชื่อ field ใน form data
	if columnName == "" {
		http.Error(w, "Column name is required", http.StatusBadRequest)
		return
	}

	// 2. อ่านข้อมูลจากไฟล์ Excel (ใช้ Utility function)
	excelFileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	excelValues, err := utils.ReadExcelColumn(excelFileBytes, columnName)
	if err != nil {
		// It might be good to include columnName in this error message too
		http.Error(w, fmt.Sprintf("Error processing Excel file for column '%s': %v", columnName, err), http.StatusBadRequest)
		return
	}
	if len(excelValues) == 0 {
		http.Error(w, fmt.Sprintf("No data found in the specified Excel column '%s'.", columnName), http.StatusBadRequest)
		return
	}

	// 3. เรียกใช้ Service เพื่อประมวลผลการเปรียบเทียบ
	response, err := h.service.CompareExcelWithDB(r.Context(), excelValues, columnName) // Pass columnName here
	if err != nil {
		// Log error server-side
		fmt.Printf("Service error during comparison for column '%s': %v\n", columnName, err)
		http.Error(w, fmt.Sprintf("Error during comparison for column '%s': %v", columnName, err), http.StatusInternalServerError)
		return
	}

	// 4. ส่งผลลัพธ์กลับไปยัง Client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
