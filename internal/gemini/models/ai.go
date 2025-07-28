package models

import (
	"time"

	"github.com/google/uuid"
)

// AiMessage represents a message in an AI conversation
type AiMessage struct {
	Role    string `json:"role"`    // 'user', 'assistant', or 'system'
	Content string `json:"content"` // Message content
}

// AiResponse represents an AI query response
type AiResponse struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// AiRequest represents an incoming AI query request
type AiRequest struct {
	Query   string `json:"query"`
	Context string `json:"context,omitempty"`
}
