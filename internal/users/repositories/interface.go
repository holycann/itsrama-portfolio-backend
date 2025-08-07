package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

// UserRepository extends BaseRepository with user-specific methods
type UserRepository interface {
	base.EnhancedRepository[models.User, models.UserDTO]

	// User-specific query methods
	FindByEmail(ctx context.Context, email string) (*models.UserDTO, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// Additional user-specific operations
	UpdateLastSignIn(ctx context.Context, userID string) error
	ChangeUserRole(ctx context.Context, userID, newRole string) error
}

// UserProfileRepository extends BaseRepository with profile-specific methods
type UserProfileRepository interface {
	base.EnhancedRepository[models.UserProfile, models.UserProfileDTO]

	// Profile-specific query methods
	FindByUserID(ctx context.Context, userID string) (*models.UserProfileDTO, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)

	// Additional profile-specific operations
	UpdateProfileImage(ctx context.Context, profileID string, imageURL string) error
	UpdateBio(ctx context.Context, profileID, bio string) error
}

// UserBadgeRepository extends BaseRepository with badge-specific methods
type UserBadgeRepository interface {
	base.EnhancedRepository[models.UserBadge, models.UserBadgeDTO]

	// Badge-specific query methods
	FindUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadgeDTO, error)
	FindUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadgeDTO, error)

	// Additional badge-specific operations
	CountUserBadges(ctx context.Context, userID string) (int, error)
	RemoveBadgeFromUser(ctx context.Context, userID, badgeID string) error
}
