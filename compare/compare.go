package compare

type ExcelItem struct {
	Value   string `json:"value"`
	IsMatch bool   `json:"isMatch"`
}

type CompareResponse struct {
	TotalExcelRows int         `json:"totalExcelRows"`
	TotalDBRows    int         `json:"totalDBRows"`
	MatchedRows    int         `json:"matchedRows"`
	ExcelItems     []ExcelItem `json:"excelItems"`
}
