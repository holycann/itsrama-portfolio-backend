package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

// EventRepository defines methods for event-related database operations
type EventRepository interface {
	repository.BaseRepository[models.Event, models.ResponseEvent]

	// Specialized query methods
	FindPopularEvents(ctx context.Context, limit int) ([]models.ResponseEvent, error)
	FindRelatedEvents(ctx context.Context, eventID string, limit int) ([]models.ResponseEvent, error)
	FindRecentEvents(ctx context.Context, limit int) ([]models.ResponseEvent, error)
	FindEventsByLocation(ctx context.Context, locationID uuid.UUID) ([]models.ResponseEvent, error)
	UpdateViews(ctx context.Context, id string) string
	GetEventViews(ctx context.Context, id string) (int, error)
}

// LocalStoryRepository defines methods for local story-related database operations
type LocalStoryRepository interface {
	repository.BaseRepository[models.LocalStory, models.LocalStory]

	// Specialized methods for local stories
	FindStoriesByLocation(ctx context.Context, locationID uuid.UUID) ([]*models.LocalStory, error)
	FindStoriesByOriginCulture(ctx context.Context, culture string) ([]*models.LocalStory, error)
}
