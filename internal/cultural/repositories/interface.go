package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type EventRepository interface {
	repository.BaseRepository[models.Event]
	ListTrendingEvent(ctx context.Context, limit int) ([]models.Event, error)
	ListRelatedEvents(ctx context.Context, eventID string, limit int) ([]models.Event, error)
	Search(ctx context.Context, query string, limit, offset int) ([]models.Event, error)
	UpdateViews(ctx context.Context, id string) string // tambah update event view
}

type LocalStoryRepository interface {
	repository.BaseRepository[models.LocalStory]
}
