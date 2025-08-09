package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/models"
	userServices "github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserBadgeRepository is a mock implementation of UserBadgeRepository
type MockUserBadgeRepository struct {
	mock.Mock
}

func (m *MockUserBadgeRepository) AddBadgeToUser(ctx context.Context, payload *models.UserBadge) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockUserBadgeRepository) RemoveBadgeFromUser(ctx context.Context, payload *models.UserBadgePayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockUserBadgeRepository) FindByID(ctx context.Context, id string) (*models.UserBadgeDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.UserBadgeDTO), args.Error(1)
}

func (m *MockUserBadgeRepository) FindUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadgeDTO, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.UserBadgeDTO), args.Error(1)
}

func (m *MockUserBadgeRepository) FindUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadgeDTO, error) {
	args := m.Called(ctx, badgeID)
	return args.Get(0).([]models.UserBadgeDTO), args.Error(1)
}

func (m *MockUserBadgeRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.UserBadgeDTO), args.Int(1), args.Error(2)
}

func (m *MockUserBadgeRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockUserBadgeRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserBadgeRepository) BulkCreate(ctx context.Context, badges []*models.UserBadge) ([]models.UserBadge, error) {
	args := m.Called(ctx, badges)
	return args.Get(0).([]models.UserBadge), args.Error(1)
}

func (m *MockUserBadgeRepository) BulkDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockUserBadgeRepository) BulkUpdate(ctx context.Context, badges []*models.UserBadge) ([]models.UserBadge, error) {
	args := m.Called(ctx, badges)
	return args.Get(0).([]models.UserBadge), args.Error(1)
}

func (m *MockUserBadgeRepository) CountUserBadges(ctx context.Context, userID string) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockUserBadgeRepository) Create(ctx context.Context, value *models.UserBadge) (*models.UserBadge, error) {
	args := m.Called(ctx, value)
	return args.Get(0).(*models.UserBadge), args.Error(1)
}

func (m *MockUserBadgeRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserBadgeRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.UserBadgeDTO, error) {
	args := m.Called(ctx, field, value)
	return args.Get(0).([]models.UserBadgeDTO), args.Error(1)
}

func (m *MockUserBadgeRepository) List(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.UserBadgeDTO), args.Error(1)
}

func (m *MockUserBadgeRepository) Update(ctx context.Context, value *models.UserBadge) (*models.UserBadge, error) {
	args := m.Called(ctx, value)
	return args.Get(0).(*models.UserBadge), args.Error(1)
}

func TestAddBadgeToUser(t *testing.T) {
	mockRepo := new(MockUserBadgeRepository)
	userBadgeService := userServices.NewUserBadgeService(mockRepo)

	userID := uuid.New()
	badgeID := uuid.New()
	payload := models.UserBadge{
		UserID:  userID,
		BadgeID: badgeID,
	}

	// Mock existing badges check
	mockRepo.On("FindUserBadgesByUser", mock.Anything, userID.String()).Return([]models.UserBadgeDTO{}, nil)

	// Mock badge assignment
	mockRepo.On("AddBadgeToUser", mock.Anything, &payload).Return(nil)

	err := userBadgeService.AddBadgeToUser(context.Background(), payload)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAddDuplicateBadgeToUser(t *testing.T) {
	mockRepo := new(MockUserBadgeRepository)
	userBadgeService := userServices.NewUserBadgeService(mockRepo)

	userID := uuid.New()
	badgeID := uuid.New()
	payload := models.UserBadge{
		UserID:  userID,
		BadgeID: badgeID,
	}

	// Mock existing badges check with the same badge
	existingBadges := []models.UserBadgeDTO{
		{
			BadgeID: badgeID,
		},
	}
	mockRepo.On("FindUserBadgesByUser", mock.Anything, userID.String()).Return(existingBadges, nil)

	err := userBadgeService.AddBadgeToUser(context.Background(), payload)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already assigned")
	mockRepo.AssertExpectations(t)
}

func TestRemoveBadgeFromUser(t *testing.T) {
	mockRepo := new(MockUserBadgeRepository)
	userBadgeService := userServices.NewUserBadgeService(mockRepo)

	userID := uuid.New()
	badgeID := uuid.New()
	payload := models.UserBadgePayload{
		UserID:  userID,
		BadgeID: badgeID,
	}

	// Mock badge removal
	mockRepo.On("RemoveBadgeFromUser", mock.Anything, &payload).Return(nil)

	err := userBadgeService.RemoveBadgeFromUser(context.Background(), payload)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUserBadgesByUser(t *testing.T) {
	mockRepo := new(MockUserBadgeRepository)
	userBadgeService := userServices.NewUserBadgeService(mockRepo)

	userID := uuid.New()
	badgeID1 := uuid.New()
	badgeID2 := uuid.New()

	expectedBadges := []models.UserBadgeDTO{
		{
			BadgeID: badgeID1,
		},
		{
			BadgeID: badgeID2,
		},
	}

	mockRepo.On("FindUserBadgesByUser", mock.Anything, userID.String()).Return(expectedBadges, nil)

	result, err := userBadgeService.GetUserBadgesByUser(context.Background(), userID.String())

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, badgeID1, result[0].BadgeID)
	assert.Equal(t, badgeID2, result[1].BadgeID)
	mockRepo.AssertExpectations(t)
}

func TestListUserBadges(t *testing.T) {
	mockRepo := new(MockUserBadgeRepository)
	userBadgeService := userServices.NewUserBadgeService(mockRepo)

	userID := uuid.New()
	listOptions := base.ListOptions{
		Page:      1,
		PerPage:   10,
		SortBy:    "created_at",
		SortOrder: base.SortDescending,
		Filters: []base.FilterOption{
			{
				Field:    "user_id",
				Operator: base.OperatorEqual,
				Value:    userID.String(),
			},
		},
	}

	expectedBadges := []models.UserBadgeDTO{
		{
			BadgeID: uuid.New(),
		},
		{
			BadgeID: uuid.New(),
		},
	}

	mockRepo.On("Search", mock.Anything, listOptions).Return(expectedBadges, len(expectedBadges), nil)

	result, total, err := userBadgeService.ListUserBadges(context.Background(), listOptions)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, len(expectedBadges), total)
	mockRepo.AssertExpectations(t)
}

func TestCountUserBadges(t *testing.T) {
	mockRepo := new(MockUserBadgeRepository)
	userBadgeService := userServices.NewUserBadgeService(mockRepo)

	userID := uuid.New()

	expectedCount := 3

	mockRepo.On("CountUserBadges", mock.Anything, userID.String()).Return(expectedCount, nil)

	result, err := userBadgeService.CountUserBadges(context.Background(), userID.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)
	mockRepo.AssertExpectations(t)
}
