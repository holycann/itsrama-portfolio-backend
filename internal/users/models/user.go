package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the system
// @Description Comprehensive model for tracking and managing user account information
// @Description Provides a structured representation of user authentication and account details
// @Tags User Management
type User struct {
	// Unique identifier for the user
	// @Description Globally unique UUID for the user account, generated automatically
	// @Description Serves as the primary key and reference for the user
	// @Example "user_123"
	// @Format uuid
	ID uuid.UUID `json:"id" db:"id"`

	// User's email address
	// @Description Unique email used for authentication and communication
	// @Description Must be a valid, unique email address
	// @Example "john.doe@example.com"
	// @Format email
	Email string `json:"email" db:"email" validate:"required,email,unique"`

	// Hashed password for account authentication
	// @Description Securely hashed user password
	// @Description Never stored or returned in plain text
	Password string `json:"-" db:"password" validate:"required,min=8"`

	// User's phone number
	// @Description Unique phone number for user contact and authentication
	// @Description Optional field for additional user identification
	// @Example "+1234567890"
	// @Format phone
	Phone string `json:"phone" db:"phone" validate:"omitempty,e164"`

	// User's role in the system
	// @Description Defines the user's access level and permissions
	// @Description Determines system-wide access control
	// @Example "user"
	// @Enum user,admin,moderator
	Role string `json:"role" db:"role" validate:"required,oneof=user admin moderator"`

	// Timestamp when the user account was created
	// @Description Precise timestamp of user account creation in UTC
	// @Description Helps track user lifecycle and origin
	// @Format date-time
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Timestamp when the user account was last updated
	// @Description Precise timestamp of the last modification to the user account in UTC
	// @Description Indicates when user information was last changed
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// UserDTO represents the data transfer object for user information
// @Description Comprehensive data transfer object for user details
// @Description Used for API responses to provide rich user information with controlled exposure
// @Tags User Management
type UserDTO struct {
	// Unique identifier for the user
	// @Description Globally unique UUID for the user account
	// @Example "user_123"
	// @Format uuid
	ID uuid.UUID `json:"id"`

	// User's email address
	// @Description Unique email used for communication
	// @Example "john.doe@example.com"
	// @Format email
	Email string `json:"email"`

	// User's phone number
	// @Description Optional phone number for user contact
	// @Example "+1234567890"
	// @Format phone
	Phone string `json:"phone,omitempty"`

	// User's role in the system
	// @Description Defines the user's access level and permissions
	// @Example "user"
	// @Enum user,admin,moderator
	Role string `json:"role"`

	// Timestamp when the user account was created
	// @Description Precise timestamp of user account creation in UTC
	// @Format date-time
	CreatedAt *time.Time `json:"created_at"`

	// Timestamp when the user account was last updated
	// @Description Precise timestamp of the last modification to the user account in UTC
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at"`
}

// ToDTO converts a User to a UserDTO
// @Description Transforms a User model into a lightweight data transfer object
// @Description Useful for API responses and data serialization
// @Return UserDTO Converted user data transfer object
func (u *User) ToDTO() UserDTO {
	return UserDTO{
		ID:        u.ID,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UserCreate represents the payload for creating a new user account
// @Description Structured payload for user account creation operations
// @Description Supports initializing new user accounts with authentication details
// @Tags User Management
type UserCreate struct {
	// User's email address
	// @Description Unique email for user registration
	// @Description Must be a valid, unique email address
	// @Example "john.doe@example.com"
	// @Format email
	Email string `json:"email" validate:"required,email,unique"`

	// User's phone number
	// @Description Optional phone number for user contact
	// @Example "+1234567890"
	// @Format phone
	Phone string `json:"phone,omitempty" validate:"omitempty,e164"`

	// User's password for account authentication
	// @Description Password for local email-based authentication
	// @Description Must meet minimum security requirements
	// @MinLength 8
	Password string `json:"password,omitempty" validate:"omitempty,min=8"`

	// User's role during account creation
	// @Description Specifies the initial role for the user account
	// @Description Optional field with a default value of "user"
	// @Example "user"
	// @Enum user,admin,moderator
	Role string `json:"role,omitempty" validate:"omitempty,oneof=user admin moderator" default:"user"`
}

// UserUpdate represents the payload for updating user account details
// @Description Structured payload for user account update operations
// @Description Supports partial updates with optional fields
// @Tags User Management
type UserUpdate struct {
	// Unique identifier for the user
	// @Description Globally unique UUID of the user account to be updated
	// @Description Must match an existing user in the system
	// @Example "user_123"
	// @Format uuid
	ID uuid.UUID `json:"id" validate:"required"`

	// User's email address
	// @Description Updated email address for the user account
	// @Description Optional field for changing contact information
	// @Example "john.updated@example.com"
	// @Format email
	Email string `json:"email,omitempty" validate:"omitempty,email,unique"`

	// User's phone number
	// @Description Optional phone number for user contact
	// @Example "+1234567890"
	// @Format phone
	Phone string `json:"phone,omitempty" validate:"omitempty,e164"`
}

// UserUpdatePassword represents the payload for updating a user's password
// @Description Payload for changing user password with current and new password validation
type UserUpdatePassword struct {
	// Unique identifier for the user (optional during creation)
	// @example "user_123"
	ID uuid.UUID `json:"id,omitempty" validate:"omitempty" example:"user_123"`

	// Current password for verification
	// @example "oldPassword123!"
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=72" example:"oldPassword123!"`

	// New password to replace the current password
	// @example "newSecurePassword456!"
	NewPassword string `json:"new_password" validate:"required,min=8,max=72,passwordStrength,nefield=CurrentPassword" example:"newSecurePassword456!"`
}

// UserRoleUpdate represents the payload for updating a user's role
// @Description Payload for changing user role with role validation
type UserRoleUpdate struct {
	// Unique identifier for the user
	// @example "user_123"
	ID uuid.UUID `json:"id" validate:"required" example:"user_123"`

	// New role for the user
	// @example "admin"
	// @enums "authenticated","admin"
	NewRole string `json:"new_role" validate:"required,oneof=authenticated admin" enums:"authenticated,admin" example:"admin"`
}
