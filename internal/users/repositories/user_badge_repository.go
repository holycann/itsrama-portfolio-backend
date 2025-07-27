package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/repository"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type userBadgeRepository struct {
	client *supabase.Client
	table  string
}

func NewUserBadgeRepository(client *supabase.Client) UserBadgeRepository {
	return &userBadgeRepository{
		client: client,
		table:  "user_badges",
	}
}

func (r *userBadgeRepository) Create(ctx context.Context, badge *models.UserBadge) error {
	_, err := r.client.
		From(r.table).
		Insert(badge, false, "", "minimal", "").
		ExecuteTo(&badge)
	return err
}

func (r *userBadgeRepository) FindByID(ctx context.Context, id string) (*models.UserBadge, error) {
	var badges []models.UserBadge
	_, err := r.client.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		ExecuteTo(&badges)

	if err != nil {
		return nil, err
	}

	if len(badges) > 0 {
		return &badges[0], nil
	}
	return nil, fmt.Errorf("user badge not found")
}

func (r *userBadgeRepository) Update(ctx context.Context, badge *models.UserBadge) error {
	_, _, err := r.client.
		From(r.table).
		Update(badge, "minimal", "").
		Eq("id", badge.ID.String()).
		Execute()
	return err
}

func (r *userBadgeRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.client.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *userBadgeRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.UserBadge, error) {
	var badges []models.UserBadge
	query := r.client.
		From(r.table).
		Select("*", "", false)

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case "=":
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case "like":
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder == repository.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	query = query.Range(opts.Offset, opts.Offset+opts.Limit-1, "")

	_, err := query.ExecuteTo(&badges)
	return badges, err
}

func (r *userBadgeRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
	query := r.client.
		From(r.table).
		Select("id", "exact", false)

	// Apply filters
	for _, filter := range filters {
		switch filter.Operator {
		case "=":
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case "like":
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	_, count, err := query.Execute()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *userBadgeRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *userBadgeRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.UserBadge, error) {
	var badges []models.UserBadge
	_, err := r.client.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&badges)
	return badges, err
}

// Specialized methods for user badges
func (r *userBadgeRepository) FindUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadge, error) {
	var badges []models.UserBadge
	_, err := r.client.
		From(r.table).
		Select("*", "", false).
		Eq("user_id", userID).
		ExecuteTo(&badges)
	return badges, err
}

func (r *userBadgeRepository) FindUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadge, error) {
	var badges []models.UserBadge
	_, err := r.client.
		From(r.table).
		Select("*", "", false).
		Eq("badge_id", badgeID).
		ExecuteTo(&badges)
	return badges, err
}
