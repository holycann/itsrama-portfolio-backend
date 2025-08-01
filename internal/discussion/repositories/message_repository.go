// Package repositories provides an implementation of repository for message data management
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

type messageRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewMessageRepository(supabaseClient *supabase.Client) MessageRepository {
	return &messageRepository{
		supabaseClient: supabaseClient,
		table:          "messages",
	}
}

func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	var participants []models.ResponseParticipant
	participantQuery := r.supabaseClient.
		From("discussion_participants").
		Select("*", "", false).
		Eq("thread_id", message.ThreadID.String()).
		Eq("user_id", message.UserID.String())

	_, err := participantQuery.ExecuteTo(&participants)
	if err != nil {
		return fmt.Errorf("error checking thread participant: %v", err)
	}

	if len(participants) == 0 {
		return fmt.Errorf("user is not a participant of this thread")
	}

	_, _, err = r.supabaseClient.
		From(r.table).
		Insert(message, false, "", "minimal", "").
		Execute()
	return err
}

func (r *messageRepository) FindByID(ctx context.Context, id string) (*models.ResponseMessage, error) {
	var responseMessage *models.ResponseMessage
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&responseMessage)
	if err != nil {
		return nil, err
	}

	return responseMessage, nil
}

func (r *messageRepository) Update(ctx context.Context, message *models.Message) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(message, "minimal", "").
		Eq("id", message.ID.String()).
		Execute()
	return err
}

func (r *messageRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *messageRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.ResponseMessage, error) {
	var responseMessages []models.ResponseMessage
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

	_, err := query.ExecuteTo(&responseMessages)
	if err != nil {
		return nil, err
	}

	return responseMessages, nil
}

func (r *messageRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
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

func (r *messageRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *messageRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.ResponseMessage, error) {
	var responseMessages []models.ResponseMessage
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&responseMessages)
	if err != nil {
		return nil, err
	}

	return responseMessages, nil
}

// Specialized methods for messages
func (r *messageRepository) FindMessagesByThread(ctx context.Context, threadID string) ([]models.ResponseMessage, error) {
	var responseMessages []models.ResponseMessage
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("thread_id", threadID).
		Order("created_at", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&responseMessages)
	if err != nil {
		return nil, err
	}

	return responseMessages, nil
}

func (r *messageRepository) FindMessagesByUser(ctx context.Context, userID string) ([]models.ResponseMessage, error) {
	var responseMessages []models.ResponseMessage
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("user_id", userID).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&responseMessages)
	if err != nil {
		return nil, err
	}

	return responseMessages, nil
}

func (r *messageRepository) FindRecentMessages(ctx context.Context, limit int) ([]models.ResponseMessage, error) {
	var responseMessages []models.ResponseMessage
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Limit(limit, "").
		ExecuteTo(&responseMessages)
	if err != nil {
		return nil, err
	}

	return responseMessages, nil
}

func (r *messageRepository) Search(ctx context.Context, opts repository.ListOptions) ([]models.ResponseMessage, int, error) {
	var responseMessages []models.ResponseMessage

	query := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

	// Apply search query if provided
	if opts.SearchQuery != "" {
		query = query.Or(
			fmt.Sprintf("content.ilike.%%%s%%", opts.SearchQuery),
			fmt.Sprintf("thread_id.ilike.%%%s%%", opts.SearchQuery),
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
	_, err := query.ExecuteTo(&responseMessages)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search messages: %w", err)
	}

	// Count total matching records
	_, count, err := query.Execute()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return responseMessages, int(count), nil
}
