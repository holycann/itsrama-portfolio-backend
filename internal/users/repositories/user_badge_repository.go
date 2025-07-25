package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/supabase-community/supabase-go"
)

// UserBadgeRepository handles database operations for user badges
type UserBadgeRepository struct {
	client *supabase.Client
	table  string
}

// UserBadgeConfig provides configuration for the user badge repository
type UserBadgeConfig struct {
	Table string
}

// DefaultUserBadgeConfig returns the default configuration for user badge repository
func DefaultUserBadgeConfig() *UserBadgeConfig {
	return &UserBadgeConfig{
		Table: "user_badges", // Adjust table name as needed
	}
}

// NewUserBadgeRepository creates a new instance of UserBadgeRepository
func NewUserBadgeRepository(client *supabase.Client, config UserBadgeConfig) *UserBadgeRepository {
	return &UserBadgeRepository{
		client: client,
		table:  config.Table,
	}
}

// Create adds a new badge to a user
func (r *UserBadgeRepository) Create(ctx context.Context, badge *models.UserBadge) (*models.UserBadge, error) {
	var result []models.UserBadge

	_, err := r.client.
		From(r.table).
		Insert(badge, false, "", "minimal", "").
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("failed to create user badge: %w", err)
	}

	if len(result) > 0 {
		return &result[0], nil
	}
	return nil, fmt.Errorf("no badge created")
}

// FindByUserID retrieves badges for a specific user
func (r *UserBadgeRepository) FindByUserID(ctx context.Context, userID string, search *models.UserBadgeSearch) ([]models.UserBadge, error) {
	query := r.client.
		From(r.table).
		Select("*", "", false).
		Eq("user_id", userID)

	if search != nil {
		if search.BadgeID != "" {
			query = query.Eq("badge_id", search.BadgeID)
		}
		query = query.Range(search.Offset, search.Offset+search.Limit-1, "")
	}

	var badges []models.UserBadge
	_, err := query.ExecuteTo(&badges)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user badges: %w", err)
	}

	return badges, nil
}

// Delete removes a specific badge from a user
func (r *UserBadgeRepository) Delete(ctx context.Context, userID, badgeID string) error {
	_, _, err := r.client.
		From(r.table).
		Delete("minimal", "").
		Eq("user_id", userID).
		Eq("badge_id", badgeID).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to delete user badge: %w", err)
	}
	return nil
}

// CountUserBadges counts the number of badges a user has
func (r *UserBadgeRepository) CountUserBadges(ctx context.Context, userID string) (int, error) {
	_, count, err := r.client.
		From(r.table).
		Select("id", "exact", true).
		Eq("user_id", userID).
		Execute()

	if err != nil {
		return 0, fmt.Errorf("failed to count user badges: %w", err)
	}

	return int(count), nil
}
