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

	"github.com/holycann/cultour-backend/internal/discussion/handlers"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	userModels "github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
)

// Mock thread service
type mockThreadService struct {
	mock.Mock
}

func (m *mockThreadService) CreateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	args := m.Called(ctx, thread)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Thread), args.Error(1)
}

func (m *mockThreadService) GetThreadByID(ctx context.Context, id string) (*models.ThreadDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ThreadDTO), args.Error(1)
}

func (m *mockThreadService) GetThreadByEvent(ctx context.Context, eventID string) (*models.ThreadDTO, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ThreadDTO), args.Error(1)
}

func (m *mockThreadService) GetActiveThreads(ctx context.Context, limit int) ([]models.ThreadDTO, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.ThreadDTO), args.Error(1)
}

func (m *mockThreadService) ListThreads(ctx context.Context, opts base.ListOptions) ([]models.ThreadDTO, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.ThreadDTO), args.Error(1)
}

func (m *mockThreadService) SearchThreads(ctx context.Context, query string, opts base.ListOptions) ([]models.ThreadDTO, int, error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).([]models.ThreadDTO), args.Int(1), args.Error(2)
}

func (m *mockThreadService) CountThreads(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *mockThreadService) UpdateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	args := m.Called(ctx, thread)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Thread), args.Error(1)
}

func (m *mockThreadService) DeleteThread(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockThreadService) JoinThread(ctx context.Context, threadID, userID string) error {
	args := m.Called(ctx, threadID, userID)
	return args.Error(0)
}

// Helper function to setup test environment
func setupThreadHandlerTest() (*gin.Engine, *mockThreadService, *handlers.ThreadHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(mockThreadService)
	appLogger := logger.DefaultLogger()
	threadHandler := handlers.NewThreadHandler(mockService, appLogger)

	return router, mockService, threadHandler
}

func TestCreateThread(t *testing.T) {
	router, mockService, threadHandler := setupThreadHandlerTest()

	// Test case: successful thread creation
	t.Run("Successful thread creation", func(t *testing.T) {
		router.POST("/threads", threadHandler.CreateThread)

		threadID := uuid.New()
		eventID := uuid.New()
		creatorID := uuid.New()

		now := time.Now()

		newThread := &models.Thread{
			ID:        threadID,
			EventID:   eventID,
			CreatorID: creatorID,
			Status:    "active",
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		mockService.On("CreateThread", mock.Anything, mock.AnythingOfType("*models.Thread")).Return(newThread, nil)

		createThread := models.CreateThread{
			EventID:   eventID,
			CreatorID: creatorID,
			Status:    "active",
		}

		jsonData, _ := json.Marshal(createThread)
		req, _ := http.NewRequest(http.MethodPost, "/threads", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetThreadByID(t *testing.T) {
	router, mockService, threadHandler := setupThreadHandlerTest()

	// Test case: successful thread retrieval
	t.Run("Successful thread retrieval", func(t *testing.T) {
		router.GET("/threads/:id", threadHandler.GetThreadByID)

		threadID := uuid.New()
		eventID := uuid.New()
		creatorID := uuid.New()

		now := time.Now()

		testThreadDTO := &models.ThreadDTO{
			ID:        threadID,
			EventID:   eventID,
			Status:    "active",
			CreatedAt: &now,
			UpdatedAt: &now,
			Creator: &userModels.User{
				ID:    creatorID,
				Email: "creator@example.com",
			},
		}

		mockService.On("GetThreadByID", mock.Anything, threadID.String()).Return(testThreadDTO, nil)

		req, _ := http.NewRequest(http.MethodGet, "/threads/"+threadID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: thread not found
	t.Run("Thread not found", func(t *testing.T) {
		router.GET("/threads/:id", threadHandler.GetThreadByID)

		threadID := uuid.New()
		mockService.On("GetThreadByID", mock.Anything, threadID.String()).Return(nil, errors.New(errors.ErrNotFound, "Thread not found", nil))

		req, _ := http.NewRequest(http.MethodGet, "/threads/"+threadID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListThreads(t *testing.T) {
	router, mockService, threadHandler := setupThreadHandlerTest()

	// Test case: successful thread listing
	t.Run("Successful thread listing", func(t *testing.T) {
		router.GET("/threads", threadHandler.ListThreads)

		threadID1 := uuid.New()
		threadID2 := uuid.New()
		eventID1 := uuid.New()
		eventID2 := uuid.New()

		now := time.Now()

		testThreads := []models.ThreadDTO{
			{
				ID:        threadID1,
				EventID:   eventID1,
				Status:    "active",
				CreatedAt: &now,
			},
			{
				ID:        threadID2,
				EventID:   eventID2,
				Status:    "active",
				CreatedAt: &now,
			},
		}

		mockService.On("ListThreads", mock.Anything, mock.AnythingOfType("base.ListOptions")).Return(testThreads, nil)
		mockService.On("CountThreads", mock.Anything, mock.AnythingOfType("[]base.FilterOption")).Return(len(testThreads), nil)

		req, _ := http.NewRequest(http.MethodGet, "/threads", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetThreadByEvent(t *testing.T) {
	router, mockService, threadHandler := setupThreadHandlerTest()

	// Test case: successful thread retrieval by event
	t.Run("Successful thread retrieval by event", func(t *testing.T) {
		router.GET("/events/:id/thread", threadHandler.GetThreadByEvent)

		threadID := uuid.New()
		eventID := uuid.New()
		creatorID := uuid.New()

		now := time.Now()

		testThreadDTO := &models.ThreadDTO{
			ID:        threadID,
			EventID:   eventID,
			Status:    "active",
			CreatedAt: &now,
			Creator: &userModels.User{
				ID:    creatorID,
				Email: "creator@example.com",
			},
		}

		mockService.On("GetThreadByEvent", mock.Anything, eventID.String()).Return(testThreadDTO, nil)

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/thread", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: thread not found for event
	t.Run("Thread not found for event", func(t *testing.T) {
		router.GET("/events/:id/thread", threadHandler.GetThreadByEvent)

		eventID := uuid.New()
		mockService.On("GetThreadByEvent", mock.Anything, eventID.String()).Return(nil, errors.New(errors.ErrNotFound, "Thread not found for event", nil))

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/thread", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteThread(t *testing.T) {
	router, mockService, threadHandler := setupThreadHandlerTest()

	// Test case: successful thread deletion
	t.Run("Successful thread deletion", func(t *testing.T) {
		router.DELETE("/threads/:id", threadHandler.DeleteThread)

		threadID := uuid.New()
		mockService.On("DeleteThread", mock.Anything, threadID.String()).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/threads/"+threadID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}
