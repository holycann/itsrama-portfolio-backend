package services_test

import (
	"context"
	"mime/multipart"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/services"
	placeModels "github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
)

// Mock event repository
type mockEventRepository struct {
	mock.Mock
}

func (m *mockEventRepository) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *mockEventRepository) FindByID(ctx context.Context, id string) (*models.EventDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.EventDTO), args.Error(1)
}

func (m *mockEventRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.EventDTO, error) {
	args := m.Called(ctx, field, value)
	return args.Get(0).([]models.EventDTO), args.Error(1)
}

func (m *mockEventRepository) Update(ctx context.Context, event *models.Event) (*models.Event, error) {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *mockEventRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockEventRepository) List(ctx context.Context, opts base.ListOptions) ([]models.EventDTO, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.EventDTO), args.Error(1)
}

func (m *mockEventRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *mockEventRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *mockEventRepository) FindPopularEvents(ctx context.Context, limit int) ([]models.EventDTO, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.EventDTO), args.Error(1)
}

func (m *mockEventRepository) FindEventsByLocation(ctx context.Context, locationID string) ([]models.EventDTO, error) {
	args := m.Called(ctx, locationID)
	return args.Get(0).([]models.EventDTO), args.Error(1)
}

func (m *mockEventRepository) FindRelatedEvents(ctx context.Context, eventID, locationID string, limit int) ([]models.EventDTO, error) {
	args := m.Called(ctx, eventID, locationID, limit)
	return args.Get(0).([]models.EventDTO), args.Error(1)
}

func (m *mockEventRepository) GetEventViews(ctx context.Context, eventID string) (int, error) {
	args := m.Called(ctx, eventID)
	return args.Int(0), args.Error(1)
}

func (m *mockEventRepository) UpdateViews(ctx context.Context, userID, eventID string) string {
	args := m.Called(ctx, userID, eventID)
	return args.String(0)
}

func (m *mockEventRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.EventDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.EventDTO), args.Int(1), args.Error(2)
}

// Mock location service
type mockLocationService struct {
	mock.Mock
}

func (m *mockLocationService) CreateLocation(ctx context.Context, location *placeModels.LocationCreate) (*placeModels.Location, error) {
	args := m.Called(ctx, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*placeModels.Location), args.Error(1)
}

func (m *mockLocationService) GetLocationByID(ctx context.Context, id string) (*placeModels.LocationDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*placeModels.LocationDTO), args.Error(1)
}

func (m *mockLocationService) GetLocationByName(ctx context.Context, name string) (*placeModels.LocationDTO, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*placeModels.LocationDTO), args.Error(1)
}

// Mock Supabase storage
type mockSupabaseStorage struct {
	mock.Mock
}

func (m *mockSupabaseStorage) UploadFile(ctx context.Context, bucket, path string, file *multipart.FileHeader) (string, error) {
	args := m.Called(ctx, bucket, path, file)
	return args.String(0), args.Error(1)
}

func (m *mockSupabaseStorage) DeleteFile(ctx context.Context, bucket, path string) error {
	args := m.Called(ctx, bucket, path)
	return args.Error(0)
}

func (m *mockSupabaseStorage) GetPublicURL(bucket, path string) string {
	args := m.Called(bucket, path)
	return args.String(0)
}

// Test setup helper
func setupEventServiceTest() (*mockEventRepository, *mockLocationService, *mockSupabaseStorage, services.EventService) {
	mockRepo := new(mockEventRepository)
	mockLocationService := new(mockLocationService)
	mockStorage := new(mockSupabaseStorage)
	eventService := services.NewEventService(mockRepo, mockLocationService, mockStorage)
	return mockRepo, mockLocationService, mockStorage, eventService
}

func TestGetEventByID(t *testing.T) {
	mockRepo, _, _, eventService := setupEventServiceTest()
	ctx := context.Background()

	// Test case: successful event retrieval
	t.Run("Successful event retrieval", func(t *testing.T) {
		eventID := uuid.New()
		userID := uuid.New()
		locationID := uuid.New()
		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)

		event := &models.EventDTO{
			ID:            eventID,
			UserID:        userID,
			LocationID:    locationID,
			Name:          "Test Event",
			Description:   "Event for testing",
			ImageURL:      "https://example.com/event.jpg",
			StartDate:     now,
			EndDate:       tomorrow,
			IsKidFriendly: true,
			CreatedAt:     &now,
		}

		// Mock repository behavior
		mockRepo.On("FindByID", ctx, eventID.String()).Return(event, nil)

		// Call the service
		result, err := eventService.GetEventByID(ctx, eventID.String())

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, eventID, result.ID)
		assert.Equal(t, userID, result.UserID)
		assert.Equal(t, locationID, result.LocationID)
		assert.Equal(t, event.Name, result.Name)
		assert.Equal(t, event.Description, result.Description)

		mockRepo.AssertExpectations(t)
	})

	// Test case: event not found
	t.Run("Event not found", func(t *testing.T) {
		eventID := uuid.New()

		// Mock repository behavior
		mockRepo.On("FindByID", ctx, eventID.String()).Return(nil, errors.New(errors.ErrNotFound, "Event not found", nil))

		// Call the service
		result, err := eventService.GetEventByID(ctx, eventID.String())

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.ErrNotFound, errors.GetType(err))

		mockRepo.AssertExpectations(t)
	})
}

func TestListEvents(t *testing.T) {
	mockRepo, _, _, eventService := setupEventServiceTest()
	ctx := context.Background()

	// Test case: successful event listing
	t.Run("Successful event listing", func(t *testing.T) {
		eventID1 := uuid.New()
		eventID2 := uuid.New()
		userID := uuid.New()
		locationID := uuid.New()
		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)

		events := []models.EventDTO{
			{
				ID:            eventID1,
				UserID:        userID,
				LocationID:    locationID,
				Name:          "Test Event 1",
				Description:   "Event 1 for testing",
				ImageURL:      "https://example.com/event1.jpg",
				StartDate:     now,
				EndDate:       tomorrow,
				IsKidFriendly: true,
				CreatedAt:     &now,
			},
			{
				ID:            eventID2,
				UserID:        userID,
				LocationID:    locationID,
				Name:          "Test Event 2",
				Description:   "Event 2 for testing",
				ImageURL:      "https://example.com/event2.jpg",
				StartDate:     now,
				EndDate:       tomorrow,
				IsKidFriendly: false,
				CreatedAt:     &now,
			},
		}

		listOpts := base.ListOptions{
			Page:     1,
			PageSize: 10,
		}

		// Mock repository behavior
		mockRepo.On("List", ctx, listOpts).Return(events, nil)

		// Call the service
		results, err := eventService.ListEvents(ctx, listOpts)

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, eventID1, results[0].ID)
		assert.Equal(t, eventID2, results[1].ID)

		mockRepo.AssertExpectations(t)
	})

	// Test case: empty list
	t.Run("Empty event list", func(t *testing.T) {
		listOpts := base.ListOptions{
			Page:     1,
			PageSize: 10,
		}

		// Mock repository behavior
		mockRepo.On("List", ctx, listOpts).Return([]models.EventDTO{}, nil)

		// Call the service
		results, err := eventService.ListEvents(ctx, listOpts)

		// Assertions
		assert.NoError(t, err)
		assert.Empty(t, results)

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteEvent(t *testing.T) {
	mockRepo, _, _, eventService := setupEventServiceTest()
	ctx := context.Background()

	// Test case: successful event deletion
	t.Run("Successful event deletion", func(t *testing.T) {
		eventID := uuid.New()

		// Mock repository behavior
		mockRepo.On("Exists", ctx, eventID.String()).Return(true, nil)
		mockRepo.On("Delete", ctx, eventID.String()).Return(nil)

		// Call the service
		err := eventService.DeleteEvent(ctx, eventID.String())

		// Assertions
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	// Test case: event not found
	t.Run("Event not found for deletion", func(t *testing.T) {
		eventID := uuid.New()

		// Mock repository behavior
		mockRepo.On("Exists", ctx, eventID.String()).Return(false, nil)

		// Call the service
		err := eventService.DeleteEvent(ctx, eventID.String())

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, errors.ErrNotFound, errors.GetType(err))

		mockRepo.AssertExpectations(t)
	})
}

func TestGetTrendingEvents(t *testing.T) {
	mockRepo, _, _, eventService := setupEventServiceTest()
	ctx := context.Background()

	// Test case: successful retrieval of trending events
	t.Run("Successful retrieval of trending events", func(t *testing.T) {
		eventID1 := uuid.New()
		eventID2 := uuid.New()
		userID := uuid.New()
		locationID := uuid.New()
		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)
		limit := 5

		events := []models.EventDTO{
			{
				ID:            eventID1,
				UserID:        userID,
				LocationID:    locationID,
				Name:          "Trending Event 1",
				Description:   "Trending event 1 for testing",
				ImageURL:      "https://example.com/event1.jpg",
				StartDate:     now,
				EndDate:       tomorrow,
				IsKidFriendly: true,
				CreatedAt:     &now,
				Views:         100,
			},
			{
				ID:            eventID2,
				UserID:        userID,
				LocationID:    locationID,
				Name:          "Trending Event 2",
				Description:   "Trending event 2 for testing",
				ImageURL:      "https://example.com/event2.jpg",
				StartDate:     now,
				EndDate:       tomorrow,
				IsKidFriendly: false,
				CreatedAt:     &now,
				Views:         80,
			},
		}

		// Mock repository behavior
		mockRepo.On("FindPopularEvents", ctx, limit).Return(events, nil)

		// Call the service
		results, err := eventService.GetTrendingEvents(ctx, limit)

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, eventID1, results[0].ID)
		assert.Equal(t, eventID2, results[1].ID)
		assert.Equal(t, 100, results[0].Views)
		assert.Equal(t, 80, results[1].Views)

		mockRepo.AssertExpectations(t)
	})
}

func TestGetRelatedEvents(t *testing.T) {
	mockRepo, _, _, eventService := setupEventServiceTest()
	ctx := context.Background()

	// Test case: successful retrieval of related events
	t.Run("Successful retrieval of related events", func(t *testing.T) {
		eventID := uuid.New()
		relatedEventID := uuid.New()
		userID := uuid.New()
		locationID := uuid.New()
		now := time.Now()
		tomorrow := now.Add(24 * time.Hour)
		limit := 3

		events := []models.EventDTO{
			{
				ID:            relatedEventID,
				UserID:        userID,
				LocationID:    locationID,
				Name:          "Related Event",
				Description:   "Related event for testing",
				ImageURL:      "https://example.com/related.jpg",
				StartDate:     now,
				EndDate:       tomorrow,
				IsKidFriendly: true,
				CreatedAt:     &now,
			},
		}

		// Mock repository behavior
		mockRepo.On("FindRelatedEvents", ctx, eventID.String(), locationID.String(), limit).Return(events, nil)

		// Call the service
		results, err := eventService.GetRelatedEvents(ctx, eventID.String(), locationID.String(), limit)

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, relatedEventID, results[0].ID)
		assert.Equal(t, locationID, results[0].LocationID)

		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateEventViews(t *testing.T) {
	mockRepo, _, _, eventService := setupEventServiceTest()
	ctx := context.Background()

	// Test case: successful update of event views
	t.Run("Successful update of event views", func(t *testing.T) {
		eventID := uuid.New().String()
		userID := uuid.New().String()
		viewID := uuid.New().String()

		// Mock repository behavior
		mockRepo.On("UpdateViews", ctx, userID, eventID).Return(viewID)

		// Call the service
		result := eventService.UpdateEventViews(ctx, userID, eventID)

		// Assertions
		assert.Equal(t, viewID, result)
		mockRepo.AssertExpectations(t)
	})
}
