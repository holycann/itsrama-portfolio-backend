// Package repositories provides an implementation of repository for event data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"
	"strings"

	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

// eventRepository is a concrete implementation of the EventRepository interface
// that manages CRUD operations for event entities in the Supabase database.
type eventRepository struct {
	supabaseClient *supabase.Client // Supabase client for interacting with the database
	table          string           // Name of the table where event data is stored
	column         string           // Columns to be selected in the query
	returning      string           // Type of data returned after an operation
}

// EventRepositoryConfig contains custom configuration for the event repository
// allowing flexibility in setting repository parameters.
type EventRepositoryConfig struct {
	Table     string // Name of the table to be used
	Column    string // Columns to be selected in the query
	Returning string // Type of data to be returned
}

// DefaultConfig returns the default configuration for the event repository
// Useful for providing standard settings if no custom configuration is provided.
func DefaultEventConfig() *EventRepositoryConfig {
	return &EventRepositoryConfig{
		Table:     "events",  // Default table for events
		Column:    "*",       // Select all columns
		Returning: "minimal", // Return minimal data
	}
}

// NewEventRepository creates a new instance of the event repository
// with the given configuration and Supabase client.
func NewEventRepository(supabaseClient *supabase.Client, cfg EventRepositoryConfig) EventRepository {
	return &eventRepository{
		supabaseClient: supabaseClient,
		table:          cfg.Table,
		column:         cfg.Column,
		returning:      cfg.Returning,
	}
}

// Create adds a new event to the database
// Accepts context and event object, returns an error if the process fails.
func (r *eventRepository) Create(ctx context.Context, event *models.Event) error {
	_, err := r.supabaseClient.
		From(r.table).
		Insert(event, false, "", "minimal", "").
		ExecuteTo(&event)
	if err != nil {
		return err
	}

	return nil
}

// FindByID searches and returns a event based on its unique ID
// Returns a event object or an error if the event is not found.
func (r *eventRepository) FindByID(ctx context.Context, id string) (*models.Event, error) {
	var event *models.Event

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// Update modifies an existing event in the database
// Accepts a modified event object, returns an error if the process fails.
func (r *eventRepository) Update(ctx context.Context, event *models.Event) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(event, r.returning, "").
		Eq("id", event.ID).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

// Delete removes a event from the database based on its ID
// Returns an error if the deletion process fails.
func (r *eventRepository) Delete(ctx context.Context, id string) error {
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

// List retrieves a list of events with limit and offset
// Useful for implementing pagination or limiting the number of data retrieved.
func (r *eventRepository) List(ctx context.Context, limit, offset int) ([]models.Event, error) {
	var events []models.Event

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Range(offset, offset+limit-1, "").
		ExecuteTo(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Count calculates the total number of events stored in the database
// Useful for determining dataset size or for pagination purposes.
func (r *eventRepository) Count(ctx context.Context) (int, error) {
	// Query to count the number of records in the event table
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

// GetTrendingEvent retrieves a list of trending events based on the highest views
func (r *eventRepository) ListTrendingEvent(ctx context.Context, limit int) ([]models.Event, error) {
	var events []models.Event

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Order("views", &postgrest.OrderOpts{Ascending: false}).
		Limit(limit, "").
		ExecuteTo(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Search searches for events based on a query string in the name or description
func (r *eventRepository) Search(ctx context.Context, query string, limit, offset int) ([]models.Event, error) {
	var events []models.Event

	// Escape % and _ in query to prevent wildcard injection
	escapedQuery := strings.ReplaceAll(strings.ReplaceAll(query, "%", "\\%"), "_", "\\_")
	likeQuery := "%" + escapedQuery + "%"
	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Or("name.ilike."+likeQuery+",description.ilike."+likeQuery, "").
		Range(offset, offset+limit-1, "").
		ExecuteTo(&events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) UpdateViews(ctx context.Context, id string) string {
	err := r.supabaseClient.
		Rpc("increment_or_create_event_views", "", map[string]interface{}{
			"event_id": id,
		})
	return err
}

// ListRelatedEvents retrieves a list of events related to a specific event
func (r *eventRepository) ListRelatedEvents(ctx context.Context, eventID string, limit int) ([]models.Event, error) {
	var events []models.Event

	_, err := r.supabaseClient.
		From(r.table).
		Select(r.column, "", false).
		Neq("id", eventID). // Exclude the current event
		Limit(limit, "").
		ExecuteTo(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}
