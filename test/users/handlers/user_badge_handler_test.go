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
	"github.com/holycann/cultour-backend/internal/users/handlers"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserBadgeService is a mock implementation of UserBadgeService
type MockUserBadgeService struct {
	mock.Mock
}

func (m *MockUserBadgeService) AddBadgeToUser(ctx context.Context, payload models.UserBadge) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockUserBadgeService) RemoveBadgeFromUser(ctx context.Context, payload models.UserBadgePayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockUserBadgeService) GetUserBadgeByID(ctx context.Context, id string) (*models.UserBadgeDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.UserBadgeDTO), args.Error(1)
}

func (m *MockUserBadgeService) GetUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadgeDTO, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.UserBadgeDTO), args.Error(1)
}

func (m *MockUserBadgeService) GetUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadgeDTO, error) {
	args := m.Called(ctx, badgeID)
	return args.Get(0).([]models.UserBadgeDTO), args.Error(1)
}

func (m *MockUserBadgeService) ListUserBadges(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.UserBadgeDTO), args.Int(1), args.Error(2)
}

func (m *MockUserBadgeService) SearchUserBadges(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.UserBadgeDTO), args.Int(1), args.Error(2)
}

func (m *MockUserBadgeService) DeleteUserBadge(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserBadgeService) CountUserBadges(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func setupUserBadgeTestRouter(userBadgeHandler *handlers.UserBadgeHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add middleware to set user context
	router.Use(func(c *gin.Context) {
		userID := uuid.New().String()
		c.Set("user_id", userID)
		c.Set("email", "test@example.com")
		c.Set("role", "authenticated")
		c.Set("token", "test_token")
		c.Next()
	})

	router.POST("/users/badges", userBadgeHandler.AssignBadge)
	router.GET("/users/badges", userBadgeHandler.GetUserBadges)
	router.DELETE("/users/badges", userBadgeHandler.RemoveBadge)
	router.GET("/users/badges/count", userBadgeHandler.CountUserBadges)
	router.GET("/users/badges/me", userBadgeHandler.GetUserBadgesByUser)

	return router
}

func TestAssignBadge(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserBadgeService := new(MockUserBadgeService)
	userBadgeHandler := handlers.NewUserBadgeHandler(mockUserBadgeService, mockLogger)

	router := setupUserBadgeTestRouter(userBadgeHandler)

	userID := uuid.New()
	badgeID := uuid.New()
	payload := models.UserBadge{
		UserID:  userID,
		BadgeID: badgeID,
	}

	// Prepare request body
	jsonBody, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/users/badges", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Mock service response
	mockUserBadgeService.On("AddBadgeToUser", mock.Anything, payload).Return(nil)
	mockUserBadgeService.On("GetUserBadgesByUser", mock.Anything, userID.String()).Return([]models.UserBadgeDTO{
		{
			BadgeID: badgeID,
		},
	}, nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "message")
	assert.Equal(t, "Badge successfully assigned", response["message"])

	mockUserBadgeService.AssertExpectations(t)
}

func TestGetUserBadges(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserBadgeService := new(MockUserBadgeService)
	userBadgeHandler := handlers.NewUserBadgeHandler(mockUserBadgeService, mockLogger)

	router := setupUserBadgeTestRouter(userBadgeHandler)

	// Prepare list options
	listOptions := base.ListOptions{
		Page:      1,
		PerPage:   10,
		SortBy:    "created_at",
		SortOrder: base.SortDescending,
	}

	// Mock service response
	expectedBadges := []models.UserBadgeDTO{
		{
			BadgeID: uuid.New(),
		},
		{
			BadgeID: uuid.New(),
		},
	}

	mockUserBadgeService.On("ListUserBadges", mock.Anything, listOptions).Return(expectedBadges, len(expectedBadges), nil)

	// Perform request
	req, _ := http.NewRequest("GET", "/users/badges?limit=10&offset=0&sort_by=created_at&sort_order=desc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "pagination")
	assert.Contains(t, response, "message")

	mockUserBadgeService.AssertExpectations(t)
}

func TestRemoveBadge(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserBadgeService := new(MockUserBadgeService)
	userBadgeHandler := handlers.NewUserBadgeHandler(mockUserBadgeService, mockLogger)

	router := setupUserBadgeTestRouter(userBadgeHandler)

	userID := uuid.New()
	badgeID := uuid.New()
	payload := models.UserBadgePayload{
		UserID:  userID,
		BadgeID: badgeID,
	}

	// Perform request
	req, _ := http.NewRequest("DELETE", "/users/badges?badge_id="+badgeID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Mock service response
	mockUserBadgeService.On("RemoveBadgeFromUser", mock.Anything, payload).Return(nil)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "message")
	assert.Equal(t, "Badge removed successfully", response["message"])

	mockUserBadgeService.AssertExpectations(t)
}

func TestCountUserBadges(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserBadgeService := new(MockUserBadgeService)
	userBadgeHandler := handlers.NewUserBadgeHandler(mockUserBadgeService, mockLogger)

	router := setupUserBadgeTestRouter(userBadgeHandler)

	// Mock service response
	expectedCount := 5

	mockUserBadgeService.On("CountUserBadges", mock.Anything, mock.Anything).Return(expectedCount, nil)

	// Perform request
	req, _ := http.NewRequest("GET", "/users/badges/count", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "message")
	assert.Equal(t, "User badges counted successfully", response["message"])
	assert.Equal(t, float64(expectedCount), response["data"])

	mockUserBadgeService.AssertExpectations(t)
}

func TestGetUserBadgesByUser(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserBadgeService := new(MockUserBadgeService)
	userBadgeHandler := handlers.NewUserBadgeHandler(mockUserBadgeService, mockLogger)

	router := setupUserBadgeTestRouter(userBadgeHandler)

	// Mock service response
	expectedBadges := []models.UserBadgeDTO{
		{
			BadgeID: uuid.New(),
		},
		{
			BadgeID: uuid.New(),
		},
	}

	mockUserBadgeService.On("GetUserBadgesByUser", mock.Anything, mock.Anything).Return(expectedBadges, nil)

	// Perform request
	req, _ := http.NewRequest("GET", "/users/badges/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "message")
	assert.Equal(t, "User badges retrieved successfully", response["message"])

	mockUserBadgeService.AssertExpectations(t)
}
