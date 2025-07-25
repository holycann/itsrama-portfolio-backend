package models

import (
	"time"

	"github.com/holycann/cultour-backend/internal/users/models"
)

// Message represents a message entity in the discussion system
type Message struct {
	ID        string    `json:"id" db:"id" example:"msg_12345"`                             // Unique ID for the message, example: "msg_12345"
	ThreadID  string    `json:"thread_id" db:"thread_id" example:"thread_12345"`            // ID of the thread the message belongs to
	UserID    string    `json:"user_id" db:"user_id" example:"user_67890"`                  // ID of the user who sent the message
	Content   string    `json:"content" db:"content" example:"This is the message content"` // Message content
	CreatedAt time.Time `json:"created_at" db:"created_at"`                                 // Message creation time
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
