package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

// BadgeService defines operations for managing badges
type BadgeService interface {
	// Badge Creation and Management
	CreateBadge(ctx context.Context, badgeCreate *models.BadgeCreate) (*models.BadgeDTO, error)
	UpdateBadge(ctx context.Context, id string, badgeUpdate *models.BadgeUpdate) (*models.BadgeDTO, error)
	DeleteBadge(ctx context.Context, id string) error

	// Badge Retrieval Operations
	GetBadgeByID(ctx context.Context, id string) (*models.BadgeDTO, error)
	GetBadgeByName(ctx context.Context, name string) (*models.BadgeDTO, error)

	// Listing and Searching
	ListBadges(ctx context.Context, opts base.ListOptions) ([]models.BadgeDTO, int, error)
	SearchBadges(ctx context.Context, opts base.ListOptions) ([]models.BadgeDTO, int, error)

	// Utility Operations
	CountBadges(ctx context.Context, filters []base.FilterOption) (int, error)
	GetPopularBadges(ctx context.Context, limit int) ([]models.BadgeDTO, error)
}
