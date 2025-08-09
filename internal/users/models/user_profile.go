package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

// UserProfile represents a user's detailed profile information
// @Description Comprehensive model for tracking and managing user profile details
// @Description Provides a structured representation of user personal information with rich metadata
// @Tags User Profiles
type UserProfile struct {
	// Unique identifier for the user profile
	// @Description Globally unique UUID for the user profile, generated automatically
	// @Description Serves as the primary key and reference for the user profile
	// @Example "profile_123"
	// @Format uuid
	ID uuid.UUID `json:"id" db:"id"`

	// Associated user ID
	// @Description Unique identifier linking the profile to a specific user account
	// @Description Ensures one-to-one relationship between user and profile
	// @Example "user_123"
	// @Format uuid
	UserID uuid.UUID `json:"user_id" db:"user_id" validate:"required"`

	// User's full name
	// @Description Official or preferred full name of the user
	// @Description Provides a clear, personal identifier for the user
	// @Example "John Doe"
	// @MinLength 2
	// @MaxLength 100
	Fullname string `json:"fullname" db:"fullname" validate:"required,min=2,max=100"`

	// User's biographical information
	// @Description Personal description, interests, or professional background
	// @Description Allows users to share more about themselves
	// @Example "Software engineer passionate about building great products"
	// @MaxLength 500
	Bio *string `json:"bio,omitempty" db:"bio" validate:"omitempty,max=500"`

	// URL to user's avatar image
	// @Description Public URL pointing to the user's profile picture
	// @Description Serves as a visual representation of the user
	// @Example "https://example.com/avatar.jpg"
	// @Format uri
	AvatarUrl *string `json:"avatar_url,omitempty" db:"avatar_url" validate:"omitempty,url" format:"uri"`

	// URL to uploaded KTP image
	// @Description Public URL of the user's official government-issued ID
	// @Description Used for identity verification purposes
	// @Example "https://example.com/ktp.jpg"
	// @Format uri
	IdentityImageUrl *string `json:"identity_image_url,omitempty" db:"identity_image_url" validate:"omitempty,url" format:"uri"`

	// Timestamp when the profile was created
	// @Description Precise timestamp of user profile creation in UTC
	// @Description Helps track profile information lifecycle and origin
	// @Format date-time
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`

	// Timestamp when the profile was last updated
	// @Description Precise timestamp of the last modification to the profile details in UTC
	// @Description Indicates when profile information was last changed
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// UserProfileDTO represents the data transfer object for user profile information
// @Description Comprehensive data transfer object for user profile details
// @Description Used for API responses to provide rich user profile information
// @Tags User Profiles
type UserProfileDTO struct {
	// Unique identifier for the user profile
	// @Description Globally unique UUID for the user profile
	// @Example "profile_123"
	// @Format uuid
	ID uuid.UUID `json:"id"`

	// User's full name
	// @Description Official or preferred full name of the user
	// @Example "John Doe"
	Fullname string `json:"fullname"`

	// User's biographical information
	// @Description Personal description, interests, or professional background
	// @Example "Software engineer passionate about building great products"
	Bio *string `json:"bio"`

	// URL to user's avatar image
	// @Description Public URL pointing to the user's profile picture
	// @Example "https://example.com/avatar.jpg"
	// @Format uri
	AvatarUrl *string `json:"avatar_url"`

	// URL to uploaded KTP image
	// @Description Public URL of the user's official government-issued ID
	// @Example "https://example.com/ktp.jpg"
	// @Format uri
	IdentityImageUrl *string `json:"identity_image_url"`

	// Timestamp when the profile was created
	// @Description Precise timestamp of user profile creation in UTC
	// @Format date-time
	CreatedAt *time.Time `json:"created_at"`

	// Timestamp when the profile was last updated
	// @Description Precise timestamp of the last modification to the profile details in UTC
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at"`

	// Associated user details
	// @Description Comprehensive information about the user account
	// @Description Provides additional context for the user profile
	User *User `json:"user,omitempty"`
}

// ToDTO converts a UserProfile to a UserProfileDTO
// @Description Transforms a UserProfile model into a lightweight data transfer object
// @Description Useful for API responses and data serialization
// @Return UserProfileDTO Converted user profile data transfer object
func (up *UserProfile) ToDTO() UserProfileDTO {
	return UserProfileDTO{
		ID:               up.ID,
		Fullname:         up.Fullname,
		Bio:              up.Bio,
		AvatarUrl:        up.AvatarUrl,
		IdentityImageUrl: up.IdentityImageUrl,
		CreatedAt:        up.CreatedAt,
		UpdatedAt:        up.UpdatedAt,
	}
}

// UserProfileCreate represents the payload for creating a user profile
// @Description Structured payload for user profile creation operations
// @Description Supports input for initializing new user profile records
// @Tags User Profiles
type UserProfileCreate struct {
	// Associated user ID (only for creation)
	// @Description Unique identifier of the user for whom the profile is being created
	// @Description Must be a valid, unique user account
	// @Example "user_123"
	// @Format uuid
	UserID uuid.UUID `json:"user_id,omitempty" validate:"omitempty"`

	// User's full name
	// @Description Official or preferred full name of the user
	// @Description Optional during creation, can be updated later
	// @Example "John Doe"
	// @MinLength 2
	// @MaxLength 100
	Fullname string `json:"fullname,omitempty" validate:"omitempty,min=2,max=100"`

	// User's biographical information
	// @Description Personal description, interests, or professional background
	// @Description Optional field for additional user context
	// @Example "Software engineer passionate about building great products"
	// @MaxLength 500
	Bio string `json:"bio,omitempty" validate:"omitempty,max=500"`

	// URL to user's avatar image
	// @Description Public URL pointing to the user's profile picture
	// @Description Optional field for profile visual representation
	// @Example "https://example.com/avatar.jpg"
	// @Format uri
	AvatarUrl string `json:"avatar_url,omitempty" validate:"omitempty,url" format:"uri"`

	// URL to uploaded KTP image
	// @Description Public URL of the user's official government-issued ID
	// @Description Optional field for identity verification
	// @Example "https://example.com/ktp.jpg"
	// @Format uri
	IdentityImageUrl string `json:"identity_image_url,omitempty" validate:"omitempty,url" format:"uri"`
}

// UserProfileUpdate represents the payload for updating a user profile
// @Description Structured payload for user profile update operations
// @Description Supports partial updates with optional fields
// @Tags User Profiles
type UserProfileUpdate struct {
	// Unique identifier for the user profile
	// @Description Globally unique UUID of the user profile to be updated
	// @Description Must match an existing user profile in the system
	// @Example "profile_123"
	// @Format uuid
	ID uuid.UUID `json:"id" validate:"required"`

	// Associated user ID for the profile update
	// @Description Unique identifier of the user associated with this profile
	// @Description Optional during update, helps ensure correct profile ownership
	// @Example "user_123"
	// @Format uuid
	UserID uuid.UUID `json:"user_id,omitempty" validate:"omitempty"`

	// User's full name
	// @Description Updated full name for the user profile
	// @Description Optional field for renaming the profile
	// @Example "John A. Doe"
	// @MinLength 2
	// @MaxLength 100
	Fullname string `json:"fullname,omitempty" validate:"omitempty,min=2,max=100"`

	// User's biographical information
	// @Description Updated personal description or background
	// @Description Optional field for refining profile information
	// @Example "Senior software engineer with a passion for innovative solutions"
	// @MaxLength 500
	Bio string `json:"bio,omitempty" validate:"omitempty,max=500"`
}

// UserProfileAvatarUpdate represents the payload for updating a user's avatar
// @Description Structured payload for user profile avatar update operations
// @Description Supports uploading a new profile picture
// @Tags User Profiles
type UserProfileAvatarUpdate struct {
	// Unique identifier for the user profile
	// @Description Globally unique UUID of the user profile to be updated
	// @Description Must match an existing user profile in the system
	// @Example "profile_123"
	// @Format uuid
	ID uuid.UUID `json:"id" validate:"required"`

	// URL to the new avatar image
	// @Description Public URL pointing to the new user profile picture
	// @Description Optional URL for avatar image
	// @Example "https://example.com/new-avatar.jpg"
	// @Format uri
	AvatarUrl string `json:"avatar_url,omitempty" validate:"omitempty,url" format:"uri" form:"avatar_url"`

	// Image file for avatar upload
	// @Description Multipart file upload for the new avatar image
	// @Description Allows direct image file upload during profile update
	Image *multipart.FileHeader `json:"-" form:"image" validate:"required"`
}

// UserProfileIdentityUpdate represents the payload for updating a user's identity verification image
// @Description Structured payload for user profile identity image update operations
// @Description Supports uploading a new government-issued ID image
// @Tags User Profiles
type UserProfileIdentityUpdate struct {
	// Unique identifier for the user profile
	// @Description Globally unique UUID of the user profile to be updated
	// @Description Must match an existing user profile in the system
	// @Example "profile_123"
	// @Format uuid
	ID uuid.UUID `json:"id" validate:"required"`

	// URL to the new identity image
	// @Description Public URL pointing to the new government-issued ID image
	// @Description Optional URL for identity verification image
	// @Example "https://example.com/new-ktp.jpg"
	// @Format uri
	IdentityImageUrl string `json:"identity_image_url,omitempty" validate:"omitempty,url" format:"uri" form:"identity_image_url"`

	// Image file for identity verification
	// @Description Multipart file upload for the new identity verification image
	// @Description Allows direct image file upload during profile update
	Image *multipart.FileHeader `json:"-" form:"image" validate:"required"`
}
