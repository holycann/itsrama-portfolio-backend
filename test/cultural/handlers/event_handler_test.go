package handlers_test

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/holycann/cultour-backend/internal/cultural/handlers"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	placeModels "github.com/holycann/cultour-backend/internal/place/models"
	userModels "github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/logger"
)

// Mock event service
type mockEventService struct {
	mock.Mock
}

func (m *mockEventService) CreateEvent(ctx context.Context, event *models.EventPayload, image *multipart.FileHeader) (*models.EventDTO, error) {
	args := m.Called(ctx, event, image)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EventDTO), args.Error(1)
}

func (m *mockEventService) GetEventByID(ctx context.Context, id string) (*models.EventDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EventDTO), args.Error(1)
}

func (m *mockEventService) ListEvents(ctx context.Context, opts base.ListOptions) ([]models.EventDTO, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.EventDTO), args.Error(1)
}

func (m *mockEventService) UpdateEvent(ctx context.Context, event *models.EventPayload, image *multipart.FileHeader) (*models.EventDTO, error) {
	args := m.Called(ctx, event, image)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EventDTO), args.Error(1)
}

func (m *mockEventService) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockEventService) CountEvents(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *mockEventService) UpdateEventViews(ctx context.Context, userID, eventID string) string {
	args := m.Called(ctx, userID, eventID)
	return args.String(0)
}

func (m *mockEventService) GetTrendingEvents(ctx context.Context, limit int) ([]models.EventDTO, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.EventDTO), args.Error(1)
}

func (m *mockEventService) GetRelatedEvents(ctx context.Context, eventID, locationID string, limit int) ([]models.EventDTO, error) {
	args := m.Called(ctx, eventID, locationID, limit)
	return args.Get(0).([]models.EventDTO), args.Error(1)
}

func (m *mockEventService) SearchEvents(ctx context.Context, query string, opts base.ListOptions) ([]models.EventDTO, error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).([]models.EventDTO), args.Error(1)
}

// Helper function to setup test environment
func setupEventHandlerTest() (*gin.Engine, *mockEventService, *handlers.EventHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(mockEventService)
	appLogger := logger.DefaultLogger()
	eventHandler := handlers.NewEventHandler(mockService, appLogger)

	return router, mockService, eventHandler
}

func TestGetEventByID(t *testing.T) {
	router, mockService, eventHandler := setupEventHandlerTest()

	// Test case: successful event retrieval
	t.Run("Successful event retrieval", func(t *testing.T) {
		router.GET("/events/:id", eventHandler.GetEventByID)

		eventID := uuid.New()
		userID := uuid.New()
		locationID := uuid.New()

		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)

		testEventDTO := &models.EventDTO{
			ID:            eventID,
			Name:          "Test Event",
			Description:   "Event for testing",
			ImageURL:      "https://example.com/event.png",
			StartDate:     now,
			EndDate:       tomorrow,
			IsKidFriendly: true,
			Creator: &userModels.User{
				ID:    userID,
				Email: "user@example.com",
			},
			Location: &placeModels.Location{
				ID:        locationID,
				Name:      "Test Location",
				Latitude:  1.234,
				Longitude: 4.567,
			},
		}

		mockService.On("GetEventByID", mock.Anything, eventID.String()).Return(testEventDTO, nil)

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: event not found
	t.Run("Event not found", func(t *testing.T) {
		router.GET("/events/:id", eventHandler.GetEventByID)

		eventID := uuid.New()
		mockService.On("GetEventByID", mock.Anything, eventID.String()).Return(nil, errors.New(errors.ErrNotFound, "Event not found", nil))

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestListEvents(t *testing.T) {
	router, mockService, eventHandler := setupEventHandlerTest()

	// Test case: successful event listing
	t.Run("Successful event listing", func(t *testing.T) {
		router.GET("/events", eventHandler.ListEvents)

		eventID1 := uuid.New()
		eventID2 := uuid.New()

		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)

		testEvents := []models.EventDTO{
			{
				ID:            eventID1,
				Name:          "Test Event 1",
				Description:   "Event 1 for testing",
				ImageURL:      "https://example.com/event1.png",
				StartDate:     now,
				EndDate:       tomorrow,
				IsKidFriendly: true,
			},
			{
				ID:            eventID2,
				Name:          "Test Event 2",
				Description:   "Event 2 for testing",
				ImageURL:      "https://example.com/event2.png",
				StartDate:     now,
				EndDate:       tomorrow,
				IsKidFriendly: false,
			},
		}

		mockService.On("ListEvents", mock.Anything, mock.AnythingOfType("base.ListOptions")).Return(testEvents, nil)

		req, _ := http.NewRequest(http.MethodGet, "/events", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteEvent(t *testing.T) {
	router, mockService, eventHandler := setupEventHandlerTest()

	// Test case: successful event deletion
	t.Run("Successful event deletion", func(t *testing.T) {
		router.DELETE("/events/:id", eventHandler.DeleteEvent)

		eventID := uuid.New()
		mockService.On("DeleteEvent", mock.Anything, eventID.String()).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/events/"+eventID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})

	// Test case: event not found
	t.Run("Event not found for deletion", func(t *testing.T) {
		router.DELETE("/events/:id", eventHandler.DeleteEvent)

		eventID := uuid.New()
		mockService.On("DeleteEvent", mock.Anything, eventID.String()).Return(errors.New(errors.ErrNotFound, "Event not found", nil))

		req, _ := http.NewRequest(http.MethodDelete, "/events/"+eventID.String(), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetTrendingEvents(t *testing.T) {
	router, mockService, eventHandler := setupEventHandlerTest()

	// Test case: successful trending events retrieval
	t.Run("Successful trending events retrieval", func(t *testing.T) {
		router.GET("/events/trending", eventHandler.GetTrendingEvents)

		eventID1 := uuid.New()
		eventID2 := uuid.New()

		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)

		testEvents := []models.EventDTO{
			{
				ID:          eventID1,
				Name:        "Trending Event 1",
				Description: "Trending event 1 for testing",
				ImageURL:    "https://example.com/trending1.png",
				StartDate:   now,
				EndDate:     tomorrow,
				Views:       100,
			},
			{
				ID:          eventID2,
				Name:        "Trending Event 2",
				Description: "Trending event 2 for testing",
				ImageURL:    "https://example.com/trending2.png",
				StartDate:   now,
				EndDate:     tomorrow,
				Views:       80,
			},
		}

		mockService.On("GetTrendingEvents", mock.Anything, mock.AnythingOfType("int")).Return(testEvents, nil)

		req, _ := http.NewRequest(http.MethodGet, "/events/trending", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertExpectations(t)
	})
}
