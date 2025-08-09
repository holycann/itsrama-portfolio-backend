package repositories

import (
	"context"
	"fmt"

	"github.com/holycann/cultour-backend/internal/achievement/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
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

func (r *badgeRepository) Create(ctx context.Context, badge *models.Badge) (*models.Badge, error) {
	_, _, err := r.client.
		From(r.table).
		Insert(badge, false, "", "minimal", "").
		Execute()

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create badge")
	}

	return badge, nil
}

func (r *badgeRepository) FindByID(ctx context.Context, id string) (*models.Badge, error) {
	var badge models.Badge

	_, err := r.client.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&badge)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to fetch badge")
	}

	return &badge, nil
}

func (r *badgeRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.Badge, error) {
	var badges []models.Badge

	_, err := r.client.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&badges)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to fetch badges")
	}

	if len(badges) == 0 {
		return nil, fmt.Errorf("no badges found")
	}

	return badges, nil
}

func (r *badgeRepository) Update(ctx context.Context, badge *models.Badge) (*models.Badge, error) {
	_, _, err := r.client.
		From(r.table).
		Update(badge, "minimal", "").
		Eq("id", badge.ID.String()).
		Execute()

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update badge")
	}

	return badge, nil
}

func (r *badgeRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.client.
		From(r.table).
		Delete("minimal", "").
		Single().
		Eq("id", id).
		Execute()

	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to delete badge")
	}

	return nil
}

func (r *badgeRepository) List(ctx context.Context, opts base.ListOptions) ([]models.Badge, error) {
	var badges []models.Badge

	query := r.client.
		From(r.table).
		Select("*", "", false)

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder != base.SortDescending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	offset := (opts.Page - 1) * opts.PerPage
	query = query.Range(offset, offset+opts.PerPage-1, "")

	_, err := query.ExecuteTo(&badges)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to list badges")
	}

	return badges, nil
}

func (r *badgeRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	query := r.client.
		From(r.table).
		Select("id", "exact", true)

	// Apply filters
	for _, filter := range filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	_, count, err := query.Execute()
	if err != nil {
		return 0, errors.Wrap(err, errors.ErrDatabase, "failed to count badges")
	}

	return int(count), nil
}

func (r *badgeRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, count, err := r.client.
		From(r.table).
		Select("id", "exact", true).
		Eq("id", id).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, errors.Wrap(err, errors.ErrDatabase, "failed to check badge existence")
	}

	return count > 0, nil
}

// Specialized methods for badges
func (r *badgeRepository) FindBadgeByName(ctx context.Context, name string) (*models.Badge, error) {
	badges, err := r.FindByField(ctx, "name", name)
	if err != nil {
		return nil, err
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

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find popular badges")
	}

	return badges, nil
}

func (r *badgeRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.Badge, int, error) {
	var badges []models.Badge

	query := r.client.
		From(r.table).
		Select("*", "", false)

	// Apply search query if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("name.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("description.ilike.%%%s%%", opts.Search),
		)
	}

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Count total results
	_, totalCount, err := query.Execute()
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to count search results")
	}

	// Apply sorting
	if opts.SortBy != "" {
		ascending := opts.SortOrder != base.SortDescending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	offset := (opts.Page - 1) * opts.PerPage
	query = query.Range(offset, offset+opts.PerPage-1, "")

	_, err = query.ExecuteTo(&badges)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to execute search")
	}

	return badges, int(totalCount), nil
}

func (r *badgeRepository) BulkCreate(ctx context.Context, values []*models.Badge) ([]models.Badge, error) {
	var results []models.Badge
	for _, badge := range values {
		createdBadge, err := r.Create(ctx, badge)
		if err != nil {
			return nil, err
		}
		results = append(results, *createdBadge)
	}
	return results, nil
}

func (r *badgeRepository) BulkUpdate(ctx context.Context, values []*models.Badge) ([]models.Badge, error) {
	var results []models.Badge
	for _, badge := range values {
		updatedBadge, err := r.Update(ctx, badge)
		if err != nil {
			return nil, err
		}
		results = append(results, *updatedBadge)
	}
	return results, nil
}

func (r *badgeRepository) BulkDelete(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *badgeRepository) BulkUpsert(ctx context.Context, values []*models.Badge) ([]models.Badge, error) {
	var results []models.Badge
	for _, badge := range values {
		updatedBadge, err := r.Update(ctx, badge)
		if err != nil {
			createdBadge, createErr := r.Create(ctx, badge)
			if createErr != nil {
				return nil, createErr
			}
			results = append(results, *createdBadge)
		} else {
			results = append(results, *updatedBadge)
		}
	}
	return results, nil
}
