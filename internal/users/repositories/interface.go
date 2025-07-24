package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type UserRepository interface {
	repository.BaseRepository[models.User]

	SoftDelete(ctx context.Context, id string) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type UserProfileRepository interface {
	repository.BaseRepository[models.UserProfile]

	SoftDelete(ctx context.Context, id string) error
	FindByUserID(ctx context.Context, userID string) (*models.UserProfile, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
}
