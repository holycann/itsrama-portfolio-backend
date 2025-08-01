package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	"github.com/holycann/cultour-backend/internal/place/services"
	"github.com/holycann/cultour-backend/internal/supabase"
	"github.com/holycann/cultour-backend/pkg/repository"
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

	// Validate all required fields
	if event.UserID == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	if event.CityID == uuid.Nil {
		return fmt.Errorf("city_id is required")
	}
	if event.ProvinceID == uuid.Nil {
		return fmt.Errorf("province_id is required")
	}
	if event.Name == "" {
		return fmt.Errorf("name is required")
	}
	if event.Description == "" {
		return fmt.Errorf("description is required")
	}
	if event.StartDate.IsZero() {
		return fmt.Errorf("start_date is required")
	}
	if event.EndDate.IsZero() {
		return fmt.Errorf("end_date is required")
	}

	// Cek apakah lokasi dengan nama tersebut sudah ada, jika belum ada maka create location
	location, err := s.locationService.GetLocationByName(ctx, event.Location.Name)
	if err != nil || location == nil {
		// Location not found, create new location
		if event.Location != nil {
			err := s.locationService.CreateLocation(ctx, event.Location)
			if err != nil {
				return fmt.Errorf("failed to create location: %v", err)
			}
			location = event.Location
		} else {
			return fmt.Errorf("location data is required")
		}
	}
	event.Location = location

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

func (s *eventService) UpdateEvent(ctx context.Context, event *models.RequestEvent) error {
	// Validate event object
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Validate required fields
	if event.ID == uuid.Nil {
		return fmt.Errorf("event ID is required for update")
	}

	// Convert RequestEvent to Event
	eventData := models.Event{
		ID:            event.ID,
		UserID:        event.UserID,
		LocationID:    event.Location.ID,
		CityID:        event.CityID,
		ProvinceID:    event.ProvinceID,
		Name:          event.Name,
		Description:   event.Description,
		StartDate:     event.StartDate,
		EndDate:       event.EndDate,
		IsKidFriendly: event.IsKidFriendly,
		UpdatedAt:     time.Now(),
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

func (s *eventService) UpdateEventViews(ctx context.Context, id string) string {
	if id == "" {
		return "event ID cannot be empty"
	}
	return s.eventRepo.UpdateViews(ctx, id)
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
