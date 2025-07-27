package models

import (
	"time"

	"github.com/google/uuid"
)

// UserBadge represents the badge earned by a user
// @Description Badge information associated with a user
type UserBadge struct {
	// Unique identifier for the user badge
	// @example "user_badge_123"
	ID uuid.UUID `json:"id" db:"id"`

	// Associated user ID
	// @example "user_123"
	UserID uuid.UUID `json:"user_id" db:"user_id" validate:"required"`

	// Badge identifier
	// @example "explorer"
	BadgeID uuid.UUID `json:"badge_id" db:"badge_id" validate:"required"`

	// Timestamp when the badge was earned
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserBadgeCreate represents the payload for creating a new user badge
// @Description Payload for assigning a badge to a user
type UserBadgeCreate struct {
	// Associated user ID
	// @example "user_123"
	UserID uuid.UUID `json:"user_id" validate:"required"`

	// Badge identifier
	// @example "explorer"
	BadgeID uuid.UUID `json:"badge_id" validate:"required"`
}

// UserBadgeSearch represents the search criteria for user badges
// @Description Search parameters for retrieving user badges
type UserBadgeSearch struct {
	// User ID to filter badges
	// @example "user_123"
	UserID uuid.UUID `json:"user_id,omitempty"`

	// Badge ID to filter
	// @example "explorer"
	BadgeID uuid.UUID `json:"badge_id,omitempty"`

	// Pagination limit
	// @example 10
	Limit int `json:"limit,omitempty" default:"10"`

	// Pagination offset
	// @example 0
	Offset int `json:"offset,omitempty" default:"0"`
}
