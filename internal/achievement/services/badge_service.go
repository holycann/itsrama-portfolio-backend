package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/internal/achievement/repositories"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
)

type badgeService struct {
	repo repositories.BadgeRepository
}

func NewBadgeService(repo repositories.BadgeRepository) BadgeService {
	return &badgeService{
		repo: repo,
	}
}

func (s *badgeService) CreateBadge(ctx context.Context, badgeCreate *models.BadgeCreate) (*models.BadgeDTO, error) {
	// Validate badge creation
	if badgeCreate.Name == "" {
		return nil, errors.New(errors.ErrValidation, "badge name is required", nil)
	}

	now := time.Now()
	badge := &models.Badge{
		ID:          uuid.New(),
		Name:        badgeCreate.Name,
		Description: badgeCreate.Description,
		IconURL:     badgeCreate.IconURL,
		CreatedAt:   &now,
	}

	badge, err := s.repo.Create(ctx, badge)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create badge")
	}

	dto := badge.ToDTO()
	return &dto, nil
}

func (s *badgeService) GetBadgeByID(ctx context.Context, id string) (*models.BadgeDTO, error) {
	badge, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to retrieve badge")
	}

	dto := badge.ToDTO()
	return &dto, nil
}

func (s *badgeService) ListBadges(ctx context.Context, opts base.ListOptions) ([]models.BadgeDTO, int, error) {
	// Set default values if not provided
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 10
	}

	badges, total, err := s.repo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to list badges")
	}

	// Convert to DTOs
	badgeDTOs := make([]models.BadgeDTO, len(badges))
	for i, badge := range badges {
		badgeDTOs[i] = badge.ToDTO()
	}

	return badgeDTOs, total, nil
}

func (s *badgeService) UpdateBadge(ctx context.Context, id string, badgeUpdate *models.BadgeUpdate) (*models.BadgeDTO, error) {
	// First, retrieve the existing badge
	existingBadge, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrNotFound, "badge not found")
	}

	// Update fields
	now := time.Now()
	existingBadge.Name = badgeUpdate.Name
	existingBadge.Description = badgeUpdate.Description
	existingBadge.IconURL = badgeUpdate.IconURL
	existingBadge.UpdatedAt = &now

	badge, err := s.repo.Update(ctx, existingBadge)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update badge")
	}

	dto := badge.ToDTO()
	return &dto, nil
}

func (s *badgeService) DeleteBadge(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *badgeService) CountBadges(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.repo.Count(ctx, filters)
}

func (s *badgeService) GetBadgeByName(ctx context.Context, name string) (*models.BadgeDTO, error) {
	badge, err := s.repo.FindBadgeByName(ctx, name)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrNotFound, "badge not found")
	}

	dto := badge.ToDTO()
	return &dto, nil
}

func (s *badgeService) GetPopularBadges(ctx context.Context, limit int) ([]models.BadgeDTO, error) {
	badges, err := s.repo.FindPopularBadges(ctx, limit)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to retrieve popular badges")
	}

	// Convert to DTOs
	badgeDTOs := make([]models.BadgeDTO, len(badges))
	for i, badge := range badges {
		badgeDTOs[i] = badge.ToDTO()
	}

	return badgeDTOs, nil
}

func (s *badgeService) SearchBadges(ctx context.Context, opts base.ListOptions) ([]models.BadgeDTO, int, error) {
	// Validate ListOptions
	if err := opts.Validate(); err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrValidation, "invalid list options")
	}

	badges, total, err := s.repo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to search badges")
	}

	// Convert to DTOs
	badgeDTOs := make([]models.BadgeDTO, len(badges))
	for i, badge := range badges {
		badgeDTOs[i] = badge.ToDTO()
	}

	return badgeDTOs, total, nil
}
