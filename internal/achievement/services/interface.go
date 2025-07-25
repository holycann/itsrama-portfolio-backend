package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/achievement/models"
)

// BadgeService defines the interface for badge business logic
type BadgeService interface {
	CreateBadge(ctx context.Context, badgeCreate *models.BadgeCreate) (*models.Badge, error)
	GetBadgeByID(ctx context.Context, id string) (*models.Badge, error)
	ListBadges(ctx context.Context, limit, offset int) ([]models.Badge, error)
	UpdateBadge(ctx context.Context, id string, badgeUpdate *models.BadgeCreate) (*models.Badge, error)
	DeleteBadge(ctx context.Context, id string) error
	CountBadges(ctx context.Context) (int, error)
}
