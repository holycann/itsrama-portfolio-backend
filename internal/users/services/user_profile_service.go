package services

import (
	"context"
	"database/sql"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/achievement/services"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/holycann/cultour-backend/pkg/supabase"
)

type userProfileService struct {
	repo      repositories.UserProfileRepository
	userRepo  repositories.UserRepository
	userBadge repositories.UserBadgeRepository
	badge     services.BadgeService
	storage   *supabase.SupabaseStorage
}

func NewUserProfileService(
	repo repositories.UserProfileRepository,
	userRepo repositories.UserRepository,
	userBadge repositories.UserBadgeRepository,
	badge services.BadgeService,
	storage *supabase.SupabaseStorage,
) UserProfileService {
	return &userProfileService{
		repo:      repo,
		userRepo:  userRepo,
		userBadge: userBadge,
		badge:     badge,
		storage:   storage,
	}
}

func (s *userProfileService) CreateProfile(ctx context.Context, userProfile *models.UserProfileCreate) (*models.UserProfileDTO, error) {
	// Validate input
	if err := base.ValidateModel(userProfile); err != nil {
		return nil, err
	}

	// Validate user existence
	if userProfile.UserID == uuid.Nil {
		return nil, errors.New(errors.ErrValidation, "user ID is required", nil)
	}

	// Check if user exists
	userExists, err := s.userRepo.Exists(ctx, userProfile.UserID.String())
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error checking user existence")
	}
	if !userExists {
		return nil, errors.New(errors.ErrNotFound, fmt.Sprintf("user %s does not exist", userProfile.UserID), nil)
	}

	// Check if user profile exists by user id
	exists, err := s.repo.ExistsByUserID(ctx, userProfile.UserID.String())
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error checking profile existence")
	}
	if exists {
		return nil, errors.New(errors.ErrConflict, fmt.Sprintf("profile for user %s already exists", userProfile.UserID), nil)
	}

	// Check if Penjelajah badge exists
	penjelajahBadge, err := s.badge.GetBadgeByName(ctx, "Penjelajah")
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal, "error retrieving penjelajah badge")
	}

	// Prepare profile for creation
	now := time.Now().UTC()
	profile := &models.UserProfile{
		ID:               uuid.New(),
		UserID:           userProfile.UserID,
		Fullname:         userProfile.Fullname,
		Bio:              &userProfile.Bio,
		AvatarUrl:        &userProfile.AvatarUrl,
		IdentityImageUrl: &userProfile.IdentityImageUrl,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}

	// Create profile
	createdProfile, err := s.repo.Create(ctx, profile)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error creating user profile")
	}

	// Assign Penjelajah badge
	userBadgePayload := &models.UserBadge{
		UserID:    profile.UserID,
		BadgeID:   penjelajahBadge.ID,
		CreatedAt: &now,
	}

	err = s.userBadge.AddBadgeToUser(ctx, userBadgePayload)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal, "error assigning penjelajah badge")
	}

	dto := createdProfile.ToDTO()
	return &dto, nil
}

func (s *userProfileService) GetProfileByID(ctx context.Context, id string) (*models.UserProfileDTO, error) {
	// Validate input
	if id == "" {
		return nil, errors.New(errors.ErrValidation, "profile ID cannot be empty", nil)
	}

	// Retrieve profile
	profile, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.ErrNotFound, fmt.Sprintf("user profile with ID %s not found", id), err)
		}
		return nil, errors.Wrap(err, errors.ErrDatabase, "error retrieving user profile")
	}

	return profile, nil
}

func (s *userProfileService) ListProfiles(ctx context.Context, opts base.ListOptions) ([]models.UserProfileDTO, int, error) {
	// Set default values if not provided
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 10
	}

	profiles, total, err := s.repo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "error listing profiles")
	}

	return profiles, total, nil
}

func (s *userProfileService) UpdateProfile(ctx context.Context, userProfile *models.UserProfileUpdate) (*models.UserProfileDTO, error) {
	// Validate input
	if err := base.ValidateModel(userProfile); err != nil {
		return nil, err
	}

	// Check if profile exists
	existingProfile, err := s.repo.Exists(ctx, userProfile.ID.String())
	if err != nil {
		return nil, err
	}

	if !existingProfile {
		return nil, errors.New(errors.ErrNotFound, fmt.Sprintf("user profile with ID %s not found", userProfile.ID), nil)
	}

	// Prepare update
	updatePayload := &models.UserProfile{
		ID:       userProfile.ID,
		UserID:   userProfile.UserID,
		Fullname: userProfile.Fullname,
		Bio:      &userProfile.Bio,
	}

	// Perform update
	updatedProfile, err := s.repo.Update(ctx, updatePayload)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error updating profile")
	}

	dto := updatedProfile.ToDTO()
	return &dto, nil
}

func (s *userProfileService) UpdateProfileAvatar(ctx context.Context, payload *models.UserProfileAvatarUpdate) (*models.UserProfileDTO, error) {
	// Validate input
	if err := base.ValidateModel(payload); err != nil {
		return nil, err
	}

	// Check if profile exists
	existingProfile, err := s.GetProfileByID(ctx, payload.ID.String())
	if err != nil {
		return nil, err
	}

	// Upload avatar
	avatarUrl, err := s.updateAvatar(ctx, existingProfile.User.ID.String(), payload.Image)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal, "failed to update avatar")
	}

	payload.AvatarUrl = avatarUrl

	// Perform update
	err = s.repo.UpdateAvatarImage(ctx, payload)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error updating avatar")
	}

	existingProfile.AvatarUrl = &payload.AvatarUrl

	return existingProfile, nil
}

func (s *userProfileService) UpdateProfileIdentity(ctx context.Context, payload *models.UserProfileIdentityUpdate) (*models.UserProfileDTO, error) {
	// Validate input
	if err := base.ValidateModel(payload); err != nil {
		return nil, err
	}

	// Check if profile exists
	existingProfile, err := s.repo.FindByID(ctx, payload.ID.String())
	if err != nil {
		return nil, err
	}

	// Upload identity image
	identityUrl, err := s.updateIdentity(ctx, existingProfile.User.ID.String(), payload.Image)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal, "failed to update identity image")
	}

	payload.IdentityImageUrl = identityUrl

	// Perform update
	err = s.repo.VerifyIdentity(ctx, payload)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error updating identity image")
	}

	// Assign Warlok badge
	badge, err := s.badge.GetBadgeByName(ctx, "Warlok")
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal, "failed to get Warlok badge")
	}

	userBadgePayload := &models.UserBadge{
		UserID:  existingProfile.User.ID,
		BadgeID: badge.ID,
	}

	err = s.userBadge.AddBadgeToUser(ctx, userBadgePayload)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrInternal, "failed to create user badge")
	}

	existingProfile.IdentityImageUrl = &payload.IdentityImageUrl

	return existingProfile, nil
}

func (s *userProfileService) DeleteProfile(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return errors.New(errors.ErrValidation, "profile ID cannot be empty", nil)
	}

	// Check if profile exists
	_, err := s.GetProfileByID(ctx, id)
	if err != nil {
		return errors.New(errors.ErrNotFound, "Profile not found", err)
	}

	// Soft delete the profile
	return s.repo.Delete(ctx, id)
}

func (s *userProfileService) CountProfiles(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.repo.Count(ctx, filters)
}

func (s *userProfileService) GetProfileByUserID(ctx context.Context, userID string) (*models.UserProfileDTO, error) {
	// Validate input
	if userID == "" {
		return nil, errors.New(errors.ErrValidation, "user ID cannot be empty", nil)
	}

	// Check if user exists
	exists, err := s.repo.ExistsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error checking user existence")
	}
	if !exists {
		return nil, errors.New(errors.ErrNotFound, fmt.Sprintf("user with ID %s not found", userID), nil)
	}

	// Retrieve profile
	profile, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error retrieving user profile")
	}

	return profile, nil
}

func (s *userProfileService) GetProfileByFullname(ctx context.Context, fullname string) ([]models.UserProfileDTO, error) {
	// Validate input
	if fullname == "" {
		return nil, errors.New(errors.ErrValidation, "fullname cannot be empty", nil)
	}

	// Retrieve profiles
	profiles, err := s.repo.FindByFullname(ctx, fullname)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error retrieving user profiles")
	}

	return profiles, nil
}

func (s *userProfileService) SearchProfiles(ctx context.Context, opts base.ListOptions) ([]models.UserProfileDTO, int, error) {
	// Validate ListOptions
	if err := opts.Validate(); err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrValidation, "invalid list options")
	}

	profiles, total, err := s.repo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "error searching profiles")
	}

	return profiles, total, nil
}

// Private helper methods for file uploads using Supabase storage
func (s *userProfileService) updateAvatar(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	// Generate a unique path for the avatar
	avatarPath := "images/avatars/" + userID + filepath.Ext(file.Filename)

	// Upload the file to Supabase storage
	_, err := s.storage.Upload(ctx, file, avatarPath)
	if err != nil {
		return "", errors.Wrap(err, errors.ErrInternal, "failed to upload avatar")
	}

	// Get the public URL for the uploaded file
	publicURL, err := s.storage.GetPublicURL(avatarPath)
	if err != nil {
		return "", errors.Wrap(err, errors.ErrInternal, "failed to get avatar URL")
	}

	return publicURL, nil
}

func (s *userProfileService) updateIdentity(ctx context.Context, userID string, file *multipart.FileHeader) (string, error) {
	// Generate a unique path for the identity image
	identityPath := "images/identity/" + userID + filepath.Ext(file.Filename)

	// Upload the file to Supabase storage
	_, err := s.storage.Upload(ctx, file, identityPath)
	if err != nil {
		return "", errors.Wrap(err, errors.ErrInternal, "failed to upload identity image")
	}

	// Get the public URL for the uploaded file
	publicURL, err := s.storage.GetPublicURL(identityPath)
	if err != nil {
		return "", errors.Wrap(err, errors.ErrInternal, "failed to get identity image URL")
	}

	return publicURL, nil
}
