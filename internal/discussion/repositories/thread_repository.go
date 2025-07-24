// Package repositories menyediakan implementasi repository untuk manajemen data thread diskusi
// menggunakan Supabase sebagai backend penyimpanan data.
package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/supabase-community/supabase-go"
)

// threadRepository adalah implementasi konkret dari interface ThreadRepository
// yang mengelola operasi CRUD untuk entitas thread pada database Supabase.
type threadRepository struct {
	supabaseClient *supabase.Client // Supabase client untuk berinteraksi dengan database
	table          string           // Nama tabel tempat data thread disimpan
	column         string           // Kolom-kolom yang dipilih dalam query
	returning      string           // Tipe data yang dikembalikan setelah operasi
}

// ThreadRepositoryConfig berisi konfigurasi kustom untuk repository thread
// memungkinkan fleksibilitas dalam pengaturan parameter repository.
type ThreadRepositoryConfig struct {
	Table     string // Nama tabel yang digunakan
	Column    string // Kolom-kolom yang dipilih dalam query
	Returning string // Tipe data yang dikembalikan
}

// DefaultThreadConfig mengembalikan konfigurasi default untuk repository thread
// Berguna untuk memberikan pengaturan standar jika tidak ada konfigurasi kustom yang diberikan.
func DefaultThreadConfig() *ThreadRepositoryConfig {
	return &ThreadRepositoryConfig{
		Table:     "threads", // Tabel default untuk thread
		Column:    "*",       // Pilih semua kolom
		Returning: "minimal", // Kembalikan data minimal
	}
}

// NewThreadRepository membuat instance baru dari repository thread
// dengan konfigurasi dan Supabase client yang diberikan.
func NewThreadRepository(supabaseClient *supabase.Client, cfg ThreadRepositoryConfig) ThreadRepository {
	return &threadRepository{
		supabaseClient: supabaseClient,
		table:          cfg.Table,
		column:         cfg.Column,
		returning:      cfg.Returning,
	}
}

// FindByTitle mencari dan mengembalikan thread berdasarkan judul
// Mengembalikan objek thread atau error jika thread tidak ditemukan.
func (r *threadRepository) FindByTitle(ctx context.Context, title string) (*models.Thread, error) {
	var thread *models.Thread

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Eq("title", title).
		Single().
		ExecuteTo(&thread)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

// Create menambahkan thread baru ke database
// Menerima context dan objek thread, mengembalikan error jika proses gagal.
func (r *threadRepository) Create(ctx context.Context, thread *models.Thread) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(thread, false, "", "minimal", "").
		ExecuteTo(&thread)
	if err != nil {
		return err
	}

	return nil
}

// FindByID mencari dan mengembalikan thread berdasarkan ID uniknya
// Mengembalikan objek thread atau error jika thread tidak ditemukan.
func (r *threadRepository) FindByID(ctx context.Context, id string) (*models.Thread, error) {
	var thread *models.Thread

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&thread)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

// Update memodifikasi thread yang sudah ada di database
// Menerima objek thread yang sudah dimodifikasi, mengembalikan error jika proses gagal.
func (r *threadRepository) Update(ctx context.Context, thread *models.Thread) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(thread, r.returning, "").
		Eq("id", thread.ID).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

// Delete menghapus thread dari database berdasarkan ID-nya
// Mengembalikan error jika proses penghapusan gagal.
func (r *threadRepository) Delete(ctx context.Context, id string) error {
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

// List mengambil daftar thread dengan limit dan offset
// Berguna untuk implementasi paginasi atau membatasi jumlah data yang diambil.
func (r *threadRepository) List(ctx context.Context, limit, offset int) ([]models.Thread, error) {
	var threads []models.Thread

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

// Count menghitung total jumlah thread yang tersimpan di database
// Berguna untuk mengetahui ukuran dataset atau untuk keperluan paginasi.
func (r *threadRepository) Count(ctx context.Context) (int, error) {
	// Query untuk menghitung jumlah record pada tabel thread
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

// ListByThreadID mengambil daftar thread berdasarkan ID thread dengan limit dan offset
// Berguna untuk mendapatkan thread spesifik dengan paginasi.
func (r *threadRepository) ListByThreadID(ctx context.Context, threadID string, limit, offset int) ([]models.Thread, error) {
	var threads []models.Thread

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("thread_id", threadID).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

// CountByThreadID menghitung total jumlah thread berdasarkan ID thread
// Berguna untuk mengetahui jumlah thread dalam suatu thread atau untuk keperluan paginasi.
func (r *threadRepository) CountByThreadID(ctx context.Context, threadID string) (int, error) {
	// Query untuk menghitung jumlah record pada tabel thread berdasarkan thread_id
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
