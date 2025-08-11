package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/internal/achievement/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
)

// Mock badge repository
type mockBadgeRepository struct {
	mock.Mock
}

func (m *mockBadgeRepository) Create(ctx context.Context, badge *models.Badge) (*models.Badge, error) {
	args := m.Called(ctx, badge)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Badge), args.Error(1)
}

func (m *mockBadgeRepository) FindByID(ctx context.Context, id string) (*models.Badge, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Badge), args.Error(1)
}

func (m *mockBadgeRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.Badge, error) {
	args := m.Called(ctx, field, value)
	return args.Get(0).([]models.Badge), args.Error(1)
}

func (m *mockBadgeRepository) Update(ctx context.Context, badge *models.Badge) (*models.Badge, error) {
	args := m.Called(ctx, badge)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Badge), args.Error(1)
}

func (m *mockBadgeRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockBadgeRepository) List(ctx context.Context, opts base.ListOptions) ([]models.Badge, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.Badge), args.Error(1)
}

func (m *mockBadgeRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *mockBadgeRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *mockBadgeRepository) FindBadgeByName(ctx context.Context, name string) (*models.Badge, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Badge), args.Error(1)
}

func (m *mockBadgeRepository) FindPopularBadges(ctx context.Context, limit int) ([]models.Badge, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]models.Badge), args.Error(1)
}

func (m *mockBadgeRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.Badge, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.Badge), args.Int(1), args.Error(2)
}

func (m *mockBadgeRepository) BulkCreate(ctx context.Context, values []*models.Badge) ([]models.Badge, error) {
	args := m.Called(ctx, values)
	return args.Get(0).([]models.Badge), args.Error(1)
}

func (m *mockBadgeRepository) BulkUpdate(ctx context.Context, values []*models.Badge) ([]models.Badge, error) {
	args := m.Called(ctx, values)
	return args.Get(0).([]models.Badge), args.Error(1)
}

func (m *mockBadgeRepository) BulkDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *mockBadgeRepository) BulkUpsert(ctx context.Context, values []*models.Badge) ([]models.Badge, error) {
	args := m.Called(ctx, values)
	return args.Get(0).([]models.Badge), args.Error(1)
}

// Test setup helper
func setupBadgeServiceTest() (*mockBadgeRepository, services.BadgeService) {
	mockRepo := new(mockBadgeRepository)
	badgeService := services.NewBadgeService(mockRepo)
	return mockRepo, badgeService
}

func TestCreateBadge(t *testing.T) {
	mockRepo, badgeService := setupBadgeServiceTest()
	ctx := context.Background()

	// Test case: successful badge creation
	t.Run("Successful badge creation", func(t *testing.T) {
		badgeID := uuid.New()
		now := time.Now()

		badgeCreate := &models.BadgeCreate{
			Name:        "Test Badge",
			Description: "Badge for testing",
			IconURL:     "https://example.com/badge.png",
		}

		createdBadge := &models.Badge{
			ID:          badgeID,
			Name:        badgeCreate.Name,
			Description: badgeCreate.Description,
			IconURL:     badgeCreate.IconURL,
			CreatedAt:   &now,
			UpdatedAt:   &now,
		}

		// Mock repository behavior
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Badge")).Return(createdBadge, nil)

		// Call the service
		result, err := badgeService.CreateBadge(ctx, badgeCreate)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, badgeID, result.ID)
		assert.Equal(t, badgeCreate.Name, result.Name)
		assert.Equal(t, badgeCreate.Description, result.Description)
		assert.Equal(t, badgeCreate.IconURL, result.IconURL)

		mockRepo.AssertExpectations(t)
	})

	// Test case: badge with duplicate name
	t.Run("Badge with duplicate name", func(t *testing.T) {
		badgeName := "Existing Badge"
		badgeCreate := &models.BadgeCreate{
			Name:        badgeName,
			Description: "Badge for testing",
			IconURL:     "https://example.com/badge.png",
		}

		// Mock repository behavior to simulate duplicate name
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Badge")).Return(nil, errors.New(errors.ErrConflict, "Badge with this name already exists", nil))

		// Call the service
		result, err := badgeService.CreateBadge(ctx, badgeCreate)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.ErrConflict, err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestGetBadgeByID(t *testing.T) {
	mockRepo, badgeService := setupBadgeServiceTest()
	ctx := context.Background()

	// Test case: successful badge retrieval
	t.Run("Successful badge retrieval", func(t *testing.T) {
		badgeID := uuid.New()
		now := time.Now()

		badge := &models.Badge{
			ID:          badgeID,
			Name:        "Test Badge",
			Description: "Badge for testing",
			IconURL:     "https://example.com/badge.png",
			CreatedAt:   &now,
			UpdatedAt:   &now,
		}

		// Mock repository behavior
		mockRepo.On("FindByID", ctx, badgeID.String()).Return(badge, nil)

		// Call the service
		result, err := badgeService.GetBadgeByID(ctx, badgeID.String())

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, badgeID, result.ID)
		assert.Equal(t, badge.Name, result.Name)
		assert.Equal(t, badge.Description, result.Description)
		assert.Equal(t, badge.IconURL, result.IconURL)

		mockRepo.AssertExpectations(t)
	})

	// Test case: badge not found
	t.Run("Badge not found", func(t *testing.T) {
		badgeID := uuid.New()

		// Mock repository behavior
		mockRepo.On("FindByID", ctx, badgeID.String()).Return(nil, errors.New(errors.ErrNotFound, "Badge not found", nil))

		// Call the service
		result, err := badgeService.GetBadgeByID(ctx, badgeID.String())

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.ErrNotFound, err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestListBadges(t *testing.T) {
	mockRepo, badgeService := setupBadgeServiceTest()
	ctx := context.Background()

	// Test case: successful badge listing
	t.Run("Successful badge listing", func(t *testing.T) {
		badgeID1 := uuid.New()
		badgeID2 := uuid.New()
		now := time.Now()

		badges := []models.Badge{
			{
				ID:          badgeID1,
				Name:        "Test Badge 1",
				Description: "Badge 1 for testing",
				IconURL:     "https://example.com/badge1.png",
				CreatedAt:   &now,
				UpdatedAt:   &now,
			},
			{
				ID:          badgeID2,
				Name:        "Test Badge 2",
				Description: "Badge 2 for testing",
				IconURL:     "https://example.com/badge2.png",
				CreatedAt:   &now,
				UpdatedAt:   &now,
			},
		}

		listOpts := base.ListOptions{
			Page:    1,
			PerPage: 10,
		}

		// Mock repository behavior
		mockRepo.On("List", ctx, listOpts).Return(badges, nil)
		mockRepo.On("Count", ctx, mock.AnythingOfType("[]base.FilterOption")).Return(len(badges), nil)

		// Call the service
		results, count, err := badgeService.ListBadges(ctx, listOpts)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Len(t, results, 2)
		assert.Equal(t, badgeID1, results[0].ID)
		assert.Equal(t, badgeID2, results[1].ID)

		mockRepo.AssertExpectations(t)
	})

	// Test case: empty list
	t.Run("Empty badge list", func(t *testing.T) {
		listOpts := base.ListOptions{
			Page:    1,
			PerPage: 10,
		}

		// Mock repository behavior
		mockRepo.On("List", ctx, listOpts).Return([]models.Badge{}, nil)
		mockRepo.On("Count", ctx, mock.AnythingOfType("[]base.FilterOption")).Return(0, nil)

		// Call the service
		results, count, err := badgeService.ListBadges(ctx, listOpts)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
		assert.Empty(t, results)

		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateBadge(t *testing.T) {
	mockRepo, badgeService := setupBadgeServiceTest()
	ctx := context.Background()

	// Test case: successful badge update
	t.Run("Successful badge update", func(t *testing.T) {
		badgeID := uuid.New()
		now := time.Now()

		existingBadge := &models.Badge{
			ID:          badgeID,
			Name:        "Original Badge",
			Description: "Original description",
			IconURL:     "https://example.com/original.png",
			CreatedAt:   &now,
			UpdatedAt:   &now,
		}

		badgeUpdate := &models.BadgeUpdate{
			ID:          badgeID,
			Name:        "Updated Badge",
			Description: "Updated description",
			IconURL:     "https://example.com/updated.png",
		}

		updatedBadge := &models.Badge{
			ID:          badgeID,
			Name:        badgeUpdate.Name,
			Description: badgeUpdate.Description,
			IconURL:     badgeUpdate.IconURL,
			CreatedAt:   existingBadge.CreatedAt,
			UpdatedAt:   &now,
		}

		// Mock repository behavior
		mockRepo.On("FindByID", ctx, badgeID.String()).Return(existingBadge, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*models.Badge")).Return(updatedBadge, nil)

		// Call the service
		result, err := badgeService.UpdateBadge(ctx, badgeID.String(), badgeUpdate)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, badgeID, result.ID)
		assert.Equal(t, badgeUpdate.Name, result.Name)
		assert.Equal(t, badgeUpdate.Description, result.Description)
		assert.Equal(t, badgeUpdate.IconURL, result.IconURL)

		mockRepo.AssertExpectations(t)
	})

	// Test case: badge not found
	t.Run("Badge not found for update", func(t *testing.T) {
		badgeID := uuid.New()

		badgeUpdate := &models.BadgeUpdate{
			ID:          badgeID,
			Name:        "Updated Badge",
			Description: "Updated description",
			IconURL:     "https://example.com/updated.png",
		}

		// Mock repository behavior
		mockRepo.On("FindByID", ctx, badgeID.String()).Return(nil, errors.New(errors.ErrNotFound, "Badge not found", nil))

		// Call the service
		result, err := badgeService.UpdateBadge(ctx, badgeID.String(), badgeUpdate)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.ErrNotFound, err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteBadge(t *testing.T) {
	mockRepo, badgeService := setupBadgeServiceTest()
	ctx := context.Background()

	// Test case: successful badge deletion
	t.Run("Successful badge deletion", func(t *testing.T) {
		badgeID := uuid.New()

		// Mock repository behavior
		mockRepo.On("Exists", ctx, badgeID.String()).Return(true, nil)
		mockRepo.On("Delete", ctx, badgeID.String()).Return(nil)

		// Call the service
		err := badgeService.DeleteBadge(ctx, badgeID.String())

		// Assertions
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	// Test case: badge not found
	t.Run("Badge not found for deletion", func(t *testing.T) {
		badgeID := uuid.New()

		// Mock repository behavior
		mockRepo.On("Exists", ctx, badgeID.String()).Return(false, nil)

		// Call the service
		err := badgeService.DeleteBadge(ctx, badgeID.String())

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, errors.ErrNotFound, err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestGetPopularBadges(t *testing.T) {
	mockRepo, badgeService := setupBadgeServiceTest()
	ctx := context.Background()

	// Test case: successful retrieval of popular badges
	t.Run("Successful retrieval of popular badges", func(t *testing.T) {
		badgeID1 := uuid.New()
		badgeID2 := uuid.New()
		now := time.Now()
		limit := 5

		badges := []models.Badge{
			{
				ID:          badgeID1,
				Name:        "Popular Badge 1",
				Description: "Popular badge 1 for testing",
				IconURL:     "https://example.com/badge1.png",
				CreatedAt:   &now,
				UpdatedAt:   &now,
			},
			{
				ID:          badgeID2,
				Name:        "Popular Badge 2",
				Description: "Popular badge 2 for testing",
				IconURL:     "https://example.com/badge2.png",
				CreatedAt:   &now,
				UpdatedAt:   &now,
			},
		}

		// Mock repository behavior
		mockRepo.On("FindPopularBadges", ctx, limit).Return(badges, nil)

		// Call the service
		results, err := badgeService.GetPopularBadges(ctx, limit)

		// Assertions
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, badgeID1, results[0].ID)
		assert.Equal(t, badgeID2, results[1].ID)

		mockRepo.AssertExpectations(t)
	})
}
