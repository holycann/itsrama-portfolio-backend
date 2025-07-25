package services

import (
	"context"
	"fmt"
	"time"

	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/internal/cultural/repositories"
	"github.com/holycann/cultour-backend/pkg/utils"
)

type localStoryService struct {
	localStoryRepo repositories.LocalStoryRepository
}

// NewLocalStoryService creates a new instance of the local story service
// with the given local story repository.
func NewLocalStoryService(localStoryRepo repositories.LocalStoryRepository) LocalStoryService {
	return &localStoryService{
		localStoryRepo: localStoryRepo,
	}
}

// CreateLocalStory adds a new local story to the database
// Validates the local story object before creating
func (s *localStoryService) CreateLocalStory(ctx context.Context, localStory *models.LocalStory) error {
	// Validate local story object
	if localStory == nil {
		return fmt.Errorf("local story cannot be nil")
	}

	// Validate required fields (example validation)
	if localStory.Title == "" {
		return fmt.Errorf("local story title is required")
	}

	// Generate UUID if not provided
	localStory.ID = utils.GenerateUUIDIfEmpty(localStory.ID)

	// Set timestamps
	localStory.CreatedAt = time.Now()
	localStory.UpdatedAt = time.Now()

	// Call repository to create local story
	return s.localStoryRepo.Create(ctx, localStory)
}

// GetLocalStories retrieves a list of local stories with pagination
func (s *localStoryService) GetLocalStories(ctx context.Context, limit, offset int) ([]*models.LocalStory, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve local stories from repository
	localStories, err := s.localStoryRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert []models.LocalStory to []*models.LocalStory
	localStoryPtrs := make([]*models.LocalStory, len(localStories))
	for i := range localStories {
		localStoryPtrs[i] = &localStories[i]
	}

	return localStoryPtrs, nil
}

// GetLocalStoryByID retrieves a single local story by its unique identifier
func (s *localStoryService) GetLocalStoryByID(ctx context.Context, id string) (*models.LocalStory, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("local story ID cannot be empty")
	}

	// Retrieve local story from repository
	return s.localStoryRepo.FindByID(ctx, id)
}

// GetLocalStoryByTitle retrieves a local story by its title
// Note: This method is not directly supported by the current repository implementation
// You might need to add a custom method in the repository or implement filtering
func (s *localStoryService) GetLocalStoryByTitle(ctx context.Context, title string) (*models.LocalStory, error) {
	// Validate title
	if title == "" {
		return nil, fmt.Errorf("local story title cannot be empty")
	}

	// Since the current repository doesn't have a direct method for this,
	// we'll use a workaround by listing all local stories and finding by title
	localStories, err := s.localStoryRepo.List(ctx, 1, 0)
	if err != nil {
		return nil, err
	}

	// Find local story by title (linear search)
	for _, localStory := range localStories {
		if localStory.Title == title {
			return &localStory, nil
		}
	}

	return nil, fmt.Errorf("local story with title %s not found", title)
}

// UpdateLocalStory updates an existing local story in the database
func (s *localStoryService) UpdateLocalStory(ctx context.Context, localStory *models.LocalStory) error {
	// Validate local story object
	if localStory == nil {
		return fmt.Errorf("local story cannot be nil")
	}

	// Validate required fields
	if localStory.ID == "" {
		return fmt.Errorf("local story ID is required for update")
	}

	// Update timestamp
	localStory.UpdatedAt = time.Now()

	// Call repository to update local story
	return s.localStoryRepo.Update(ctx, localStory)
}

// DeleteLocalStory removes a local story from the database by its ID
func (s *localStoryService) DeleteLocalStory(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("local story ID cannot be empty")
	}

	// Call repository to delete local story
	return s.localStoryRepo.Delete(ctx, id)
}

// Count calculates the total number of local stories stored
func (s *localStoryService) Count(ctx context.Context) (int, error) {
	return s.localStoryRepo.Count(ctx)
}
