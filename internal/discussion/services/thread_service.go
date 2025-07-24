package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
)

type threadService struct {
	threadRepo repositories.ThreadRepository
}

// NewThreadService membuat instance baru dari thread service
func NewThreadService(threadRepo repositories.ThreadRepository) ThreadService {
	return &threadService{
		threadRepo: threadRepo,
	}
}

// CreateThread menambahkan thread baru ke database
// Melakukan validasi objek thread sebelum membuat
func (s *threadService) CreateThread(ctx context.Context, thread *models.Thread) error {
	// Validasi objek thread
	if thread == nil {
		return fmt.Errorf("thread tidak boleh nil")
	}

	// Validasi field yang wajib diisi (contoh: Title)
	if thread.Title == "" {
		return fmt.Errorf("judul thread wajib diisi")
	}

	thread.ID = uuid.NewString()

	// Panggil repository untuk membuat thread
	return s.threadRepo.Create(ctx, thread)
}

// GetThreads mengambil daftar thread dengan paginasi
func (s *threadService) GetThreads(ctx context.Context, limit, offset int) ([]*models.Thread, error) {
	// Validasi parameter paginasi
	if limit <= 0 {
		limit = 10 // Limit default
	}
	if offset < 0 {
		offset = 0
	}

	// Ambil thread dari repository
	threads, err := s.threadRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Konversi []models.Thread ke []*models.Thread
	threadPtrs := make([]*models.Thread, len(threads))
	for i := range threads {
		threadPtrs[i] = &threads[i]
	}

	return threadPtrs, nil
}

// GetThreadByID mengambil satu thread berdasarkan ID uniknya
func (s *threadService) GetThreadByID(ctx context.Context, id string) (*models.Thread, error) {
	// Validasi ID
	if id == "" {
		return nil, fmt.Errorf("thread ID tidak boleh kosong")
	}

	// Ambil thread dari repository
	return s.threadRepo.FindByID(ctx, id)
}

// GetThreadByTitle mengambil satu thread berdasarkan judul uniknya
func (s *threadService) GetThreadByTitle(ctx context.Context, title string) (*models.Thread, error) {
	// Validasi title
	if title == "" {
		return nil, fmt.Errorf("judul thread tidak boleh kosong")
	}

	// Ambil thread dari repository
	return s.threadRepo.FindByTitle(ctx, title)
}

// UpdateThread memperbarui thread yang sudah ada di database
func (s *threadService) UpdateThread(ctx context.Context, thread *models.Thread) error {
	// Validasi objek thread
	if thread == nil {
		return fmt.Errorf("thread tidak boleh nil")
	}

	// Validasi field yang wajib diisi
	if thread.ID == "" {
		return fmt.Errorf("thread ID wajib diisi untuk update")
	}

	// Panggil repository untuk update thread
	return s.threadRepo.Update(ctx, thread)
}

// DeleteThread menghapus thread dari database berdasarkan ID-nya
func (s *threadService) DeleteThread(ctx context.Context, id string) error {
	// Validasi ID
	if id == "" {
		return fmt.Errorf("thread ID tidak boleh kosong")
	}

	// Panggil repository untuk menghapus thread
	return s.threadRepo.Delete(ctx, id)
}

// Count menghitung jumlah total thread yang tersimpan
func (s *threadService) Count(ctx context.Context) (int, error) {
	return s.threadRepo.Count(ctx)
}
