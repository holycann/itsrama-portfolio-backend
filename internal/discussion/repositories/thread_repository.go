// Package repositories provides an implementation of repository for discussion thread data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/base"
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

func (r *threadRepository) Create(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(thread, false, "", "minimal", "").
		Execute()
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (r *threadRepository) FindByID(ctx context.Context, id string) (*models.ThreadDTO, error) {
	var threadDTO models.ThreadDTO

	// Use join to fetch thread and participant in a single query
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*), creator:users_view!threads_creator_id_fkey(*)", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&threadDTO)

	if err != nil {
		return nil, err
	}

	return &threadDTO, nil
}

func (r *threadRepository) Update(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(thread, "minimal", "").
		Eq("id", thread.ID.String()).
		Execute()
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (r *threadRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *threadRepository) List(ctx context.Context, opts base.ListOptions) ([]models.ThreadDTO, error) {
	var threads []models.ThreadDTO

	query := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*), creator:users_view!threads_creator_id_fkey(*)", "", false)

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

	_, err := query.ExecuteTo(&threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (r *threadRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
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

func (r *threadRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", true).
		Eq("id", id).
		Limit(1, "").
		Execute()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *threadRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.ThreadDTO, error) {
	var threads []models.ThreadDTO

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*), creator:users_view!threads_creator_id_fkey(*)", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (r *threadRepository) FindThreadByEvent(ctx context.Context, eventID string) (*models.ThreadDTO, error) {
	var thread *models.ThreadDTO

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*), creator:users_view!threads_creator_id_fkey(*)", "", false).
		Eq("event_id", eventID).
		Single().
		ExecuteTo(&thread)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (r *threadRepository) FindActiveThreads(ctx context.Context, limit int) ([]models.ThreadDTO, error) {
	var threads []models.ThreadDTO

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*), creator:users_view!threads_creator_id_fkey(*)", "", false).
		Eq("status", "active").
		Limit(limit, "").
		ExecuteTo(&threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (r *threadRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.ThreadDTO, int, error) {
	var threads []models.ThreadDTO

	query := r.supabaseClient.
		From(r.table).
		Select("*, discussion_participants(*), creator:users_view!threads_creator_id_fkey(*)", "", false)

	// Apply search query if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("title.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("description.ilike.%%%s%%", opts.Search),
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

func (r *threadRepository) BulkCreate(ctx context.Context, threads []*models.Thread) ([]models.Thread, error) {
	var createdThreads []models.Thread

	for _, thread := range threads {
		// Generate UUID if not provided
		if base.IsZero(thread.ID) {
			thread.ID = uuid.New()
		}

		_, _, err := r.supabaseClient.
			From(r.table).
			Insert(thread, false, "", "minimal", "").
			Execute()
		if err != nil {
			return nil, err
		}

		createdThreads = append(createdThreads, *thread)
	}

	return createdThreads, nil
}

func (r *threadRepository) BulkUpdate(ctx context.Context, threads []*models.Thread) ([]models.Thread, error) {
	var updatedThreads []models.Thread

	for _, thread := range threads {
		_, _, err := r.supabaseClient.
			From(r.table).
			Update(thread, "minimal", "").
			Eq("id", thread.ID.String()).
			Execute()
		if err != nil {
			return nil, err
		}

		updatedThreads = append(updatedThreads, *thread)
	}

	return updatedThreads, nil
}

func (r *threadRepository) BulkDelete(ctx context.Context, ids []string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		In("id", ids).
		Execute()
	return err
}

func (r *threadRepository) BulkUpsert(ctx context.Context, threads []*models.Thread) ([]models.ThreadDTO, error) {
	var upsertedThreads []models.ThreadDTO

	for _, thread := range threads {
		// Generate UUID if not provided
		if base.IsZero(thread.ID) {
			thread.ID = uuid.New()
		}

		var upsertedThread models.ThreadDTO
		_, err := r.supabaseClient.
			From(r.table).
			Upsert(thread, "id", "minimal", "").
			Single().
			ExecuteTo(&upsertedThread)
		if err != nil {
			return nil, err
		}

		upsertedThreads = append(upsertedThreads, upsertedThread)
	}

	return upsertedThreads, nil
}
