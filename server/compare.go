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

	r.ParseMultipartForm(10 << 20)

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
		http.Error(w, fmt.Sprintf("No data found in the specified Excel column '%s'.", columnName), http.StatusBadRequest)
		return
	}

	response, err := h.service.CompareExcelWithDB(r.Context(), excelValues, columnName)
	if err != nil {
		fmt.Printf("Service error during comparison for column '%s': %v\n", columnName, err)
		http.Error(w, fmt.Sprintf("Error during comparison for column '%s': %v", columnName, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
