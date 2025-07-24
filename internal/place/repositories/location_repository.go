// Package repositories provides an implementation of repository for location data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/supabase-community/supabase-go"
)

// locationRepository is a concrete implementation of the LocationRepository interface
// that manages CRUD operations for location entities in the Supabase database.
type locationRepository struct {
	supabaseClient *supabase.Client // Supabase client for interacting with the database
	table          string           // Name of the table where location data is stored
	column         string           // Columns to be selected in the query
	returning      string           // Type of data returned after an operation
}

// LocationRepositoryConfig contains custom configuration for the location repository
// allowing flexibility in setting repository parameters.
type LocationRepositoryConfig struct {
	Table     string // Name of the table to be used
	Column    string // Columns to be selected in the query
	Returning string // Type of data to be returned
}

// DefaultConfig returns the default configuration for the location repository
// Useful for providing standard settings if no custom configuration is provided.
func DefaultLocationConfig() *LocationRepositoryConfig {
	return &LocationRepositoryConfig{
		Table:     "locations", // Default table for locations
		Column:    "*",         // Select all columns
		Returning: "minimal",   // Return minimal data
	}
}

// NewLocationRepository creates a new instance of the location repository
// with the given configuration and Supabase client.
func NewLocationRepository(supabaseClient *supabase.Client, cfg LocationRepositoryConfig) LocationRepository {
	return &locationRepository{
		supabaseClient: supabaseClient,
		table:          cfg.Table,
		column:         cfg.Column,
		returning:      cfg.Returning,
	}
}

// Create adds a new location to the database
// Accepts context and location object, returns an error if the process fails.
func (r *locationRepository) Create(ctx context.Context, location *models.Location) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(location, false, "", "minimal", "").
		ExecuteTo(&location)
	if err != nil {
		return err
	}

	return nil
}

// FindByID searches and returns a location based on its unique ID
// Returns a location object or an error if the location is not found.
func (r *locationRepository) FindByID(ctx context.Context, id string) (*models.Location, error) {
	var location *models.Location

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&location)
	if err != nil {
		return nil, err
	}

	return location, nil
}

// Update modifies an existing location in the database
// Accepts a modified location object, returns an error if the process fails.
func (r *locationRepository) Update(ctx context.Context, location *models.Location) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(location, r.returning, "").
		Eq("id", location.ID).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a location from the database based on its ID
// Returns an error if the deletion process fails.
func (r *locationRepository) Delete(ctx context.Context, id string) error {
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

// List retrieves a list of locations with limit and offset
// Useful for implementing pagination or limiting the number of data retrieved.
func (r *locationRepository) List(ctx context.Context, limit, offset int) ([]models.Location, error) {
	var locations []models.Location

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

// Count calculates the total number of locations stored in the database
// Useful for determining dataset size or for pagination purposes.
func (r *locationRepository) Count(ctx context.Context) (int, error) {
	// Query to count the number of records in the location table
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
