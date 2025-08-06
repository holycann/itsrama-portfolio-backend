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

type locationService struct {
	locationRepo repositories.LocationRepository
}

func NewLocationService(locationRepo repositories.LocationRepository) LocationService {
	return &locationService{
		locationRepo: locationRepo,
	}
}

func (s *locationService) CreateLocation(ctx context.Context, location *models.Location) error {
	// Validate location object
	if location == nil {
		return fmt.Errorf("location cannot be nil")
	}

	// Validate required fields
	if location.Name == "" {
		return fmt.Errorf("location name is required")
	}

	// Set default values
	location.ID = uuid.New()
	now := time.Now()
	location.CreatedAt = now
	location.UpdatedAt = now

	// Call repository to create location
	return s.locationRepo.Create(ctx, location)
}

func (s *locationService) GetLocationByID(ctx context.Context, id string) (*models.Location, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("location ID cannot be empty")
	}

	// Retrieve location from repository
	return s.locationRepo.FindByID(ctx, id)
}

func (s *locationService) ListLocations(ctx context.Context, opts repository.ListOptions) ([]models.Location, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.locationRepo.List(ctx, opts)
}

func (s *locationService) UpdateLocation(ctx context.Context, location *models.Location) error {
	// Validate location object
	if location == nil {
		return fmt.Errorf("location cannot be nil")
	}

	// Validate required fields
	if location.ID == uuid.Nil {
		return fmt.Errorf("location ID is required for update")
	}

	// Update timestamp
	location.UpdatedAt = time.Now()

	// Call repository to update location
	return s.locationRepo.Update(ctx, location)
}

func (s *locationService) DeleteLocation(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("location ID cannot be empty")
	}

	// Call repository to delete location
	return s.locationRepo.Delete(ctx, id)
}

func (s *locationService) CountLocations(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.locationRepo.Count(ctx, filters)
}

func (s *locationService) GetLocationByName(ctx context.Context, name string) (*models.Location, error) {
	// Validate name
	if name == "" {
		return nil, fmt.Errorf("location name cannot be empty")
	}

	// Use repository's search method
	locations, err := s.locationRepo.FindByField(ctx, "name", name)
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, fmt.Errorf("location with name %s not found", name)
	}

	return &locations[0], nil
}

func (s *locationService) GetLocationsByCity(ctx context.Context, cityID string) ([]models.Location, error) {
	// Convert string to UUID
	cityUUID, err := uuid.Parse(cityID)
	if err != nil {
		return nil, fmt.Errorf("invalid city ID: %w", err)
	}

	return s.locationRepo.FindLocationsByCity(ctx, cityUUID.String())
}

// func (s *locationService) GetLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.Location, error) {
// 	return s.locationRepo.FindLocationsByProximity(ctx, latitude, longitude, radius)
// }

func (s *locationService) SearchLocations(ctx context.Context, query string, opts repository.ListOptions) ([]models.Location, error) {
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

	return s.locationRepo.List(ctx, opts)
}
