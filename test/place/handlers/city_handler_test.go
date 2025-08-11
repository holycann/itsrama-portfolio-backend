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

// Mock city service
type mockCityService struct {
	mock.Mock
}

func (m *mockCityService) CreateCity(ctx context.Context, city *models.City) (*models.City, error) {
	args := m.Called(ctx, city)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.City), args.Error(1)
}

func (m *mockCityService) GetCityByID(ctx context.Context, id string) (*models.CityDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CityDTO), args.Error(1)
}

func (m *mockCityService) GetCityByName(ctx context.Context, name string) (*models.CityDTO, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CityDTO), args.Error(1)
}

func (m *mockCityService) GetCityByCode(ctx context.Context, code string) (*models.CityDTO, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CityDTO), args.Error(1)
}

func (m *mockCityService) GetCitiesByProvince(ctx context.Context, provinceID string) ([]models.CityDTO, error) {
	args := m.Called(ctx, provinceID)
	return args.Get(0).([]models.CityDTO), args.Error(1)
}

func (m *mockCityService) ListCities(ctx context.Context, opts base.ListOptions) ([]models.CityDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.CityDTO), args.Int(1), args.Error(2)
}

func (m *mockCityService) ListCitiesByPopulation(ctx context.Context, minPopulation, maxPopulation string) ([]models.CityDTO, error) {
	args := m.Called(ctx, minPopulation, maxPopulation)
	return args.Get(0).([]models.CityDTO), args.Error(1)
}

func (m *mockCityService) CountCities(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *mockCityService) SearchCities(ctx context.Context, opts base.ListOptions) ([]models.CityDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.CityDTO), args.Int(1), args.Error(2)
}

func (m *mockCityService) UpdateCity(ctx context.Context, city *models.City) (*models.City, error) {
	args := m.Called(ctx, city)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.City), args.Error(1)
}

func (m *mockCityService) DeleteCity(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper function to setup test environment
func setupCityHandlerTest() (*gin.Engine, *mockCityService, *handlers.CityHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(mockCityService)
	appLogger := logger.DefaultLogger()
	cityHandler := handlers.NewCityHandler(mockService, appLogger)

	return router, mockService, cityHandler
}

func TestCreateCity(t *testing.T) {
	router, mockService, cityHandler := setupCityHandlerTest()

	// Test case: successful city creation
	t.Run("Successful city creation", func(t *testing.T) {
		router.POST("/cities", cityHandler.CreateCity)

		cityID := uuid.New()
		provinceID := uuid.New()
		now := time.Now()

		newCity := &models.City{
			ID:          cityID,
			Name:        "Test City",
			Description: "Test city description",
			ProvinceID:  provinceID,
			ImageURL:    "https://example.com/city.jpg",
			CreatedAt:   &now,
			UpdatedAt:   &now,
		}

		cityCreate := &models.CityCreate{
			Name:        "Test City",
			Description: "Test city description",
			ProvinceID:  provinceID,
			ImageURL:    "https://example.com/city.jpg",
		}

		mockService.On("CreateCity", mock.Anything, mock.AnythingOfType("*models.City")).Return(newCity, nil)

		jsonData, _ := json.Marshal(cityCreate)
		req, _ := http.NewRequest(http.MethodPost, "/cities", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: validation failure
	t.Run("Validation failure", func(t *testing.T) {
		router.POST("/cities", cityHandler.CreateCity)

		// Missing required fields
		cityCreate := &models.CityCreate{
			// Name and ProvinceID are required but missing
		}

		jsonData, _ := json.Marshal(cityCreate)
		req, _ := http.NewRequest(http.MethodPost, "/cities", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestGetCityByID(t *testing.T) {
	router, mockService, cityHandler := setupCityHandlerTest()

	// Test case: successful city retrieval
	t.Run("Successful city retrieval", func(t *testing.T) {
		router.GET("/cities/:id", cityHandler.GetCityByID)

		cityID := uuid.New()
		provinceID := uuid.New()

		city := &models.CityDTO{
			ID:          cityID,
			Name:        "Test City",
			Description: "Test city description",
			ProvinceID:  provinceID,
			ImageURL:    "https://example.com/city.jpg",
		}

		mockService.On("GetCityByID", mock.Anything, cityID.String()).Return(city, nil)

		req, _ := http.NewRequest(http.MethodGet, "/cities/"+cityID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: city not found
	t.Run("City not found", func(t *testing.T) {
		router.GET("/cities/:id", cityHandler.GetCityByID)

		cityID := uuid.New()

		mockService.On("GetCityByID", mock.Anything, cityID.String()).Return(nil, errors.New(errors.ErrNotFound, "City not found", nil))

		req, _ := http.NewRequest(http.MethodGet, "/cities/"+cityID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListCities(t *testing.T) {
	router, mockService, cityHandler := setupCityHandlerTest()

	// Test case: successful cities listing
	t.Run("Successful cities listing", func(t *testing.T) {
		router.GET("/cities", cityHandler.ListCities)

		cityID1 := uuid.New()
		cityID2 := uuid.New()
		provinceID := uuid.New()

		cities := []models.CityDTO{
			{
				ID:          cityID1,
				Name:        "Test City 1",
				Description: "Test city 1 description",
				ProvinceID:  provinceID,
				ImageURL:    "https://example.com/city1.jpg",
			},
			{
				ID:          cityID2,
				Name:        "Test City 2",
				Description: "Test city 2 description",
				ProvinceID:  provinceID,
				ImageURL:    "https://example.com/city2.jpg",
			},
		}

		mockService.On("ListCities", mock.Anything, mock.AnythingOfType("base.ListOptions")).Return(cities, len(cities), nil)
		mockService.On("CountCities", mock.Anything, mock.AnythingOfType("[]base.FilterOption")).Return(len(cities), nil)

		req, _ := http.NewRequest(http.MethodGet, "/cities", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateCity(t *testing.T) {
	router, mockService, cityHandler := setupCityHandlerTest()

	// Test case: successful city update
	t.Run("Successful city update", func(t *testing.T) {
		router.PUT("/cities/:id", cityHandler.UpdateCity)

		cityID := uuid.New()
		provinceID := uuid.New()
		now := time.Now()

		updatedCity := &models.City{
			ID:          cityID,
			Name:        "Updated City",
			Description: "Updated city description",
			ProvinceID:  provinceID,
			ImageURL:    "https://example.com/updated.jpg",
			UpdatedAt:   &now,
		}

		cityUpdate := &models.CityUpdate{
			ID:          cityID,
			Name:        "Updated City",
			Description: "Updated city description",
			ImageURL:    "https://example.com/updated.jpg",
		}

		mockService.On("UpdateCity", mock.Anything, mock.AnythingOfType("*models.City")).Return(updatedCity, nil)

		jsonData, _ := json.Marshal(cityUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/cities/"+cityID.String(), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteCity(t *testing.T) {
	router, mockService, cityHandler := setupCityHandlerTest()

	// Test case: successful city deletion
	t.Run("Successful city deletion", func(t *testing.T) {
		router.DELETE("/cities/:id", cityHandler.DeleteCity)

		cityID := uuid.New()

		mockService.On("DeleteCity", mock.Anything, cityID.String()).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/cities/"+cityID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}
