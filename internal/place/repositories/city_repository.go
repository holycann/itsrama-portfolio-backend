// Package repositories provides an implementation of repository for city data management
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

type cityRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewCityRepository(supabaseClient *supabase.Client) CityRepository {
	return &cityRepository{
		supabaseClient: supabaseClient,
		table:          "cities",
	}
}

func (r *cityRepository) Create(ctx context.Context, city *models.City) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(city, false, "", "minimal", "").
		ExecuteTo(&city)
	return err
}

func (r *cityRepository) FindByID(ctx context.Context, id string) (*models.City, error) {
	var city *models.City
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&city)
	return city, err
}

func (r *cityRepository) Update(ctx context.Context, city *models.City) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(city, "minimal", "").
		Eq("id", city.ID.String()).
		Execute()
	return err
}

func (r *cityRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *cityRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.City, error) {
	var cities []models.City
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

	_, err := query.ExecuteTo(&cities)
	return cities, err
}

func (r *cityRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
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

func (r *cityRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *cityRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.City, error) {
	var cities []models.City
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&cities)
	return cities, err
}

// Specialized methods for cities
func (r *cityRepository) FindCitiesByProvince(ctx context.Context, provinceID string) ([]models.City, error) {
	var cities []models.City
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("province_id", provinceID).
		ExecuteTo(&cities)
	return cities, err
}

func (r *cityRepository) FindCityByName(ctx context.Context, name string) (*models.City, error) {
	cities, err := r.FindByField(ctx, "name", name)
	if err != nil {
		return nil, err
	}
	if len(cities) == 0 {
		return nil, fmt.Errorf("city not found")
	}
	return &cities[0], nil
}

func (r *cityRepository) Search(ctx context.Context, opts repository.ListOptions) ([]models.City, int, error) {
	var cities []models.City

	query := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

	// Apply search query if provided
	if opts.SearchQuery != "" {
		query = query.Or(
			fmt.Sprintf("name.ilike.%%%s%%", opts.SearchQuery),
			fmt.Sprintf("province_id.ilike.%%%s%%", opts.SearchQuery),
		)
	}

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case "=":
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case "like":
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Execute query to get results
	_, err := query.ExecuteTo(&cities)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search cities: %w", err)
	}

	// Count total matching records
	_, count, err := query.Execute()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cities: %w", err)
	}

	return cities, int(count), nil
}
