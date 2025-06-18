package compare

import (
	"context"
	"fmt"
	"time"

	"log"

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

	log.Printf("Connected to database: %s", db.Options().Database)

	ctxQuery, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	if columnName == "" {
		return nil, fmt.Errorf("columnName cannot be empty")
	}

	// Map allowed column names to ensure correct matching
	allowedColumns := map[string]string{
		"goods_en": "goods_en",
		"hs_code":  "hs_code",
	}
	mappedColumn, ok := allowedColumns[columnName]
	if !ok {
		return nil, fmt.Errorf("invalid column name: %s", columnName)
	}

	query := fmt.Sprintf("SELECT %s FROM public.tbl_compare_goods WHERE %s IS NOT NULL AND %s != ''", pg.Ident(mappedColumn), pg.Ident(mappedColumn), pg.Ident(mappedColumn))
	log.Printf("Executing query: %s", query)

	var dbValues []string
	_, err := db.WithContext(ctxQuery).Query(&dbValues, query)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, fmt.Errorf("failed to query database for column %s: %w", columnName, err)
	}

	log.Printf("Retrieved values: %+v", dbValues)
	return dbValues, nil
}
