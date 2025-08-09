package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/repositories"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
)

type cityService struct {
	cityRepo repositories.CityRepository
}

func NewCityService(cityRepo repositories.CityRepository) CityService {
	return &cityService{
		cityRepo: cityRepo,
	}
}

func (s *cityService) CreateCity(ctx context.Context, city *models.City) (*models.City, error) {
	// Validate city object
	if city == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"City cannot be nil",
			nil,
			errors.WithContext("input", "nil city"),
		)
	}

	// Validate required fields
	if err := base.ValidateModel(city); err != nil {
		return nil, err
	}

	// Set default values
	city.ID = uuid.New()
	now := time.Now()
	city.CreatedAt = &now
	city.UpdatedAt = &now

	// Call repository to create city
	return s.cityRepo.Create(ctx, city)
}

func (s *cityService) GetCityByID(ctx context.Context, id string) (*models.CityDTO, error) {
	// Validate ID
	if id == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"City ID cannot be empty",
			nil,
			errors.WithContext("input", "empty ID"),
		)
	}

	// Retrieve city from repository
	return s.cityRepo.FindByID(ctx, id)
}

func (s *cityService) ListCities(ctx context.Context, opts base.ListOptions) ([]models.CityDTO, int, error) {
	// Retrieve cities from repository
	cities, total, err := s.cityRepo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "Failed to list cities")
	}

	return cities, total, nil
}

func (s *cityService) UpdateCity(ctx context.Context, city *models.City) (*models.City, error) {
	// Validate city object
	if city == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"City cannot be nil",
			nil,
			errors.WithContext("input", "nil city"),
		)
	}

	// Validate required fields
	if city.ID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"City ID is required for update",
			nil,
			errors.WithContext("input", "missing ID"),
		)
	}

	// Validate model
	if err := base.ValidateModel(city); err != nil {
		return nil, err
	}

	// Update timestamp
	now := time.Now()
	city.UpdatedAt = &now

	return s.cityRepo.Update(ctx, city)
}

func (s *cityService) DeleteCity(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"City ID cannot be empty",
			nil,
			errors.WithContext("input", "empty ID"),
		)
	}

	// Call repository to delete city
	return s.cityRepo.Delete(ctx, id)
}

func (s *cityService) CountCities(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.cityRepo.Count(ctx, filters)
}

func (s *cityService) GetCityByName(ctx context.Context, name string) (*models.CityDTO, error) {
	// Validate name
	if name == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"City name cannot be empty",
			nil,
			errors.WithContext("input", "empty name"),
		)
	}

	return s.cityRepo.FindCityByName(ctx, name)
}

func (s *cityService) GetCityByCode(ctx context.Context, code string) (*models.CityDTO, error) {
	// Validate code
	if code == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"City code cannot be empty",
			nil,
			errors.WithContext("input", "empty code"),
		)
	}

	return s.cityRepo.FindCityByCode(ctx, code)
}

func (s *cityService) GetCitiesByProvince(ctx context.Context, provinceID string) ([]models.CityDTO, error) {
	// Validate province ID
	if provinceID == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Province ID cannot be empty",
			nil,
			errors.WithContext("input", "empty province ID"),
		)
	}

	return s.cityRepo.FindCitiesByProvince(ctx, provinceID)
}

func (s *cityService) ListCitiesByPopulation(ctx context.Context, minPopulation, maxPopulation string) ([]models.CityDTO, error) {
	return s.cityRepo.ListCitiesByPopulation(ctx, minPopulation, maxPopulation)
}

func (s *cityService) SearchCities(ctx context.Context, opts base.ListOptions) ([]models.CityDTO, int, error) {
	return s.cityRepo.Search(ctx, opts)
}
