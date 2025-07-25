// Package repositories provides an implementation of repository for discussion thread data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/supabase-community/supabase-go"
)

// threadRepository is a concrete implementation of the ThreadRepository interface
// that manages CRUD operations for thread entities in the Supabase database.
type threadRepository struct {
	supabaseClient *supabase.Client // Supabase client for interacting with the database
	table          string           // Name of the table where thread data is stored
	column         string           // Columns to be selected in the query
	returning      string           // Type of data returned after an operation
}

// ThreadRepositoryConfig contains custom configuration for the thread repository
// allowing flexibility in setting repository parameters.
type ThreadRepositoryConfig struct {
	Table     string // Name of the table to be used
	Column    string // Columns to be selected in the query
	Returning string // Type of data to be returned
}

// DefaultThreadConfig returns the default configuration for the thread repository
// Useful for providing standard settings if no custom configuration is provided.
func DefaultThreadConfig() *ThreadRepositoryConfig {
	return &ThreadRepositoryConfig{
		Table:     "threads", // Default table for threads
		Column:    "*",       // Select all columns
		Returning: "minimal", // Return minimal data
	}
}

// NewThreadRepository creates a new instance of the thread repository
// with the given configuration and Supabase client.
func NewThreadRepository(supabaseClient *supabase.Client, cfg ThreadRepositoryConfig) ThreadRepository {
	return &threadRepository{
		supabaseClient: supabaseClient,
		table:          cfg.Table,
		column:         cfg.Column,
		returning:      cfg.Returning,
	}
}

// FindByTitle searches and returns a thread by title
// Returns a thread object or an error if the thread is not found.
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

// Create adds a new thread to the database
// Accepts context and thread object, returns an error if the process fails.
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

// FindByID searches and returns a thread based on its unique ID
// Returns a thread object or an error if the thread is not found.
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

// Update modifies an existing thread in the database
// Accepts a modified thread object, returns an error if the process fails.
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

// Delete removes a thread from the database based on its ID
// Returns an error if the deletion process fails.
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

// List retrieves a list of threads with limit and offset
// Useful for implementing pagination or limiting the number of data retrieved.
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

// Count calculates the total number of threads stored in the database
// Useful for determining dataset size or for pagination purposes.
func (r *threadRepository) Count(ctx context.Context) (int, error) {
	// Query to count the number of records in the thread table
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", false).
		Execute()
	if err != nil {
		return 0, err
	}

	// Check if the response contains a count
	if count <= 0 {
		return 0, nil
	}

	return int(count), nil
}

// ListByThreadID retrieves a list of threads by thread ID with limit and offset
// Useful for getting specific threads with pagination.
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

// CountByThreadID calculates the total number of threads by thread ID
// Useful for determining the number of threads within a specific thread or for pagination purposes.
func (r *threadRepository) CountByThreadID(ctx context.Context, threadID string) (int, error) {
	// Query to count the number of records in the thread table by thread_id
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", false).
		Eq("thread_id", threadID).
		Execute()
	if err != nil {
		return 0, err
	}

	// Check if the response contains a count
	if count <= 0 {
		return 0, nil
	}

	return int(count), nil
}
