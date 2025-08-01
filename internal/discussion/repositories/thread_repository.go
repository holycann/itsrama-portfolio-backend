// Package repositories provides an implementation of repository for discussion thread data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/repository"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type threadRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewThreadRepository(supabaseClient *supabase.Client) ThreadRepository {
	return &threadRepository{
		supabaseClient: supabaseClient,
		table:          "threads",
	}
}

func (r *threadRepository) Create(ctx context.Context, thread *models.Thread) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(thread, false, "", "minimal", "").
		ExecuteTo(&thread)
	return err
}

func (r *threadRepository) FindByID(ctx context.Context, id string) (*models.Thread, error) {
	var thread *models.Thread
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&thread)
	return thread, err
}

func (r *threadRepository) Update(ctx context.Context, thread *models.Thread) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(thread, "minimal", "").
		Eq("id", thread.ID.String()).
		Execute()
	return err
}

func (r *threadRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *threadRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.Thread, error) {
	var threads []models.Thread
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

	_, err := query.ExecuteTo(&threads)
	return threads, err
}

func (r *threadRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
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

func (r *threadRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *threadRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.Thread, error) {
	var threads []models.Thread
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&threads)
	return threads, err
}

// Specialized methods for threads
func (r *threadRepository) FindByTitle(ctx context.Context, title string) (*models.Thread, error) {
	threads, err := r.FindByField(ctx, "title", title)
	if err != nil {
		return nil, err
	}
	if len(threads) == 0 {
		return nil, fmt.Errorf("thread not found")
	}
	return &threads[0], nil
}

func (r *threadRepository) FindThreadsByEvent(ctx context.Context, eventID string) ([]models.Thread, error) {
	var threads []models.Thread
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("event_id", eventID).
		ExecuteTo(&threads)
	return threads, err
}

func (r *threadRepository) FindActiveThreads(ctx context.Context, limit int) ([]models.Thread, error) {
	var threads []models.Thread
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("status", "active").
		Limit(limit, "").
		ExecuteTo(&threads)
	return threads, err
}

func (r *threadRepository) Search(ctx context.Context, opts repository.ListOptions) ([]models.Thread, int, error) {
	var threads []models.Thread

	query := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

	// Apply search query if provided
	if opts.SearchQuery != "" {
		query = query.Or(
			fmt.Sprintf("title.ilike.%%%s%%", opts.SearchQuery),
			fmt.Sprintf("description.ilike.%%%s%%", opts.SearchQuery),
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
	_, err := query.ExecuteTo(&threads)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search threads: %w", err)
	}

	// Count total matching records
	_, count, err := query.Execute()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count threads: %w", err)
	}

	return threads, int(count), nil
}
