// Package repositories provides an implementation of repository for discussion participant data management
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

type participantRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewParticipantRepository(supabaseClient *supabase.Client) ParticipantRepository {
	return &participantRepository{
		supabaseClient: supabaseClient,
		table:          "discussion_participants",
	}
}

func (r *participantRepository) Create(ctx context.Context, participant *models.Participant) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(participant, false, "", "minimal", "").
		Execute()
	return err
}

func (r *participantRepository) FindByID(ctx context.Context, id string) (*models.ResponseParticipant, error) {
	var participant *models.ResponseParticipant
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&participant)
	return participant, err
}

func (r *participantRepository) Update(ctx context.Context, participant *models.Participant) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(participant, "minimal", "").
		Eq("thread_id", participant.ThreadID.String()).
		Eq("user_id", participant.UserID.String()).
		Execute()
	return err
}

// Delete method is deprecated and will always return an error
func (r *participantRepository) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("delete method is deprecated for participant repository")
}

func (r *participantRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.ResponseParticipant, error) {
	var participants []models.ResponseParticipant
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

	_, err := query.ExecuteTo(&participants)
	return participants, err
}

func (r *participantRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
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

// Deprecated: This method is no longer recommended for use
func (r *participantRepository) Exists(ctx context.Context, id string) (bool, error) {
	// This method is deprecated and will be removed in future versions
	return false, fmt.Errorf("method Exists is deprecated")
}

func (r *participantRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.ResponseParticipant, error) {
	var participants []models.ResponseParticipant
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&participants)
	return participants, err
}

func (r *participantRepository) FindParticipantsByThread(ctx context.Context, threadID string) ([]models.ResponseParticipant, error) {
	var responseParticipants []models.ResponseParticipant
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("thread_id", threadID).
		ExecuteTo(&responseParticipants)
	if err != nil {
		return nil, err
	}

	return responseParticipants, nil
}

func (r *participantRepository) FindThreadParticipants(ctx context.Context, threadID string) ([]models.ResponseParticipant, error) {
	var responseParticipants []models.ResponseParticipant
	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("thread_id", threadID).
		ExecuteTo(&responseParticipants)
	if err != nil {
		return nil, err
	}

	return responseParticipants, nil
}

func (r *participantRepository) RemoveParticipant(ctx context.Context, threadID, userID string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("thread_id", threadID).
		Eq("user_id", userID).
		Execute()
	if err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	return nil
}

func (r *participantRepository) Search(ctx context.Context, opts repository.ListOptions) ([]models.ResponseParticipant, int, error) {
	var responseParticipants []models.ResponseParticipant
	query := r.supabaseClient.
		From(r.table).
		Select("*", "", false)

	// Apply search query if provided
	if opts.SearchQuery != "" {
		query = query.Or(
			fmt.Sprintf("users_profile.fullname.ilike.%%%s%%", opts.SearchQuery),
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
	_, err := query.ExecuteTo(&responseParticipants)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search participants: %w", err)
	}

	// Count total matching records
	_, count, err := query.Execute()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count participants: %w", err)
	}

	return responseParticipants, int(count), nil
}
