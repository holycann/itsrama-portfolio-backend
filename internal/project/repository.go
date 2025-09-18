package project

import (
	"context"
	"fmt"

	"github.com/holycann/itsrama-portfolio-backend/internal/base"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/supabase"
	postgrest "github.com/supabase-community/postgrest-go"
)

type ProjectRepository interface {
	base.BaseRepository[Project, ProjectDTO]
	CreateProjectTechStack(ctx context.Context, project *ProjectTechStack) (*ProjectTechStack, error)
	DeleteProjectTechStack(ctx context.Context, projectID string) error
}

type projectRepository struct {
	supabaseClient *supabase.SupabaseClient
	storage        supabase.SupabaseStorage
	table          string
}

func NewProjectRepository(supabaseClient *supabase.SupabaseClient, storage supabase.SupabaseStorage) ProjectRepository {
	return &projectRepository{
		supabaseClient: supabaseClient,
		storage:        storage,
		table:          "project",
	}
}

func (r *projectRepository) Create(ctx context.Context, project *Project) (*Project, error) {
	_, _, err := r.supabaseClient.GetClient().
		From(r.table).
		Insert(project, false, "", "minimal", "").
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create project")
	}
	return project, nil
}

func (r *projectRepository) CreateProjectTechStack(ctx context.Context, project *ProjectTechStack) (*ProjectTechStack, error) {
	_, _, err := r.supabaseClient.GetClient().
		From("project_tech_stack").
		Insert(project, false, "", "minimal", "").
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create project tech stack")
	}
	return project, nil
}

func (r *projectRepository) Update(ctx context.Context, project *Project) (*Project, error) {
	_, _, err := r.supabaseClient.GetClient().
		From(r.table).
		Update(project, "minimal", "").
		Eq("id", project.ID.String()).
		Execute()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update project")
	}
	return project, nil
}

func (r *projectRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.GetClient().
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to delete project")
	}
	return nil
}

func (r *projectRepository) DeleteProjectTechStack(ctx context.Context, projectID string) error {
	_, _, err := r.supabaseClient.GetClient().
		From("project_tech_stack").
		Delete("minimal", "").
		Eq("project_id", projectID).
		Execute()
	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to delete project tech stack")
	}
	return nil
}

func (r *projectRepository) List(ctx context.Context, opts base.ListOptions) ([]ProjectDTO, error) {
	var projects []ProjectDTO
	query := r.supabaseClient.GetClient().
		From(r.table).
		Select("*, project_tech_stack(tech_stack_id, tech_stack(id, name))", "", false)

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
			fmt.Sprintf("title.ilike.%%%s%%", opts.Search),
			"",
		)
		query = query.Or(
			fmt.Sprintf("category.ilike.%%%s%%", opts.Search),
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

	_, err := query.ExecuteTo(&projects)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to list projects")
	}

	return projects, nil
}

func (r *projectRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
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
		return 0, errors.Wrap(err, errors.ErrDatabase, "failed to count projects")
	}

	return int(count), nil
}

func (r *projectRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, count, err := r.supabaseClient.GetClient().
		From(r.table).
		Select("id", "exact", true).
		Eq("id", id).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, errors.Wrap(err, errors.ErrDatabase, "failed to check project existence")
	}

	return count > 0, nil
}

func (r *projectRepository) FindByField(ctx context.Context, field string, value interface{}) ([]ProjectDTO, error) {
	var projects []ProjectDTO
	_, err := r.supabaseClient.GetClient().
		From(r.table).
		Select("*, project_tech_stack(tech_stack_id, tech_stack(id, name))", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&projects)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find projects by field")
	}
	return projects, nil
}

func (r *projectRepository) Search(ctx context.Context, opts base.ListOptions) ([]ProjectDTO, int, error) {
	var projects []ProjectDTO
	query := r.supabaseClient.GetClient().
		From(r.table).
		Select("*, project_tech_stack(tech_stack_id, tech_stack(id, name))", "", false)

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
			fmt.Sprintf("title.ilike.%%%s%%", opts.Search),
			"",
		)
		query = query.Or(
			fmt.Sprintf("category.ilike.%%%s%%", opts.Search),
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

	_, err := query.ExecuteTo(&projects)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to search projects")
	}

	// Count total results
	count, err := r.Count(ctx, opts.Filters)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to count projects")
	}

	return projects, count, nil
}
