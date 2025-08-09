// Package repositories provides an implementation of repository for message data management
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

func (r *messageRepository) Create(ctx context.Context, message *models.Message) (*models.Message, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(message, false, "", "minimal", "").
		Execute()
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (r *messageRepository) FindByID(ctx context.Context, id string) (*models.MessageDTO, error) {
	var messageDTO models.MessageDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, sender:users_view(*, profile:users_profile(*))", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&messageDTO)
	if err != nil {
		return nil, err
	}

	return &messageDTO, nil
}

func (r *messageRepository) Update(ctx context.Context, message *models.Message) (*models.Message, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(message, "minimal", "").
		Eq("id", message.ID.String()).
		Execute()
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (r *messageRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *messageRepository) List(ctx context.Context, opts base.ListOptions) ([]models.MessageDTO, error) {
	var messageDTOs []models.MessageDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*, sender:users_view(*, profile:users_profile(*))", "", false)

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

	_, err := query.ExecuteTo(&messageDTOs)
	if err != nil {
		return nil, err
	}

	return messageDTOs, nil
}

func (r *messageRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
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

func (r *messageRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *messageRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.MessageDTO, error) {
	var messageDTOs []models.MessageDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, sender:users_view(*, profile:users_profile(*))", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&messageDTOs)
	if err != nil {
		return nil, err
	}

	return messageDTOs, nil
}

func (r *messageRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.MessageDTO, int, error) {
	var messageDTOs []models.MessageDTO

	query := r.supabaseClient.
		From(r.table).
		Select("*, sender:users_view(*, profile:users_profile(*))", "", false)

	// Apply search query if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("content.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("thread_id.ilike.%%%s%%", opts.Search),
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
	_, err := query.ExecuteTo(&messageDTOs)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search messages: %w", err)
	}

	// Count total matching records
	_, count, err := query.Execute()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return messageDTOs, int(count), nil
}

func (r *messageRepository) BulkCreate(ctx context.Context, messages []*models.Message) ([]models.Message, error) {
	var createdMessages []models.Message

	for _, message := range messages {
		// Generate UUID if not provided
		if base.IsZero(message.ID) {
			message.ID = uuid.New()
		}

		_, _, err := r.supabaseClient.
			From(r.table).
			Insert(message, false, "", "minimal", "").
			Execute()
		if err != nil {
			return nil, err
		}

		createdMessages = append(createdMessages, *message)
	}

	return createdMessages, nil
}

func (r *messageRepository) BulkUpdate(ctx context.Context, messages []*models.Message) ([]models.Message, error) {
	var updatedMessages []models.Message

	for _, message := range messages {
		_, _, err := r.supabaseClient.
			From(r.table).
			Update(message, "minimal", "").
			Eq("id", message.ID.String()).
			Execute()
		if err != nil {
			return nil, err
		}

		updatedMessages = append(updatedMessages, *message)
	}

	return updatedMessages, nil
}

func (r *messageRepository) BulkDelete(ctx context.Context, ids []string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		In("id", ids).
		Execute()
	return err
}

func (r *messageRepository) FindMessagesByThread(ctx context.Context, threadID string) ([]models.MessageDTO, error) {
	var messages []models.MessageDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, sender:users_view(*, profile:users_profile(*))", "", false).
		Eq("thread_id", threadID).
		Order("created_at", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&messages)
	return messages, err
}

func (r *messageRepository) FindMessagesByUser(ctx context.Context, userID string) ([]models.MessageDTO, error) {
	var messages []models.MessageDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, sender:users_view(*, profile:users_profile(*))", "", false).
		Eq("sender_id", userID).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&messages)
	return messages, err
}

func (r *messageRepository) FindRecentMessages(ctx context.Context, limit int) ([]models.MessageDTO, error) {
	var messageDTOs []models.MessageDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, sender:users_view(*, profile:users_profile(*))", "", false).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Limit(limit, "").
		ExecuteTo(&messageDTOs)
	if err != nil {
		return nil, err
	}

	return messageDTOs, nil
}
