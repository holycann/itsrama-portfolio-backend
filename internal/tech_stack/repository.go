package tech_stack

import (
	"context"
	"fmt"

	"github.com/holycann/itsrama-portfolio-backend/internal/base"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/supabase"
	postgrest "github.com/supabase-community/postgrest-go"
)

type TechStackRepository interface {
	base.BaseRepository[TechStack, TechStack]
}

type techStackRepository struct {
	supabaseClient *supabase.SupabaseClient
	table          string
}

func NewTechStackRepository(supabaseClient *supabase.SupabaseClient) TechStackRepository {
	return &techStackRepository{
		supabaseClient: supabaseClient,
		table:          "tech_stack",
	}
}

func (r *techStackRepository) Create(ctx context.Context, techStack *TechStack) (*TechStack, error) {
	_, _, err := r.supabaseClient.GetClient().
		From(r.table).
		Insert(techStack, false, "", "minimal", "").
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create tech stack")
	}
	return techStack, nil
}

func (r *techStackRepository) Update(ctx context.Context, techStack *TechStack) (*TechStack, error) {
	_, _, err := r.supabaseClient.GetClient().
		From(r.table).
		Update(techStack, "minimal", "").
		Eq("id", techStack.ID.String()).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update tech stack")
	}
	return techStack, nil
}

func (r *techStackRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.GetClient().
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to delete tech stack")
	}
	return nil
}

func (r *techStackRepository) List(ctx context.Context, opts base.ListOptions) ([]TechStack, error) {
	var techStacks []TechStack
	query := r.supabaseClient.GetClient().
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

	// Apply search if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("name.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("category.ilike.%%%s%%", opts.Search),
		)
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == base.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	offset := (opts.Page - 1) * opts.PerPage
	query = query.Range(offset, offset+opts.PerPage-1, "")

	_, err := query.ExecuteTo(&techStacks)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to list tech stacks")
	}

	return techStacks, nil
}

func (r *techStackRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	query := r.supabaseClient.GetClient().
		From(r.table).
		Select("id", "exact", true)

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
		return 0, errors.Wrap(err, errors.ErrDatabase, "failed to count tech stacks")
	}

	return int(count), nil
}

func (r *techStackRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, count, err := r.supabaseClient.GetClient().
		From(r.table).
		Select("id", "exact", true).
		Eq("id", id).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, errors.Wrap(err, errors.ErrDatabase, "failed to check tech stack existence")
	}

	return count > 0, nil
}

func (r *techStackRepository) FindByField(ctx context.Context, field string, value interface{}) ([]TechStack, error) {
	var techStacks []TechStack
	_, err := r.supabaseClient.GetClient().
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&techStacks)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find tech stacks by field")
	}
	return techStacks, nil
}

func (r *techStackRepository) Search(ctx context.Context, opts base.ListOptions) ([]TechStack, int, error) {
	var techStacks []TechStack
	query := r.supabaseClient.GetClient().
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

	// Apply search if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("name.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("category.ilike.%%%s%%", opts.Search),
		)
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == base.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	offset := (opts.Page - 1) * opts.PerPage
	query = query.Range(offset, offset+opts.PerPage-1, "")

	_, err := query.ExecuteTo(&techStacks)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to search tech stacks")
	}

	// Count total results
	count, err := r.Count(ctx, opts.Filters)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to count tech stacks")
	}

	return techStacks, count, nil
}
