package models

import (
	"time"

	"github.com/google/uuid"
	userModels "github.com/holycann/cultour-backend/internal/users/models"
)

// Participant represents a user's participation in a discussion thread
// @Description Detailed participant entry with thread and user references
type Participant struct {
	// Unique identifier for the thread
	// @example "thread_123"
	ThreadID uuid.UUID `json:"thread_id" db:"thread_id" validate:"required"`

	// Reference to the user participating in the thread
	// @example "user_789"
	UserID uuid.UUID `json:"user_id" db:"user_id" validate:"required"`

	// Timestamp when the user joined the thread
	JoinedAt *time.Time `json:"joined_at" db:"joined_at"`

	// Timestamp when the participant entry was last updated
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// ParticipantDTO represents a detailed data transfer object for participants
// @Description Comprehensive participant data transfer object with additional details
type ParticipantDTO struct {
	// Unique identifier for the thread
	// @example "thread_123"
	ThreadID uuid.UUID `json:"thread_id"`

	// Reference to the user participating in the thread
	// @example "user_789"
	UserID uuid.UUID `json:"user_id"`

	// Timestamp when the user joined the thread
	JoinedAt *time.Time `json:"joined_at"`

	// Timestamp when the participant entry was last updated
	UpdatedAt *time.Time `json:"updated_at"`

	// User's profile details
	User *userModels.User `json:"user,omitempty"`
}

// ToDTO converts a Participant to a ParticipantDTO
func (p *Participant) ToDTO() ParticipantDTO {
	return ParticipantDTO{
		ThreadID:  p.ThreadID,
		UserID:    p.UserID,
		JoinedAt:  p.JoinedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// CreateParticipant represents the payload for creating or updating a participant
// @Description Data transfer object for participant creation or update
type CreateParticipant struct {
	// Unique identifier for the thread
	// @example "thread_123"
	ThreadID uuid.UUID `json:"thread_id" validate:"required"`

	// Reference to the user participating in the thread
	// @example "user_789"
	UserID uuid.UUID `json:"user_id" validate:"required"`
}
