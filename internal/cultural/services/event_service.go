package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	discussionModels "github.com/holycann/cultour-backend/internal/discussion/models"
	discussionServices "github.com/holycann/cultour-backend/internal/discussion/services"
	placeModels "github.com/holycann/cultour-backend/internal/place/models"
	placeServices "github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/supabase"
	storage_go "github.com/supabase-community/storage-go"
)

type eventService struct {
	eventRepo       repositories.EventRepository
	locationService placeServices.LocationService
	storage         supabase.SupabaseStorage
	threadService   discussionServices.ThreadService
}

func NewEventService(
	eventRepo repositories.EventRepository,
	service placeServices.LocationService,
	storage supabase.SupabaseStorage,
	threadService discussionServices.ThreadService,
) EventService {
	return &eventService{
		eventRepo:       eventRepo,
		storage:         storage,
		locationService: service,
		threadService:   threadService,
	}
}

func (s *eventService) findOrCreateLocation(ctx context.Context, locationPayload *placeModels.LocationCreate) (*placeModels.Location, error) {
	const tolerance = 0.0001

	// Search for existing locations with matching name
	existingLocations, err := s.locationService.ListLocations(ctx, base.ListOptions{
		Filters: []base.FilterOption{
			{
				Field:    "name",
				Operator: "eq",
				Value:    locationPayload.Name,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Check for location with similar coordinates
	for _, loc := range existingLocations {
		if abs(loc.Latitude-locationPayload.Latitude) <= tolerance &&
			abs(loc.Longitude-locationPayload.Longitude) <= tolerance {
			location := &placeModels.Location{
				ID:        loc.ID,
				CityID:    loc.CityID,
				Name:      loc.Name,
				Latitude:  loc.Latitude,
				Longitude: loc.Longitude,
			}
			return location, nil
		}
	}
	// Convert LocationPayload to Location for creation
	locationToCreate := &placeModels.LocationCreate{
		Name:      locationPayload.Name,
		Latitude:  locationPayload.Latitude,
		Longitude: locationPayload.Longitude,
		CityID:    locationPayload.CityID,
	}

	// Create new location if no match found
	location, err := s.locationService.CreateLocation(ctx, locationToCreate)
	if err != nil {
		return nil, err
	}

	return location, nil
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func (s *eventService) CreateEvent(ctx context.Context, event *models.EventPayload, image *multipart.FileHeader) (*models.EventDTO, error) {
	// Validate input
	if err := base.ValidateModel(event); err != nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Invalid event payload",
			err,
			errors.WithContext("payload", event),
		)
	}

	// Find or create location
	location, err := s.findOrCreateLocation(ctx, event.Location)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to process event location",
			errors.WithContext("location_name", event.Location.Name),
		)
	}

	now := time.Now().UTC()
	eventData := &models.Event{
		ID:            uuid.New(),
		UserID:        event.UserID,
		LocationID:    location.ID,
		Name:          event.Name,
		Description:   event.Description,
		StartDate:     event.StartDate,
		EndDate:       event.EndDate,
		IsKidFriendly: event.IsKidFriendly,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	// Upload image if provided
	if image != nil {
		imageURL, err := s.uploadEventImage(ctx, eventData.ID.String(), image)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload event image",
				errors.WithContext("event_id", eventData.ID),
			)
		}
		eventData.ImageURL = imageURL
	}

	// Create event in repository
	createdEvent, err := s.eventRepo.Create(ctx, eventData)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to create event",
			errors.WithContext("event_name", event.Name),
		)
	}

	// Create a thread for the event
	threadInput := &discussionModels.CreateThread{
		EventID:   createdEvent.ID,
		CreatorID: event.UserID,
		Status:    "active",
	}

	_, err = s.threadService.CreateThread(ctx, threadInput)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to create thread for event",
			errors.WithContext("event_id", createdEvent.ID),
		)
	}

	// Fetch and return created event details
	eventDTO, err := s.getEventDetails(ctx, createdEvent.ID.String())
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to retrieve created event details",
			errors.WithContext("event_id", createdEvent.ID),
		)
	}

	return eventDTO, nil
}

func (s *eventService) uploadEventImage(ctx context.Context, eventID string, file *multipart.FileHeader) (string, error) {
	if eventID == "" {
		return "", fmt.Errorf("event ID cannot be empty")
	}
	if file == nil {
		return "", fmt.Errorf("file data is required")
	}

	destPath := "images/events/" + eventID + filepath.Ext(file.Filename)

	_, err := s.storage.Upload(ctx, file, destPath, storage_go.FileOptions{
		ContentType: func(s string) *string { return &s }("image"),
		Upsert:      func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		return "", errors.Wrap(err,
			errors.ErrInternal,
			"Failed to upload event image",
			errors.WithContext("event_id", eventID),
		)
	}

	signedURL, err := s.storage.GetPublicURL(destPath)
	if err != nil {
		return "", errors.Wrap(err,
			errors.ErrInternal,
			"Failed to get public URL for event image",
			errors.WithContext("dest_path", destPath),
		)
	}

	return signedURL, nil
}

func (s *eventService) getEventDetails(ctx context.Context, eventID string) (*models.EventDTO, error) {
	if eventID == "" {
		return nil, fmt.Errorf("event ID cannot be empty")
	}

	return s.eventRepo.FindByID(ctx, eventID)
}

func (s *eventService) GetEventByID(ctx context.Context, id string) (*models.EventDTO, error) {
	return s.getEventDetails(ctx, id)
}

func (s *eventService) ListEvents(ctx context.Context, opts base.ListOptions) ([]models.EventDTO, int, error) {
	events, total, err := s.eventRepo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to list events",
			errors.WithContext("options", opts),
		)
	}

	return events, total, nil
}

func (s *eventService) UpdateEvent(ctx context.Context, event *models.EventPayload, image *multipart.FileHeader) (*models.EventDTO, error) {
	// Validate input
	if err := base.ValidateModel(event); err != nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Invalid event payload",
			err,
			errors.WithContext("payload", event),
		)
	}

	fmt.Println("event:", event.ID)

	// Retrieve existing event
	existingEvent, err := s.getEventDetails(ctx, event.ID.String())
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to retrieve existing event",
			errors.WithContext("event_id", event.ID),
		)
	}

	// Find or create location
	location, err := s.findOrCreateLocation(ctx, event.Location)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to process event location",
			errors.WithContext("location_name", event.Location.Name),
		)
	}

	now := time.Now().UTC()
	eventData := &models.Event{
		ID:            event.ID,
		UserID:        event.UserID,
		LocationID:    location.ID,
		Name:          s.getUpdatedValue(event.Name, existingEvent.Name),
		Description:   s.getUpdatedValue(event.Description, existingEvent.Description),
		StartDate:     s.getUpdatedTime(event.StartDate, existingEvent.StartDate),
		EndDate:       s.getUpdatedTime(event.EndDate, existingEvent.EndDate),
		IsKidFriendly: s.getUpdatedBool(event.IsKidFriendly, existingEvent.IsKidFriendly),
		ImageURL:      existingEvent.ImageURL,
		UpdatedAt:     &now,
	}

	// Upload new image if provided
	if image != nil {
		imageURL, err := s.uploadEventImage(ctx, eventData.ID.String(), image)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload event image",
				errors.WithContext("event_id", eventData.ID),
			)
		}
		eventData.ImageURL = imageURL
	}

	// Update event in repository
	updatedEvent, err := s.eventRepo.Update(ctx, eventData)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to update event",
			errors.WithContext("event_id", eventData.ID),
		)
	}

	// Fetch and return updated event details
	updatedEventDTO, err := s.getEventDetails(ctx, updatedEvent.ID.String())
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to retrieve updated event details",
			errors.WithContext("event_id", updatedEvent.ID),
		)
	}

	return updatedEventDTO, nil
}

// Helper methods for updating values
func (s *eventService) getUpdatedValue(newValue, existingValue string) string {
	if newValue != "" && newValue != existingValue {
		return newValue
	}
	return existingValue
}

func (s *eventService) getUpdatedTime(newTime, existingTime time.Time) time.Time {
	if !newTime.IsZero() && newTime != existingTime {
		return newTime
	}
	return existingTime
}

func (s *eventService) getUpdatedBool(newBool, existingBool bool) bool {
	if newBool != existingBool {
		return newBool
	}
	return existingBool
}

func (s *eventService) DeleteEvent(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"Event ID cannot be empty",
			nil,
		)
	}

	return s.eventRepo.Delete(ctx, id)
}

func (s *eventService) CountEvents(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.eventRepo.Count(ctx, filters)
}

func (s *eventService) UpdateEventViews(ctx context.Context, userID, eventID string) string {
	if userID == "" || eventID == "" {
		return "user ID and event ID cannot be empty"
	}
	return s.eventRepo.UpdateViews(ctx, userID, eventID)
}

func (s *eventService) GetTrendingEvents(ctx context.Context, limit int) ([]models.EventDTO, error) {
	return s.eventRepo.FindPopularEvents(ctx, limit)
}

func (s *eventService) GetRelatedEvents(ctx context.Context, eventID, locationID string, limit int) ([]models.EventDTO, error) {
	return s.eventRepo.FindRelatedEvents(ctx, eventID, locationID, limit)
}

func (s *eventService) SearchEvents(ctx context.Context, query string, opts base.ListOptions) ([]models.EventDTO, int, error) {
	// Attach search term
	opts.Search = query

	events, total, err := s.eventRepo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to search events",
			errors.WithContext("query", query),
		)
	}

	return events, total, nil
}
