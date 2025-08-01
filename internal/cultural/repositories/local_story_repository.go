// Package repositories provides an implementation of repository for local story data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/repository"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type localStoryRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewLocalStoryRepository(supabaseClient *supabase.Client) LocalStoryRepository {
	return &localStoryRepository{
		supabaseClient: supabaseClient,
		table:          "local_stories",
	}
}

func (r *localStoryRepository) Create(ctx context.Context, localStory *models.LocalStory) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(*localStory, false, "", "minimal", "").
		ExecuteTo(localStory)
	return err
}

func (r *localStoryRepository) FindByID(ctx context.Context, id string) (*models.LocalStory, error) {
	var localStory *models.LocalStory
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&localStory)
	return localStory, err
}

func (r *localStoryRepository) Update(ctx context.Context, localStory *models.LocalStory) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(*localStory, "minimal", "").
		Eq("id", (*localStory).ID.String()).
		Execute()
	return err
}

func (r *localStoryRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *localStoryRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.LocalStory, error) {
	var localStories []models.LocalStory
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

	_, err := query.ExecuteTo(&localStories)
	return localStories, err
}

func (r *localStoryRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
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

func (r *localStoryRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *localStoryRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.LocalStory, error) {
	var localStories []models.LocalStory
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&localStories)
	return localStories, err
}

// Specialized methods for local stories
func (r *localStoryRepository) FindStoriesByLocation(ctx context.Context, locationID uuid.UUID) ([]*models.LocalStory, error) {
	var localStories []models.LocalStory
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("location_id", locationID.String()).
		ExecuteTo(&localStories)

	// Convert to slice of pointers
	result := make([]*models.LocalStory, len(localStories))
	for i := range localStories {
		result[i] = &localStories[i]
	}

	return result, err
}

func (r *localStoryRepository) FindStoriesByOriginCulture(ctx context.Context, culture string) ([]*models.LocalStory, error) {
	var localStories []models.LocalStory
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("origin_culture", culture).
		ExecuteTo(&localStories)

	// Convert to slice of pointers
	result := make([]*models.LocalStory, len(localStories))
	for i := range localStories {
		result[i] = &localStories[i]
	}

	return result, err
}

// Modify Search method to match base repository interface
func (r *localStoryRepository) Search(ctx context.Context, opts repository.ListOptions) ([]models.LocalStory, int, error) {
	var localStories []models.LocalStory
	query := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

	// Apply search query if provided
	if opts.SearchQuery != "" {
		escapedQuery := strings.ReplaceAll(strings.ReplaceAll(opts.SearchQuery, "%", "\\%"), "_", "\\_")
		likeQuery := "%" + escapedQuery + "%"
		query = query.Or(fmt.Sprintf("title.ilike.%s,story_text.ilike.%s", likeQuery, likeQuery), "")
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

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == repository.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	query = query.Range(opts.Offset, opts.Offset+opts.Limit-1, "")

	// Execute query
	_, err := query.ExecuteTo(&localStories)
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	count, err := r.Count(ctx, opts.Filters)
	if err != nil {
		return nil, 0, err
	}

	return localStories, count, nil
}
