package server

import (
	"encoding/json"
	"fmt"
	"hpc-express-service/compare"
	"hpc-express-service/utils"
	"io"
	"net/http"
)

type ExcelHandlerInterface interface {
	CompareExcel(w http.ResponseWriter, r *http.Request)
}

type excelHandler struct {
	service compare.ExcelServiceInterface
}

func NewExcelHandler(svc compare.ExcelServiceInterface) ExcelHandlerInterface {
	return &excelHandler{service: svc}
}

func (h *excelHandler) CompareExcel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing multipart form: %v", err), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("excelFile")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving the file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	columnName := r.FormValue("columnName")
	if columnName == "" {
		http.Error(w, "Column name is required", http.StatusBadRequest)
		return
	}

	// จำกัดคอลัมน์ที่อนุญาต
	allowedColumns := map[string]bool{
		"goods_en":  true,
		"goods_th":  true,
		"hs_code":   true,
		"tariff":    true,
		"unit_code": true,
		"duty_rate": true,
	}
	if !allowedColumns[columnName] {
		http.Error(w, fmt.Sprintf("Column '%s' is not allowed for comparison", columnName), http.StatusBadRequest)
		return
	}

	excelFileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading file: %v", err), http.StatusInternalServerError)
		return
	}

	excelValues, err := utils.ReadExcelColumn(excelFileBytes, columnName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing Excel file for column '%s': %v", columnName, err), http.StatusBadRequest)
		return
	}
	if len(excelValues) == 0 {
		http.Error(w, fmt.Sprintf("No data found in the specified Excel column '%s'", columnName), http.StatusBadRequest)
		return
	}

	response, err := h.service.CompareExcelWithDB(r.Context(), excelValues, columnName)
	if err != nil {
		fmt.Printf("Service error during comparison for column '%s': %v\n", columnName, err)
		http.Error(w, fmt.Sprintf("Error during comparison for column '%s': %v", columnName, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
