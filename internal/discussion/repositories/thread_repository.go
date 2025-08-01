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
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(thread, false, "", "minimal", "").
		Execute()
	return err
}

func (r *threadRepository) FindByID(ctx context.Context, id string) (*models.ResponseThread, error) {
	var responseThread *models.ResponseThread

	// Use join to fetch thread and participant in a single query
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*)", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&responseThread)

	if err != nil {
		return nil, err
	}

	return responseThread, nil
}

func (r *threadRepository) Update(ctx context.Context, thread *models.Thread) error {
	// Update only the thread part
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(&thread, "minimal", "").
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

func (r *threadRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.ResponseThread, error) {
	var responseThreads []models.ResponseThread

	query := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*)", "", false)

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

	_, err := query.ExecuteTo(&responseThreads)
	if err != nil {
		return nil, err
	}

	return responseThreads, nil
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

func (r *threadRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.ResponseThread, error) {
	var responseThreads []models.ResponseThread

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*)", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&responseThreads)
	if err != nil {
		return nil, err
	}

	return responseThreads, nil
}

func (r *threadRepository) FindThreadByEvent(ctx context.Context, eventID string) (*models.ResponseThread, error) {
	var responseThread *models.ResponseThread

	// Find the thread using a join to get the first participant
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*)", "", false).
		Eq("event_id", eventID).
		Single().
		ExecuteTo(&responseThread)

	if err != nil {
		return nil, err
	}

	return responseThread, nil
}

func (r *threadRepository) FindActiveThreads(ctx context.Context, limit int) ([]models.ResponseThread, error) {
	var responseThreads []models.ResponseThread

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*)", "", false).
		Eq("status", "active").
		Limit(limit, "").
		ExecuteTo(&responseThreads)
	if err != nil {
		return nil, err
	}

	return responseThreads, nil
}

func (r *threadRepository) Search(ctx context.Context, opts repository.ListOptions) ([]models.ResponseThread, int, error) {
	var responseThreads []models.ResponseThread

	query := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*)", "", false)

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
	_, err := query.ExecuteTo(&responseThreads)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search threads: %w", err)
	}

	// Count total matching records
	_, count, err := query.Execute()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count threads: %w", err)
	}

	return responseThreads, int(count), nil
}

func (r *threadRepository) JoinThread(ctx context.Context, threadID, userID string) error {
	// Prepare the data for insertion into discussion_participants table
	data := map[string]interface{}{
		"thread_id": threadID,
		"user_id":   userID,
	}

	// Insert the participant into the discussion_participants table
	_, _, err := r.supabaseClient.
		From("discussion_participants").
		Insert(data, false, "", "minimal", "").
		Execute()

	if err != nil {
		return fmt.Errorf("failed to join thread: %w", err)
	}

	return nil
}
