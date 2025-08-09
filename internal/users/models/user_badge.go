package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/achievement/models"
)

// UserBadge represents the association between a user and their earned badges
// @Description Comprehensive model for tracking user achievements and badge acquisitions
// @Description Provides a structured representation of user-badge relationships with metadata
// @Tags User Achievements
type UserBadge struct {
	// Associated user ID
	// @Description Globally unique UUID of the user who earned the badge
	// @Description Links the badge to a specific user account
	// @Example "user_123"
	// @Format uuid
	UserID uuid.UUID `json:"user_id" db:"user_id" validate:"required" example:"user_123"`

	// Badge identifier
	// @Description Unique identifier of the badge earned by the user
	// @Description Represents the specific achievement milestone
	// @Example "explorer"
	// @Format uuid
	BadgeID uuid.UUID `json:"badge_id" db:"badge_id" validate:"required" example:"explorer"`

	// Timestamp when the badge was earned
	// @Description Precise timestamp of when the user acquired the badge
	// @Description Helps track the user's achievement progression
	// @Format date-time
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Timestamp of the last user badge update
	// @Description Precise timestamp of the last modification to the user badge record
	// @Description Indicates any changes or updates to the badge status
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// UserBadgeDTO represents the data transfer object for user badge information
// @Description Comprehensive data transfer object for user badge details
// @Description Used for API responses to provide rich user achievement information
// @Tags User Achievements
type UserBadgeDTO struct {
	// Badge identifier
	// @Description Unique identifier of the badge earned by the user
	// @Example "explorer"
	// @Format uuid
	BadgeID uuid.UUID `json:"badge_id" example:"explorer"`

	// Full badge details
	// @Description Comprehensive information about the earned badge
	// @Description Provides context and details of the specific achievement
	Badge *models.Badge `json:"badge,omitempty"`

	// Timestamp when the badge was earned
	// @Description Precise timestamp of when the user acquired the badge
	// @Description Helps track the user's achievement progression
	// @Format date-time
	CreatedAt *time.Time `json:"created_at"`

	// Timestamp of the last user badge update
	// @Description Precise timestamp of the last modification to the user badge record
	// @Description Indicates any changes or updates to the badge status
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// ToDTO converts a UserBadge to a UserBadgeDTO
// @Description Transforms a UserBadge model into a lightweight data transfer object
// @Description Useful for API responses and data serialization
// @Return UserBadgeDTO Converted user badge data transfer object
func (ub *UserBadge) ToDTO() UserBadgeDTO {
	return UserBadgeDTO{
		BadgeID:   ub.BadgeID,
		CreatedAt: ub.CreatedAt,
		UpdatedAt: ub.UpdatedAt,
	}
}

// UserBadgePayload represents the payload for assigning or removing a badge from a user
// @Description Structured payload for user badge management operations
// @Description Supports creating or removing user badge associations
// @Tags User Achievements
type UserBadgePayload struct {
	// Associated user ID
	// @Description Unique identifier of the user receiving the badge
	// @Description Must be a valid user account
	// @Example "user_123"
	// @Format uuid
	UserID uuid.UUID `json:"user_id" validate:"required" example:"user_123"`

	// Badge identifier
	// @Description Unique identifier of the badge to be assigned or removed
	// @Description Must be a valid badge in the system
	// @Example "explorer"
	// @Format uuid
	BadgeID uuid.UUID `json:"badge_id" validate:"required" example:"explorer"`
}
