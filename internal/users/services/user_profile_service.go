package services

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/supabase"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
	storage_go "github.com/supabase-community/storage-go"
)

type userProfileService struct {
	repo     repositories.UserProfileRepository
	userRepo repositories.UserRepository
	storage  *supabase.SupabaseStorage
}

func NewUserProfileService(
	repo repositories.UserProfileRepository,
	userRepo repositories.UserRepository,
	storage *supabase.SupabaseStorage,
) UserProfileService {
	return &userProfileService{
		repo:     repo,
		userRepo: userRepo,
		storage:  storage,
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

func (s *userProfileService) UpdateProfile(ctx context.Context, userProfile *models.UserProfile, avatar *multipart.FileHeader, identity *multipart.FileHeader) error {
	// Validate input
	if userProfile == nil {
		return fmt.Errorf("user profile cannot be nil")
	}
	if userProfile.ID == uuid.Nil || userProfile.UserID == uuid.Nil {
		return fmt.Errorf("profile ID and user ID is required for update")
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

	// Merge: if value ada di userProfile, update, jika tidak ada pake existingProfile
	mergedProfile := *existingProfile // start with existing

	// Only update fields if userProfile has a non-zero value
	if userProfile.Fullname != "" {
		mergedProfile.Fullname = userProfile.Fullname
	}
	if userProfile.Bio != "" {
		mergedProfile.Bio = userProfile.Bio
	}
	if userProfile.AvatarUrl != "" {
		mergedProfile.AvatarUrl = userProfile.AvatarUrl
	}
	if userProfile.IdentityImageUrl != "" {
		mergedProfile.IdentityImageUrl = userProfile.IdentityImageUrl
	}
	// Always update UpdatedAt
	mergedProfile.UpdatedAt = time.Now().UTC()

	// If avatar file is provided, upload and update AvatarUrl
	if avatar != nil {
		avatarUrl, err := s.UpdateAvatar(ctx, mergedProfile.UserID.String(), avatar)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return fmt.Errorf("failed to update avatar")
		}
		mergedProfile.AvatarUrl = avatarUrl
	}

	// If identity file is provided, upload and update IdentityImageUrl
	if identity != nil {
		identityUrl, err := s.UpdateIdentity(ctx, mergedProfile.UserID.String(), identity)
		if err != nil {
			fmt.Println("Error:", err.Error())
			return fmt.Errorf("failed to update identity image")
		}
		mergedProfile.IdentityImageUrl = identityUrl
	}

	// Perform update
	return s.repo.Update(ctx, &mergedProfile)
}

// UpdateAvatar uploads the avatar file and returns the public URL
func (s *userProfileService) UpdateAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	// Validate input
	if userID == "" {
		return "", fmt.Errorf("user ID is required")
	}
	if file == nil {
		return "", fmt.Errorf("file data is required")
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Get file extension from the uploaded file's filename
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg" // default to jpg if extension is missing
	}

	// Build the destination path
	destPath := s.storage.GetDefaultFolder() + "/avatar/" + userID + ext

	// Upload file and get the response
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

// UpdateIdentity uploads the identity image file and returns the public URL
func (s *userProfileService) UpdateIdentity(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	// Validate input
	if userID == "" {
		return "", fmt.Errorf("user ID is required")
	}
	if file == nil {
		return "", fmt.Errorf("file data is required")
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Get file extension from the uploaded file's filename
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg" // default to jpg if extension is missing
	}

	// Build the destination path
	destPath := s.storage.GetDefaultFolder() + "/identity/" + userID + ext

	// Upload file and get the response
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
