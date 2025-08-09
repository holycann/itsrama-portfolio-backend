// Package repositories provides an implementation of repository for city data management
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

func (r *cityRepository) Create(ctx context.Context, city *models.City) (*models.City, error) {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(city, false, "", "minimal", "").
		ExecuteTo(&city)
	return city, err
}

func (r *cityRepository) FindByID(ctx context.Context, id string) (*models.CityDTO, error) {
	var city *models.CityDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, province:provinces(*)", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&city)
	return city, err
}

func (r *cityRepository) Update(ctx context.Context, city *models.City) (*models.City, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(city, "minimal", "").
		Eq("id", city.ID.String()).
		Execute()
	return city, err
}

func (r *cityRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *cityRepository) List(ctx context.Context, opts base.ListOptions) ([]models.CityDTO, error) {
	var cities []models.CityDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*, province:provinces(*)", "", false)

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

	_, err := query.ExecuteTo(&cities)
	return cities, err
}

func (r *cityRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
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

func (r *cityRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *cityRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.CityDTO, error) {
	var cities []models.CityDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, province:provinces(*)", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&cities)
	return cities, err
}

func (r *cityRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.CityDTO, int, error) {
	var cities []models.CityDTO

	query := r.supabaseClient.
		From(r.table).
		Select("*, province:provinces(*)", "", false)

	// Apply search query if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("name.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("province_id.ilike.%%%s%%", opts.Search),
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

func (r *cityRepository) BulkCreate(ctx context.Context, cities []*models.City) ([]models.City, error) {
	var createdCities []models.City
	for _, city := range cities {
		_, err := r.supabaseClient.
			From(r.table).
			Insert(city, false, "", "minimal", "").
			ExecuteTo(&city)
		if err != nil {
			return nil, err
		}
		createdCities = append(createdCities, *city)
	}
	return createdCities, nil
}

func (r *cityRepository) BulkUpdate(ctx context.Context, cities []*models.City) ([]models.City, error) {
	var updatedCities []models.City
	for _, city := range cities {
		_, _, err := r.supabaseClient.
			From(r.table).
			Update(city, "minimal", "").
			Eq("id", city.ID.String()).
			Execute()
		if err != nil {
			return nil, err
		}
		updatedCities = append(updatedCities, *city)
	}
	return updatedCities, nil
}

func (r *cityRepository) BulkDelete(ctx context.Context, ids []string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		In("id", ids).
		Execute()
	return err
}

// Specialized methods for cities
func (r *cityRepository) FindCitiesByProvince(ctx context.Context, provinceID string) ([]models.CityDTO, error) {
	var cities []models.CityDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, province:provinces(*)", "", false).
		Eq("province_id", provinceID).
		ExecuteTo(&cities)
	return cities, err
}

func (r *cityRepository) FindCityByName(ctx context.Context, name string) (*models.CityDTO, error) {
	cities, err := r.FindByField(ctx, "name", name)
	if err != nil {
		return nil, err
	}
	if len(cities) == 0 {
		return nil, fmt.Errorf("city not found")
	}
	return &cities[0], nil
}

func (r *cityRepository) FindCityByCode(ctx context.Context, code string) (*models.CityDTO, error) {
	cities, err := r.FindByField(ctx, "code", code)
	if err != nil {
		return nil, err
	}
	if len(cities) == 0 {
		return nil, fmt.Errorf("city not found")
	}
	return &cities[0], nil
}

func (r *cityRepository) ListCitiesByPopulation(ctx context.Context, minPopulation, maxPopulation string) ([]models.CityDTO, error) {
	var cities []models.CityDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, province:provinces(*)", "", false).
		Gte("population", minPopulation).
		Lte("population", maxPopulation).
		ExecuteTo(&cities)
	return cities, err
}
