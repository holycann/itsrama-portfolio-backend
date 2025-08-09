package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/place/models"
	placeServices "github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProvinceRepository is a mock implementation of ProvinceRepository
type MockProvinceRepository struct {
	mock.Mock
}

func (m *MockProvinceRepository) Create(ctx context.Context, province *models.Province) (*models.Province, error) {
	args := m.Called(ctx, province)
	return args.Get(0).(*models.Province), args.Error(1)
}

func (m *MockProvinceRepository) FindByID(ctx context.Context, id string) (*models.ProvinceDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.ProvinceDTO), args.Error(1)
}

func (m *MockProvinceRepository) FindByName(ctx context.Context, name string) (*models.ProvinceDTO, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.ProvinceDTO), args.Error(1)
}

func (m *MockProvinceRepository) FindByCode(ctx context.Context, code string) (*models.ProvinceDTO, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*models.ProvinceDTO), args.Error(1)
}

func (m *MockProvinceRepository) FindProvinceByCode(ctx context.Context, code string) (*models.ProvinceDTO, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*models.ProvinceDTO), args.Error(1)
}

func (m *MockProvinceRepository) FindProvinceByName(ctx context.Context, name string) (*models.ProvinceDTO, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*models.ProvinceDTO), args.Error(1)
}

func (m *MockProvinceRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.ProvinceDTO, error) {
	args := m.Called(ctx, field, value)
	return args.Get(0).([]models.ProvinceDTO), args.Error(1)
}

func (m *MockProvinceRepository) List(ctx context.Context, opts base.ListOptions) ([]models.ProvinceDTO, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.ProvinceDTO), args.Error(1)
}

func (m *MockProvinceRepository) ListByRegion(ctx context.Context, region string) ([]models.ProvinceDTO, error) {
	args := m.Called(ctx, region)
	return args.Get(0).([]models.ProvinceDTO), args.Error(1)
}

func (m *MockProvinceRepository) ListProvincesByRegion(ctx context.Context, region string) ([]models.ProvinceDTO, error) {
	args := m.Called(ctx, region)
	return args.Get(0).([]models.ProvinceDTO), args.Error(1)
}

func (m *MockProvinceRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

func (m *MockProvinceRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.ProvinceDTO, int, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]models.ProvinceDTO), args.Int(1), args.Error(2)
}

func (m *MockProvinceRepository) Update(ctx context.Context, province *models.Province) (*models.Province, error) {
	args := m.Called(ctx, province)
	return args.Get(0).(*models.Province), args.Error(1)
}

func (m *MockProvinceRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProvinceRepository) BulkCreate(ctx context.Context, provinces []*models.Province) ([]models.Province, error) {
	args := m.Called(ctx, provinces)
	return args.Get(0).([]models.Province), args.Error(1)
}

func (m *MockProvinceRepository) BulkDelete(ctx context.Context, ids []string) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *MockProvinceRepository) BulkUpdate(ctx context.Context, provinces []*models.Province) ([]models.Province, error) {
	args := m.Called(ctx, provinces)
	return args.Get(0).([]models.Province), args.Error(1)
}

func (m *MockProvinceRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func TestCreateProvince(t *testing.T) {
	mockRepo := new(MockProvinceRepository)
	provinceService := placeServices.NewProvinceService(mockRepo)

	provinceToCreate := &models.Province{
		Name:        "Test Province",
		Description: "Test Description",
	}

	mockRepo.On("Create", mock.Anything, provinceToCreate).Return(provinceToCreate, nil)

	result, err := provinceService.CreateProvince(context.Background(), provinceToCreate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, provinceToCreate.Name, result.Name)
	assert.Equal(t, provinceToCreate.Description, result.Description)

	mockRepo.AssertExpectations(t)
}

func TestGetProvinceByID(t *testing.T) {
	mockRepo := new(MockProvinceRepository)
	provinceService := placeServices.NewProvinceService(mockRepo)

	provinceID := uuid.New()
	expectedProvince := &models.ProvinceDTO{
		ID:          provinceID,
		Name:        "Test Province",
		Description: "Test Description",
	}

	mockRepo.On("FindByID", mock.Anything, provinceID.String()).Return(expectedProvince, nil)

	result, err := provinceService.GetProvinceByID(context.Background(), provinceID.String())

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProvince.Name, result.Name)
	assert.Equal(t, expectedProvince.Description, result.Description)

	mockRepo.AssertExpectations(t)
}

func TestListProvinces(t *testing.T) {
	mockRepo := new(MockProvinceRepository)
	provinceService := placeServices.NewProvinceService(mockRepo)

	listOptions := base.ListOptions{
		Page:    1,
		PerPage: 10,
		SortBy:  "name",
	}

	expectedProvinces := []models.ProvinceDTO{
		{
			Name:        "Province 1",
			Description: "Description 1",
		},
		{
			Name:        "Province 2",
			Description: "Description 2",
		},
	}

	mockRepo.On("List", mock.Anything, listOptions).Return(expectedProvinces, nil)

	result, err := provinceService.ListProvinces(context.Background(), listOptions)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedProvinces[0].Name, result[0].Name)
	assert.Equal(t, expectedProvinces[1].Name, result[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestSearchProvinces(t *testing.T) {
	mockRepo := new(MockProvinceRepository)
	provinceService := placeServices.NewProvinceService(mockRepo)

	searchOptions := base.ListOptions{
		Page:    1,
		PerPage: 10,
		SortBy:  "name",
	}

	expectedProvinces := []models.ProvinceDTO{
		{
			Name:        "Province 1",
			Description: "Description 1",
		},
		{
			Name:        "Province 2",
			Description: "Description 2",
		},
	}

	mockRepo.On("Search", mock.Anything, searchOptions).Return(expectedProvinces, len(expectedProvinces), nil)

	result, total, err := provinceService.SearchProvinces(context.Background(), searchOptions)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, len(expectedProvinces), total)
	assert.Equal(t, expectedProvinces[0].Name, result[0].Name)
	assert.Equal(t, expectedProvinces[1].Name, result[1].Name)

	mockRepo.AssertExpectations(t)
}
