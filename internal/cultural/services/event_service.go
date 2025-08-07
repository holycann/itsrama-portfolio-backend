package services

import (
	"context"
	"fmt"
	"math"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	placeModels "github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/internal/supabase"
	"github.com/holycann/cultour-backend/pkg/repository"
	"github.com/holycann/cultour-backend/pkg/validator"
	storage_go "github.com/supabase-community/storage-go"
)

type eventService struct {
	eventRepo       repositories.EventRepository
	locationService services.LocationService
	storage         *supabase.SupabaseStorage
}

func NewEventService(eventRepo repositories.EventRepository, service services.LocationService, storage *supabase.SupabaseStorage) EventService {
	return &eventService{
		eventRepo:       eventRepo,
		storage:         storage,
		locationService: service,
	}
}

func (s *eventService) CreateEvent(ctx context.Context, event *models.RequestEvent, image *multipart.FileHeader) error {
	// Validate event object
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Use centralized validator
	if err := validator.ValidateStruct(event); err != nil {
		return err
	}

	// Additional specific validations
	if event.StartDate.IsZero() {
		return fmt.Errorf("start_date is required")
	}
	if event.EndDate.IsZero() {
		return fmt.Errorf("end_date is required")
	}

	tolerance := 0.0001
	existingLocations, err := s.locationService.ListLocations(ctx, repository.ListOptions{
		Filters: []repository.FilterOption{
			{
				Field:    "name",
				Operator: "eq",
				Value:    event.Location.Name,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to check existing locations: %v", err)
	}

	var matchedLocation *placeModels.Location
	for _, loc := range existingLocations {
		if math.Abs(loc.Latitude-event.Location.Latitude) <= tolerance &&
			math.Abs(loc.Longitude-event.Location.Longitude) <= tolerance {
			matchedLocation = &loc
			break
		}
	}

	if matchedLocation == nil {
		event.Location.ID = uuid.New()
		event.Location.CityID = event.CityID
		if err := s.locationService.CreateLocation(ctx, event.Location); err != nil {
			return fmt.Errorf("failed to create new location: %v", err)
		}

	} else {
		event.Location = matchedLocation
	}

	eventData := models.Event{
		ID:            uuid.New(),
		UserID:        event.UserID,
		LocationID:    event.Location.ID,
		CityID:        event.CityID,
		ProvinceID:    event.ProvinceID,
		Name:          event.Name,
		Description:   event.Description,
		StartDate:     event.StartDate,
		EndDate:       event.EndDate,
		IsKidFriendly: event.IsKidFriendly,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	imageUrl, err := s.uploadEventImage(ctx, eventData.ID.String(), image)
	if err != nil {
		return fmt.Errorf("failed to upload event image: %v", err)
	}

	eventData.ImageURL = imageUrl
	event.ImageURL = imageUrl
	event.ID = eventData.ID

	// Call repository to create event
	return s.eventRepo.Create(ctx, &eventData)
}

func (s *eventService) uploadEventImage(ctx context.Context, eventID string, file *multipart.FileHeader) (string, error) {
	// Validate input
	if eventID == "" {
		return "", fmt.Errorf("event ID cannot be empty")
	}
	if file == nil {
		return "", fmt.Errorf("file data is required")
	}

	// Open the file
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Get file extension from the uploaded file's filename
	// (default to .jpg if missing)
	ext := ".jpg"
	if file.Filename != "" {
		for i := len(file.Filename) - 1; i >= 0; i-- {
			if file.Filename[i] == '.' {
				ext = file.Filename[i:]
				break
			}
		}
	}

	// Build the destination path (assuming you have a storage service, adjust as needed)
	destPath := s.storage.GetDefaultFolder() + "/images/" + "event/" + eventID + ext

	result, _ := s.storage.GetClient().UploadFile(s.storage.GetBucketID(), destPath, f, storage_go.FileOptions{
		ContentType: func(s string) *string { return &s }("image"),
		Upsert:      func(b bool) *bool { return &b }(true),
	})
	if result.Key == "" {
		return "", fmt.Errorf("failed to upload file: %v", result)
	}
	url := s.storage.GetClient().GetPublicUrl(s.storage.GetBucketID(), destPath)
	if url.SignedURL == "" {
		return "", fmt.Errorf("failed to get public url: %v", url)
	}

	return url.SignedURL, nil
}

func (s *eventService) deleteEventImage(ctx context.Context, path string) error {
	_, err := s.storage.GetClient().RemoveFile(s.storage.GetBucketID(), []string{path})
	if err != nil {
		return fmt.Errorf("failed to delete event image: %v", err)
	}
	return nil
}

func (s *eventService) GetEventByID(ctx context.Context, id string) (*models.ResponseEvent, error) {
	if id == "" {
		return nil, fmt.Errorf("event ID cannot be empty")
	}
	event, err := s.eventRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *eventService) ListEvents(ctx context.Context, opts repository.ListOptions) ([]models.ResponseEvent, error) {
	// Modify the list options to include full details
	opts.Limit = opts.Limit
	opts.Offset = opts.Offset
	opts.SortBy = opts.SortBy
	opts.SortOrder = opts.SortOrder

	// Use the new method to fetch events with full details
	events, _, err := s.eventRepo.Search(ctx, opts)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *eventService) UpdateEvent(ctx context.Context, event *models.RequestEvent, image *multipart.FileHeader) error {
	// Validate event object
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Validate required fields
	if event.ID == uuid.Nil || event.UserID == uuid.Nil {
		return fmt.Errorf("event ID and user ID are required for update")
	}

	existingEvent, err := s.eventRepo.FindByID(ctx, event.ID.String())
	if err != nil {
		return fmt.Errorf("failed to get existing event: %v", err)
	}

	// Convert RequestEvent to Event
	eventData := models.Event{
		ID:     event.ID,
		UserID: event.UserID,
		CityID: func() uuid.UUID {
			if event.CityID != uuid.Nil && event.CityID != existingEvent.CityID {
				return event.CityID
			}
			return existingEvent.CityID
		}(),
		ProvinceID: func() uuid.UUID {
			if event.ProvinceID != uuid.Nil && event.ProvinceID != existingEvent.ProvinceID {
				return event.ProvinceID
			}
			return existingEvent.ProvinceID
		}(),
		Name: func() string {
			if event.Name != "" && event.Name != existingEvent.Name {
				return event.Name
			}
			return existingEvent.Name
		}(),
		Description: func() string {
			if event.Description != "" && event.Description != existingEvent.Description {
				return event.Description
			}
			return existingEvent.Description
		}(),
		StartDate: func() time.Time {
			if !event.StartDate.IsZero() && event.StartDate != existingEvent.StartDate {
				return event.StartDate
			}
			return existingEvent.StartDate
		}(),
		EndDate: func() time.Time {
			if !event.EndDate.IsZero() && event.EndDate != existingEvent.EndDate {
				return event.EndDate
			}
			return existingEvent.EndDate
		}(),
		IsKidFriendly: func() bool {
			if event.IsKidFriendly != existingEvent.IsKidFriendly {
				return event.IsKidFriendly
			}
			return existingEvent.IsKidFriendly
		}(),
		UpdatedAt: time.Now(),
		ImageURL:  existingEvent.ImageURL,
	}

	if image != nil {
		var imagePath string
		parts := strings.SplitN(existingEvent.ImageURL, "public/cultour/", 2)
		if len(parts) == 2 {
			imagePath = parts[1]
		} else {
			return fmt.Errorf("prefix not found")
		}

		err := s.deleteEventImage(ctx, imagePath)
		if err != nil {
			return fmt.Errorf("failed to delete existing event image: %v", err)
		}
		imageUrl, err := s.uploadEventImage(ctx, eventData.ID.String(), image)
		if err != nil {
			return fmt.Errorf("failed to upload event image: %v", err)
		}
		eventData.ImageURL = imageUrl
	}

	tolerance := 0.0001

	if event.Location == nil {
		return fmt.Errorf("location data is required for updating LocationID")
	}

	existingLocations, err := s.locationService.ListLocations(ctx, repository.ListOptions{
		Filters: []repository.FilterOption{
			{
				Field:    "name",
				Operator: "eq",
				Value:    event.Location.Name,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to check existing location: %v", err)
	}

	var matchedLocation *placeModels.Location
	for _, loc := range existingLocations {
		if math.Abs(loc.Latitude-event.Location.Latitude) <= tolerance &&
			math.Abs(loc.Longitude-event.Location.Longitude) <= tolerance {
			matchedLocation = &loc
			break
		}
	}

	if matchedLocation == nil {
		event.Location.ID = uuid.New()
		event.Location.CityID = event.CityID
		if err := s.locationService.CreateLocation(ctx, event.Location); err != nil {
			return fmt.Errorf("failed to create new location: %v", err)
		}
		eventData.LocationID = event.Location.ID
	} else {
		eventData.LocationID = matchedLocation.ID
	}

	// Call repository to update event
	return s.eventRepo.Update(ctx, &eventData)
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

func (s *eventService) UpdateEventViews(ctx context.Context, userID, eventID string) string {
	if userID == "" || eventID == "" {
		return "user ID and event ID cannot be empty"
	}
	return s.eventRepo.UpdateViews(ctx, userID, eventID)
}

func (s *eventService) GetTrendingEvents(ctx context.Context, limit int) ([]models.ResponseEvent, error) {
	return s.eventRepo.FindPopularEvents(ctx, limit)
}

func (s *eventService) GetRelatedEvents(ctx context.Context, eventID string, limit int) ([]models.ResponseEvent, error) {
	return s.eventRepo.FindRelatedEvents(ctx, eventID, limit)
}

func (s *eventService) SearchEvents(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseEvent, error) {
	// Modify the list options to include search query
	opts.SearchQuery = query

	events, _, err := s.eventRepo.Search(ctx, opts)
	if err != nil {
		return nil, err
	}

	return events, nil
}
