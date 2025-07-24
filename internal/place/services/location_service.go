package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/repositories"
)

type locationService struct {
	locationRepo repositories.LocationRepository
}

// NewLocationService creates a new instance of the location service
// with the given location repository.
func NewLocationService(locationRepo repositories.LocationRepository) LocationService {
	return &locationService{
		locationRepo: locationRepo,
	}
}

// CreateLocation adds a new location to the database
// Validates the location object before creating
func (s *locationService) CreateLocation(ctx context.Context, location *models.Location) error {
	// Validate location object
	if location == nil {
		return fmt.Errorf("location cannot be nil")
	}

	// Validate required fields (example validation)
	if location.Name == "" {
		return fmt.Errorf("location name is required")
	}

	location.ID = uuid.NewString()

	// Call repository to create location
	return s.locationRepo.Create(ctx, location)
}

// GetLocations retrieves a list of locations with pagination
func (s *locationService) GetLocations(ctx context.Context, limit, offset int) ([]*models.Location, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve locations from repository
	locations, err := s.locationRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert []models.Location to []*models.Location
	locationPtrs := make([]*models.Location, len(locations))
	for i := range locations {
		locationPtrs[i] = &locations[i]
	}

	return locationPtrs, nil
}

// GetLocationByID retrieves a single location by its unique identifier
func (s *locationService) GetLocationByID(ctx context.Context, id string) (*models.Location, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("location ID cannot be empty")
	}

	// Retrieve location from repository
	return s.locationRepo.FindByID(ctx, id)
}

// GetLocationByName retrieves a location by its name
// Note: This method is not directly supported by the current repository implementation
// You might need to add a custom method in the repository or implement filtering
func (s *locationService) GetLocationByName(ctx context.Context, name string) (*models.Location, error) {
	// Validate name
	if name == "" {
		return nil, fmt.Errorf("location name cannot be empty")
	}

	// Since the current repository doesn't have a direct method for this,
	// we'll use a workaround by listing all locations and finding by name
	locations, err := s.locationRepo.List(ctx, 1, 0)
	if err != nil {
		return nil, err
	}

	// Find location by name (linear search)
	for _, location := range locations {
		if location.Name == name {
			return &location, nil
		}
	}

	return nil, fmt.Errorf("location with name %s not found", name)
}

// UpdateLocation updates an existing location in the database
func (s *locationService) UpdateLocation(ctx context.Context, location *models.Location) error {
	// Validate location object
	if location == nil {
		return fmt.Errorf("location cannot be nil")
	}

	// Validate required fields
	if location.ID == "" {
		return fmt.Errorf("location ID is required for update")
	}

	// Call repository to update location
	return s.locationRepo.Update(ctx, location)
}

// DeleteLocation removes a location from the database by its ID
func (s *locationService) DeleteLocation(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("location ID cannot be empty")
	}

	// Call repository to delete location
	return s.locationRepo.Delete(ctx, id)
}

// Count menghitung jumlah total lokasi yang tersimpan
func (s *locationService) Count(ctx context.Context) (int, error) {
	return s.locationRepo.Count(ctx)
}
