package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/repository"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type provinceRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewProvinceRepository(supabaseClient *supabase.Client) ProvinceRepository {
	return &provinceRepository{
		supabaseClient: supabaseClient,
		table:          "provinces",
	}
}

func (r *provinceRepository) Create(ctx context.Context, province *models.Province) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(province, false, "", "minimal", "").
		ExecuteTo(&province)
	return err
}

func (r *provinceRepository) FindByID(ctx context.Context, id string) (*models.Province, error) {
	var province *models.Province
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&province)
	return province, err
}

func (r *provinceRepository) Update(ctx context.Context, province *models.Province) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(province, "minimal", "").
		Eq("id", province.ID.String()).
		Execute()
	return err
}

func (r *provinceRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *provinceRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.Province, error) {
	var provinces []models.Province
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

	_, err := query.ExecuteTo(&provinces)
	return provinces, err
}

func (r *provinceRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
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

func (r *provinceRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *provinceRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.Province, error) {
	var provinces []models.Province
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&provinces)
	return provinces, err
}

// Specialized methods for provinces
func (r *provinceRepository) FindProvinceByName(ctx context.Context, name string) (*models.Province, error) {
	provinces, err := r.FindByField(ctx, "name", name)
	if err != nil {
		return nil, err
	}
	if len(provinces) == 0 {
		return nil, fmt.Errorf("province not found")
	}
	return &provinces[0], nil
}
