package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type BadgeService interface {
	CreateBadge(ctx context.Context, badgeCreate *models.BadgeCreate) error
	GetBadgeByID(ctx context.Context, id string) (*models.Badge, error)
	GetBadgeByName(ctx context.Context, name string) (*models.Badge, error)
	ListBadges(ctx context.Context, opts repository.ListOptions) ([]models.Badge, error)
	UpdateBadge(ctx context.Context, id string, badgeUpdate *models.BadgeCreate) error
	DeleteBadge(ctx context.Context, id string) error
	CountBadges(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetPopularBadges(ctx context.Context, limit int) ([]models.Badge, error)
}
