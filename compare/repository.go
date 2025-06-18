package compare

import (
	"database/sql"
	"fmt"
	// ยังคงใช้ model เดิม
)

type ExcelRepositoryInterface interface {
	GetAllValuesFromDB() ([]string, error)
}

type excelRepository struct {
	db *sql.DB
}

func NewExcelRepository(db *sql.DB) ExcelRepositoryInterface {
	return &excelRepository{db: db}
}

func (r *excelRepository) GetAllValuesFromDB() ([]string, error) {
	var dbValues []string

	rows, err := r.db.Query("SELECT goods_en FROM goods_hs_code")
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			return nil, fmt.Errorf("failed to scan row from database: %w", err)
		}
		dbValues = append(dbValues, val)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating database rows: %w", err)
	}

	return dbValues, nil
}
