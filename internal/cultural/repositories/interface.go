package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

// EventRepository defines methods for event-related database operations
type EventRepository interface {
	repository.BaseRepository[models.Event]

	// Advanced search and filtering methods
	Search(ctx context.Context, filter repository.ListOptions) ([]models.Event, int, error)

	// Specialized query methods
	FindPopularEvents(ctx context.Context, limit int) ([]models.Event, error)
	FindRelatedEvents(ctx context.Context, eventID string, limit int) ([]models.Event, error)
	FindRecentEvents(ctx context.Context, limit int) ([]models.Event, error)
	FindEventsByLocation(ctx context.Context, locationID uuid.UUID) ([]models.Event, error)
	UpdateViews(ctx context.Context, id string) string
}

// LocalStoryRepository defines methods for local story-related database operations
type LocalStoryRepository interface {
	repository.BaseRepository[models.LocalStory]

	// Specialized methods for local stories
	FindStoriesByLocation(ctx context.Context, locationID uuid.UUID) ([]models.LocalStory, error)
	FindStoriesByOriginCulture(ctx context.Context, culture string) ([]models.LocalStory, error)
	FindByID(ctx context.Context, id string) (*models.LocalStory, error)
	Delete(ctx context.Context, id string) error
}
