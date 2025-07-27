package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type userBadgeService struct {
	repo repositories.UserBadgeRepository
}

func NewUserBadgeService(repo repositories.UserBadgeRepository) UserBadgeService {
	return &userBadgeService{
		repo: repo,
	}
}

func (s *userBadgeService) CreateUserBadge(ctx context.Context, userBadgeCreate *models.UserBadgeCreate) error {
	// Validate input
	if userBadgeCreate == nil {
		return fmt.Errorf("user badge creation data cannot be nil")
	}

	// Validate user ID
	if userBadgeCreate.UserID == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	// Validate badge ID
	if userBadgeCreate.BadgeID == uuid.Nil {
		return fmt.Errorf("badge ID is required")
	}

	// Check if the user already has this badge
	existingBadges, err := s.repo.FindByField(ctx, "user_id", userBadgeCreate.UserID)
	if err != nil {
		return fmt.Errorf("error checking existing badge: %w", err)
	}

	// If badge already exists, return an error
	for _, badge := range existingBadges {
		if badge.BadgeID == userBadgeCreate.BadgeID {
			return fmt.Errorf("badge %s already assigned to user", userBadgeCreate.BadgeID)
		}
	}

	// Create new user badge
	now := time.Now()
	newBadge := &models.UserBadge{
		ID:        uuid.New(),
		UserID:    userBadgeCreate.UserID,
		BadgeID:   userBadgeCreate.BadgeID,
		CreatedAt: now,
	}

	return s.repo.Create(ctx, newBadge)
}

func (s *userBadgeService) GetUserBadgeByID(ctx context.Context, id string) (*models.UserBadge, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("user badge ID cannot be empty")
	}

	return s.repo.FindByID(ctx, id)
}

func (s *userBadgeService) ListUserBadges(ctx context.Context, opts repository.ListOptions) ([]models.UserBadge, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.repo.List(ctx, opts)
}

func (s *userBadgeService) DeleteUserBadge(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return fmt.Errorf("user badge ID cannot be empty")
	}

	return s.repo.Delete(ctx, id)
}

func (s *userBadgeService) CountUserBadges(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.repo.Count(ctx, filters)
}

func (s *userBadgeService) GetUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadge, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// Convert string to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return s.repo.FindUserBadgesByUser(ctx, userUUID.String())
}

func (s *userBadgeService) GetUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadge, error) {
	// Validate input
	if badgeID == "" {
		return nil, fmt.Errorf("badge ID cannot be empty")
	}

	// Convert string to UUID
	badgeUUID, err := uuid.Parse(badgeID)
	if err != nil {
		return nil, fmt.Errorf("invalid badge ID: %w", err)
	}

	return s.repo.FindUserBadgesByBadge(ctx, badgeUUID.String())
}

func (s *userBadgeService) SearchUserBadges(ctx context.Context, query string, opts repository.ListOptions) ([]models.UserBadge, error) {
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
			Field:    "user_id",
			Operator: "like",
			Value:    query,
		},
		repository.FilterOption{
			Field:    "badge_id",
			Operator: "like",
			Value:    query,
		},
	)

	return s.repo.List(ctx, opts)
}
