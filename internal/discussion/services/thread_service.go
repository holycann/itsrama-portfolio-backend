package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type threadService struct {
	threadRepo repositories.ThreadRepository
}

func NewThreadService(threadRepo repositories.ThreadRepository) ThreadService {
	return &threadService{
		threadRepo: threadRepo,
	}
}

func (s *threadService) CreateThread(ctx context.Context, thread *models.Thread) error {
	// Validasi objek thread
	if thread == nil {
		return fmt.Errorf("thread tidak boleh nil")
	}

	// Set default values
	thread.ID = uuid.New()
	thread.CreatedAt = time.Now()

	// Set default status if not provided
	if thread.Status == "" {
		thread.Status = "active"
	}

	// Panggil repository untuk membuat thread
	return s.threadRepo.Create(ctx, thread)
}

func (s *threadService) GetThreadByID(ctx context.Context, id string) (*models.ResponseThread, error) {
	// Validasi ID
	if id == "" {
		return nil, fmt.Errorf("thread ID tidak boleh kosong")
	}

	// Ambil thread dari repository
	thread, err := s.threadRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (s *threadService) ListThreads(ctx context.Context, opts repository.ListOptions) ([]models.ResponseThread, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	threads, err := s.threadRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (s *threadService) UpdateThread(ctx context.Context, thread *models.Thread) error {
	// Validasi objek thread
	if thread == nil {
		return fmt.Errorf("thread tidak boleh nil")
	}

	// Validasi field yang wajib diisi
	if thread.ID == uuid.Nil {
		return fmt.Errorf("thread ID wajib diisi untuk update")
	}

	// Update timestamp
	thread.UpdatedAt = time.Now()

	// Panggil repository untuk update thread
	return s.threadRepo.Update(ctx, thread)
}

func (s *threadService) DeleteThread(ctx context.Context, id string) error {
	// Validasi ID
	if id == "" {
		return fmt.Errorf("thread ID tidak boleh kosong")
	}

	// Panggil repository untuk menghapus thread
	return s.threadRepo.Delete(ctx, id)
}

func (s *threadService) CountThreads(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.threadRepo.Count(ctx, filters)
}

func (s *threadService) GetThreadByEvent(ctx context.Context, eventID string) (*models.ResponseThread, error) {
	thread, err := s.threadRepo.FindThreadByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (s *threadService) GetActiveThreads(ctx context.Context, limit int) ([]models.ResponseThread, error) {
	if limit <= 0 {
		limit = 10
	}
	threads, err := s.threadRepo.FindActiveThreads(ctx, limit)
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (s *threadService) SearchThreads(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseThread, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	// Add search query to filters
	opts.Filters = append(opts.Filters,
		repository.FilterOption{
			Field:    "title",
			Operator: "like",
			Value:    query,
		},
	)

	threads, err := s.threadRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (s *threadService) JoinThread(ctx context.Context, threadID, userID string) error {
	// Validate input parameters
	if threadID == "" {
		return fmt.Errorf("thread ID cannot be empty")
	}
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Call repository to join thread
	return s.threadRepo.JoinThread(ctx, threadID, userID)
}
