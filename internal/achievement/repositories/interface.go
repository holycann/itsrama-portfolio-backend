package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

type BadgeRepository interface {
	base.BaseRepository[models.Badge, models.Badge]
	FindBadgeByName(ctx context.Context, name string) (*models.Badge, error)
	FindPopularBadges(ctx context.Context, limit int) ([]models.Badge, error)
}
