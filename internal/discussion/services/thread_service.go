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

type threadService struct {
	threadRepo repositories.ThreadRepository
}

func NewThreadService(threadRepo repositories.ThreadRepository) ThreadService {
	return &threadService{
		threadRepo: threadRepo,
	}
}

func (s *threadService) CreateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error) {
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

	// Set default values
	if thread.ID == uuid.Nil {
		thread.ID = uuid.New()
	}

	now := time.Now()
	thread.CreatedAt = &now

	// Set default status if not provided
	if thread.Status == "" {
		thread.Status = "active"
	}

	// Call repository to create thread
	return s.threadRepo.Create(ctx, thread)
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

	// Call repository to join thread
	return s.threadRepo.JoinThread(ctx, threadID, userID)
}
