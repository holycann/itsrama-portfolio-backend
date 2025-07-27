package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/models"
)

// MessageType represents the type of message in the discussion system
type MessageType string

// Enum values for MessageType
const (
	DiscussionMessageType MessageType = "discussion"
	AIMessageType         MessageType = "ai"
)

// Message represents a message entity in the discussion system
type Message struct {
	ID        uuid.UUID   `json:"id" db:"id"`                               // Unique ID for the message
	ThreadID  uuid.UUID   `json:"thread_id" db:"thread_id"`                 // ID of the thread the message belongs to
	UserID    uuid.UUID   `json:"user_id" db:"user_id"`                     // ID of the user who sent the message
	Content   string      `json:"content" db:"content" validate:"required"` // Message content
	Type      MessageType `json:"type" db:"type"`                           // Type of message
	CreatedAt time.Time   `json:"created_at" db:"created_at"`               // Message creation time
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`               // Message last update time
}

// RequestMessage is used for message creation or update requests
type RequestMessage struct {
	Message
}

// ResponseMessage is used for returning message data to the client
type ResponseMessage struct {
	Message
	User *models.User `json:"user,omitempty"`
}
