package models

import (
	"time"

	"github.com/google/uuid"
)

// UserDTO represents the data transfer object for user information
type UserDTO struct {
	// Unique identifier for the user
	ID string `json:"id" example:"user_123"`

	// User's email address
	Email string `json:"email" example:"user@example.com"`

	// User's phone number
	Phone string `json:"phone,omitempty" example:"+1234567890"`

	// User's role in the system
	Role string `json:"role" enums:"user,admin,moderator" example:"user"`

	// Timestamp of the last user sign-in
	LastSignInAt *time.Time `json:"last_sign_in_at,omitempty"`

	// Timestamp when the user account was created
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// Timestamp of the last user account update
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// UserProfileDTO represents the data transfer object for user profile
type UserProfileDTO struct {
	// Unique identifier for the user profile
	ID uuid.UUID `json:"id"`

	// Associated user ID
	UserID uuid.UUID `json:"user_id"`

	// User's full name
	Fullname string `json:"fullname"`

	// User's biographical information
	Bio string `json:"bio,omitempty"`

	// URL to user's avatar image
	AvatarURL string `json:"avatar_url,omitempty" format:"uri"`

	// URL to uploaded identity image
	IdentityImageURL string `json:"identity_image_url,omitempty" format:"uri"`

	// Timestamp when the profile was created
	CreatedAt time.Time `json:"created_at"`

	// Timestamp when the profile was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// UserBadgeDTO represents the data transfer object for user badges
type UserBadgeDTO struct {
	// Unique identifier for the user badge
	ID uuid.UUID `json:"id"`

	// Associated user ID
	UserID uuid.UUID `json:"user_id"`

	// Badge identifier
	BadgeID uuid.UUID `json:"badge_id"`

	// Timestamp when the badge was earned
	CreatedAt time.Time `json:"created_at"`
}

// Conversion methods to map between models and DTOs

// ToDTO converts a User model to UserDTO
func (u *User) ToDTO() UserDTO {
	return UserDTO{
		ID:           u.ID,
		Email:        u.Email,
		Phone:        u.Phone,
		Role:         u.Role,
		LastSignInAt: u.LastSignInAt,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// ToDTO converts a UserProfile model to UserProfileDTO
func (up *UserProfile) ToDTO() UserProfileDTO {
	return UserProfileDTO{
		ID:               up.ID,
		UserID:           up.UserID,
		Fullname:         up.Fullname,
		Bio:              up.Bio,
		AvatarURL:        up.AvatarUrl,
		IdentityImageURL: up.IdentityImageUrl,
		CreatedAt:        up.CreatedAt,
		UpdatedAt:        up.UpdatedAt,
	}
}

// ToDTO converts a UserBadge model to UserBadgeDTO
func (ub *UserBadge) ToDTO() UserBadgeDTO {
	return UserBadgeDTO{
		ID:        ub.ID,
		UserID:    ub.UserID,
		BadgeID:   ub.BadgeID,
		CreatedAt: ub.CreatedAt,
	}
}

// FromDTO methods to convert DTOs back to models (if needed)

// FromUserDTO converts UserDTO to User model
func FromUserDTO(dto UserDTO) User {
	return User{
		ID:           dto.ID,
		Email:        dto.Email,
		Phone:        dto.Phone,
		Role:         dto.Role,
		LastSignInAt: dto.LastSignInAt,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
}

// FromUserProfileDTO converts UserProfileDTO to UserProfile model
func FromUserProfileDTO(dto UserProfileDTO) UserProfile {
	return UserProfile{
		ID:               dto.ID,
		UserID:           dto.UserID,
		Fullname:         dto.Fullname,
		Bio:              dto.Bio,
		AvatarUrl:        dto.AvatarURL,
		IdentityImageUrl: dto.IdentityImageURL,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
	}
}

// FromUserBadgeDTO converts UserBadgeDTO to UserBadge model
func FromUserBadgeDTO(dto UserBadgeDTO) UserBadge {
	return UserBadge{
		ID:        dto.ID,
		UserID:    dto.UserID,
		BadgeID:   dto.BadgeID,
		CreatedAt: dto.CreatedAt,
	}
}
