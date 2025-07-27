package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type userProfileService struct {
	repo     repositories.UserProfileRepository
	userRepo repositories.UserRepository
}

func NewUserProfileService(
	repo repositories.UserProfileRepository,
	userRepo repositories.UserRepository,
) UserProfileService {
	return &userProfileService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *userProfileService) CreateProfile(ctx context.Context, userProfile *models.UserProfile) error {
	// Validate input
	if userProfile == nil {
		return fmt.Errorf("user profile cannot be nil")
	}

	// Validate user existence
	if userProfile.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userProfile.UserID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with ID %s not found", userProfile.UserID)
		}
		return fmt.Errorf("error checking user existence: %w", err)
	}

	// Check if profile already exists for this user
	exists, err := s.repo.ExistsByUserID(ctx, userProfile.UserID.String())
	if err != nil {
		return fmt.Errorf("error checking profile existence: %w", err)
	}
	if exists {
		return fmt.Errorf("profile for user %s already exists", userProfile.UserID)
	}

	// Set default values
	userProfile.ID = uuid.New()
	now := time.Now().UTC()
	userProfile.CreatedAt = now
	userProfile.UpdatedAt = now

	// Create profile
	return s.repo.Create(ctx, userProfile)
}

func (s *userProfileService) GetProfileByID(ctx context.Context, id string) (*models.UserProfile, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("profile ID cannot be empty")
	}

	// Retrieve profile
	profile, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user profile with ID %s not found", id)
		}
		return nil, fmt.Errorf("error retrieving user profile: %w", err)
	}

	return profile, nil
}

func (s *userProfileService) ListProfiles(ctx context.Context, opts repository.ListOptions) ([]models.UserProfile, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.repo.List(ctx, opts)
}

func (s *userProfileService) UpdateProfile(ctx context.Context, userProfile *models.UserProfile) error {
	// Validate input
	if userProfile == nil {
		return fmt.Errorf("user profile cannot be nil")
	}
	if userProfile.ID == uuid.Nil {
		return fmt.Errorf("profile ID is required for update")
	}

	// Check if profile exists
	existingProfile, err := s.GetProfileByID(ctx, userProfile.ID.String())
	if err != nil {
		return err
	}

	// Validate user existence if user ID is changed
	if userProfile.UserID != uuid.Nil && userProfile.UserID != existingProfile.UserID {
		_, err := s.userRepo.FindByID(ctx, userProfile.UserID.String())
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("user with ID %s not found", userProfile.UserID)
			}
			return fmt.Errorf("error checking user existence: %w", err)
		}
	}

	// Update timestamps
	now := time.Now().UTC()
	userProfile.UpdatedAt = now

	// Preserve creation timestamp
	userProfile.CreatedAt = existingProfile.CreatedAt

	// Perform update
	return s.repo.Update(ctx, userProfile)
}

func (s *userProfileService) DeleteProfile(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return fmt.Errorf("profile ID cannot be empty")
	}

	// Check if profile exists
	_, err := s.GetProfileByID(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete the profile
	return s.repo.Delete(ctx, id)
}

func (s *userProfileService) CountProfiles(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.repo.Count(ctx, filters)
}

func (s *userProfileService) GetProfileByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// Convert string to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Check if user exists
	_, err = s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", userID)
		}
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}

	// Retrieve profile
	profile, err := s.repo.FindByField(ctx, "user_id", userUUID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user profile: %w", err)
	}

	if len(profile) == 0 {
		return nil, fmt.Errorf("profile for user %s not found", userID)
	}

	return &profile[0], nil
}

func (s *userProfileService) SearchProfiles(ctx context.Context, query string, opts repository.ListOptions) ([]models.UserProfile, error) {
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
			Field:    "fullname",
			Operator: "like",
			Value:    query,
		},
		repository.FilterOption{
			Field:    "bio",
			Operator: "like",
			Value:    query,
		},
	)

	return s.repo.List(ctx, opts)
}
