package services

import (
	"context"
	"mime/multipart"

	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

// EventService defines operations for managing events
type EventService interface {
	// Event Creation and Management
	CreateEvent(ctx context.Context, event *models.EventPayload, image *multipart.FileHeader) (*models.EventDTO, error)
	UpdateEvent(ctx context.Context, event *models.EventPayload, image *multipart.FileHeader) (*models.EventDTO, error)
	DeleteEvent(ctx context.Context, id string) error

	// Event Retrieval Operations
	GetEventByID(ctx context.Context, id string) (*models.EventDTO, error)
	ListEvents(ctx context.Context, opts base.ListOptions) ([]models.EventDTO, int, error)
	SearchEvents(ctx context.Context, query string, opts base.ListOptions) ([]models.EventDTO, int, error)

	// Specialized Event Operations
	GetTrendingEvents(ctx context.Context, limit int) ([]models.EventDTO, error)
	GetRelatedEvents(ctx context.Context, eventID, locationID string, limit int) ([]models.EventDTO, error)
	UpdateEventViews(ctx context.Context, userID, eventID string) string
	CountEvents(ctx context.Context, filters []base.FilterOption) (int, error)
}
