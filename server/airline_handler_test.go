package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	// "time" // Not directly used in this version, but often useful

	"hpc-express-service/airline" // Assuming airline types are in this package
	"hpc-express-service/factory" // For ServiceFactory if needed, or mock directly
	"hpc-express-service/constant" // For response codes like CodeSuccess

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAirlineService is a mock type for the airline.Service interface
type MockAirlineService struct {
	mock.Mock
}

// GetAllAirlines provides a mock function for Service.GetAllAirlines
func (m *MockAirlineService) GetAllAirlines(ctx context.Context) ([]*airline.Airline, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		if args.Error(1) != nil {
			return nil, args.Error(1)
		}
		return nil, nil // Should not happen if mock is set up correctly
	}
	return args.Get(0).([]*airline.Airline), args.Error(1)
}

func TestAirlineHandler_GetAllAirlines(t *testing.T) {
	// Helper function to set up the router with the mock service
	setupRouter := func(airlineSvc airline.Service) *chi.Mux {
		r := chi.NewRouter()

		// We are unit testing the handler, so we inject the mock service directly.
		// No need for the full service factory here.
		airlineCtrl := airlineHandler{svc: airlineSvc}

		// Create a /v1 sub-router and mount the airline handler as it's done in server.go
		// This ensures the path prefix is correctly handled in tests.
		v1Router := chi.NewRouter()
		v1Router.Mount("/airlines", airlineCtrl.router()) // airlineCtrl.router() returns the sub-router for airline
		r.Mount("/v1", v1Router)
		
		return r
	}

	t.Run("successful retrieval", func(t *testing.T) {
		mockService := new(MockAirlineService)
		expectedAirlines := []*airline.Airline{
			{UUID: "uuid1", Name: "Airline One", LogoURL: "http://logo.url/1"},
		}
		// This is what the JSON response for a single airline object should look like
		expectedAirlineJSON := map[string]interface{}{
			"uuid":    "uuid1",
			"name":    "Airline One",
			"logoUrl": "http://logo.url/1", // JSON tag in airline.Airline struct is `json:"logoUrl"`
		}

		mockService.On("GetAllAirlines", mock.Anything).Return(expectedAirlines, nil).Once()

		router := setupRouter(mockService)
		req, _ := http.NewRequest("GET", "/v1/airlines", nil) // Request to the endpoint
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		// Assertions based on the structure of SuccessResponse from server.go
		assert.Equal(t, float64(constant.CodeSuccess), responseBody["code"], "Response code should match CodeSuccess")
		assert.Equal(t, "success", responseBody["message"], "Response message should be 'success'")
		
		responseData, ok := responseBody["data"].([]interface{})
		assert.True(t, ok, "Response data should be a slice")
		assert.Len(t, responseData, 1, "Response data should contain one airline")

		airlineData, ok := responseData[0].(map[string]interface{})
		assert.True(t, ok, "Airline data should be a map")
		
		// Compare individual fields of the airline
		assert.Equal(t, expectedAirlineJSON["uuid"], airlineData["uuid"])
		assert.Equal(t, expectedAirlineJSON["name"], airlineData["name"])
		assert.Equal(t, expectedAirlineJSON["logoUrl"], airlineData["logoUrl"])

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(MockAirlineService)
		serviceErr := errors.New("service layer error") // The actual error message

		mockService.On("GetAllAirlines", mock.Anything).Return(nil, serviceErr).Once()

		router := setupRouter(mockService)
		req, _ := http.NewRequest("GET", "/v1/airlines", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)
		
		assert.Equal(t, http.StatusBadRequest, rr.Code, "HTTP status code should be Bad Request on service error") 

		// Assuming ErrInvalidRequest is used, which populates ErrResponse
		var errorResponse ErrResponse 
		err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
		assert.NoError(t, err, "Failed to unmarshal error response")

		// Check fields of ErrResponse based on server.go's ErrInvalidRequest
		assert.Equal(t, int64(constant.CodeError), errorResponse.AppCode, "AppCode should match constant.CodeError")
		assert.Equal(t, serviceErr.Error(), errorResponse.Message, "Error message should match the service error")

		mockService.AssertExpectations(t)
	})

	// Test for empty list of airlines
	t.Run("successful retrieval with empty list", func(t *testing.T) {
		mockService := new(MockAirlineService)
		expectedAirlines := []*airline.Airline{} // Empty slice

		mockService.On("GetAllAirlines", mock.Anything).Return(expectedAirlines, nil).Once()

		router := setupRouter(mockService)
		req, _ := http.NewRequest("GET", "/v1/airlines", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		assert.Equal(t, float64(constant.CodeSuccess), responseBody["code"])
		assert.Equal(t, "success", responseBody["message"])
		
		responseData, ok := responseBody["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, responseData, 0, "Data should be an empty slice for no airlines")

		mockService.AssertExpectations(t)
	})
}

// Note: The test assumes that `server.ErrResponse`, `server.SuccessResponse`, 
// `server.ErrInvalidRequest`, and `constant.CodeSuccess`/`constant.CodeError` are accessible and behave as expected.
// If `ErrResponse` is not exported (e.g. starts with lowercase `errResponse`), 
// a local struct mirroring its JSON structure would be needed for unmarshalling and assertion.
// The provided code already had `ErrResponse` as exported.
// The factory import was removed from setupRouter as it wasn't strictly necessary for this unit test.
// The simplified `setupRouter` directly uses the `airlineHandler` with the mock service.
```
