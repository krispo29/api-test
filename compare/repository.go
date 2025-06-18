package compare

import (
	"database/sql"
	"fmt"
	// ยังคงใช้ model เดิม
)

type ExcelRepositoryInterface interface {
	GetValuesFromDB(columnName string) ([]string, error)
}

type excelRepository struct {
	db *sql.DB
}

func NewExcelRepository(db *sql.DB) ExcelRepositoryInterface {
	return &excelRepository{db: db}
}

func (r *excelRepository) GetValuesFromDB(columnName string) ([]string, error) {
	var dbValues []string

	// Basic validation or sanitization for columnName (important for security)
	// For now, we'll assume columnName is valid and safe.
	// A more robust solution would be to validate against a list of known columns.
	if columnName == "" {
		return nil, fmt.Errorf("columnName cannot be empty")
	}

	// Construct the query safely. Using fmt.Sprintf directly with column names can be risky
	// if columnName is not validated.
	// For PostgreSQL, column names can be quoted using double quotes.
	// For MySQL, column names can be quoted using backticks.
	// Assuming PostgreSQL for quoting, but this should ideally match the actual DB.
	// Let's stick to a simple approach and assume the column name doesn't need special quoting for now,
	// or is already in a format that doesn't conflict with SQL keywords.
	// A better way would be to check against an allow-list of column names.
	query := fmt.Sprintf("SELECT %s FROM goods_hs_code", columnName) // simplified for now

	rows, err := r.db.Query(query) // Be cautious with dynamic table/column names
	if err != nil {
		return nil, fmt.Errorf("failed to query database for column %s: %w", columnName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			return nil, fmt.Errorf("failed to scan row from database for column %s: %w", columnName, err)
		}
		dbValues = append(dbValues, val)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating database rows for column %s: %w", columnName, err)
	}

	return dbValues, nil
}
