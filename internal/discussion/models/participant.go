package models

import (
	// "time"

	"github.com/google/uuid"
)

// Participant represents a user's participation in a discussion thread
type Participant struct {
	ThreadID uuid.UUID `json:"thread_id" db:"thread_id"` // Reference to the thread
	UserID   uuid.UUID `json:"user_id" db:"user_id"`     // Reference to the user
	// JoinedAt  time.Time `json:"joined_at" db:"joined_at"`   // Time when the user joined the thread
	// UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // Last update time for the participant entry
}

// RequestParticipant is used for creating or updating participant entries
type RequestParticipant struct {
	Participant
}

// ResponseParticipant is used for returning participant data to the client
type ResponseParticipant struct {
	Participant
}
