package services

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
	"github.com/holycann/cultour-backend/pkg/base"
)

type userBadgeService struct {
	repo repositories.UserBadgeRepository
}

func NewUserBadgeService(repo repositories.UserBadgeRepository) UserBadgeService {
	return &userBadgeService{
		repo: repo,
	}
}

func (s *userBadgeService) AddBadgeToUser(ctx context.Context, payload models.UserBadge) error {
	// Validate input
	if err := base.ValidateModel(&payload); err != nil {
		return err
	}

	// Check if the user already has this badge
	existingBadges, err := s.repo.FindUserBadgesByUser(ctx, payload.UserID.String())
	if err != nil {
		return fmt.Errorf("error checking existing badge: %w", err)
	}

	// If badge already exists, return an error
	for _, badge := range existingBadges {
		if badge.BadgeID == payload.BadgeID {
			return fmt.Errorf("badge %s already assigned to user", payload.BadgeID)
		}
	}

	return s.repo.AddBadgeToUser(ctx, &payload)
}

func (s *userBadgeService) RemoveBadgeFromUser(ctx context.Context, payload models.UserBadgePayload) error {
	// Validate input
	if err := base.ValidateModel(&payload); err != nil {
		return err
	}

	return s.repo.RemoveBadgeFromUser(ctx, &payload)
}

func (s *userBadgeService) GetUserBadgeByID(ctx context.Context, id string) (*models.UserBadgeDTO, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("user badge ID cannot be empty")
	}

	userBadge, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return userBadge, nil
}

func (s *userBadgeService) GetUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadgeDTO, error) {
	// Validate input
	if userID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	userBadges, err := s.repo.FindUserBadgesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return userBadges, nil
}

func (s *userBadgeService) GetUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadgeDTO, error) {
	// Validate input
	if badgeID == "" {
		return nil, fmt.Errorf("badge ID cannot be empty")
	}

	userBadges, err := s.repo.FindUserBadgesByBadge(ctx, badgeID)
	if err != nil {
		return nil, err
	}

	return userBadges, nil
}

func (s *userBadgeService) ListUserBadges(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, int, error) {
	userBadges, total, err := s.repo.Search(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	return userBadges, total, nil
}

func (s *userBadgeService) SearchUserBadges(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, int, error) {
	return s.ListUserBadges(ctx, opts)
}

func (s *userBadgeService) DeleteUserBadge(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return fmt.Errorf("user badge ID cannot be empty")
	}

	return s.repo.Delete(ctx, id)
}

func (s *userBadgeService) CountUserBadges(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.repo.Count(ctx, filters)
}
