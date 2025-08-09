package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
)

type participantService struct {
	participantRepo repositories.ParticipantRepository
}

func NewParticipantService(participantRepo repositories.ParticipantRepository) ParticipantService {
	return &participantService{
		participantRepo: participantRepo,
	}
}

func (s *participantService) CreateParticipant(ctx context.Context, participant *models.Participant) (*models.Participant, error) {
	// Validate participant object
	if participant == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Participant cannot be nil",
			nil,
		)
	}

	// Validate model
	if err := base.ValidateModel(participant); err != nil {
		return nil, err
	}

	// Validate thread and user IDs
	if participant.ThreadID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread ID is required and must be a valid UUID",
			nil,
		)
	}

	if participant.UserID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"User ID is required and must be a valid UUID",
			nil,
		)
	}

	// Perform an atomic check for existing participant
	existingParticipant, err := s.GetParticipantByThread(ctx, participant.ThreadID.String(), participant.UserID.String())
	if err != nil {
		return nil, errors.Wrap(
			err,
			errors.ErrDatabase,
			"Failed to check existing participant",
		)
	}

	// If participant already exists, return an error
	if existingParticipant != nil {
		return nil, errors.New(
			errors.ErrConflict,
			"User is already a participant in this thread",
			nil,
		)
	}

	// Set timestamps
	now := time.Now().UTC()
	participant.JoinedAt = &now
	participant.UpdatedAt = &now

	// Attempt to create participant
	createdParticipant, err := s.participantRepo.Create(ctx, participant)
	if err != nil {
		return nil, errors.Wrap(
			err,
			errors.ErrDatabase,
			"Failed to create participant",
		)
	}

	return createdParticipant, nil
}

func (s *participantService) GetParticipantByID(ctx context.Context, id string) (*models.ParticipantDTO, error) {
	// Validate ID
	if id == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Participant ID cannot be empty",
			nil,
		)
	}

	// Retrieve participant from repository
	return s.participantRepo.FindByID(ctx, id)
}

func (s *participantService) ListParticipants(ctx context.Context, opts base.ListOptions) ([]models.ParticipantDTO, error) {
	// Set default pagination
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 10
	}

	return s.participantRepo.List(ctx, opts)
}

func (s *participantService) UpdateParticipant(ctx context.Context, participant *models.Participant) (*models.Participant, error) {
	// Validate participant object
	if participant == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Participant cannot be nil",
			nil,
		)
	}

	// Validate model
	if err := base.ValidateModel(participant); err != nil {
		return nil, err
	}

	// Validate required fields
	if participant.ThreadID == uuid.Nil || participant.UserID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread ID and User ID are required for update",
			nil,
		)
	}

	// Update timestamp
	now := time.Now()
	participant.UpdatedAt = &now

	// Call repository to update participant
	return s.participantRepo.Update(ctx, participant)
}

func (s *participantService) CountParticipants(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.participantRepo.Count(ctx, filters)
}

func (s *participantService) GetParticipantsByThread(ctx context.Context, threadID string) ([]models.ParticipantDTO, error) {
	// Validate thread ID
	if threadID == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread ID cannot be empty",
			nil,
		)
	}

	return s.participantRepo.FindParticipantsByThread(ctx, threadID)
}

func (s *participantService) GetParticipantByThread(ctx context.Context, threadID, userID string) (*models.ParticipantDTO, error) {
	// Validate thread ID
	if threadID == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread ID cannot be empty",
			nil,
		)
	}

	// Validate user ID
	if userID == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"User ID cannot be empty",
			nil,
		)
	}

	return s.participantRepo.FindParticipantByThread(ctx, threadID, userID)
}

func (s *participantService) GetThreadParticipants(ctx context.Context, threadID string) ([]models.ParticipantDTO, error) {
	// Validate thread ID
	if threadID == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread ID cannot be empty",
			nil,
		)
	}

	return s.participantRepo.FindThreadParticipants(ctx, threadID)
}

func (s *participantService) RemoveParticipant(ctx context.Context, threadID, userID string) error {
	// Validate input parameters
	if threadID == "" {
		return errors.New(
			errors.ErrValidation,
			"Thread ID cannot be empty",
			nil,
		)
	}
	if userID == "" {
		return errors.New(
			errors.ErrValidation,
			"User ID cannot be empty",
			nil,
		)
	}

	// Call repository to remove participant
	return s.participantRepo.RemoveParticipant(ctx, threadID, userID)
}

func (s *participantService) SearchParticipants(ctx context.Context, query string, opts base.ListOptions) ([]models.ParticipantDTO, int, error) {
	// Set default pagination
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 10
	}

	// Search participants
	participants, count, err := s.participantRepo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(
			err,
			errors.ErrDatabase,
			"Failed to search participants",
		)
	}
	return participants, count, nil
}
