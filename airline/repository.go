package airline

import (
	"context"
	"time"

	"github.com/go-pg/pg/v9"
)

// Repository defines the interface for airline data operations.
type Repository interface {
	GetAllAirlines(ctx context.Context) ([]*Airline, error)
}

type repository struct {
	contextTimeout time.Duration
}

// NewRepository creates a new airline repository.
func NewRepository(timeout time.Duration) Repository {
	return &repository{
		contextTimeout: timeout,
	}
}

// GetAllAirlines retrieves all airlines from the database.
func (r *repository) GetAllAirlines(ctx context.Context) ([]*Airline, error) {
	db := ctx.Value("postgreSQLConn").(*pg.DB)
	// Create a new context with a timeout for the database query.
	// Note: Using context.Background() here creates a new root context,
	// detaching it from the incoming request's context.
	// Consider whether this is the desired behavior or if the incoming context should be preferred.
	queryCtx, cancel := context.WithTimeout(context.Background(), r.contextTimeout)
	defer cancel()

	sqlStr := `
		SELECT 
			"uuid",
			"name",
			"logo_url"  -- Assuming the column name is logo_url
		FROM public.tbl_airlines
		WHERE deleted_at IS NULL
		ORDER BY name ASC; -- Optional: Order by name
	`

	var airlines []*Airline
	_, err := db.QueryContext(queryCtx, &airlines, sqlStr)

	if err != nil {
		// Consider logging the error here or returning a more specific error type
		return nil, err
	}

	return airlines, nil
}
