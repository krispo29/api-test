package compare

import (
	"context"
	"fmt"
	"time" // Added for context timeout consistency

	"github.com/go-pg/pg/v9" // Added for go-pg
)

type ExcelRepositoryInterface interface {
	GetValuesFromDB(ctx context.Context, columnName string) ([]string, error)
}

type excelRepository struct {
	contextTimeout time.Duration // Added for consistency
}

func NewExcelRepository(timeout time.Duration) ExcelRepositoryInterface {
	return &excelRepository{
		contextTimeout: timeout,
	}
}

func (r *excelRepository) GetValuesFromDB(ctx context.Context, columnName string) ([]string, error) {
	db := ctx.Value("postgreSQLConn").(*pg.DB)
	// It's good practice to ensure db is not nil, though in a well-structured app it should always be there.
	if db == nil {
		return nil, fmt.Errorf("database connection not found in context")
	}

	ctxQuery, cancel := context.WithTimeout(ctx, r.contextTimeout)
	defer cancel()

	if columnName == "" {
		return nil, fmt.Errorf("columnName cannot be empty")
	}

	// IMPORTANT: columnName is used directly in the query.
	// This is a SQL injection risk if columnName is not validated/sanitized.
	// For now, proceeding as per previous plan's assumption.
	// Consider using `?` placeholder for column name if go-pg supports it,
	// or `pg.Ident` for safe identifier quoting.
	// query := fmt.Sprintf("SELECT %s FROM goods_hs_code", columnName) // Original approach
	// Safer with pg.Ident if columnName is simple:
	query := fmt.Sprintf("SELECT %s FROM goods_hs_code", pg.Ident(columnName))


	var dbValues []string
	// Using db.ModelContext for selecting into a slice of strings from a single column
	// might be tricky. A raw query is often simpler for this specific case.
	// Let's try with a simple query execution that scans results.
	// The most direct way with go-pg to get a slice of single values from a query
	// is often to query into a slice of structs and then extract, or use a simpler model.
	// However, for a single column of strings:
	_, err := db.QueryContext(ctxQuery, pg.Scan(&dbValues), query)

	if err != nil {
		// Check for pg.ErrNoRows specifically if needed, though QueryContext might not return it for scans into slices.
		return nil, fmt.Errorf("failed to query database for column %s: %w", columnName, err)
	}

	return dbValues, nil
}
