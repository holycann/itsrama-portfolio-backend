// Package repositories provides an implementation of repository for location data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/repository"
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

func (r *locationRepository) Create(ctx context.Context, location *models.Location) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(location, false, "", "minimal", "").
		ExecuteTo(&location)
	return err
}

func (r *locationRepository) FindByID(ctx context.Context, id string) (*models.Location, error) {
	var location *models.Location
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&location)
	return location, err
}

func (r *locationRepository) Update(ctx context.Context, location *models.Location) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(location, "minimal", "").
		Eq("id", location.ID.String()).
		Execute()
	return err
}

func (r *locationRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *locationRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.Location, error) {
	var locations []models.Location
	query := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case "=":
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case "like":
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == repository.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	query = query.Range(opts.Offset, opts.Offset+opts.Limit-1, "")

	_, err := query.ExecuteTo(&locations)
	return locations, err
}

func (r *locationRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
	query := r.supabaseClient.
		From(r.table).
		Select("id", "exact", false)

	// Apply filters
	for _, filter := range filters {
		switch filter.Operator {
		case "=":
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case "like":
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

func (r *locationRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.Location, error) {
	var locations []models.Location
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&locations)
	return locations, err
}

// Specialized methods for locations
func (r *locationRepository) FindLocationsByCity(ctx context.Context, cityID string) ([]models.Location, error) {
	var locations []models.Location
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("city_id", cityID).
		ExecuteTo(&locations)
	return locations, err
}

func (r *locationRepository) FindLocationsByProximity(ctx context.Context, latitude, longitude float64, radius float64) ([]models.Location, error) {
	// Note: This is a placeholder. Actual implementation would depend on Supabase's geospatial query capabilities
	var locations []models.Location
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		// Add geospatial filtering logic here
		ExecuteTo(&locations)
	return locations, err
}
