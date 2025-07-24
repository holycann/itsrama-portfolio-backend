package models

import "time"

// UserProfile represents the extended user profile information
// @Description Detailed user profile with additional personal information
type UserProfile struct {
	// Unique identifier for the user profile
	// @example "profile_123"
	ID string `json:"id" db:"id" example:"profile_123"`

	// Associated user ID
	// @example "user_123"
	UserID string `json:"user_id" db:"user_id" validate:"required" example:"user_123"`

	// User's full name
	// @example "John Doe"
	Fullname string `json:"fullname" db:"fullname" validate:"required" example:"John Doe"`

	// User's biographical information
	// @example "Software engineer passionate about building great products"
	Bio string `json:"bio" db:"bio"`

	// URL to user's avatar image
	// @example "https://example.com/avatar.jpg"
	AvatarUrl string `json:"avatar_url" db:"avatar_url" format:"uri"`

	// Timestamp when the profile was created
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Timestamp of the last profile update
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`

	// Timestamp of profile soft deletion
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// UserProfileCreate represents the payload for creating a new user profile
// @Description Payload for creating a user profile
type UserProfileCreate struct {
	// Associated user ID
	// @example "user_123"
	UserID string `json:"user_id" validate:"required" example:"user_123"`

	// User's full name
	// @example "John Doe"
	Fullname string `json:"fullname" validate:"required" example:"John Doe"`

	// User's biographical information
	// @example "Software engineer passionate about building great products"
	Bio string `json:"bio,omitempty"`

	// URL to user's avatar image
	// @example "https://example.com/avatar.jpg"
	AvatarUrl string `json:"avatar_url,omitempty" format:"uri"`
}
