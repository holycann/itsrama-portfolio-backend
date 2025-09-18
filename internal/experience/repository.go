package experience

import (
	"context"
	"fmt"

	"github.com/holycann/itsrama-portfolio-backend/internal/base"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/supabase"
	postgrest "github.com/supabase-community/postgrest-go"
)

type ExperienceRepository interface {
	base.BaseRepository[Experience, ExperienceDTO]
	CreateExperienceTechStack(ctx context.Context, experienceTechStack *ExperienceTechStack) (*ExperienceTechStack, error)
	DeleteExperienceTechStack(ctx context.Context, experienceID string) error
}

type experienceRepository struct {
	supabaseClient *supabase.SupabaseClient
	storage        supabase.SupabaseStorage
	table          string
}

func NewExperienceRepository(supabaseClient *supabase.SupabaseClient, storage supabase.SupabaseStorage) ExperienceRepository {
	return &experienceRepository{
		supabaseClient: supabaseClient,
		storage:        storage,
		table:          "experience",
	}
}

func (r *experienceRepository) Create(ctx context.Context, experience *Experience) (*Experience, error) {
	_, _, err := r.supabaseClient.GetClient().
		From(r.table).
		Insert(experience, false, "", "minimal", "").
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create experience")
	}
	return experience, nil
}

func (r *experienceRepository) CreateExperienceTechStack(ctx context.Context, experienceTechStack *ExperienceTechStack) (*ExperienceTechStack, error) {
	_, _, err := r.supabaseClient.GetClient().
		From("experience_tech_stack").
		Insert(experienceTechStack, false, "", "minimal", "").
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create experience tech stack")
	}
	return experienceTechStack, nil
}

func (r *experienceRepository) Update(ctx context.Context, experience *Experience) (*Experience, error) {
	_, _, err := r.supabaseClient.GetClient().
		From(r.table).
		Update(experience, "minimal", "").
		Eq("id", experience.ID.String()).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update experience")
	}
	return experience, nil
}

func (r *experienceRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.GetClient().
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to delete experience")
	}
	return nil
}

func (r *experienceRepository) DeleteExperienceTechStack(ctx context.Context, experienceID string) error {
	_, _, err := r.supabaseClient.GetClient().
		From("experience_tech_stack").
		Delete("minimal", "").
		Eq("experience_id", experienceID).
		Execute()
	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to delete experience tech stack")
	}
	return nil
}

func (r *experienceRepository) List(ctx context.Context, opts base.ListOptions) ([]ExperienceDTO, error) {
	var experience []ExperienceDTO
	query := r.supabaseClient.GetClient().
		From(r.table).
		Select("*, experience_tech_stack(tech_stack_id, tech_stack(id, name))", "", false)

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
			fmt.Sprintf("role.ilike.%%%s%%", opts.Search),
			"",
		)
		query = query.Or(
			fmt.Sprintf("company.ilike.%%%s%%", opts.Search),
			"",
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

	_, err := query.ExecuteTo(&experience)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to list experience")
	}

	return experience, nil
}

func (r *experienceRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
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
		return 0, errors.Wrap(err, errors.ErrDatabase, "failed to count experience")
	}

	return int(count), nil
}

func (r *experienceRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, count, err := r.supabaseClient.GetClient().
		From(r.table).
		Select("id", "exact", true).
		Eq("id", id).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, errors.Wrap(err, errors.ErrDatabase, "failed to check experience existence")
	}

	return count > 0, nil
}

func (r *experienceRepository) FindByField(ctx context.Context, field string, value interface{}) ([]ExperienceDTO, error) {
	var experience []ExperienceDTO
	_, err := r.supabaseClient.GetClient().
		From(r.table).
		Select("*, experience_tech_stack(tech_stack_id, tech_stack(id, name))", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&experience)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find experience by field")
	}
	return experience, nil
}

func (r *experienceRepository) Search(ctx context.Context, opts base.ListOptions) ([]ExperienceDTO, int, error) {
	var experience []ExperienceDTO
	query := r.supabaseClient.GetClient().
		From(r.table).
		Select("*, experience_tech_stack(tech_stack_id, tech_stack(id, name))", "", false)

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
			fmt.Sprintf("role.ilike.%%%s%%", opts.Search),
			"",
		)
		query = query.Or(
			fmt.Sprintf("tech_stack.name.ilike.%%%s%%", opts.Search),
			"",
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

	_, err := query.ExecuteTo(&experience)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to search experience")
	}

	// Count total results
	count, err := r.Count(ctx, opts.Filters)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to count experience")
	}

	return experience, count, nil
}
