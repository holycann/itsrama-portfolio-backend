// Package repositories menyediakan implementasi repository untuk manajemen data pesan (message)
// menggunakan Supabase sebagai backend penyimpanan data.
package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/supabase-community/supabase-go"
)

// messageRepository adalah implementasi konkret dari interface MessageRepository
// yang mengelola operasi CRUD untuk entitas message pada database Supabase.
type messageRepository struct {
	supabaseClient *supabase.Client // Supabase client untuk berinteraksi dengan database
	table          string           // Nama tabel tempat data message disimpan
	column         string           // Kolom-kolom yang dipilih dalam query
	returning      string           // Tipe data yang dikembalikan setelah operasi
}

// MessageRepositoryConfig berisi konfigurasi kustom untuk repository message
// memungkinkan fleksibilitas dalam pengaturan parameter repository.
type MessageRepositoryConfig struct {
	Table     string // Nama tabel yang digunakan
	Column    string // Kolom-kolom yang dipilih dalam query
	Returning string // Tipe data yang dikembalikan
}

// DefaultMessageConfig mengembalikan konfigurasi default untuk repository message
// Berguna untuk memberikan pengaturan standar jika tidak ada konfigurasi kustom yang diberikan.
func DefaultMessageConfig() *MessageRepositoryConfig {
	return &MessageRepositoryConfig{
		Table:     "messages", // Tabel default untuk message
		Column:    "*",        // Pilih semua kolom
		Returning: "minimal",  // Kembalikan data minimal
	}
}

// NewMessageRepository membuat instance baru dari repository message
// dengan konfigurasi dan Supabase client yang diberikan.
func NewMessageRepository(supabaseClient *supabase.Client, cfg MessageRepositoryConfig) MessageRepository {
	return &messageRepository{
		supabaseClient: supabaseClient,
		table:          cfg.Table,
		column:         cfg.Column,
		returning:      cfg.Returning,
	}
}

// Create menambahkan message baru ke database
// Menerima context dan objek message, mengembalikan error jika proses gagal.
func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(message, false, "", "minimal", "").
		ExecuteTo(&message)
	if err != nil {
		return err
	}

	return nil
}

// FindByID mencari dan mengembalikan message berdasarkan ID uniknya
// Mengembalikan objek message atau error jika message tidak ditemukan.
func (r *messageRepository) FindByID(ctx context.Context, id string) (*models.Message, error) {
	var message *models.Message

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Update memodifikasi message yang sudah ada di database
// Menerima objek message yang sudah dimodifikasi, mengembalikan error jika proses gagal.
func (r *messageRepository) Update(ctx context.Context, message *models.Message) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(message, r.returning, "").
		Eq("id", message.ID).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

// Delete menghapus message dari database berdasarkan ID-nya
// Mengembalikan error jika proses penghapusan gagal.
func (r *messageRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete(r.returning, "").
		Eq("id", id).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

// List mengambil daftar message dengan limit dan offset
// Berguna untuk implementasi paginasi atau membatasi jumlah data yang diambil.
func (r *messageRepository) List(ctx context.Context, limit, offset int) ([]models.Message, error) {
	var messages []models.Message

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// Count menghitung total jumlah message yang tersimpan di database
// Berguna untuk mengetahui ukuran dataset atau untuk keperluan paginasi.
func (r *messageRepository) Count(ctx context.Context) (int, error) {
	// Query untuk menghitung jumlah record pada tabel message
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", false).
		Execute()
	if err != nil {
		return 0, err
	}

	// Cek apakah response berisi count
	if count <= 0 {
		return 0, nil
	}

	return int(count), nil
}

// CountByThreadID menghitung jumlah pesan dalam suatu thread
// Menerima ID thread dan mengembalikan jumlah pesan yang terkait dengan thread tersebut.
func (r *messageRepository) CountByThreadID(ctx context.Context, threadID string) (int, error) {
	// Query untuk menghitung jumlah pesan berdasarkan thread ID
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", false).
		Eq("thread_id", threadID).
		Execute()
	if err != nil {
		return 0, err
	}

	// Cek apakah response berisi count
	if count <= 0 {
		return 0, nil
	}

	return int(count), nil
}

// ListByThreadID mengambil daftar pesan berdasarkan ID thread dengan limit dan offset
// Berguna untuk mendapatkan pesan-pesan dalam suatu thread dengan paginasi.
func (r *messageRepository) ListByThreadID(ctx context.Context, threadID string, limit, offset int) ([]models.Message, error) {
	var messages []models.Message

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("thread_id", threadID).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
