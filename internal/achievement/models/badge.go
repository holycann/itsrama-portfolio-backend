package models

import (
	"time"

	"github.com/google/uuid"
)

// Badge represents an achievement badge in the system
// @Description Comprehensive model for tracking and managing user achievements
// @Description Provides a structured representation of badges with unique identifiers, metadata, and timestamps
// @Tags Achievements
type Badge struct {
	// Unique identifier for the badge
	// @Description Globally unique UUID for the badge, generated automatically
	// @Example "550e8400-e29b-41d4-a716-446655440000"
	// @Format uuid
	ID uuid.UUID `json:"id" db:"id"`

	// Badge name
	// @Description Human-readable, distinctive name for the badge
	// @Description Represents the achievement or milestone
	// @Example "Master Explorer"
	// @Required true
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name" db:"name" validate:"required,min=2,max=100"`

	// Badge description
	// @Description Detailed explanation of how to earn the badge
	// @Description Provides context and motivation for users to achieve this badge
	// @Example "Discovered and visited 50 unique locations across the platform"
	// @MaxLength 500
	Description string `json:"description" db:"description" validate:"omitempty,max=500"`

	// URL to badge icon
	// @Description Full URL pointing to the visual representation of the badge
	// @Description Should be a high-quality, recognizable icon or image
	// @Example "https://cdn.example.com/badges/master-explorer.png"
	// @Format uri
	IconURL string `json:"icon_url" db:"icon_url" validate:"omitempty,url" format:"uri"`

	// Timestamp when the badge was created
	// @Description Precise timestamp of badge creation in UTC
	// @Description Helps track the badge's lifecycle and origin
	// @Example "2023-06-15T14:30:00Z"
	// @Format date-time
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Timestamp when the badge was last updated
	// @Description Precise timestamp of the last modification to the badge in UTC
	// @Description Indicates when badge details were last changed
	// @Example "2023-06-16T10:15:00Z"
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// BadgeDTO represents the data transfer object for badge information
// @Description Lightweight data transfer object for badge details
// @Description Used for API responses to provide a clean, minimal representation of badges
// @Tags Achievements
type BadgeDTO struct {
	// Unique identifier for the badge
	// @Description Globally unique UUID for the badge
	// @Example "550e8400-e29b-41d4-a716-446655440000"
	// @Format uuid
	ID uuid.UUID `json:"id"`

	// Badge name
	// @Description Human-readable name of the badge
	// @Example "Master Explorer"
	Name string `json:"name"`

	// Badge description
	// @Description Detailed explanation of how to earn the badge
	// @Example "Discovered and visited 50 unique locations across the platform"
	Description string `json:"description"`

	// URL to badge icon
	// @Description Full URL pointing to the visual representation of the badge
	// @Example "https://cdn.example.com/badges/master-explorer.png"
	// @Format uri
	IconURL string `json:"icon_url"`

	// Timestamp when the badge was created
	// @Description Precise timestamp of badge creation in UTC
	// @Example "2023-06-15T14:30:00Z"
	// @Format date-time
	CreatedAt *time.Time `json:"created_at"`

	// Timestamp when the badge was last updated
	// @Description Precise timestamp of the last modification to the badge in UTC
	// @Example "2023-06-16T10:15:00Z"
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at"`
}

// ToDTO converts a Badge to a BadgeDTO
// @Description Transforms a Badge model into a lightweight data transfer object
// @Description Useful for API responses and data serialization
// @Return BadgeDTO Converted badge data transfer object
func (b *Badge) ToDTO() BadgeDTO {
	return BadgeDTO{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		IconURL:     b.IconURL,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

// BadgeCreate represents the payload for creating a badge
// @Description Payload structure for creating a new badge in the system
// @Description Provides the necessary details to mint a new achievement badge
// @Tags Achievements
type BadgeCreate struct {
	// Badge name
	// @Description Human-readable name for the new badge
	// @Description Must be unique and descriptive
	// @Example "Mountain Conqueror"
	// @Required true
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name" validate:"required,min=2,max=100"`

	// Badge description
	// @Description Detailed explanation of how to earn the badge
	// @Description Provides clear guidance on achievement criteria
	// @Example "Climbed peaks in 5 different mountain ranges"
	// @MaxLength 500
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`

	// URL to badge icon
	// @Description Full URL pointing to the visual representation of the badge
	// @Description Should be a high-quality, recognizable icon
	// @Example "https://cdn.example.com/badges/mountain-conqueror.png"
	// @Format uri
	IconURL string `json:"icon_url,omitempty" validate:"omitempty,url" format:"uri"`
}

// BadgeUpdate represents the payload for updating a badge
// @Description Payload structure for updating existing badge details
// @Description Supports partial updates with optional fields
// @Tags Achievements
type BadgeUpdate struct {
	// Unique identifier for the badge
	// @Description Globally unique UUID of the badge to be updated
	// @Description Must match an existing badge in the system
	// @Example "550e8400-e29b-41d4-a716-446655440000"
	// @Required true
	// @Format uuid
	ID uuid.UUID `json:"id" validate:"required"`

	// Badge name
	// @Description Updated human-readable name for the badge
	// @Description Optional field for renaming the badge
	// @Example "Advanced Explorer"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`

	// Badge description
	// @Description Updated explanation of how to earn the badge
	// @Description Optional field for refining achievement criteria
	// @Example "Discovered and visited 75 unique locations across the platform"
	// @MaxLength 500
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`

	// URL to badge icon
	// @Description Updated URL pointing to the visual representation of the badge
	// @Description Optional field for updating badge visual
	// @Example "https://cdn.example.com/badges/advanced-explorer.png"
	// @Format uri
	IconURL string `json:"icon_url,omitempty" validate:"omitempty,url" format:"uri"`
}
