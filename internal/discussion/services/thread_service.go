package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
)

type threadService struct {
	threadRepo         repositories.ThreadRepository
	participantService ParticipantService
}

func NewThreadService(threadRepo repositories.ThreadRepository, participantService ParticipantService) ThreadService {
	return &threadService{
		threadRepo:         threadRepo,
		participantService: participantService,
	}
}

func (s *threadService) CreateThread(ctx context.Context, thread *models.CreateThread) (*models.Thread, error) {
	// Validate thread object
	if thread == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread cannot be nil",
			nil,
		)
	}

	// Validate model
	if err := base.ValidateModel(thread); err != nil {
		return nil, err
	}

	// Additional validation for event ID
	if thread.EventID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Event ID is required and must be a valid UUID",
			nil,
		)
	}

	// Check if thread for event already exists
	existingThreads, _ := s.threadRepo.FindByField(ctx, "event_id", thread.EventID.String())

	// If thread already exists, return an error with more context
	if len(existingThreads) > 0 {
		return nil, errors.New(
			errors.ErrConflict,
			fmt.Sprintf("A discussion thread already exists for event %s", thread.EventID),
			nil,
		)
	}

	// Set default status if not provided
	if thread.Status == "" {
		thread.Status = "active"
	}

	// Prepare thread model for creation
	now := time.Now().UTC()
	threadModel := &models.Thread{
		ID:        uuid.New(),
		EventID:   thread.EventID,
		CreatorID: thread.CreatorID,
		Status:    thread.Status,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	// Call repository to create thread
	return s.threadRepo.Create(ctx, threadModel)
}

func (s *threadService) GetThreadByID(ctx context.Context, id string) (*models.ThreadDTO, error) {
	// Validate ID
	if id == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread ID cannot be empty",
			nil,
		)
	}

	// Retrieve thread from repository
	thread, err := s.threadRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(
			err,
			errors.ErrNotFound,
			"Failed to retrieve thread",
		)
	}
	return thread, nil
}

func (s *threadService) ListThreads(ctx context.Context, opts base.ListOptions) ([]models.ThreadDTO, error) {
	// Set default pagination
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 10
	}

	threads, err := s.threadRepo.List(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(
			err,
			errors.ErrDatabase,
			"Failed to list threads",
		)
	}
	return threads, nil
}

func (s *threadService) UpdateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
	// Validate thread object
	if thread == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread cannot be nil",
			nil,
		)
	}

	// Validate model
	if err := base.ValidateModel(thread); err != nil {
		return nil, err
	}

	// Validate required fields
	if thread.ID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread ID is required for update",
			nil,
		)
	}

	// Update timestamp
	now := time.Now()
	thread.UpdatedAt = &now

	// Call repository to update thread
	return s.threadRepo.Update(ctx, thread)
}

func (s *threadService) DeleteThread(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"Thread ID cannot be empty",
			nil,
		)
	}

	// Call repository to delete thread
	return s.threadRepo.Delete(ctx, id)
}

func (s *threadService) CountThreads(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.threadRepo.Count(ctx, filters)
}

func (s *threadService) GetThreadByEvent(ctx context.Context, eventID string) (*models.ThreadDTO, error) {
	// Validate event ID
	if eventID == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Event ID cannot be empty",
			nil,
		)
	}

	thread, err := s.threadRepo.FindThreadByEvent(ctx, eventID)
	if err != nil {
		return nil, errors.Wrap(
			err,
			errors.ErrNotFound,
			"Failed to retrieve thread by event",
		)
	}
	return thread, nil
}

func (s *threadService) GetActiveThreads(ctx context.Context, limit int) ([]models.ThreadDTO, error) {
	// Set default limit
	if limit <= 0 {
		limit = 10
	}

	threads, err := s.threadRepo.FindActiveThreads(ctx, limit)
	if err != nil {
		return nil, errors.Wrap(
			err,
			errors.ErrDatabase,
			"Failed to retrieve active threads",
		)
	}
	return threads, nil
}

func (s *threadService) SearchThreads(ctx context.Context, query string, opts base.ListOptions) ([]models.ThreadDTO, int, error) {
	// Set default pagination
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 10
	}

	// Search threads
	threads, count, err := s.threadRepo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(
			err,
			errors.ErrDatabase,
			"Failed to search threads",
		)
	}
	return threads, count, nil
}

func (s *threadService) JoinThread(ctx context.Context, threadID, userID string) error {
	// Validate input parameters
	threadUUID, err := uuid.Parse(threadID)
	if err != nil {
		return errors.New(
			errors.ErrValidation,
			"Invalid thread ID format",
			nil,
		)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New(
			errors.ErrValidation,
			"Invalid user ID format",
			nil,
		)
	}

	// Retrieve the thread to ensure it exists and is active
	thread, err := s.threadRepo.FindByID(ctx, threadID)
	if err != nil {
		return errors.Wrap(
			err,
			errors.ErrNotFound,
			"Thread not found",
		)
	}

	// Check thread status
	if thread.Status != "active" {
		return errors.New(
			errors.ErrValidation,
			"Cannot join an inactive thread",
			nil,
		)
	}

	// Create a participant model for validation
	participant := &models.Participant{
		ThreadID: threadUUID,
		UserID:   userUUID,
	}

	// Validate the model
	if err := base.ValidateModel(participant); err != nil {
		return err
	}

	// Check if participant already exists (atomic operation)
	existingParticipant, err := s.participantService.GetParticipantByThread(ctx, userID, threadID)
	if err != nil {
		return errors.Wrap(
			err,
			errors.ErrDatabase,
			"Failed to check existing participant",
		)
	}

	if existingParticipant != nil {
		return errors.New(
			errors.ErrConflict,
			"User is already a participant in this thread",
			nil,
		)
	}

	// Create new participant with proper timestamps
	now := time.Now().UTC()
	newParticipant := &models.Participant{
		ThreadID:  threadUUID,
		UserID:    userUUID,
		JoinedAt:  &now,
		UpdatedAt: &now,
	}

	// Attempt to create participant
	_, err = s.participantService.CreateParticipant(ctx, newParticipant)
	if err != nil {
		return errors.Wrap(
			err,
			errors.ErrDatabase,
			"Failed to join thread",
		)
	}

	return nil
}
