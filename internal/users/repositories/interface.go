package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type UserRepository interface {
	repository.BaseRepository[models.User]

	// Specialized methods for users
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type UserProfileRepository interface {
	repository.BaseRepository[models.UserProfile]

	// Specialized methods for user profiles
	FindByUserID(ctx context.Context, userID string) (*models.UserProfile, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
}

type UserBadgeRepository interface {
	repository.BaseRepository[models.UserBadge]

	// Specialized methods for user badges
	FindUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadge, error)
	FindUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadge, error)
}
