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

// GetTrendingEvents retrieves a list of trending events based on the highest views
func (s *eventService) GetTrendingEvents(ctx context.Context, limit int) ([]*models.Event, error) {
	if limit <= 0 {
		limit = 10
	}
	events, err := s.eventRepo.ListTrendingEvent(ctx, limit)
	if err != nil {
		return nil, err
	}
	eventPtrs := make([]*models.Event, len(events))
	for i := range events {
		eventPtrs[i] = &events[i]
	}
	return eventPtrs, nil
}

// ListRelatedEvents retrieves a list of events related to a specific event
func (s *eventService) ListRelatedEvents(ctx context.Context, eventID string, limit int) ([]*models.Event, error) {
	// Validate input parameters
	if eventID == "" {
		return nil, fmt.Errorf("event ID cannot be empty")
	}
	if limit <= 0 {
		limit = 3 // Default limit
	}

	// Retrieve related events from repository
	events, err := s.eventRepo.ListRelatedEvents(ctx, eventID, limit)
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

// UpdateEventViews increases the view count for an event by its ID
func (s *eventService) UpdateEventViews(ctx context.Context, id string) string {
	if id == "" {
		return "event ID cannot be empty"
	}
	return s.eventRepo.UpdateViews(ctx, id)
}

// GetEventByID retrieves a single event by its unique identifier
// and updates its view count
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

// Count calculates the total number of stored locations
func (s *eventService) Count(ctx context.Context) (int, error) {
	return s.eventRepo.Count(ctx)
}

// SearchEvents searches for events based on a query string in the name or description
func (s *eventService) SearchEvents(ctx context.Context, query string, limit, offset int) ([]*models.Event, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	events, err := s.eventRepo.Search(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	eventPtrs := make([]*models.Event, len(events))
	for i := range events {
		eventPtrs[i] = &events[i]
	}
	return eventPtrs, nil
}
