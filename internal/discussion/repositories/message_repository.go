// Package repositories provides an implementation of repository for message data management
// using Supabase as the data storage backend.
package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/supabase-community/supabase-go"
)

// messageRepository is a concrete implementation of the MessageRepository interface
// that manages CRUD operations for message entities in the Supabase database.
type messageRepository struct {
	supabaseClient *supabase.Client // Supabase client for interacting with the database
	table          string           // Name of the table where message data is stored
	column         string           // Columns to be selected in the query
	returning      string           // Type of data returned after an operation
}

// MessageRepositoryConfig contains custom configuration for the message repository
// allowing flexibility in setting repository parameters.
type MessageRepositoryConfig struct {
	Table     string // Name of the table to be used
	Column    string // Columns to be selected in the query
	Returning string // Type of data to be returned
}

// DefaultMessageConfig returns the default configuration for the message repository
// Useful for providing standard settings if no custom configuration is provided.
func DefaultMessageConfig() *MessageRepositoryConfig {
	return &MessageRepositoryConfig{
		Table:     "messages", // Default table for messages
		Column:    "*",        // Select all columns
		Returning: "minimal",  // Return minimal data
	}
}

// NewMessageRepository creates a new instance of the message repository
// with the given configuration and Supabase client.
func NewMessageRepository(supabaseClient *supabase.Client, cfg MessageRepositoryConfig) MessageRepository {
	return &messageRepository{
		supabaseClient: supabaseClient,
		table:          cfg.Table,
		column:         cfg.Column,
		returning:      cfg.Returning,
	}
}

// Create adds a new message to the database
// Accepts context and message object, returns an error if the process fails.
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

// FindByID searches and returns a message based on its unique ID
// Returns a message object or an error if the message is not found.
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

// Update modifies an existing message in the database
// Accepts a modified message object, returns an error if the process fails.
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

// Delete removes a message from the database based on its ID
// Returns an error if the deletion process fails.
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

// List retrieves a list of messages with limit and offset
// Useful for implementing pagination or limiting the number of data retrieved.
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

// Count calculates the total number of messages stored in the database
// Useful for determining dataset size or for pagination purposes.
func (r *messageRepository) Count(ctx context.Context) (int, error) {
	// Query to count the number of records in the message table
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

// CountByThreadID counts the number of messages in a thread
// Accepts a thread ID and returns the count of messages associated with that thread.
func (r *messageRepository) CountByThreadID(ctx context.Context, threadID string) (int, error) {
	// Query to count messages by thread ID
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

// ListByThreadID retrieves a list of messages by thread ID with limit and offset
// Useful for getting paginated messages within a thread.
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
