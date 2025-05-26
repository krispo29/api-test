package airline

import (
	"context"
	"time"
)

// Service defines the interface for airline business logic.
type Service interface {
	GetAllAirlines(ctx context.Context) ([]*Airline, error)
}

type service struct {
	repo           Repository // Changed from selfRepo to repo for clarity
	contextTimeout time.Duration
}

// NewService creates a new airline service.
func NewService(repo Repository, timeout time.Duration) Service {
	return &service{
		repo:           repo,
		contextTimeout: timeout,
	}
}

// GetAllAirlines retrieves all airlines.
func (s *service) GetAllAirlines(ctx context.Context) ([]*Airline, error) {
	// Create a new context with timeout, derived from the incoming context.
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	airlines, err := s.repo.GetAllAirlines(ctx)
	if err != nil {
		// Consider logging the error here or returning a more specific error type
		return nil, err
	}

	return airlines, nil
}
