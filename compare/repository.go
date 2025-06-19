package compare

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-pg/pg/v9"
)

type ExcelRepositoryInterface interface {
	GetValuesFromDB(ctx context.Context, columnName string) ([]string, error)
}

type excelRepository struct {
	contextTimeout time.Duration
}

func NewExcelRepository(timeout time.Duration) ExcelRepositoryInterface {
	return &excelRepository{
		contextTimeout: timeout,
	}
}

func (r *excelRepository) GetValuesFromDB(ctx context.Context, columnName string) ([]string, error) {
	db := ctx.Value("postgreSQLConn").(*pg.DB)
	if db == nil {
		return nil, fmt.Errorf("database connection not found in context")
	}

	ctxQuery, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	if columnName == "" {
		return nil, fmt.Errorf("columnName cannot be empty")
	}

	// ตรวจสอบว่าคอลัมน์มีอยู่ในตาราง
	var exists bool
	_, err := db.QueryOneContext(ctxQuery, &exists, `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_schema = 'public' 
			AND table_name = 'tbl_compare_goods' 
			AND column_name = ?
		)`, columnName)
	if err != nil {
		log.Printf("Failed to check column existence: %v", err)
		return nil, fmt.Errorf("failed to check column existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("column '%s' does not exist in table tbl_compare_goods", columnName)
	}

	query := fmt.Sprintf("SELECT %s FROM public.tbl_compare_goods WHERE %s IS NOT NULL AND %s != ''", pg.Ident(columnName), pg.Ident(columnName), pg.Ident(columnName))

	var dbValues []string
	_, err = db.WithContext(ctxQuery).Query(&dbValues, query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, fmt.Errorf("failed to query database for column %s: %w", columnName, err)
	}
	return dbValues, nil
}
