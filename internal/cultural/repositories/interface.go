package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

// EventRepository defines methods for managing events
type EventRepository interface {
	base.BaseRepository[models.Event, models.EventDTO]

	// Specialized event query methods
	FindPopularEvents(ctx context.Context, limit int) ([]models.EventDTO, error)
	FindRelatedEvents(ctx context.Context, eventID string, locationID string, limit int) ([]models.EventDTO, error)
	FindEventsByLocation(ctx context.Context, locationID string) ([]models.EventDTO, error)
	UpdateViews(ctx context.Context, userID, eventID string) string
	GetEventViews(ctx context.Context, id string) (int, error)
}
