package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

// UserBadgeService defines operations for managing user badges
type UserBadgeService interface {
	// Badge Creation and Management
	AddBadgeToUser(ctx context.Context, payload models.UserBadge) error
	RemoveBadgeFromUser(ctx context.Context, payload models.UserBadgePayload) error

	// Badge Retrieval Operations
	GetUserBadgeByID(ctx context.Context, id string) (*models.UserBadgeDTO, error)
	GetUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadgeDTO, error)
	GetUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadgeDTO, error)

	// Listing and Searching
	ListUserBadges(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, int, error)
	SearchUserBadges(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, int, error)

	// Utility Operations
	DeleteUserBadge(ctx context.Context, id string) error
	CountUserBadges(ctx context.Context, filters []base.FilterOption) (int, error)
}

// UserService defines operations for managing user accounts
type UserService interface {
	// User Creation and Management
	CreateUser(ctx context.Context, user *models.UserCreate) (*models.UserDTO, error)
	UpdateUser(ctx context.Context, user *models.UserUpdate) (*models.UserDTO, error)
	DeleteUser(ctx context.Context, id string) error

	// User Retrieval Operations
	GetUserByID(ctx context.Context, id string) (*models.UserDTO, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserDTO, error)

	// Listing and Searching
	ListUsers(ctx context.Context, opts base.ListOptions) ([]models.UserDTO, int, error)
	SearchUsers(ctx context.Context, opts base.ListOptions) ([]models.UserDTO, int, error)

	// Utility Operations
	CountUsers(ctx context.Context, filters []base.FilterOption) (int, error)
}

// UserProfileService defines operations for managing user profiles
type UserProfileService interface {
	// Profile Creation and Management
	CreateProfile(ctx context.Context, userProfile *models.UserProfileCreate) (*models.UserProfileDTO, error)
	UpdateProfile(ctx context.Context, userProfile *models.UserProfileUpdate) (*models.UserProfileDTO, error)
	UpdateProfileAvatar(ctx context.Context, payload *models.UserProfileAvatarUpdate) (*models.UserProfileDTO, error)
	UpdateProfileIdentity(ctx context.Context, payload *models.UserProfileIdentityUpdate) (*models.UserProfileDTO, error)
	DeleteProfile(ctx context.Context, id string) error

	// Profile Retrieval Operations
	GetProfileByID(ctx context.Context, id string) (*models.UserProfileDTO, error)
	GetProfileByUserID(ctx context.Context, userID string) (*models.UserProfileDTO, error)
	GetProfileByFullname(ctx context.Context, fullname string) ([]models.UserProfileDTO, error)

	// Listing and Searching
	ListProfiles(ctx context.Context, opts base.ListOptions) ([]models.UserProfileDTO, int, error)
	SearchProfiles(ctx context.Context, opts base.ListOptions) ([]models.UserProfileDTO, int, error)

	// Utility Operations
	CountProfiles(ctx context.Context, filters []base.FilterOption) (int, error)
}
