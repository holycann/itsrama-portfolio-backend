package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
)

type eventService struct {
	eventRepo repositories.EventRepository
}

// NewEventService creates a new instance of the event service
// with the given event repository.
func NewEventService(eventRepo repositories.EventRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
	}
}

// CreateEvent adds a new event to the database
// Validates the event object before creating
func (s *eventService) CreateEvent(ctx context.Context, event *models.Event) error {
	// Validate event object
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Validate required fields (example validation)
	if event.Name == "" {
		return fmt.Errorf("event name is required")
	}

	event.ID = uuid.NewString()

	// Call repository to create event
	return s.eventRepo.Create(ctx, event)
}

// GetEvents retrieves a list of events with pagination
func (s *eventService) GetEvents(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve events from repository
	events, err := s.eventRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert []models.Event to []*models.Event
	eventPtrs := make([]*models.Event, len(events))
	for i := range events {
		eventPtrs[i] = &events[i]
	}

	return eventPtrs, nil
}

// GetEventByID retrieves a single event by its unique identifier
func (s *eventService) GetEventByID(ctx context.Context, id string) (*models.Event, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("event ID cannot be empty")
	}

	// Retrieve event from repository
	return s.eventRepo.FindByID(ctx, id)
}

// GetEventByName retrieves a event by its name
// Note: This method is not directly supported by the current repository implementation
// You might need to add a custom method in the repository or implement filtering
func (s *eventService) GetEventByName(ctx context.Context, name string) (*models.Event, error) {
	// Validate name
	if name == "" {
		return nil, fmt.Errorf("event name cannot be empty")
	}

	// Since the current repository doesn't have a direct method for this,
	// we'll use a workaround by listing all events and finding by name
	events, err := s.eventRepo.List(ctx, 1, 0)
	if err != nil {
		return nil, err
	}

	// Find event by name (linear search)
	for _, event := range events {
		if event.Name == name {
			return &event, nil
		}
	}

	return nil, fmt.Errorf("event with name %s not found", name)
}

// UpdateEvent updates an existing event in the database
func (s *eventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	// Validate event object
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Validate required fields
	if event.ID == "" {
		return fmt.Errorf("event ID is required for update")
	}

	// Call repository to update event
	return s.eventRepo.Update(ctx, event)
}

// DeleteEvent removes a event from the database by its ID
func (s *eventService) DeleteEvent(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("event ID cannot be empty")
	}

	// Call repository to delete event
	return s.eventRepo.Delete(ctx, id)
}

// Count menghitung jumlah total lokasi yang tersimpan
func (s *eventService) Count(ctx context.Context) (int, error) {
	return s.eventRepo.Count(ctx)
}
