package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/models"
	userServices "github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserProfileRepository is a mock implementation of UserProfileRepository
type MockUserProfileRepository struct {
	mock.Mock
}

func (m *MockUserProfileRepository) Create(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error) {
	args := m.Called(ctx, profile)
	return args.Get(0).(*models.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) FindByID(ctx context.Context, id string) (*models.UserProfile, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) Update(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error) {
	args := m.Called(ctx, profile)
	return args.Get(0).(*models.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserProfileRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.UserProfile, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.UserProfile), args.Int(1), args.Error(2)
}

func (m *MockUserProfileRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockUserProfileRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserProfileRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserProfileRepository) FindByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) UpdateAvatarImage(ctx context.Context, payload *models.UserProfileAvatarUpdate) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func (m *MockUserProfileRepository) VerifyIdentity(ctx context.Context, payload *models.UserProfileIdentityUpdate) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// MockUserBadgeRepository is a mock implementation of UserBadgeRepository
type MockUserBadgeRepository struct {
	mock.Mock
}

func (m *MockUserBadgeRepository) AddBadgeToUser(ctx context.Context, payload *models.UserBadgePayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

// MockBadgeService is a mock implementation of BadgeService
type MockBadgeService struct {
	mock.Mock
}

func (m *MockBadgeService) GetBadgeByName(ctx context.Context, name string) (*models.BadgeDTO, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.BadgeDTO), args.Error(1)
}

// MockSupabaseStorage is a mock implementation of SupabaseStorage
type MockSupabaseStorage struct {
	mock.Mock
}

func (m *MockSupabaseStorage) Upload(ctx context.Context, file interface{}, path string) (string, error) {
	args := m.Called(ctx, file, path)
	return args.String(0), args.Error(1)
}

func (m *MockSupabaseStorage) GetPublicURL(filepath string) (string, error) {
	args := m.Called(filepath)
	return args.String(0), args.Error(1)
}

func TestCreateUserProfile(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockProfileRepo := new(MockUserProfileRepository)
	mockUserBadgeRepo := new(MockUserBadgeRepository)
	mockBadgeService := new(MockBadgeService)
	mockStorage := new(MockSupabaseStorage)

	userProfileService := userServices.NewUserProfileService(
		mockProfileRepo,
		mockUserRepo,
		mockUserBadgeRepo,
		mockBadgeService,
		mockStorage,
	)

	userID := uuid.New()
	testProfile := &models.UserProfileCreate{
		UserID:   userID,
		Fullname: "John Doe",
		Bio:      "Test bio",
	}

	// Mock user existence check
	mockUserRepo.On("Exists", mock.Anything, userID.String()).Return(true, nil)

	// Mock profile existence check
	mockProfileRepo.On("ExistsByUserID", mock.Anything, userID.String()).Return(false, nil)

	// Mock profile creation
	mockProfileRepo.On("Create", mock.Anything, mock.Anything).Return(&models.UserProfile{
		ID:       uuid.New(),
		UserID:   userID,
		Fullname: testProfile.Fullname,
		Bio:      testProfile.Bio,
	}, nil)

	// Mock badge retrieval
	mockBadgeService.On("GetBadgeByName", mock.Anything, "Penjelajah").Return(&models.BadgeDTO{
		ID:   uuid.New(),
		Name: "Penjelajah",
	}, nil)

	// Mock badge assignment
	mockUserBadgeRepo.On("AddBadgeToUser", mock.Anything, mock.Anything).Return(nil)

	// Execute the service method
	result, err := userProfileService.CreateProfile(context.Background(), testProfile)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testProfile.Fullname, result.Fullname)
	assert.Equal(t, testProfile.Bio, result.Bio)

	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockProfileRepo.AssertExpectations(t)
	mockUserBadgeRepo.AssertExpectations(t)
	mockBadgeService.AssertExpectations(t)
}

func TestCreateUserProfileWithExistingProfile(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockProfileRepo := new(MockUserProfileRepository)
	mockUserBadgeRepo := new(MockUserBadgeRepository)
	mockBadgeService := new(MockBadgeService)
	mockStorage := new(MockSupabaseStorage)

	userProfileService := userServices.NewUserProfileService(
		mockProfileRepo,
		mockUserRepo,
		mockUserBadgeRepo,
		mockBadgeService,
		mockStorage,
	)

	userID := uuid.New()
	testProfile := &models.UserProfileCreate{
		UserID:   userID,
		Fullname: "John Doe",
		Bio:      "Test bio",
	}

	// Mock user existence check
	mockUserRepo.On("Exists", mock.Anything, userID.String()).Return(true, nil)

	// Mock profile existence check (profile already exists)
	mockProfileRepo.On("ExistsByUserID", mock.Anything, userID.String()).Return(true, nil)

	// Execute the service method
	result, err := userProfileService.CreateProfile(context.Background(), testProfile)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")

	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockProfileRepo.AssertExpectations(t)
}

func TestGetUserProfileByID(t *testing.T) {
	mockProfileRepo := new(MockUserProfileRepository)

	userProfileService := userServices.NewUserProfileService(
		mockProfileRepo,
		nil,
		nil,
		nil,
		nil,
	)

	profileID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()

	expectedProfile := &models.UserProfile{
		ID:        profileID,
		UserID:    userID,
		Fullname:  "John Doe",
		Bio:       "Test bio",
		CreatedAt: &now,
	}

	mockProfileRepo.On("FindByID", mock.Anything, profileID.String()).Return(expectedProfile, nil)

	result, err := userProfileService.GetProfileByID(context.Background(), profileID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProfile.Fullname, result.Fullname)
	assert.Equal(t, expectedProfile.Bio, result.Bio)

	mockProfileRepo.AssertExpectations(t)
}

func TestUpdateUserProfile(t *testing.T) {
	mockProfileRepo := new(MockUserProfileRepository)

	userProfileService := userServices.NewUserProfileService(
		mockProfileRepo,
		nil,
		nil,
		nil,
		nil,
	)

	profileID := uuid.New()
	updateProfile := &models.UserProfileUpdate{
		ID:       profileID,
		Fullname: "Updated Name",
		Bio:      "Updated bio",
	}

	// Mock profile existence check
	mockProfileRepo.On("Exists", mock.Anything, profileID.String()).Return(true, nil)

	// Mock profile update
	mockProfileRepo.On("Update", mock.Anything, mock.Anything).Return(&models.UserProfile{
		ID:       profileID,
		Fullname: updateProfile.Fullname,
		Bio:      updateProfile.Bio,
	}, nil)

	result, err := userProfileService.UpdateProfile(context.Background(), updateProfile)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updateProfile.Fullname, result.Fullname)
	assert.Equal(t, updateProfile.Bio, result.Bio)

	mockProfileRepo.AssertExpectations(t)
}

func TestDeleteUserProfile(t *testing.T) {
	mockProfileRepo := new(MockUserProfileRepository)

	userProfileService := userServices.NewUserProfileService(
		mockProfileRepo,
		nil,
		nil,
		nil,
		nil,
	)

	profileID := uuid.New()

	// Mock profile existence check
	mockProfileRepo.On("FindByID", mock.Anything, profileID.String()).Return(&models.UserProfile{
		ID: profileID,
	}, nil)

	// Mock profile deletion
	mockProfileRepo.On("Delete", mock.Anything, profileID.String()).Return(nil)

	err := userProfileService.DeleteProfile(context.Background(), profileID.String())

	assert.NoError(t, err)

	mockProfileRepo.AssertExpectations(t)
}
