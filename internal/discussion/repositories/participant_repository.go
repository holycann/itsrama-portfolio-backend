// Package repositories provides an implementation of repository for discussion participant data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/base"
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

func (r *participantRepository) Create(ctx context.Context, participant *models.Participant) (*models.Participant, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(participant, false, "", "minimal", "").
		Execute()
	if err != nil {
		return nil, err
	}

	return participant, nil
}

func (r *participantRepository) FindByID(ctx context.Context, id string) (*models.ParticipantDTO, error) {
	var participantDTO models.ParticipantDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, user:users(id,username,profile_picture)", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&participantDTO)
	return &participantDTO, err
}

func (r *participantRepository) Update(ctx context.Context, participant *models.Participant) (*models.Participant, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(participant, "minimal", "").
		Eq("thread_id", participant.ThreadID.String()).
		Eq("user_id", participant.UserID.String()).
		Execute()
	if err != nil {
		return nil, err
	}
	return participant, nil
}

func (r *participantRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *participantRepository) List(ctx context.Context, opts base.ListOptions) ([]models.ParticipantDTO, error) {
	var participants []models.ParticipantDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*, user:users(id,username,profile_picture)", "", false)

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

	_, err := query.ExecuteTo(&participants)
	return participants, err
}

func (r *participantRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
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

func (r *participantRepository) Exists(ctx context.Context, id string) (bool, error) {
	count, err := r.Count(ctx, []base.FilterOption{
		{Field: "id", Operator: base.OperatorEqual, Value: id},
	})
	return count > 0, err
}

func (r *participantRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.ParticipantDTO, error) {
	var participants []models.ParticipantDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, user:users(id,username,profile_picture)", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&participants)
	return participants, err
}

func (r *participantRepository) FindParticipantsByThread(ctx context.Context, threadID string) ([]models.ParticipantDTO, error) {
	var participants []models.ParticipantDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, user:users(id,username,profile_picture)", "", false).
		Eq("thread_id", threadID).
		ExecuteTo(&participants)
	return participants, err
}

func (r *participantRepository) FindThreadParticipants(ctx context.Context, threadID string) ([]models.ParticipantDTO, error) {
	return r.FindParticipantsByThread(ctx, threadID)
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

func (r *participantRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.ParticipantDTO, int, error) {
	var participants []models.ParticipantDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*, user:users(id,username,profile_picture)", "", false)

	// Apply search query if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("user.username.ilike.%%%s%%", opts.Search),
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
	_, err := query.ExecuteTo(&participants)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search participants: %w", err)
	}

	// Count total matching records
	_, count, err := query.Execute()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count participants: %w", err)
	}

	return participants, int(count), nil
}

func (r *participantRepository) BulkCreate(ctx context.Context, participants []*models.Participant) ([]models.Participant, error) {
	var createdParticipants []models.Participant

	for _, participant := range participants {
		_, _, err := r.supabaseClient.
			From(r.table).
			Insert(participant, false, "", "minimal", "").
			Execute()
		if err != nil {
			return nil, err
		}

		createdParticipants = append(createdParticipants, *participant)
	}

	return createdParticipants, nil
}

func (r *participantRepository) BulkUpdate(ctx context.Context, participants []*models.Participant) ([]models.Participant, error) {
	var updatedParticipants []models.Participant

	for _, participant := range participants {
		_, _, err := r.supabaseClient.
			From(r.table).
			Update(participant, "minimal", "").
			Eq("thread_id", participant.ThreadID.String()).
			Eq("user_id", participant.UserID.String()).
			Execute()
		if err != nil {
			return nil, err
		}

		updatedParticipants = append(updatedParticipants, *participant)
	}

	return updatedParticipants, nil
}

func (r *participantRepository) BulkDelete(ctx context.Context, ids []string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		In("id", ids).
		Execute()
	return err
}

func (r *participantRepository) BulkUpsert(ctx context.Context, participants []*models.Participant) ([]models.ParticipantDTO, error) {
	var upsertedParticipants []models.ParticipantDTO

	for _, participant := range participants {
		var upsertedParticipant models.ParticipantDTO
		_, err := r.supabaseClient.
			From(r.table).
			Upsert(participant, "id", "minimal", "").
			Single().
			ExecuteTo(&upsertedParticipant)
		if err != nil {
			return nil, err
		}

		upsertedParticipants = append(upsertedParticipants, upsertedParticipant)
	}

	return upsertedParticipants, nil
}
