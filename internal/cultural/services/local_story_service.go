package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
	"github.com/holycann/cultour-backend/pkg/validator"
)

type localStoryService struct {
	localStoryRepo repositories.LocalStoryRepository
}

func NewLocalStoryService(localStoryRepo repositories.LocalStoryRepository) LocalStoryService {
	return &localStoryService{
		localStoryRepo: localStoryRepo,
	}
}

func (s *localStoryService) CreateLocalStory(ctx context.Context, localStory *models.LocalStory) error {
	// Validate local story object
	if localStory == nil {
		return fmt.Errorf("local story cannot be nil")
	}

	// Use centralized validator
	if err := validator.ValidateStruct(localStory); err != nil {
		return err
	}

	// Additional specific validations
	if err := validator.ValidateString(localStory.Title, "Title", 3, 100); err != nil {
		return err
	}
	if err := validator.ValidateString(localStory.Summary, "Summary", 10, 500); err != nil {
		return err
	}
	if err := validator.ValidateString(localStory.StoryText, "StoryText", 50, 5000); err != nil {
		return err
	}

	// Set default values
	localStory.ID = uuid.New()
	now := time.Now()
	localStory.CreatedAt = now
	localStory.UpdatedAt = now

	// Call repository to create local story
	return s.localStoryRepo.Create(ctx, localStory)
}

func (s *localStoryService) GetLocalStoryByID(ctx context.Context, id string) (*models.LocalStory, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("local story ID cannot be empty")
	}

	// Parse the ID to UUID
	storyID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid local story ID: %w", err)
	}

	// Retrieve local story from repository
	return s.localStoryRepo.FindByID(ctx, storyID.String())
}

func (s *localStoryService) ListLocalStories(ctx context.Context, opts repository.ListOptions) ([]models.LocalStory, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.localStoryRepo.List(ctx, opts)
}

func (s *localStoryService) UpdateLocalStory(ctx context.Context, localStory *models.LocalStory) error {
	// Validate local story object
	if localStory == nil {
		return fmt.Errorf("local story cannot be nil")
	}

	// Validate required fields
	if localStory.ID == uuid.Nil {
		return fmt.Errorf("local story ID is required for update")
	}

	// Update timestamp
	localStory.UpdatedAt = time.Now()

	// Call repository to update local story
	return s.localStoryRepo.Update(ctx, localStory)
}

func (s *localStoryService) DeleteLocalStory(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("local story ID cannot be empty")
	}

	// Parse the ID to UUID
	storyID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid local story ID: %w", err)
	}

	// Call repository to delete local story
	return s.localStoryRepo.Delete(ctx, storyID.String())
}

func (s *localStoryService) CountLocalStories(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.localStoryRepo.Count(ctx, filters)
}

func (s *localStoryService) GetLocalStoriesByLocation(ctx context.Context, locationID string) ([]*models.LocalStory, error) {
	// Convert string to UUID
	locUUID, err := uuid.Parse(locationID)
	if err != nil {
		return nil, fmt.Errorf("invalid location ID: %w", err)
	}

	return s.localStoryRepo.FindStoriesByLocation(ctx, locUUID)
}

func (s *localStoryService) GetLocalStoriesByOriginCulture(ctx context.Context, culture string) ([]*models.LocalStory, error) {
	return s.localStoryRepo.FindStoriesByOriginCulture(ctx, culture)
}

func (s *localStoryService) SearchLocalStories(ctx context.Context, query string, opts repository.ListOptions) ([]models.LocalStory, error) {
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
			Field:    "title",
			Operator: "like",
			Value:    query,
		},
		repository.FilterOption{
			Field:    "story_text",
			Operator: "like",
			Value:    query,
		},
	)

	return s.localStoryRepo.List(ctx, opts)
}
