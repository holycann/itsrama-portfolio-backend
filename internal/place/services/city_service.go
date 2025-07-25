package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/repositories"
)

type cityService struct {
	cityRepo repositories.CityRepository
}

// NewCityService creates a new instance of the city service
// with the given city repository.
func NewCityService(cityRepo repositories.CityRepository) CityService {
	return &cityService{
		cityRepo: cityRepo,
	}
}

// CreateCity adds a new city to the database
// Validates the city object before creating
func (s *cityService) CreateCity(ctx context.Context, city *models.City) error {
	// Validate city object
	if city == nil {
		return fmt.Errorf("city cannot be nil")
	}

	// Validate required fields (example validation)
	if city.Name == "" {
		return fmt.Errorf("city name is required")
	}

	city.ID = uuid.NewString()

	// Call repository to create city
	return s.cityRepo.Create(ctx, city)
}

// GetCitys retrieves a list of cities with pagination
func (s *cityService) GetCities(ctx context.Context, limit, offset int) ([]*models.City, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve cities from repository
	cities, err := s.cityRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert []models.City to []*models.City
	cityPtrs := make([]*models.City, len(cities))
	for i := range cities {
		cityPtrs[i] = &cities[i]
	}

	return cityPtrs, nil
}

// GetCityByID retrieves a single city by its unique identifier
func (s *cityService) GetCityByID(ctx context.Context, id string) (*models.City, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("city ID cannot be empty")
	}

	// Retrieve city from repository
	return s.cityRepo.FindByID(ctx, id)
}

// GetCityByName retrieves a city by its name
// Note: This method is not directly supported by the current repository implementation
// You might need to add a custom method in the repository or implement filtering
func (s *cityService) GetCityByName(ctx context.Context, name string) (*models.City, error) {
	// Validate name
	if name == "" {
		return nil, fmt.Errorf("city name cannot be empty")
	}

	// Since the current repository doesn't have a direct method for this,
	// we'll use a workaround by listing all cities and finding by name
	cities, err := s.cityRepo.List(ctx, 1, 0)
	if err != nil {
		return nil, err
	}

	// Find city by name (linear search)
	for _, city := range cities {
		if city.Name == name {
			return &city, nil
		}
	}

	return nil, fmt.Errorf("city with name %s not found", name)
}

// UpdateCity updates an existing city in the database
func (s *cityService) UpdateCity(ctx context.Context, city *models.City) error {
	// Validate city object
	if city == nil {
		return fmt.Errorf("city cannot be nil")
	}

	// Validate required fields
	if city.ID == "" {
		return fmt.Errorf("city ID is required for update")
	}

	// Call repository to update city
	return s.cityRepo.Update(ctx, city)
}

// DeleteCity removes a city from the database by its ID
func (s *cityService) DeleteCity(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("city ID cannot be empty")
	}

	// Call repository to delete city
	return s.cityRepo.Delete(ctx, id)
}

// Count calculates the total number of stored locations
func (s *cityService) Count(ctx context.Context) (int, error) {
	return s.cityRepo.Count(ctx)
}
