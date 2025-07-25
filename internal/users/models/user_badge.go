package models

import "time"

// UserBadge represents the badge earned by a user
// @Description Badge information associated with a user
type UserBadge struct {
	// Unique identifier for the user badge
	// @example "user_badge_123"
	ID string `json:"id" db:"id" example:"user_badge_123"`

	// Associated user ID
	// @example "user_123"
	UserID string `json:"user_id" db:"user_id" validate:"required" example:"user_123"`

	// Badge identifier
	// @example "explorer"
	BadgeID string `json:"badge_id" db:"badge_id" validate:"required" example:"explorer"`

	// Badge name
	// @example "Penjelajah"
	BadgeName string `json:"badge_name" db:"badge_name" validate:"required" example:"Penjelajah"`

	// Badge description
	// @example "Explored multiple cultural events"
	BadgeDescription string `json:"badge_description" db:"badge_description"`

	// Badge icon URL
	// @example "https://example.com/badges/explorer.png"
	BadgeIconUrl string `json:"badge_icon_url" db:"badge_icon_url" format:"uri"`

	// Timestamp when the badge was earned
	EarnedAt *time.Time `json:"earned_at" db:"earned_at"`
}

// UserBadgeCreate represents the payload for creating a new user badge
// @Description Payload for assigning a badge to a user
type UserBadgeCreate struct {
	// Associated user ID
	// @example "user_123"
	UserID string `json:"user_id" validate:"required" example:"user_123"`

	// Badge identifier
	// @example "explorer"
	BadgeID string `json:"badge_id" validate:"required" example:"explorer"`
}

// UserBadgeSearch represents the search criteria for user badges
// @Description Search parameters for retrieving user badges
type UserBadgeSearch struct {
	// User ID to filter badges
	// @example "user_123"
	UserID string `json:"user_id,omitempty"`

	// Badge ID to filter
	// @example "explorer"
	BadgeID string `json:"badge_id,omitempty"`

	// Pagination limit
	// @example 10
	Limit int `json:"limit,omitempty" default:"10"`

	// Pagination offset
	// @example 0
	Offset int `json:"offset,omitempty" default:"0"`
}
