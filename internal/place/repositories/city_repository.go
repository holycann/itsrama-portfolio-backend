// Package repositories provides an implementation of repository for city data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/supabase-community/supabase-go"
)

// cityRepository is a concrete implementation of the CityRepository interface
// that manages CRUD operations for city entities in the Supabase database.
type cityRepository struct {
	supabaseClient *supabase.Client // Supabase client for interacting with the database
	table          string           // Name of the table where city data is stored
	column         string           // Columns to be selected in the query
	returning      string           // Type of data returned after an operation
}

// CityRepositoryConfig contains custom configuration for the city repository
// allowing flexibility in setting repository parameters.
type CityRepositoryConfig struct {
	Table     string // Name of the table to be used
	Column    string // Columns to be selected in the query
	Returning string // Type of data to be returned
}

// DefaultConfig returns the default configuration for the city repository
// Useful for providing standard settings if no custom configuration is provided.
func DefaultCityConfig() *CityRepositoryConfig {
	return &CityRepositoryConfig{
		Table:     "cities",  // Default table for cities
		Column:    "*",       // Select all columns
		Returning: "minimal", // Return minimal data
	}
}

// NewCityRepository creates a new instance of the city repository
// with the given configuration and Supabase client.
func NewCityRepository(supabaseClient *supabase.Client, cfg CityRepositoryConfig) CityRepository {
	return &cityRepository{
		supabaseClient: supabaseClient,
		table:          cfg.Table,
		column:         cfg.Column,
		returning:      cfg.Returning,
	}
}

// Create adds a new city to the database
// Accepts context and city object, returns an error if the process fails.
func (r *cityRepository) Create(ctx context.Context, city *models.City) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(city, false, "", "minimal", "").
		ExecuteTo(&city)
	if err != nil {
		return err
	}

	return nil
}

// FindByID searches and returns a city based on its unique ID
// Returns a city object or an error if the city is not found.
func (r *cityRepository) FindByID(ctx context.Context, id string) (*models.City, error) {
	var city *models.City

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&city)
	if err != nil {
		return nil, err
	}

	return city, nil
}

// Update modifies an existing city in the database
// Accepts a modified city object, returns an error if the process fails.
func (r *cityRepository) Update(ctx context.Context, city *models.City) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(city, r.returning, "").
		Eq("id", city.ID).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a city from the database based on its ID
// Returns an error if the deletion process fails.
func (r *cityRepository) Delete(ctx context.Context, id string) error {
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

// List retrieves a list of cities with limit and offset
// Useful for implementing pagination or limiting the number of data retrieved.
func (r *cityRepository) List(ctx context.Context, limit, offset int) ([]models.City, error) {
	var cities []models.City

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		ExecuteTo(&cities)
	if err != nil {
		return nil, err
	}

	return cities, nil
}

// Count calculates the total number of cities stored in the database
// Useful for determining dataset size or for pagination purposes.
func (r *cityRepository) Count(ctx context.Context) (int, error) {
	// Query to count the number of records in the city table
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
