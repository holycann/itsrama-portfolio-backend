package models

import (
	"time"
)

// User represents the user account information
// @Description User account details with authentication and metadata
type User struct {
	// Unique identifier for the user
	// @example "user_123"
	ID string `json:"id" db:"id" example:"user_123"`

	// User's email address
	// @example "user@example.com"
	Email string `json:"email" db:"email" validate:"required,email" example:"admin@gmail.com"`

	// User's password (never returned in responses)
	// @example "securePassword123!"
	Password string `json:"-" db:"password" swaggerignore:"true"`

	// User's phone number
	// @example "+1234567890"
	Phone string `json:"phone" db:"phone" validate:"omitempty,e164" example:"+1234567890"`

	// User's role in the system
	// @example "user"
	// @enums "user","admin","moderator"
	Role string `json:"role" db:"role" validate:"required,oneof=user admin moderator" enums:"user,admin,moderator" example:"user"`

	// Timestamp of the last user sign-in
	LastSignInAt *time.Time `json:"last_sign_in_at" db:"last_sign_in_at"`

	// Timestamp when the user account was created
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Timestamp of the last user account update
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`

	// Timestamp of user account soft deletion
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// UserCreate represents the payload for creating a new user
// @Description Payload for user registration
type UserCreate struct {
	// User's email address
	// @example "user@example.com"
	Email string `json:"email" validate:"required,email" example:"admin@gmail.com"`

	// User's password
	// @example "securePassword123!"
	Password string `json:"password" validate:"required,min=8,max=72,password" example:"admin123"`

	// User's phone number (optional)
	// @example "+1234567890"
	Phone string `json:"phone,omitempty" validate:"omitempty,e164" example:"+1234567890"`

	// User's role in the system
	// @example "user"
	// @enums "user","admin","moderator"
	Role string `json:"role" validate:"required,oneof=user admin moderator" enums:"user,admin,moderator" example:"user"`
}

// UserUpdate represents the payload for updating user details
// @Description Payload for updating user information
type UserUpdate struct {
	// User's email address (optional)
	// @example "newemail@example.com"
	Email string `json:"email,omitempty" validate:"omitempty,email" example:"newemail@example.com"`

	// User's phone number (optional)
	// @example "+1234567890"
	Phone string `json:"phone,omitempty" validate:"omitempty,e164" example:"+1234567890"`

	// User's role in the system (optional)
	// @example "admin"
	// @enums "user","admin","moderator"
	Role string `json:"role,omitempty" validate:"omitempty,oneof=user admin moderator" enums:"user,admin,moderator" example:"admin"`
}
