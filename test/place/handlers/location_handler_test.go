package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/holycann/cultour-backend/internal/place/handlers"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
)

// Mock location service
type mockLocationService struct {
	mock.Mock
}

func (m *mockLocationService) CreateLocation(ctx context.Context, location *models.LocationCreate) (*models.Location, error) {
	args := m.Called(ctx, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *mockLocationService) GetLocationByID(ctx context.Context, id string) (*models.LocationDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LocationDTO), args.Error(1)
}

func (m *mockLocationService) GetLocationByName(ctx context.Context, name string) (*models.LocationDTO, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LocationDTO), args.Error(1)
}

func (m *mockLocationService) GetLocationsByCity(ctx context.Context, cityID string) ([]models.LocationDTO, error) {
	args := m.Called(ctx, cityID)
	return args.Get(0).([]models.LocationDTO), args.Error(1)
}

func (m *mockLocationService) ListLocations(ctx context.Context, opts base.ListOptions) ([]models.LocationDTO, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.LocationDTO), args.Error(1)
}

func (m *mockLocationService) GetLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.LocationDTO, error) {
	args := m.Called(ctx, latitude, longitude, radius)
	return args.Get(0).([]models.LocationDTO), args.Error(1)
}

func (m *mockLocationService) CountLocations(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *mockLocationService) SearchLocations(ctx context.Context, opts base.ListOptions) ([]models.LocationDTO, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.LocationDTO), args.Error(1)
}

func (m *mockLocationService) UpdateLocation(ctx context.Context, location *models.LocationUpdate) (*models.Location, error) {
	args := m.Called(ctx, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *mockLocationService) DeleteLocation(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to setup test environment
func setupLocationHandlerTest() (*gin.Engine, *mockLocationService, *handlers.LocationHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(mockLocationService)
	appLogger := logger.DefaultLogger()
	locationHandler := handlers.NewLocationHandler(mockService, appLogger)

	return router, mockService, locationHandler
}

func TestCreateLocation(t *testing.T) {
	router, mockService, locationHandler := setupLocationHandlerTest()

	// Test case: successful location creation
	t.Run("Successful location creation", func(t *testing.T) {
		router.POST("/locations", locationHandler.CreateLocation)

		locationID := uuid.New()
		cityID := uuid.New()
		now := time.Now()

		newLocation := &models.Location{
			ID:        locationID,
			Name:      "Test Location",
			CityID:    cityID,
			Latitude:  -6.175392,
			Longitude: 106.827153,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		locationCreate := models.LocationCreate{
			Name:      "Test Location",
			CityID:    cityID,
			Latitude:  -6.175392,
			Longitude: 106.827153,
		}

		mockService.On("CreateLocation", mock.Anything, mock.AnythingOfType("*models.LocationCreate")).Return(newLocation, nil)

		jsonData, _ := json.Marshal(locationCreate)
		req, _ := http.NewRequest(http.MethodPost, "/locations", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: validation failure
	t.Run("Validation failure", func(t *testing.T) {
		router.POST("/locations", locationHandler.CreateLocation)

		// Missing required fields (Name, CityID, Latitude, Longitude)
		locationCreate := models.LocationCreate{}

		jsonData, _ := json.Marshal(locationCreate)
		req, _ := http.NewRequest(http.MethodPost, "/locations", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestGetLocationByID(t *testing.T) {
	router, mockService, locationHandler := setupLocationHandlerTest()

	// Test case: successful location retrieval
	t.Run("Successful location retrieval", func(t *testing.T) {
		router.GET("/locations/:id", locationHandler.GetLocationByID)

		locationID := uuid.New()
		cityID := uuid.New()

		location := &models.LocationDTO{
			ID:        locationID,
			Name:      "Test Location",
			CityID:    cityID,
			Latitude:  -6.175392,
			Longitude: 106.827153,
		}

		mockService.On("GetLocationByID", mock.Anything, locationID.String()).Return(location, nil)

		req, _ := http.NewRequest(http.MethodGet, "/locations/"+locationID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: location not found
	t.Run("Location not found", func(t *testing.T) {
		router.GET("/locations/:id", locationHandler.GetLocationByID)

		locationID := uuid.New()

		mockService.On("GetLocationByID", mock.Anything, locationID.String()).Return(nil, errors.New(errors.ErrNotFound, "Location not found", nil))

		req, _ := http.NewRequest(http.MethodGet, "/locations/"+locationID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListLocations(t *testing.T) {
	router, mockService, locationHandler := setupLocationHandlerTest()

	// Test case: successful locations listing
	t.Run("Successful locations listing", func(t *testing.T) {
		router.GET("/locations", locationHandler.ListLocations)

		locationID1 := uuid.New()
		locationID2 := uuid.New()
		cityID := uuid.New()

		locations := []models.LocationDTO{
			{
				ID:        locationID1,
				Name:      "Test Location 1",
				CityID:    cityID,
				Latitude:  -6.175392,
				Longitude: 106.827153,
			},
			{
				ID:        locationID2,
				Name:      "Test Location 2",
				CityID:    cityID,
				Latitude:  -6.175393,
				Longitude: 106.827154,
			},
		}

		mockService.On("ListLocations", mock.Anything, mock.AnythingOfType("base.ListOptions")).Return(locations, nil)

		req, _ := http.NewRequest(http.MethodGet, "/locations", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetLocationsByProximity(t *testing.T) {
	router, mockService, locationHandler := setupLocationHandlerTest()

	// Test case: successful proximity search
	t.Run("Successful proximity search", func(t *testing.T) {
		// Define the route with the parameters in the query string
		router.GET("/locations/proximity", locationHandler.SearchLocations)

		locationID1 := uuid.New()
		locationID2 := uuid.New()
		cityID := uuid.New()

		// Define locations that are near the target coordinates
		locations := []models.LocationDTO{
			{
				ID:        locationID1,
				Name:      "Nearby Location 1",
				CityID:    cityID,
				Latitude:  -6.175392,
				Longitude: 106.827153,
			},
			{
				ID:        locationID2,
				Name:      "Nearby Location 2",
				CityID:    cityID,
				Latitude:  -6.175393,
				Longitude: 106.827154,
			},
		}

		// Set up the mock to return these locations when proximity search is called
		mockService.On("SearchLocations", mock.Anything, mock.AnythingOfType("base.ListOptions")).Return(locations, nil)

		// Create a request with latitude, longitude, and radius in the query parameters
		req, _ := http.NewRequest(http.MethodGet, "/locations/proximity?latitude=-6.175390&longitude=106.827150&radius=1000", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateLocation(t *testing.T) {
	router, mockService, locationHandler := setupLocationHandlerTest()

	// Test case: successful location update
	t.Run("Successful location update", func(t *testing.T) {
		router.PUT("/locations/:id", locationHandler.UpdateLocation)

		locationID := uuid.New()
		cityID := uuid.New()
		now := time.Now()

		updatedLocation := &models.Location{
			ID:        locationID,
			Name:      "Updated Location",
			CityID:    cityID,
			Latitude:  -6.175392,
			Longitude: 106.827153,
			UpdatedAt: &now,
		}

		locationUpdate := models.LocationUpdate{
			ID:   locationID,
			Name: "Updated Location",
		}

		mockService.On("UpdateLocation", mock.Anything, mock.AnythingOfType("*models.LocationUpdate")).Return(updatedLocation, nil)

		jsonData, _ := json.Marshal(locationUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/locations/"+locationID.String(), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteLocation(t *testing.T) {
	router, mockService, locationHandler := setupLocationHandlerTest()

	// Test case: successful location deletion
	t.Run("Successful location deletion", func(t *testing.T) {
		router.DELETE("/locations/:id", locationHandler.DeleteLocation)

		locationID := uuid.New()

		mockService.On("DeleteLocation", mock.Anything, locationID.String()).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/locations/"+locationID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}
