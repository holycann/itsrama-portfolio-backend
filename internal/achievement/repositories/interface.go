package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/achievement/models"
)

// BadgeRepository defines the interface for badge data operations
type BadgeRepository interface {
	Create(ctx context.Context, badge *models.Badge) (*models.Badge, error)
	FindByID(ctx context.Context, id string) (*models.Badge, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.Badge, error)
	Update(ctx context.Context, badge *models.Badge) (*models.Badge, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}
