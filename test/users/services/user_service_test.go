package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, opts base.ListOptions) ([]models.User, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.User, error) {
	args := m.Called(ctx, field, value)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.User, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	testUser := &models.UserCreate{
		Email:    "test@example.com",
		Password: "password123",
		Role:     "authenticated",
	}

	mockRepo.On("ExistsByEmail", mock.Anything, testUser.Email).Return(false, nil)
	mockRepo.On("Create", mock.Anything, mock.Anything).Return(&models.User{
		ID:    uuid.New(),
		Email: testUser.Email,
		Role:  testUser.Role,
	}, nil)

	result, err := userService.CreateUser(context.Background(), testUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.Email, result.Email)
	assert.Equal(t, testUser.Role, result.Role)

	mockRepo.AssertExpectations(t)
}

func TestCreateUserWithExistingEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	testUser := &models.UserCreate{
		Email:    "existing@example.com",
		Password: "password123",
		Role:     "authenticated",
	}

	mockRepo.On("ExistsByEmail", mock.Anything, testUser.Email).Return(true, nil)

	result, err := userService.CreateUser(context.Background(), testUser)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "already exists")

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	userID := uuid.New()
	expectedUser := &models.User{
		ID:    userID,
		Email: "test@example.com",
		Role:  "authenticated",
	}

	mockRepo.On("FindByID", mock.Anything, userID.String()).Return(expectedUser, nil)

	result, err := userService.GetUserByID(context.Background(), userID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	userID := uuid.New()
	updateUser := &models.UserUpdate{
		ID:    userID,
		Email: "updated@example.com",
	}

	mockRepo.On("Update", mock.Anything, mock.Anything).Return(&models.User{
		ID:    userID,
		Email: updateUser.Email,
	}, nil)

	result, err := userService.UpdateUser(context.Background(), updateUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updateUser.Email, result.Email)

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	userID := uuid.New()

	mockRepo.On("FindByID", mock.Anything, userID.String()).Return(&models.User{ID: userID}, nil)
	mockRepo.On("Delete", mock.Anything, userID.String()).Return(nil)

	err := userService.DeleteUser(context.Background(), userID.String())

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := services.NewUserService(mockRepo)

	listOptions := base.ListOptions{
		Page:    1,
		PerPage: 10,
	}

	expectedUsers := []models.User{
		{ID: uuid.New(), Email: "user1@example.com"},
		{ID: uuid.New(), Email: "user2@example.com"},
	}

	mockRepo.On("Search", mock.Anything, listOptions).Return(expectedUsers, len(expectedUsers), nil)

	results, total, err := userService.ListUsers(context.Background(), listOptions)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, len(expectedUsers), total)

	mockRepo.AssertExpectations(t)
}
