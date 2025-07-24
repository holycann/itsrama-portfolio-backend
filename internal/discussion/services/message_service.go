package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
)

type messageService struct {
	messageRepo repositories.MessageRepository
}

// NewMessageService membuat instance baru dari message service
func NewMessageService(messageRepo repositories.MessageRepository) MessageService {
	return &messageService{
		messageRepo: messageRepo,
	}
}

// CreateMessage menambahkan pesan baru ke database
// Melakukan validasi objek message sebelum membuat
func (s *messageService) CreateMessage(ctx context.Context, message *models.Message) error {
	// Validasi objek message
	if message == nil {
		return fmt.Errorf("message tidak boleh nil")
	}

	// Validasi field yang wajib diisi (contoh: Content)
	if message.Content == "" {
		return fmt.Errorf("isi pesan wajib diisi")
	}

	message.ID = uuid.NewString()

	// Panggil repository untuk membuat message
	return s.messageRepo.Create(ctx, message)
}

// GetMessages mengambil daftar pesan dengan paginasi
func (s *messageService) GetMessages(ctx context.Context, limit, offset int) ([]*models.Message, error) {
	// Validasi parameter paginasi
	if limit <= 0 {
		limit = 10 // Limit default
	}
	if offset < 0 {
		offset = 0
	}

	// Ambil pesan dari repository
	messages, err := s.messageRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Konversi []models.Message ke []*models.Message
	messagePtrs := make([]*models.Message, len(messages))
	for i := range messages {
		messagePtrs[i] = &messages[i]
	}

	return messagePtrs, nil
}

// GetMessageByID mengambil satu pesan berdasarkan ID uniknya
func (s *messageService) GetMessageByID(ctx context.Context, id string) (*models.Message, error) {
	// Validasi ID
	if id == "" {
		return nil, fmt.Errorf("message ID tidak boleh kosong")
	}

	// Ambil pesan dari repository
	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Konversi ke slice pointer
	return message, nil
}

// GetMessageByContent mengambil pesan berdasarkan kontennya
// Catatan: Metode ini tidak didukung langsung oleh repository saat ini
// Anda mungkin perlu menambah metode custom di repository atau melakukan filter manual
func (s *messageService) GetMessageByContent(ctx context.Context, content string) ([]*models.Message, error) {
	// Validasi content
	if content == "" {
		return nil, fmt.Errorf("isi pesan tidak boleh kosong")
	}

	// Karena repository belum punya metode langsung, workaround: ambil semua pesan lalu cari yang sesuai
	messages, err := s.messageRepo.List(ctx, 1000, 0) // Naikkan limit untuk mencari lebih banyak pesan
	if err != nil {
		return nil, err
	}

	// Cari pesan berdasarkan content (linear search)
	matchedMessages := []*models.Message{}
	for i := range messages {
		if messages[i].Content == content {
			matchedMessages = append(matchedMessages, &messages[i])
		}
	}

	if len(matchedMessages) == 0 {
		return nil, fmt.Errorf("pesan dengan konten %s tidak ditemukan", content)
	}

	return matchedMessages, nil
}

// GetMessagesByThreadID mengambil daftar pesan berdasarkan ID thread dengan paginasi
func (s *messageService) GetMessagesByThreadID(ctx context.Context, threadID string, limit, offset int) ([]*models.Message, error) {
	// Validasi parameter
	if threadID == "" {
		return nil, fmt.Errorf("thread ID tidak boleh kosong")
	}

	// Validasi parameter paginasi
	if limit <= 0 {
		limit = 10 // Limit default
	}
	if offset < 0 {
		offset = 0
	}

	// Ambil pesan dari repository berdasarkan thread ID
	messages, err := s.messageRepo.ListByThreadID(ctx, threadID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Konversi []models.Message ke []*models.Message
	messagePtrs := make([]*models.Message, len(messages))
	for i := range messages {
		messagePtrs[i] = &messages[i]
	}

	return messagePtrs, nil
}

// CountByThreadID menghitung jumlah pesan dalam suatu thread
func (s *messageService) CountByThreadID(ctx context.Context, threadID string) (int, error) {
	// Validasi thread ID
	if threadID == "" {
		return 0, fmt.Errorf("thread ID tidak boleh kosong")
	}

	// Hitung jumlah pesan dari repository
	return s.messageRepo.CountByThreadID(ctx, threadID)
}

// UpdateMessage memperbarui pesan yang sudah ada di database
func (s *messageService) UpdateMessage(ctx context.Context, message *models.Message) error {
	// Validasi objek message
	if message == nil {
		return fmt.Errorf("message tidak boleh nil")
	}

	// Validasi field yang wajib diisi
	if message.ID == "" {
		return fmt.Errorf("message ID wajib diisi untuk update")
	}

	// Panggil repository untuk update message
	return s.messageRepo.Update(ctx, message)
}

// DeleteMessage menghapus pesan dari database berdasarkan ID-nya
func (s *messageService) DeleteMessage(ctx context.Context, id string) error {
	// Validasi ID
	if id == "" {
		return fmt.Errorf("message ID tidak boleh kosong")
	}

	// Panggil repository untuk hapus message
	return s.messageRepo.Delete(ctx, id)
}

// Count menghitung jumlah total pesan yang tersimpan
func (s *messageService) Count(ctx context.Context) (int, error) {
	return s.messageRepo.Count(ctx)
}
