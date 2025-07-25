package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/supabase-community/supabase-go"
)

// provinceRepository is a concrete implementation of the ProvinceRepository interface
// that manages CRUD operations for province entities in the Supabase database.
type provinceRepository struct {
	supabaseClient *supabase.Client // Supabase client for interacting with the database
	table          string           // Name of the table where province data is stored
	column         string           // Columns to be selected in the query
	returning      string           // Type of data returned after an operation
}

// ProvinceRepositoryConfig contains custom configuration for the province repository
// allowing flexibility in setting repository parameters.
type ProvinceRepositoryConfig struct {
	Table     string // Name of the table to be used
	Column    string // Columns to be selected in the query
	Returning string // Type of data to be returned
}

// DefaultProvinceConfig returns the default configuration for the province repository
// Useful for providing standard settings if no custom configuration is provided.
func DefaultProvinceConfig() *ProvinceRepositoryConfig {
	return &ProvinceRepositoryConfig{
		Table:     "provinces", // Default table for provinces
		Column:    "*",         // Select all columns
		Returning: "minimal",   // Return minimal data
	}
}

// NewProvinceRepository creates a new instance of the province repository
// with the given configuration and Supabase client.
func NewProvinceRepository(supabaseClient *supabase.Client, cfg ProvinceRepositoryConfig) ProvinceRepository {
	return &provinceRepository{
		supabaseClient: supabaseClient,
		table:          cfg.Table,
		column:         cfg.Column,
		returning:      cfg.Returning,
	}
}

// Create adds a new province to the database
// Accepts context and province object, returns an error if the process fails.
func (r *provinceRepository) Create(ctx context.Context, province *models.Province) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(province, false, "", "minimal", "").
		ExecuteTo(&province)
	if err != nil {
		return err
	}

	return nil
}

// FindByID searches and returns a province based on its unique ID
// Returns a province object or an error if the province is not found.
func (r *provinceRepository) FindByID(ctx context.Context, id string) (*models.Province, error) {
	var province *models.Province

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&province)
	if err != nil {
		return nil, err
	}

	return province, nil
}

// Update modifies an existing province in the database
// Accepts a modified province object, returns an error if the process fails.
func (r *provinceRepository) Update(ctx context.Context, province *models.Province) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(province, r.returning, "").
		Eq("id", province.ID).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a province from the database based on its ID
// Returns an error if the deletion process fails.
func (r *provinceRepository) Delete(ctx context.Context, id string) error {
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

// List retrieves a list of provinces with limit and offset
// Useful for implementing pagination or limiting the number of data retrieved.
func (r *provinceRepository) List(ctx context.Context, limit, offset int) ([]models.Province, error) {
	var provinces []models.Province

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		ExecuteTo(&provinces)
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

// Count calculates the total number of provinces stored in the database
// Useful for determining dataset size or for pagination purposes.
func (r *provinceRepository) Count(ctx context.Context) (int, error) {
	// Query to count the number of records in the province table
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
