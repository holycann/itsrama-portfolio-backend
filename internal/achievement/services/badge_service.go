package services

import (
	"context"
	"fmt"
	"time"

	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/internal/achievement/repositories"
)

// DefaultBadgeService implements BadgeService
type DefaultBadgeService struct {
	repo repositories.BadgeRepository
}

// NewBadgeService creates a new instance of BadgeService
func NewBadgeService(repo repositories.BadgeRepository) BadgeService {
	return &DefaultBadgeService{
		repo: repo,
	}
}

// CreateBadge creates a new badge in the system
func (s *DefaultBadgeService) CreateBadge(ctx context.Context, badgeCreate *models.BadgeCreate) (*models.Badge, error) {
	// Validate badge creation
	if badgeCreate.Name == "" {
		return nil, fmt.Errorf("badge name is required")
	}

	// Generate unique ID (you might want to use a more robust ID generation method)
	now := time.Now()
	badge := &models.Badge{
		ID:          fmt.Sprintf("badge_%s", badgeCreate.Name),
		Name:        badgeCreate.Name,
		Description: badgeCreate.Description,
		IconURL:     badgeCreate.IconURL,
		CreatedAt:   &now,
	}

	return s.repo.Create(ctx, badge)
}

// GetBadgeByID retrieves a badge by its unique identifier
func (s *DefaultBadgeService) GetBadgeByID(ctx context.Context, id string) (*models.Badge, error) {
	return s.repo.FindByID(ctx, id)
}

// ListBadges retrieves badges with pagination
func (s *DefaultBadgeService) ListBadges(ctx context.Context, limit, offset int) ([]models.Badge, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.FindAll(ctx, limit, offset)
}

// UpdateBadge updates an existing badge
func (s *DefaultBadgeService) UpdateBadge(ctx context.Context, id string, badgeUpdate *models.BadgeCreate) (*models.Badge, error) {
	// First, retrieve the existing badge
	existingBadge, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("badge not found: %w", err)
	}

	// Update fields
	now := time.Now()
	existingBadge.Name = badgeUpdate.Name
	existingBadge.Description = badgeUpdate.Description
	existingBadge.IconURL = badgeUpdate.IconURL
	existingBadge.UpdatedAt = &now

	return s.repo.Update(ctx, existingBadge)
}

// DeleteBadge removes a badge from the system
func (s *DefaultBadgeService) DeleteBadge(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// CountBadges returns the total number of badges
func (s *DefaultBadgeService) CountBadges(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}
