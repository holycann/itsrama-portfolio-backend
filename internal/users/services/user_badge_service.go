package services

import (
	"context"
	"fmt"
	"time"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
)

// UserBadgeService handles business logic for user badges
type UserBadgeService struct {
	repo *repositories.UserBadgeRepository
}

// NewUserBadgeService creates a new instance of UserBadgeService
func NewUserBadgeService(repo *repositories.UserBadgeRepository) *UserBadgeService {
	return &UserBadgeService{
		repo: repo,
	}
}

// AssignBadge assigns a new badge to a user
func (s *UserBadgeService) AssignBadge(ctx context.Context, userID, badgeID, badgeName, badgeDescription, badgeIconUrl string) (*models.UserBadge, error) {
	// Check if the user already has this badge
	existingBadges, err := s.repo.FindByUserID(ctx, userID, &models.UserBadgeSearch{
		UserID:  userID,
		BadgeID: badgeID,
		Limit:   1,
		Offset:  0,
	})
	if err != nil {
		return nil, fmt.Errorf("error checking existing badge: %w", err)
	}

	// If badge already exists, return the existing badge
	if len(existingBadges) > 0 {
		return &existingBadges[0], nil
	}

	// Create new user badge
	now := time.Now()
	newBadge := &models.UserBadge{
		UserID:           userID,
		BadgeID:          badgeID,
		BadgeName:        badgeName,
		BadgeDescription: badgeDescription,
		BadgeIconUrl:     badgeIconUrl,
		EarnedAt:         &now,
	}

	return s.repo.Create(ctx, newBadge)
}

// GetUserBadges retrieves badges for a specific user
func (s *UserBadgeService) GetUserBadges(ctx context.Context, userID string, search *models.UserBadgeSearch) ([]models.UserBadge, error) {
	if search == nil {
		search = &models.UserBadgeSearch{
			Limit:  10,
			Offset: 0,
		}
	}

	return s.repo.FindByUserID(ctx, userID, search)
}

// RemoveBadge removes a specific badge from a user
func (s *UserBadgeService) RemoveBadge(ctx context.Context, userID, badgeID string) error {
	return s.repo.Delete(ctx, userID, badgeID)
}

// CountUserBadges counts the number of badges a user has
func (s *UserBadgeService) CountUserBadges(ctx context.Context, userID string) (int, error) {
	return s.repo.CountUserBadges(ctx, userID)
}
