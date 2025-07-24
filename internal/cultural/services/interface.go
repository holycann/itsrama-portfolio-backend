package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/cultural/models"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvents(ctx context.Context, limit, offset int) ([]*models.Event, error)
	GetEventByID(ctx context.Context, id string) (*models.Event, error)
	GetEventByName(ctx context.Context, name string) (*models.Event, error)
	UpdateEvent(ctx context.Context, user *models.Event) error
	DeleteEvent(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type LocalStoryService interface {
	CreateLocalStory(ctx context.Context, localStory *models.LocalStory) error
	GetLocalStories(ctx context.Context, limit, offset int) ([]*models.LocalStory, error)
	GetLocalStoryByID(ctx context.Context, id string) (*models.LocalStory, error)
	GetLocalStoryByTitle(ctx context.Context, title string) (*models.LocalStory, error)
	UpdateLocalStory(ctx context.Context, localStory *models.LocalStory) error
	DeleteLocalStory(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}
