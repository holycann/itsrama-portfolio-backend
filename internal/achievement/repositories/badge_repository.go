package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/pkg/repository"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type badgeRepository struct {
	client *supabase.Client
	table  string
}

func NewBadgeRepository(client *supabase.Client) BadgeRepository {
	return &badgeRepository{
		client: client,
		table:  "badges",
	}
}

func (r *badgeRepository) Create(ctx context.Context, badge *models.Badge) error {
	_, err := r.client.
		From(r.table).
		Insert(badge, false, "", "minimal", "").
		ExecuteTo(&badge)
	return err
}

func (r *badgeRepository) FindByID(ctx context.Context, id string) (*models.Badge, error) {
	var badges []models.Badge
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
	return nil, fmt.Errorf("badge not found")
}

func (r *badgeRepository) Update(ctx context.Context, badge *models.Badge) error {
	_, _, err := r.client.
		From(r.table).
		Update(badge, "minimal", "").
		Eq("id", badge.ID.String()).
		Execute()
	return err
}

func (r *badgeRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.client.
		From(r.table).
		Delete("minimal", "").
		Eq("id", id).
		Execute()
	return err
}

func (r *badgeRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.Badge, error) {
	var badges []models.Badge
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

func (r *badgeRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
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

func (r *badgeRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *badgeRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.Badge, error) {
	var badges []models.Badge
	_, err := r.client.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&badges)
	return badges, err
}

// Specialized methods for badges
func (r *badgeRepository) FindBadgeByName(ctx context.Context, name string) (*models.Badge, error) {
	badges, err := r.FindByField(ctx, "name", name)
	if err != nil {
		return nil, err
	}
	if len(badges) == 0 {
		return nil, fmt.Errorf("badge not found")
	}
	return &badges[0], nil
}

func (r *badgeRepository) FindPopularBadges(ctx context.Context, limit int) ([]models.Badge, error) {
	var badges []models.Badge
	_, err := r.client.
		From(r.table).
		Select("*", "", false).
		Order("created_at", &postgrest.OrderOpts{Ascending: false}).
		Limit(limit, "").
		ExecuteTo(&badges)
	return badges, err
}
