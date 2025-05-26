package airline

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAirlineRepository is a mock type for the airline.Repository interface
type MockAirlineRepository struct {
	mock.Mock
}

// GetAllAirlines provides a mock function for Repository.GetAllAirlines
func (m *MockAirlineRepository) GetAllAirlines(ctx context.Context) ([]*Airline, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Airline), args.Error(1)
}

func TestService_GetAllAirlines(t *testing.T) {
	defaultTimeout := 5 * time.Second

	t.Run("successful retrieval", func(t *testing.T) {
		mockRepo := new(MockAirlineRepository)
		expectedAirlines := []*Airline{
			{UUID: "uuid1", Name: "Airline One", LogoURL: "http://logo.url/1"},
			{UUID: "uuid2", Name: "Airline Two", LogoURL: "http://logo.url/2"},
		}

		// Setup expectations
		mockRepo.On("GetAllAirlines", mock.AnythingOfType("*context.timerCtx")).Return(expectedAirlines, nil).Once()

		service := NewService(mockRepo, defaultTimeout)
		ctx := context.Background() // Use a background context for the test call

		airlines, err := service.GetAllAirlines(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, airlines)
		assert.Equal(t, expectedAirlines, airlines)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockAirlineRepository)
		expectedError := errors.New("repository error")

		// Setup expectations
		mockRepo.On("GetAllAirlines", mock.AnythingOfType("*context.timerCtx")).Return(nil, expectedError).Once()

		service := NewService(mockRepo, defaultTimeout)
		ctx := context.Background() // Use a background context for the test call

		airlines, err := service.GetAllAirlines(ctx)

		assert.Error(t, err)
		assert.Nil(t, airlines)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("context timeout", func(t *testing.T) {
		mockRepo := new(MockAirlineRepository)
		// Simulate a scenario where the repository call hangs until context timeout
		// The service has its own timeout, which should be shorter or equal to this for the test to be meaningful.
		veryShortTimeout := 1 * time.Millisecond
		service := NewService(mockRepo, veryShortTimeout) // Service with a very short timeout

		// Mock the repository to return an error that indicates a delay,
		// or simply don't expect it to return before the service's context times out.
		// Here, we expect GetAllAirlines to be called, but the service's context should cancel it.
		mockRepo.On("GetAllAirlines", mock.AnythingOfType("*context.timerCtx")).
			Return(nil, context.DeadlineExceeded). // Simulate repo itself respecting timeout
			Maybe() // Maybe, because the service's own timeout might fire first

		ctx := context.Background()
		_, err := service.GetAllAirlines(ctx)

		assert.Error(t, err)
		// Check if the error is context.DeadlineExceeded or contains it.
		// The error might be wrapped by the service or pg driver, so direct equality might fail.
		// For this test, we expect the service's own context timeout.
		assert.EqualError(t, err, context.DeadlineExceeded.Error(), "Error should be context.DeadlineExceeded")
		// mockRepo.AssertExpectations(t) // May or may not be called depending on timing
	})
}
