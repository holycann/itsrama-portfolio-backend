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

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.UserCreate) (*models.UserDTO, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.UserDTO), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, user *models.UserUpdate) (*models.UserDTO, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.UserDTO), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.UserDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.UserDTO), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.UserDTO, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.UserDTO), args.Error(1)
}

func (m *MockUserService) ListUsers(ctx context.Context, opts base.ListOptions) ([]models.UserDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.UserDTO), args.Int(1), args.Error(2)
}

func (m *MockUserService) SearchUsers(ctx context.Context, opts base.ListOptions) ([]models.UserDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.UserDTO), args.Int(1), args.Error(2)
}

func (m *MockUserService) CountUsers(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func setupTestRouter(userHandler *handlers.UserHandler) *gin.Engine {
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

	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.ListUsers)
	router.GET("/users/search", userHandler.SearchUsers)
	router.GET("/users/:id", userHandler.GetUserByID)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	return router
}

func TestCreateUser(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockUserService, mockLogger)

	router := setupTestRouter(userHandler)

	userCreate := models.UserCreate{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "authenticated",
	}

	// Prepare request body
	jsonBody, _ := json.Marshal(userCreate)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Mock service response
	mockUserService.On("CreateUser", mock.Anything, &userCreate).Return(&models.UserDTO{
		ID:    uuid.New(),
		Email: userCreate.Email,
		Role:  userCreate.Role,
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
	assert.Equal(t, "User created successfully", response["message"])

	mockUserService.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockUserService, mockLogger)

	router := setupTestRouter(userHandler)

	// Prepare list options
	listOptions := base.ListOptions{
		Page:      1,
		PerPage:   10,
		SortBy:    "created_at",
		SortOrder: base.SortDescending,
	}

	// Mock service response
	expectedUsers := []models.UserDTO{
		{
			ID:    uuid.New(),
			Email: "user1@example.com",
			Role:  "authenticated",
		},
		{
			ID:    uuid.New(),
			Email: "user2@example.com",
			Role:  "authenticated",
		},
	}

	mockUserService.On("ListUsers", mock.Anything, listOptions).Return(expectedUsers, len(expectedUsers), nil)

	// Perform request
	req, _ := http.NewRequest("GET", "/users?limit=10&offset=0&sort_by=created_at&sort_order=desc", nil)
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

	mockUserService.AssertExpectations(t)
}

func TestSearchUsers(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockUserService, mockLogger)

	router := setupTestRouter(userHandler)

	// Prepare search options
	searchOptions := base.ListOptions{
		Page:    1,
		PerPage: 10,
		Search:  "test",
	}

	// Mock service response
	expectedUsers := []models.UserDTO{
		{
			ID:    uuid.New(),
			Email: "test1@example.com",
			Role:  "authenticated",
		},
		{
			ID:    uuid.New(),
			Email: "test2@example.com",
			Role:  "authenticated",
		},
	}

	mockUserService.On("SearchUsers", mock.Anything, searchOptions).Return(expectedUsers, len(expectedUsers), nil)

	// Perform request
	req, _ := http.NewRequest("GET", "/users/search?query=test&limit=10&offset=0", nil)
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

	mockUserService.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockUserService, mockLogger)

	router := setupTestRouter(userHandler)

	userID := uuid.New()
	updateUser := models.UserUpdate{
		ID:    userID,
		Email: "updated@example.com",
	}

	// Prepare request body
	jsonBody, _ := json.Marshal(updateUser)
	req, _ := http.NewRequest("PUT", "/users/"+userID.String(), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Mock service response
	mockUserService.On("UpdateUser", mock.Anything, &updateUser).Return(&models.UserDTO{
		ID:    userID,
		Email: updateUser.Email,
	}, nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "message")
	assert.Equal(t, "User updated successfully", response["message"])

	mockUserService.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockUserService, mockLogger)

	router := setupTestRouter(userHandler)

	userID := uuid.New()

	// Mock service response
	mockUserService.On("DeleteUser", mock.Anything, userID.String()).Return(nil)

	// Perform request
	req, _ := http.NewRequest("DELETE", "/users/"+userID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "message")
	assert.Equal(t, "User deleted successfully", response["message"])

	mockUserService.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockLogger := &logger.Logger{}
	mockUserService := new(MockUserService)
	userHandler := handlers.NewUserHandler(mockUserService, mockLogger)

	router := setupTestRouter(userHandler)

	userID := uuid.New()

	// Mock service response
	mockUserService.On("GetUserByID", mock.Anything, userID.String()).Return(&models.UserDTO{
		ID:    userID,
		Email: "test@example.com",
		Role:  "authenticated",
	}, nil)

	// Perform request
	req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "data")
	assert.Contains(t, response, "message")
	assert.Equal(t, "User retrieved successfully", response["message"])

	mockUserService.AssertExpectations(t)
}
