package server

import (
	"encoding/json"
	"fmt"
	"hpc-express-service/tools/compare"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *excelHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CompareExcel)

	return r
}

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

	// Excel processing is now handled by the service
	response, err := h.service.CompareExcelWithDB(r.Context(), excelFileBytes, columnName)
	if err != nil {
		// Check if the error is something the user should see as BadRequest vs InternalServerError
		// For now, let's assume service errors are generally internal or propagated input errors become bad requests.
		// The service function readExcelColumnFromBytes already returns specific errors that can be
		// interpreted as bad requests (e.g., column not found, no data in column).
		// We might need more sophisticated error handling here to distinguish.
		// For simplicity, if the service returns an error, we'll treat it as a potential internal server error
		// or a bad request if it's due to bad input that the service detected.
		// The service should ideally wrap errors or return specific error types.
		// For now, let's keep it as InternalServerError and log the actual error.
		fmt.Printf("Service error during comparison for column '%s': %v\n", columnName, err)           // Log for debugging
		http.Error(w, fmt.Sprintf("Error during comparison: %v", err), http.StatusInternalServerError) // User-facing error
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
