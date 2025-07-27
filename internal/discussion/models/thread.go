package models

import (
	"time"

	"github.com/google/uuid"
)

// Thread represents a discussion thread entity in the system
type Thread struct {
	ID        uuid.UUID `json:"id" db:"id"`                          // Unique ID for the thread
	Title     string    `json:"title" db:"title"`                    // Thread title
	EventID   uuid.UUID `json:"event_id" db:"event_id"`              // Reference to related event
	Status    string    `json:"status" db:"status" example:"active"` // Thread status
	CreatedAt time.Time `json:"created_at" db:"created_at"`          // Thread creation time
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`          // Thread last update time
}

// RequestThread is used for thread creation or update requests
type RequestThread struct {
	Thread
}

// ResponseThread is used for returning thread data to the client
type ResponseThread struct {
	Thread
}
