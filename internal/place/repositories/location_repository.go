// Package repositories provides an implementation of repository for location data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type locationRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewLocationRepository(supabaseClient *supabase.Client) LocationRepository {
	return &locationRepository{
		supabaseClient: supabaseClient,
		table:          "locations",
	}
}

func (r *locationRepository) Create(ctx context.Context, location *models.Location) (*models.Location, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(location, false, "", "minimal", "").
		Execute()
	return location, err
}

func (r *locationRepository) FindByID(ctx context.Context, id string) (*models.LocationDTO, error) {
	var location *models.LocationDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, city:cities(*)", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&location)
	return location, err
}

func (r *locationRepository) Update(ctx context.Context, location *models.Location) (*models.Location, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(location, "minimal", "").
		Eq("id", location.ID.String()).
		Execute()
	return location, err
}

func (r *locationRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *locationRepository) List(ctx context.Context, opts base.ListOptions) ([]models.LocationDTO, error) {
	var locations []models.LocationDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*, city:cities(*)", "", false)

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == base.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	limit, offset := opts.LimitOffset()
	query = query.Range(offset, offset+limit-1, "")

	_, err := query.ExecuteTo(&locations)
	return locations, err
}

func (r *locationRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	query := r.supabaseClient.
		From(r.table).
		Select("id", "exact", false)

	// Apply filters
	for _, filter := range filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	_, count, err := query.Execute()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *locationRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *locationRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.LocationDTO, error) {
	var locations []models.LocationDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, city:cities(*)", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&locations)
	return locations, err
}

func (r *locationRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.LocationDTO, int, error) {
	var locations []models.LocationDTO

	query := r.supabaseClient.
		From(r.table).
		Select("*, city:cities(*)", "", false)

	// Apply search query if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("name.ilike.%%%s%%", opts.Search),
			"",
		)
	}

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Execute query to get results
	_, err := query.ExecuteTo(&locations)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search locations: %w", err)
	}

	// Count total matching records
	_, count, err := query.Execute()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count locations: %w", err)
	}

	return locations, int(count), nil
}

func (r *locationRepository) BulkCreate(ctx context.Context, locations []*models.Location) ([]models.Location, error) {
	var createdLocations []models.Location
	for _, location := range locations {
		_, err := r.supabaseClient.
			From(r.table).
			Insert(location, false, "", "minimal", "").
			ExecuteTo(&location)
		if err != nil {
			return nil, err
		}
		createdLocations = append(createdLocations, *location)
	}
	return createdLocations, nil
}

func (r *locationRepository) BulkUpdate(ctx context.Context, locations []*models.Location) ([]models.Location, error) {
	var updatedLocations []models.Location
	for _, location := range locations {
		_, _, err := r.supabaseClient.
			From(r.table).
			Update(location, "minimal", "").
			Eq("id", location.ID.String()).
			Execute()
		if err != nil {
			return nil, err
		}
		updatedLocations = append(updatedLocations, *location)
	}
	return updatedLocations, nil
}

func (r *locationRepository) BulkDelete(ctx context.Context, ids []string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		In("id", ids).
		Execute()
	return err
}

// Specialized methods for locations
func (r *locationRepository) FindLocationsByCity(ctx context.Context, cityID string) ([]models.LocationDTO, error) {
	var locations []models.LocationDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, city:cities(*)", "", false).
		Eq("city_id", cityID).
		ExecuteTo(&locations)
	return locations, err
}

func (r *locationRepository) FindLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.LocationDTO, error) {
	var locations []models.LocationDTO

	// Use the PostGIS ST_DWithin function that works with the location field (geography type)
	// ST_DWithin checks if locations are within the specified radius (in meters)
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, city:cities(*)", "", false).
		Or(fmt.Sprintf("ST_DWithin(location, ST_MakePoint(%f, %f)::geography, %f)", longitude, latitude, radius), "").
		ExecuteTo(&locations)

	return locations, err
}

func (r *locationRepository) SearchLocationsByName(ctx context.Context, queryStr string) ([]models.LocationDTO, error) {
	var locations []models.LocationDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, city:cities(*)", "", false).
		Like("name", fmt.Sprintf("%%%s%%", queryStr)).
		ExecuteTo(&locations)
	return locations, err
}
