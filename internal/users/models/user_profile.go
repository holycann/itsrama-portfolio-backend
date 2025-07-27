package models

import (
	"time"

	"github.com/google/uuid"
)

// UserProfile represents the extended user profile information
// @Description Detailed user profile with additional personal information, including identity (KTP) data
type UserProfile struct {
	// Unique identifier for the user profile
	// @example "profile_123"
	ID uuid.UUID `json:"id" db:"id"`

	// Associated user ID
	// @example "user_123"
	UserID uuid.UUID `json:"user_id" db:"user_id" validate:"required"`

	// User's full name
	// @example "John Doe"
	Fullname string `json:"fullname" db:"fullname" validate:"required"`

	// User's biographical information
	// @example "Software engineer passionate about building great products"
	Bio string `json:"bio" db:"bio"`

	// URL to user's avatar image
	// @example "https://example.com/avatar.jpg"
	AvatarUrl string `json:"avatar_url" db:"avatar_url" format:"uri"`

	// URL to uploaded KTP image
	// @example "https://example.com/ktp.jpg"
	IdentityImageUrl string `json:"identity_image_url" db:"identity_image_url" format:"uri"`

	// Timestamp when the profile was created
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Timestamp when the profile was last updated
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserProfileCreate represents the payload for creating a new user profile
// @Description Payload for creating a user profile, including identity (KTP) data
type UserProfileCreate struct {
	// Associated user ID
	// @example "user_123"
	UserID uuid.UUID `json:"user_id" validate:"required"`

	// User's full name
	// @example "John Doe"
	Fullname string `json:"fullname" validate:"required"`

	// User's biographical information
	// @example "Software engineer passionate about building great products"
	Bio string `json:"bio,omitempty"`

	// URL to user's avatar image
	// @example "https://example.com/avatar.jpg"
	AvatarUrl string `json:"avatar_url,omitempty" format:"uri"`

	// URL to uploaded KTP image
	// @example "https://example.com/ktp.jpg"
	IdentityImageUrl string `json:"identity_image_url,omitempty" format:"uri"`
}
