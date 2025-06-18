package compare

type DatabaseItem struct {
	Value string `json:"value"`
}
type CompareRequest struct {
	ColumnName string `json:"columnName"`
}

type CompareResponse struct {
	TotalExcelRows int      `json:"totalExcelRows"`
	TotalDBRows    int      `json:"totalDBRows"`
	MatchedRows    int      `json:"matchedRows"`
	MismatchedRows int      `json:"mismatchedRows"`
	MissingInExcel []string `json:"missingInExcel"` // ข้อมูลที่มีใน DB แต่ไม่มีใน Excel
	MissingInDB    []string `json:"missingInDB"`    // ข้อมูลที่มีใน Excel แต่ไม่มีใน DB
}
