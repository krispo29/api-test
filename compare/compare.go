package compare

// DataModel แทนโครงสร้างข้อมูลในฐานข้อมูล
// สมมติว่าใน DB มีตารางที่เก็บค่าที่เราจะเปรียบเทียบ
type DatabaseItem struct {
	Value string `json:"value"` // Column ที่จะใช้เปรียบเทียบ เช่น "product_code", "id_card"
}

// CompareRequest สำหรับรับข้อมูลจาก Frontend
type CompareRequest struct {
	ColumnName string `json:"columnName"` // ชื่อ Column ใน Excel ที่ Frontend เลือกมา
}

// CompareResponse สำหรับส่งผลลัพธ์กลับไปยัง Frontend
type CompareResponse struct {
	TotalExcelRows int      `json:"totalExcelRows"`
	TotalDBRows    int      `json:"totalDBRows"`
	MatchedRows    int      `json:"matchedRows"`
	MismatchedRows int      `json:"mismatchedRows"`
	MissingInExcel []string `json:"missingInExcel"` // ข้อมูลที่มีใน DB แต่ไม่มีใน Excel
	MissingInDB    []string `json:"missingInDB"`    // ข้อมูลที่มีใน Excel แต่ไม่มีใน DB
}
