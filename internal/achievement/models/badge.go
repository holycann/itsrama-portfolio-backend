package models

import (
	"time"

	"github.com/google/uuid"
)

// Badge represents a simple achievement badge
type Badge struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" validate:"required"`
	Description string     `json:"description" db:"description"`
	IconURL     string     `json:"icon_url" db:"icon_url"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at" db:"updated_at"`
}

// BadgeCreate represents the data needed to create a new badge
type BadgeCreate struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
}
