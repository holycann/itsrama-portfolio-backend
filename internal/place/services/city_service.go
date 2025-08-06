package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type cityService struct {
	cityRepo repositories.CityRepository
}

func NewCityService(cityRepo repositories.CityRepository) CityService {
	return &cityService{
		cityRepo: cityRepo,
	}
}

func (s *cityService) CreateCity(ctx context.Context, city *models.City) error {
	// Validate city object
	if city == nil {
		return fmt.Errorf("city cannot be nil")
	}

	// Validate required fields
	if city.Name == "" {
		return fmt.Errorf("city name is required")
	}

	// Set default values
	city.ID = uuid.New()
	now := time.Now()
	city.CreatedAt = now
	city.UpdatedAt = now

	// Call repository to create city
	return s.cityRepo.Create(ctx, city)
}

func (s *cityService) GetCityByID(ctx context.Context, id string) (*models.ResponseCity, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("city ID cannot be empty")
	}

	// Retrieve city from repository
	return s.cityRepo.FindByID(ctx, id)
}

func (s *cityService) ListCities(ctx context.Context, opts repository.ListOptions) ([]models.ResponseCity, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.cityRepo.List(ctx, opts)
}

func (s *cityService) UpdateCity(ctx context.Context, city *models.City) error {
	// Validate city object
	if city == nil {
		return fmt.Errorf("city cannot be nil")
	}

	// Validate required fields
	if city.ID == uuid.Nil {
		return fmt.Errorf("city ID is required for update")
	}

	// Update timestamp
	city.UpdatedAt = time.Now()

	// Call repository to update city
	return s.cityRepo.Update(ctx, city)
}

func (s *cityService) DeleteCity(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("city ID cannot be empty")
	}

	// Call repository to delete city
	return s.cityRepo.Delete(ctx, id)
}

func (s *cityService) CountCities(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.cityRepo.Count(ctx, filters)
}

func (s *cityService) GetCityByName(ctx context.Context, name string) (*models.ResponseCity, error) {
	return s.cityRepo.FindCityByName(ctx, name)
}

func (s *cityService) GetCitiesByProvince(ctx context.Context, provinceID string) ([]models.ResponseCity, error) {
	// Convert string to UUID
	provUUID, err := uuid.Parse(provinceID)
	if err != nil {
		return nil, fmt.Errorf("invalid province ID: %w", err)
	}

	return s.cityRepo.FindCitiesByProvince(ctx, provUUID.String())
}

func (s *cityService) SearchCities(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseCity, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	// Add search query to filters
	opts.Filters = append(opts.Filters,
		repository.FilterOption{
			Field:    "name",
			Operator: "like",
			Value:    query,
		},
	)

	return s.cityRepo.List(ctx, opts)
}
