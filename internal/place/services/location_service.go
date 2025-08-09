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

type locationService struct {
	locationRepo repositories.LocationRepository
}

func NewLocationService(locationRepo repositories.LocationRepository) LocationService {
	return &locationService{
		locationRepo: locationRepo,
	}
}

// CreateLocation creates a new location
func (s *locationService) CreateLocation(ctx context.Context, location *models.LocationCreate) (*models.Location, error) {
	// Validate location object
	if location == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Location cannot be nil",
			nil,
			errors.WithContext("input", "nil location"),
		)
	}

	// Validate required fields
	if err := base.ValidateModel(location); err != nil {
		return nil, err
	}

	// Create new Location model from LocationCreate
	newLocation := &models.Location{
		ID:        uuid.New(),
		Name:      location.Name,
		CityID:    location.CityID,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}

	// Set timestamps
	now := time.Now()
	newLocation.CreatedAt = &now
	newLocation.UpdatedAt = &now

	// Call repository to create location
	return s.locationRepo.Create(ctx, newLocation)
}

func (s *locationService) GetLocationByID(ctx context.Context, id string) (*models.LocationDTO, error) {
	// Validate ID
	if id == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Location ID cannot be empty",
			nil,
			errors.WithContext("input", "empty ID"),
		)
	}

	// Retrieve location from repository
	return s.locationRepo.FindByID(ctx, id)
}

func (s *locationService) ListLocations(ctx context.Context, opts base.ListOptions) ([]models.LocationDTO, error) {
	return s.locationRepo.List(ctx, opts)
}

// UpdateLocation updates an existing location
func (s *locationService) UpdateLocation(ctx context.Context, location *models.LocationUpdate) (*models.Location, error) {
	// Validate location object
	if location == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Location cannot be nil",
			nil,
			errors.WithContext("input", "nil location"),
		)
	}

	// Validate required fields
	if err := base.ValidateModel(location); err != nil {
		return nil, err
	}

	// Retrieve existing location
	existingLocation, err := s.locationRepo.FindByID(ctx, location.ID.String())
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrNotFound, "Location not found")
	}

	// Update location fields
	updatedLocation := &models.Location{
		ID:        location.ID,
		Name:      location.Name,
		CityID:    location.CityID,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
		CreatedAt: existingLocation.CreatedAt,
	}

	// Update timestamp
	now := time.Now()
	updatedLocation.UpdatedAt = &now

	// Call repository to update location
	return s.locationRepo.Update(ctx, updatedLocation)
}

func (s *locationService) DeleteLocation(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"Location ID cannot be empty",
			nil,
			errors.WithContext("input", "empty ID"),
		)
	}

	// Call repository to delete location
	return s.locationRepo.Delete(ctx, id)
}

func (s *locationService) CountLocations(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.locationRepo.Count(ctx, filters)
}

func (s *locationService) GetLocationByName(ctx context.Context, name string) (*models.LocationDTO, error) {
	// Validate name
	if name == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Location name cannot be empty",
			nil,
			errors.WithContext("input", "empty name"),
		)
	}

	// Use repository's search method
	locations, err := s.locationRepo.SearchLocationsByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, errors.New(
			errors.ErrNotFound,
			"Location not found",
			nil,
			errors.WithContext("name", name),
		)
	}

	return &locations[0], nil
}

func (s *locationService) GetLocationsByCity(ctx context.Context, cityID string) ([]models.LocationDTO, error) {
	// Convert string to UUID
	cityUUID, err := uuid.Parse(cityID)
	if err != nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Invalid city ID",
			err,
			errors.WithContext("cityID", cityID),
		)
	}

	return s.locationRepo.FindLocationsByCity(ctx, cityUUID.String())
}

func (s *locationService) SearchLocations(ctx context.Context, opts base.ListOptions) ([]models.LocationDTO, error) {

	// Perform search using repository
	locations, err := s.locationRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return empty result if no locations found
	if len(locations) == 0 {
		return nil, errors.New(
			errors.ErrNotFound,
			"No locations found matching the search query",
			nil,
		)
	}

	return locations, nil
}

func (s *locationService) GetLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.LocationDTO, error) {
	// Validate input coordinates
	if latitude < -90 || latitude > 90 {
		return nil, errors.New(
			errors.ErrValidation,
			"Invalid latitude",
			nil,
			errors.WithContext("latitude", latitude),
		)
	}

	if longitude < -180 || longitude > 180 {
		return nil, errors.New(
			errors.ErrValidation,
			"Invalid longitude",
			nil,
			errors.WithContext("longitude", longitude),
		)
	}

	// Radius must be positive and is specified in meters
	if radius <= 0 {
		return nil, errors.New(
			errors.ErrValidation,
			"Radius must be positive",
			nil,
			errors.WithContext("radius", radius),
		)
	}

	// Call repository method to find locations within proximity (using PostGIS ST_DWithin)
	return s.locationRepo.FindLocationsByProximity(ctx, latitude, longitude, radius)
}
