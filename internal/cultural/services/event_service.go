package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type eventService struct {
	eventRepo repositories.EventRepository
}

func NewEventService(eventRepo repositories.EventRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, event *models.Event) error {
	// Validate event object
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Validate required fields
	if event.Name == "" {
		return fmt.Errorf("event name is required")
	}

	// Set default values
	event.ID = uuid.New()
	event.CreatedAt = time.Now()

	// Call repository to create event
	return s.eventRepo.Create(ctx, event)
}

func (s *eventService) GetEventByID(ctx context.Context, id string) (*models.Event, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("event ID cannot be empty")
	}

	// Retrieve event from repository
	event, err := s.eventRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update views count
	_ = s.UpdateEventViews(ctx, id)

	return event, nil
}

func (s *eventService) ListEvents(ctx context.Context, opts repository.ListOptions) ([]models.Event, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.eventRepo.List(ctx, opts)
}

func (s *eventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	// Validate event object
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Validate required fields
	if event.ID == uuid.Nil {
		return fmt.Errorf("event ID is required for update")
	}

	// Call repository to update event
	return s.eventRepo.Update(ctx, event)
}

func (s *eventService) DeleteEvent(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("event ID cannot be empty")
	}

	// Call repository to delete event
	return s.eventRepo.Delete(ctx, id)
}

func (s *eventService) CountEvents(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.eventRepo.Count(ctx, filters)
}

func (s *eventService) UpdateEventViews(ctx context.Context, id string) string {
	if id == "" {
		return "event ID cannot be empty"
	}
	return s.eventRepo.UpdateViews(ctx, id)
}

func (s *eventService) GetTrendingEvents(ctx context.Context, limit int) ([]models.Event, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.eventRepo.FindPopularEvents(ctx, limit)
}

func (s *eventService) GetRelatedEvents(ctx context.Context, eventID string, limit int) ([]models.Event, error) {
	// Validate input parameters
	if eventID == "" {
		return nil, fmt.Errorf("event ID cannot be empty")
	}
	if limit <= 0 {
		limit = 3 // Default limit
	}

	return s.eventRepo.FindRelatedEvents(ctx, eventID, limit)
}

func (s *eventService) SearchEvents(ctx context.Context, query string, opts repository.ListOptions) ([]models.Event, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	// Add search query to filters
	opts.Filters = append(opts.Filters,
		repository.FilterOption{
			Field:    "name",
			Operator: "like",
			Value:    query,
		},
		repository.FilterOption{
			Field:    "description",
			Operator: "like",
			Value:    query,
		},
	)

	return s.eventRepo.List(ctx, opts)
}
