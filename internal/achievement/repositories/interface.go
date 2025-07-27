package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type BadgeRepository interface {
	repository.BaseRepository[models.Badge]

	// Specialized methods for badges
	FindBadgeByName(ctx context.Context, name string) (*models.Badge, error)
	FindPopularBadges(ctx context.Context, limit int) ([]models.Badge, error)
}
