// Package repositories provides an implementation of repository for local story data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/supabase-community/supabase-go"
)

// localStoryRepository is a concrete implementation of the LocalStoryRepository interface
// that manages CRUD operations for local story entities in the Supabase database.
type localStoryRepository struct {
	supabaseClient *supabase.Client // Supabase client for interacting with the database
	table          string           // Name of the table where local story data is stored
	column         string           // Columns to be selected in the query
	returning      string           // Type of data returned after an operation
}

// LocalStoryRepositoryConfig contains custom configuration for the local story repository
// allowing flexibility in setting repository parameters.
type LocalStoryRepositoryConfig struct {
	Table     string // Name of the table to be used
	Column    string // Columns to be selected in the query
	Returning string // Type of data to be returned
}

// DefaultLocalStoryConfig returns the default configuration for the local story repository
// Useful for providing standard settings if no custom configuration is provided.
func DefaultLocalStoryConfig() *LocalStoryRepositoryConfig {
	return &LocalStoryRepositoryConfig{
		Table:     "local_stories", // Default table for local stories
		Column:    "*",             // Select all columns
		Returning: "minimal",       // Return minimal data
	}
}

// NewLocalStoryRepository creates a new instance of the local story repository
// with the given configuration and Supabase client.
func NewLocalStoryRepository(supabaseClient *supabase.Client, cfg LocalStoryRepositoryConfig) LocalStoryRepository {
	return &localStoryRepository{
		supabaseClient: supabaseClient,
		table:          cfg.Table,
		column:         cfg.Column,
		returning:      cfg.Returning,
	}
}

// Create adds a new local story to the database
// Accepts context and local story object, returns an error if the process fails.
func (r *localStoryRepository) Create(ctx context.Context, localStory *models.LocalStory) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(localStory, false, "", "minimal", "").
		ExecuteTo(&localStory)
	if err != nil {
		return err
	}

	return nil
}

// FindByID searches and returns a local story based on its unique ID
// Returns a local story object or an error if the local story is not found.
func (r *localStoryRepository) FindByID(ctx context.Context, id string) (*models.LocalStory, error) {
	var localStory *models.LocalStory

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&localStory)
	if err != nil {
		return nil, err
	}

	return localStory, nil
}

// Update modifies an existing local story in the database
// Accepts a modified local story object, returns an error if the process fails.
func (r *localStoryRepository) Update(ctx context.Context, localStory *models.LocalStory) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(localStory, r.returning, "").
		Eq("id", localStory.ID).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a local story from the database based on its ID
// Returns an error if the deletion process fails.
func (r *localStoryRepository) Delete(ctx context.Context, id string) error {
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

// List retrieves a list of local stories with limit and offset
// Useful for implementing pagination or limiting the number of data retrieved.
func (r *localStoryRepository) List(ctx context.Context, limit, offset int) ([]models.LocalStory, error) {
	var localStories []models.LocalStory

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&localStories)
	if err != nil {
		return nil, err
	}

	return localStories, nil
}

// Count calculates the total number of local stories stored in the database
// Useful for determining dataset size or for pagination purposes.
func (r *localStoryRepository) Count(ctx context.Context) (int, error) {
	// Query to count the number of records in the local stories table
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
