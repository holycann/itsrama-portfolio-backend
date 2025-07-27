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

	// Validasi field yang wajib diisi
	if thread.Title == "" {
		return fmt.Errorf("judul thread wajib diisi")
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

func (s *threadService) GetThreadByID(ctx context.Context, id string) (*models.Thread, error) {
	// Validasi ID
	if id == "" {
		return nil, fmt.Errorf("thread ID tidak boleh kosong")
	}

	// Ambil thread dari repository
	return s.threadRepo.FindByID(ctx, id)
}

func (s *threadService) ListThreads(ctx context.Context, opts repository.ListOptions) ([]models.Thread, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.threadRepo.List(ctx, opts)
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

func (s *threadService) GetThreadByTitle(ctx context.Context, title string) (*models.Thread, error) {
	// Validasi title
	if title == "" {
		return nil, fmt.Errorf("judul thread tidak boleh kosong")
	}

	// Ambil thread dari repository
	return s.threadRepo.FindByTitle(ctx, title)
}

func (s *threadService) GetThreadsByEvent(ctx context.Context, eventID string) ([]models.Thread, error) {
	return s.threadRepo.FindThreadsByEvent(ctx, eventID)
}

func (s *threadService) GetActiveThreads(ctx context.Context, limit int) ([]models.Thread, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.threadRepo.FindActiveThreads(ctx, limit)
}

func (s *threadService) SearchThreads(ctx context.Context, query string, opts repository.ListOptions) ([]models.Thread, error) {
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

	return s.threadRepo.List(ctx, opts)
}
