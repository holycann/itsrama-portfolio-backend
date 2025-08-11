package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

// UserRepository extends BaseRepository with user-specific methods
type UserRepository interface {
	base.BaseRepository[models.User, models.User]

	// User-specific query methods
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// Role and permission management
	ChangeUserRole(ctx context.Context, payload *models.UserRoleUpdate) error
}

// UserProfileRepository extends BaseRepository with profile-specific methods
type UserProfileRepository interface {
	base.BaseRepository[models.UserProfile, models.UserProfileDTO]

	// Profile-specific query methods
	FindByUserID(ctx context.Context, userID string) (*models.UserProfileDTO, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
	FindByFullname(ctx context.Context, fullname string) ([]models.UserProfileDTO, error)

	// Additional profile-specific operations
	UpdateAvatarImage(ctx context.Context, payload *models.UserProfileAvatarUpdate) error
	UpdatePersonalInfo(ctx context.Context, payload *models.UserProfileUpdate) error
	VerifyIdentity(ctx context.Context, payload *models.UserProfileIdentityUpdate) error
}

// UserBadgeRepository extends BaseRepository with badge-specific methods
type UserBadgeRepository interface {
	base.BaseRepository[models.UserBadge, models.UserBadgeDTO]

	// Badge-specific query methods
	FindUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadgeDTO, error)
	FindUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadgeDTO, error)

	// Additional badge-specific operations
	CountUserBadges(ctx context.Context, userID string) (int, error)
	AddBadgeToUser(ctx context.Context, payload *models.UserBadge) error
	RemoveBadgeFromUser(ctx context.Context, payload *models.UserBadgePayload) error
}
