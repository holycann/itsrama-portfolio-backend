package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	discussionServices "github.com/holycann/cultour-backend/internal/discussion/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockThreadRepository is a mock implementation of ThreadRepository
type MockThreadRepository struct {
	mock.Mock
}

func (m *MockThreadRepository) Create(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	args := m.Called(ctx, thread)
	return args.Get(0).(*models.Thread), args.Error(1)
}

func (m *MockThreadRepository) FindByID(ctx context.Context, id string) (*models.ThreadDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.ThreadDTO), args.Error(1)
}

func (m *MockThreadRepository) FindByEvent(ctx context.Context, eventID string) (*models.ThreadDTO, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).(*models.ThreadDTO), args.Error(1)
}

func (m *MockThreadRepository) List(ctx context.Context, opts base.ListOptions) ([]models.ThreadDTO, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.ThreadDTO), args.Error(1)
}

func (m *MockThreadRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.ThreadDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.ThreadDTO), args.Int(1), args.Error(2)
}

func (m *MockThreadRepository) Update(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	args := m.Called(ctx, thread)
	return args.Get(0).(*models.Thread), args.Error(1)
}

func (m *MockThreadRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockThreadRepository) GetActiveThreads(ctx context.Context, limit int) ([]models.ThreadDTO, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.ThreadDTO), args.Error(1)
}

func TestCreateThread(t *testing.T) {
	mockRepo := new(MockThreadRepository)
	threadService := discussionServices.NewThreadService(mockRepo)

	eventID := uuid.New()
	userID := uuid.New()
	threadToCreate := &models.Thread{
		EventID: eventID,
		UserID:  userID,
		Title:   "Event Discussion Thread",
	}

	mockRepo.On("Create", mock.Anything, threadToCreate).Return(threadToCreate, nil)

	result, err := threadService.CreateThread(context.Background(), threadToCreate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, threadToCreate.EventID, result.EventID)
	assert.Equal(t, threadToCreate.UserID, result.UserID)
	assert.Equal(t, threadToCreate.Title, result.Title)

	mockRepo.AssertExpectations(t)
}

func TestGetThreadByID(t *testing.T) {
	mockRepo := new(MockThreadRepository)
	threadService := discussionServices.NewThreadService(mockRepo)

	threadID := uuid.New()
	eventID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	expectedThread := &models.ThreadDTO{
		ID:        threadID,
		EventID:   eventID,
		UserID:    userID,
		Title:     "Event Discussion Thread",
		CreatedAt: &now,
	}

	mockRepo.On("FindByID", mock.Anything, threadID.String()).Return(expectedThread, nil)

	result, err := threadService.GetThreadByID(context.Background(), threadID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedThread.ID, result.ID)
	assert.Equal(t, expectedThread.EventID, result.EventID)
	assert.Equal(t, expectedThread.Title, result.Title)

	mockRepo.AssertExpectations(t)
}

func TestGetThreadByEvent(t *testing.T) {
	mockRepo := new(MockThreadRepository)
	threadService := discussionServices.NewThreadService(mockRepo)

	eventID := uuid.New()
	threadID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	expectedThread := &models.ThreadDTO{
		ID:        threadID,
		EventID:   eventID,
		UserID:    userID,
		Title:     "Event Discussion Thread",
		CreatedAt: &now,
	}

	mockRepo.On("FindByEvent", mock.Anything, eventID.String()).Return(expectedThread, nil)

	result, err := threadService.GetThreadByEvent(context.Background(), eventID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedThread.ID, result.ID)
	assert.Equal(t, expectedThread.EventID, result.EventID)
	assert.Equal(t, expectedThread.Title, result.Title)

	mockRepo.AssertExpectations(t)
}

func TestListThreads(t *testing.T) {
	mockRepo := new(MockThreadRepository)
	threadService := discussionServices.NewThreadService(mockRepo)

	listOptions := base.ListOptions{
		Page:    1,
		PerPage: 10,
		SortBy:  "created_at",
	}

	now := time.Now()
	expectedThreads := []models.ThreadDTO{
		{
			ID:        uuid.New(),
			EventID:   uuid.New(),
			UserID:    uuid.New(),
			Title:     "Thread 1",
			CreatedAt: &now,
		},
		{
			ID:        uuid.New(),
			EventID:   uuid.New(),
			UserID:    uuid.New(),
			Title:     "Thread 2",
			CreatedAt: &now,
		},
	}

	mockRepo.On("List", mock.Anything, listOptions).Return(expectedThreads, nil)

	result, err := threadService.ListThreads(context.Background(), listOptions)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedThreads[0].Title, result[0].Title)
	assert.Equal(t, expectedThreads[1].Title, result[1].Title)

	mockRepo.AssertExpectations(t)
}

func TestSearchThreads(t *testing.T) {
	mockRepo := new(MockThreadRepository)
	threadService := discussionServices.NewThreadService(mockRepo)

	searchOptions := base.ListOptions{
		Page:    1,
		PerPage: 10,
		SortBy:  "created_at",
	}

	now := time.Now()
	expectedThreads := []models.ThreadDTO{
		{
			ID:        uuid.New(),
			EventID:   uuid.New(),
			UserID:    uuid.New(),
			Title:     "Thread 1",
			CreatedAt: &now,
		},
		{
			ID:        uuid.New(),
			EventID:   uuid.New(),
			UserID:    uuid.New(),
			Title:     "Thread 2",
			CreatedAt: &now,
		},
	}

	mockRepo.On("Search", mock.Anything, searchOptions).Return(expectedThreads, len(expectedThreads), nil)

	result, total, err := threadService.SearchThreads(context.Background(), "", searchOptions)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, len(expectedThreads), total)
	assert.Equal(t, expectedThreads[0].Title, result[0].Title)
	assert.Equal(t, expectedThreads[1].Title, result[1].Title)

	mockRepo.AssertExpectations(t)
}

func TestGetActiveThreads(t *testing.T) {
	mockRepo := new(MockThreadRepository)
	threadService := discussionServices.NewThreadService(mockRepo)

	limit := 5
	now := time.Now()
	expectedThreads := []models.ThreadDTO{
		{
			ID:        uuid.New(),
			EventID:   uuid.New(),
			UserID:    uuid.New(),
			Title:     "Active Thread 1",
			CreatedAt: &now,
		},
		{
			ID:        uuid.New(),
			EventID:   uuid.New(),
			UserID:    uuid.New(),
			Title:     "Active Thread 2",
			CreatedAt: &now,
		},
	}

	mockRepo.On("GetActiveThreads", mock.Anything, limit).Return(expectedThreads, nil)

	result, err := threadService.GetActiveThreads(context.Background(), limit)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedThreads[0].Title, result[0].Title)
	assert.Equal(t, expectedThreads[1].Title, result[1].Title)

	mockRepo.AssertExpectations(t)
}
