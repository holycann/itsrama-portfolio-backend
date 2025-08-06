package services

import (
	"context"
	"mime/multipart"

	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *models.RequestEvent, image *multipart.FileHeader) error
	GetEventByID(ctx context.Context, id string) (*models.ResponseEvent, error)
	ListEvents(ctx context.Context, opts repository.ListOptions) ([]models.ResponseEvent, error)
	UpdateEvent(ctx context.Context, event *models.RequestEvent, image *multipart.FileHeader) error
	DeleteEvent(ctx context.Context, id string) error
	CountEvents(ctx context.Context, filters []repository.FilterOption) (int, error)
	UpdateEventViews(ctx context.Context, userID, eventID string) string
	GetTrendingEvents(ctx context.Context, limit int) ([]models.ResponseEvent, error)
	GetRelatedEvents(ctx context.Context, eventID string, limit int) ([]models.ResponseEvent, error)
	SearchEvents(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseEvent, error)
}

type LocalStoryService interface {
	CreateLocalStory(ctx context.Context, localStory *models.LocalStory) error
	GetLocalStoryByID(ctx context.Context, id string) (*models.LocalStory, error)
	ListLocalStories(ctx context.Context, opts repository.ListOptions) ([]models.LocalStory, error)
	UpdateLocalStory(ctx context.Context, localStory *models.LocalStory) error
	DeleteLocalStory(ctx context.Context, id string) error
	CountLocalStories(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetLocalStoriesByLocation(ctx context.Context, locationID string) ([]*models.LocalStory, error)
	GetLocalStoriesByOriginCulture(ctx context.Context, culture string) ([]*models.LocalStory, error)
	SearchLocalStories(ctx context.Context, query string, opts repository.ListOptions) ([]models.LocalStory, error)
}
