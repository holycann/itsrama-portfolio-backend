package models

import (
	"time"

	"github.com/google/uuid"
	userModels "github.com/holycann/cultour-backend/internal/users/models"
)

// Thread represents a discussion thread in the system
// @Description Comprehensive model for tracking discussion threads associated with events
// @Description Provides a structured representation of threads with user context and status management
// @Tags Discussion Threads
type Thread struct {
	// Unique identifier for the thread
	// @Description Globally unique UUID for the thread, generated automatically
	// @Description Serves as the primary key and reference for the thread
	// @Example "thread_123"
	// @Format uuid
	ID uuid.UUID `json:"id" db:"id" validate:"required"`

	// Reference to the related event
	// @Description Unique identifier linking the thread to a specific event
	// @Description Enables contextual discussion and event-specific communication
	// @Example "event_456"
	// @Format uuid
	EventID uuid.UUID `json:"event_id" db:"event_id" validate:"required"`

	// Reference to the user who created the thread
	// @Description Unique identifier of the thread creator
	// @Description Provides attribution and tracks thread ownership
	// @Example "user_789"
	// @Format uuid
	CreatorID uuid.UUID `json:"creator_id" db:"creator_id" validate:"required"`

	// Thread status
	// @Description Current state of the discussion thread
	// @Description Allows lifecycle management and access control
	// @Example "active"
	// @Enum active,closed,archived
	Status string `json:"status" db:"status" validate:"required,oneof=active closed archived" example:"active"`

	// Timestamp when the thread was created
	// @Description Precise timestamp of thread creation in UTC
	// @Description Helps track thread lifecycle and origin
	// @Format date-time
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Timestamp when the thread was last updated
	// @Description Precise timestamp of the last modification to the thread in UTC
	// @Description Indicates when thread details or status were last changed
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// ThreadDTO represents the data transfer object for thread information
// @Description Comprehensive data transfer object for thread details
// @Description Used for API responses to provide rich thread information with related entities
// @Tags Discussion Threads
type ThreadDTO struct {
	// Unique identifier for the thread
	// @Description Globally unique UUID for the thread
	// @Example "thread_123"
	// @Format uuid
	ID uuid.UUID `json:"id"`

	// Reference to the related event
	// @Description Unique identifier linking the thread to a specific event
	// @Description Enables contextual discovery and event-based communication
	// @Example "event_456"
	// @Format uuid
	EventID uuid.UUID `json:"event_id"`

	// Thread status
	// @Description Current state of the discussion thread
	// @Description Indicates thread accessibility and lifecycle stage
	// @Example "active"
	// @Enum active,closed,archived
	Status string `json:"status"`

	// Timestamp when the thread was created
	// @Description Precise timestamp of thread creation in UTC
	// @Format date-time
	CreatedAt *time.Time `json:"created_at"`

	// Timestamp when the thread was last updated
	// @Description Precise timestamp of the last modification to the thread in UTC
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at"`

	// Creator's profile details
	// @Description Comprehensive information about the thread creator
	// @Description Provides context and attribution for the thread
	Creator *userModels.User `json:"creator,omitempty"`

	// List of participants in the thread
	// @Description Detailed list of users participating in the discussion
	// @Description Helps track engagement and thread membership
	Participants []Participant `json:"discussion_participants,omitempty"`
}

// ToDTO converts a Thread to a ThreadDTO
// @Description Transforms a Thread model into a lightweight data transfer object
// @Description Useful for API responses and data serialization
// @Return ThreadDTO Converted thread data transfer object
func (t *Thread) ToDTO() ThreadDTO {
	return ThreadDTO{
		ID:        t.ID,
		EventID:   t.EventID,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// CreateThread represents the payload for creating or updating a thread
// @Description Structured payload for thread creation and update operations
// @Description Supports flexible input for thread initialization and modification
// @Tags Discussion Threads
type CreateThread struct {
	// Unique identifier for the thread (optional for creation)
	// @Description Optional UUID for the thread during creation or update
	// @Description Used to identify specific threads during updates
	// @Example "thread_123"
	// @Format uuid
	ID uuid.UUID `json:"id,omitempty" validate:"omitempty"`

	// Reference to the related event
	// @Description Unique identifier linking the thread to a specific event
	// @Description Required for thread context and organization
	// @Example "event_456"
	// @Format uuid
	EventID uuid.UUID `json:"event_id" validate:"required"`

	// Reference to the user who created the thread
	// @Description Optional user ID for thread creator
	// @Description Typically set automatically during thread creation
	// @Example "user_789"
	// @Format uuid
	CreatorID uuid.UUID `json:"creator_id,omitempty" validate:"omitempty"`

	// Thread status
	// @Description Optional thread status specification
	// @Description Defaults to 'active' if not specified
	// @Example "active"
	// @Enum active,closed,archived
	Status string `json:"status,omitempty" validate:"omitempty,oneof=active closed archived"`
}
