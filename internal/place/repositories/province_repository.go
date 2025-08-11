package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/pkg/base"
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

func (r *provinceRepository) Create(ctx context.Context, province *models.Province) (*models.Province, error) {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(province, false, "", "minimal", "").
		ExecuteTo(&province)
	return province, err
}

func (r *provinceRepository) FindByID(ctx context.Context, id string) (*models.ProvinceDTO, error) {
	var province *models.ProvinceDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&province)
	return province, err
}

func (r *provinceRepository) Update(ctx context.Context, province *models.Province) (*models.Province, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(province, "minimal", "").
		Eq("id", province.ID.String()).
		Execute()
	return province, err
}

func (r *provinceRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *provinceRepository) List(ctx context.Context, opts base.ListOptions) ([]models.ProvinceDTO, error) {
	var provinces []models.ProvinceDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

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

	_, err := query.ExecuteTo(&provinces)
	return provinces, err
}

func (r *provinceRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
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

func (r *provinceRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *provinceRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.ProvinceDTO, error) {
	var provinces []models.ProvinceDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&provinces)
	return provinces, err
}

func (r *provinceRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.ProvinceDTO, int, error) {
	var provinces []models.ProvinceDTO

	searchQuery := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

	// Apply search
	if opts.Search != "" {
		searchQuery = searchQuery.Or(
			fmt.Sprintf("name.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("description.ilike.%%%s%%", opts.Search),
		)
	}

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case base.OperatorEqual:
			searchQuery = searchQuery.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			searchQuery = searchQuery.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == base.SortAscending
		searchQuery = searchQuery.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	limit, offset := opts.LimitOffset()
	searchQuery = searchQuery.Range(offset, offset+limit-1, "")

	_, err := searchQuery.ExecuteTo(&provinces)

	// Get total count for pagination
	totalCount, countErr := r.Count(ctx, opts.Filters)
	if countErr != nil {
		return nil, 0, countErr
	}

	return provinces, totalCount, err
}

func (r *provinceRepository) BulkCreate(ctx context.Context, provinces []*models.Province) ([]models.Province, error) {
	var createdProvinces []models.Province
	_, err := r.supabaseClient.
		From(r.table).
		Insert(provinces, false, "", "minimal", "").
		ExecuteTo(&createdProvinces)
	return createdProvinces, err
}

func (r *provinceRepository) BulkUpdate(ctx context.Context, provinces []*models.Province) ([]models.Province, error) {
	var updatedProvinces []models.Province
	for _, province := range provinces {
		_, _, err := r.supabaseClient.
			From(r.table).
			Update(province, "minimal", "").
			Eq("id", province.ID.String()).
			Execute()
		if err != nil {
			return nil, err
		}
		updatedProvinces = append(updatedProvinces, *province)
	}
	return updatedProvinces, nil
}

func (r *provinceRepository) BulkDelete(ctx context.Context, ids []string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		In("id", ids).
		Execute()
	return err
}

func (r *provinceRepository) FindProvinceByName(ctx context.Context, name string) (*models.ProvinceDTO, error) {
	var provinces []models.ProvinceDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("name", name).
		ExecuteTo(&provinces)

	if err != nil {
		return nil, err
	}

	if len(provinces) == 0 {
		return nil, fmt.Errorf("province not found")
	}

	return &provinces[0], nil
}

func (r *provinceRepository) FindProvinceByCode(ctx context.Context, code string) (*models.ProvinceDTO, error) {
	var provinces []models.ProvinceDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("code", code).
		ExecuteTo(&provinces)

	if err != nil {
		return nil, err
	}

	if len(provinces) == 0 {
		return nil, fmt.Errorf("province not found")
	}

	return &provinces[0], nil
}

func (r *provinceRepository) ListProvincesByRegion(ctx context.Context, region string) ([]models.ProvinceDTO, error) {
	var provinces []models.ProvinceDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("region", region).
		ExecuteTo(&provinces)

	return provinces, err
}
