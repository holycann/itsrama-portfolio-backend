package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/internal/achievement/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type badgeService struct {
	repo repositories.BadgeRepository
}

func NewBadgeService(repo repositories.BadgeRepository) BadgeService {
	return &badgeService{
		repo: repo,
	}
}

func (s *badgeService) CreateBadge(ctx context.Context, badgeCreate *models.BadgeCreate) error {
	// Validate badge creation
	if badgeCreate.Name == "" {
		return fmt.Errorf("badge name is required")
	}

	now := time.Now()
	badge := &models.Badge{
		ID:          uuid.New(),
		Name:        badgeCreate.Name,
		Description: badgeCreate.Description,
		IconURL:     badgeCreate.IconURL,
		CreatedAt:   now,
	}

	return s.repo.Create(ctx, badge)
}

func (s *badgeService) GetBadgeByID(ctx context.Context, id string) (*models.Badge, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *badgeService) ListBadges(ctx context.Context, opts repository.ListOptions) ([]models.Badge, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.repo.List(ctx, opts)
}

func (s *badgeService) UpdateBadge(ctx context.Context, id string, badgeUpdate *models.BadgeCreate) error {
	// First, retrieve the existing badge
	existingBadge, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("badge not found: %w", err)
	}

	// Update fields
	now := time.Now()
	existingBadge.Name = badgeUpdate.Name
	existingBadge.Description = badgeUpdate.Description
	existingBadge.IconURL = badgeUpdate.IconURL
	existingBadge.UpdatedAt = &now

	return s.repo.Update(ctx, existingBadge)
}

func (s *badgeService) DeleteBadge(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *badgeService) CountBadges(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.repo.Count(ctx, filters)
}

func (s *badgeService) GetBadgeByName(ctx context.Context, name string) (*models.Badge, error) {
	return s.repo.FindBadgeByName(ctx, name)
}

func (s *badgeService) GetPopularBadges(ctx context.Context, limit int) ([]models.Badge, error) {
	return s.repo.FindPopularBadges(ctx, limit)
}
