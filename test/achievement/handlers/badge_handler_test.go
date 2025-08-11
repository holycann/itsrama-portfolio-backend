package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/holycann/cultour-backend/internal/achievement/handlers"
	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
)

// Mock badge service
type mockBadgeService struct {
	mock.Mock
}

func (m *mockBadgeService) CreateBadge(ctx context.Context, badgeCreate *models.BadgeCreate) (*models.BadgeDTO, error) {
	args := m.Called(ctx, badgeCreate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BadgeDTO), args.Error(1)
}

func (m *mockBadgeService) UpdateBadge(ctx context.Context, id string, badgeUpdate *models.BadgeUpdate) (*models.BadgeDTO, error) {
	args := m.Called(ctx, id, badgeUpdate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BadgeDTO), args.Error(1)
}

func (m *mockBadgeService) GetBadgeByID(ctx context.Context, id string) (*models.BadgeDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BadgeDTO), args.Error(1)
}

func (m *mockBadgeService) GetBadgeByName(ctx context.Context, name string) (*models.BadgeDTO, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BadgeDTO), args.Error(1)
}

func (m *mockBadgeService) ListBadges(ctx context.Context, opts base.ListOptions) ([]models.BadgeDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.BadgeDTO), args.Int(1), args.Error(2)
}

func (m *mockBadgeService) SearchBadges(ctx context.Context, opts base.ListOptions) ([]models.BadgeDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.BadgeDTO), args.Int(1), args.Error(2)
}

func (m *mockBadgeService) DeleteBadge(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockBadgeService) CountBadges(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *mockBadgeService) GetPopularBadges(ctx context.Context, limit int) ([]models.BadgeDTO, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.BadgeDTO), args.Error(1)
}

// Helper function to setup test environment
func setupBadgeHandlerTest() (*gin.Engine, *mockBadgeService, *handlers.BadgeHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(mockBadgeService)
	appLogger := logger.DefaultLogger()
	badgeHandler := handlers.NewBadgeHandler(mockService, appLogger)

	return router, mockService, badgeHandler
}

func TestCreateBadge(t *testing.T) {
	router, mockService, badgeHandler := setupBadgeHandlerTest()

	// Test case: successful badge creation
	t.Run("Successful badge creation", func(t *testing.T) {
		router.POST("/badges", badgeHandler.CreateBadge)

		badgeID := uuid.New()
		testBadgeDTO := &models.BadgeDTO{
			ID:          badgeID,
			Name:        "Test Badge",
			Description: "Badge for testing",
			IconURL:     "https://example.com/badge.png",
		}

		badgeCreate := &models.BadgeCreate{
			Name:        "Test Badge",
			Description: "Badge for testing",
			IconURL:     "https://example.com/badge.png",
		}

		mockService.On("CreateBadge", mock.Anything, mock.AnythingOfType("*models.BadgeCreate")).Return(testBadgeDTO, nil)

		jsonData, _ := json.Marshal(badgeCreate)
		req, _ := http.NewRequest(http.MethodPost, "/badges", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: validation failure
	t.Run("Validation failure", func(t *testing.T) {
		router.POST("/badges", badgeHandler.CreateBadge)

		badgeCreate := &models.BadgeCreate{
			// Missing required fields
		}

		jsonData, _ := json.Marshal(badgeCreate)
		req, _ := http.NewRequest(http.MethodPost, "/badges", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

func TestGetBadgeByID(t *testing.T) {
	router, mockService, badgeHandler := setupBadgeHandlerTest()

	// Test case: successful badge retrieval
	t.Run("Successful badge retrieval", func(t *testing.T) {
		router.GET("/badges/:id", badgeHandler.GetBadgeByID)

		badgeID := uuid.New()
		testBadgeDTO := &models.BadgeDTO{
			ID:          badgeID,
			Name:        "Test Badge",
			Description: "Badge for testing",
			IconURL:     "https://example.com/badge.png",
		}

		mockService.On("GetBadgeByID", mock.Anything, badgeID.String()).Return(testBadgeDTO, nil)

		req, _ := http.NewRequest(http.MethodGet, "/badges/"+badgeID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: badge not found
	t.Run("Badge not found", func(t *testing.T) {
		router.GET("/badges/:id", badgeHandler.GetBadgeByID)

		badgeID := uuid.New()
		mockService.On("GetBadgeByID", mock.Anything, badgeID.String()).Return(nil, errors.New(errors.ErrNotFound, "Badge not found", nil))

		req, _ := http.NewRequest(http.MethodGet, "/badges/"+badgeID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListBadges(t *testing.T) {
	router, mockService, badgeHandler := setupBadgeHandlerTest()

	// Test case: successful badge listing
	t.Run("Successful badge listing", func(t *testing.T) {
		router.GET("/badges", badgeHandler.ListBadges)

		badgeID1 := uuid.New()
		badgeID2 := uuid.New()
		testBadges := []models.BadgeDTO{
			{
				ID:          badgeID1,
				Name:        "Test Badge 1",
				Description: "Badge 1 for testing",
				IconURL:     "https://example.com/badge1.png",
			},
			{
				ID:          badgeID2,
				Name:        "Test Badge 2",
				Description: "Badge 2 for testing",
				IconURL:     "https://example.com/badge2.png",
			},
		}

		mockService.On("ListBadges", mock.Anything, mock.AnythingOfType("base.ListOptions")).Return(testBadges, len(testBadges), nil)
		mockService.On("CountBadges", mock.Anything, mock.AnythingOfType("[]base.FilterOption")).Return(len(testBadges), nil)

		req, _ := http.NewRequest(http.MethodGet, "/badges", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateBadge(t *testing.T) {
	router, mockService, badgeHandler := setupBadgeHandlerTest()

	// Test case: successful badge update
	t.Run("Successful badge update", func(t *testing.T) {
		router.PUT("/badges/:id", badgeHandler.UpdateBadge)

		badgeID := uuid.New()
		updatedBadgeDTO := &models.BadgeDTO{
			ID:          badgeID,
			Name:        "Updated Badge",
			Description: "Updated description",
			IconURL:     "https://example.com/updated.png",
		}

		badgeUpdate := &models.BadgeUpdate{
			Name:        "Updated Badge",
			Description: "Updated description",
			IconURL:     "https://example.com/updated.png",
		}

		mockService.On("UpdateBadge", mock.Anything, badgeID.String(), mock.AnythingOfType("*models.BadgeUpdate")).Return(updatedBadgeDTO, nil)

		jsonData, _ := json.Marshal(badgeUpdate)
		req, _ := http.NewRequest(http.MethodPut, "/badges/"+badgeID.String(), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteBadge(t *testing.T) {
	router, mockService, badgeHandler := setupBadgeHandlerTest()

	// Test case: successful badge deletion
	t.Run("Successful badge deletion", func(t *testing.T) {
		router.DELETE("/badges/:id", badgeHandler.DeleteBadge)

		badgeID := uuid.New()
		mockService.On("DeleteBadge", mock.Anything, badgeID.String()).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/badges/"+badgeID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}
