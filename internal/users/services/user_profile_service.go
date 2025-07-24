package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
)

type userProfileService struct {
	repo     repositories.UserProfileRepository
	userRepo repositories.UserRepository
}

func NewUserProfileService(
	repo repositories.UserProfileRepository,
) UserProfileService {
	return &userProfileService{
		repo: repo,
	}
}

func (s *userProfileService) CreateProfile(ctx context.Context, userProfile *models.UserProfile) error {
	// Validate input
	if userProfile == nil {
		return fmt.Errorf("user profile cannot be nil")
	}

	// Validate user existence
	if userProfile.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userProfile.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with ID %s not found", userProfile.UserID)
		}
		return fmt.Errorf("error checking user existence: %w", err)
	}

	// Check if profile already exists for this user
	exists, err := s.repo.ExistsByUserID(ctx, userProfile.UserID)
	if err != nil {
		return fmt.Errorf("error checking profile existence: %w", err)
	}
	if exists {
		return fmt.Errorf("profile for user %s already exists", userProfile.UserID)
	}

	// Set default values
	now := time.Now().UTC()
	userProfile.CreatedAt = &now
	userProfile.UpdatedAt = &now

	// Create profile
	return s.repo.Create(ctx, userProfile)
}

func (s *userProfileService) GetProfiles(ctx context.Context, limit, offset int) ([]*models.UserProfile, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve profiles
	profiles, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user profiles: %w", err)
	}

	fmt.Println("User Profiles:", profiles)

	// Convert to pointer slice
	profilePtrs := make([]*models.UserProfile, len(profiles))
	for i := range profiles {
		profilePtrs[i] = &profiles[i]
	}

	return profilePtrs, nil
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

func (s *userProfileService) GetProfileByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", userID)
		}
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}

	// Retrieve profile
	profile, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile for user %s not found", userID)
		}
		return nil, fmt.Errorf("error retrieving user profile: %w", err)
	}

	return profile, nil
}

func (s *userProfileService) UpdateProfile(ctx context.Context, userProfile *models.UserProfile) error {
	// Validate input
	if userProfile == nil {
		return fmt.Errorf("user profile cannot be nil")
	}
	if userProfile.ID == "" {
		return fmt.Errorf("profile ID is required for update")
	}

	// Check if profile exists
	existingProfile, err := s.GetProfileByID(ctx, userProfile.ID)
	if err != nil {
		return err
	}

	// Validate user existence if user ID is changed
	if userProfile.UserID != "" && userProfile.UserID != existingProfile.UserID {
		_, err := s.userRepo.FindByID(ctx, userProfile.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("user with ID %s not found", userProfile.UserID)
			}
			return fmt.Errorf("error checking user existence: %w", err)
		}
	}

	// Update timestamps
	now := time.Now().UTC()
	userProfile.UpdatedAt = &now

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
	return s.repo.SoftDelete(ctx, id)
}
