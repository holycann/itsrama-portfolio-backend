package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type userBadgeRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewUserBadgeRepository(client *supabase.Client) UserBadgeRepository {
	return &userBadgeRepository{
		supabaseClient: client,
		table:          "users_badge",
	}
}

func (r *userBadgeRepository) Create(ctx context.Context, value *models.UserBadge) (*models.UserBadge, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(value, false, "", "minimal", "").
		Execute()

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create user badge")
	}

	return value, nil
}

func (r *userBadgeRepository) FindByID(ctx context.Context, id string) (*models.UserBadgeDTO, error) {
	var userBadges models.UserBadgeDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, badge:badges(*)", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&userBadges)

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to fetch user badge")
	}

	return &userBadges, nil
}

func (r *userBadgeRepository) Update(ctx context.Context, value *models.UserBadge) (*models.UserBadge, error) {

	_, _, err := r.supabaseClient.
		From(r.table).
		Update(value, "minimal", "").
		Eq("user_id", value.UserID.String()).
		Eq("badge_id", value.BadgeID.String()).
		Execute()

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update user badge")
	}

	return value, nil
}

func (r *userBadgeRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Single().
		Eq("id", id).
		Execute()

	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to delete user badge")
	}

	return nil
}

func (r *userBadgeRepository) List(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, error) {
	var userBadges []models.UserBadgeDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*, badge:badges(*)", "", false)

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
		ascending := opts.SortOrder == base.SortAscending
		query = query.Order(opts.SortBy, &postgrest.OrderOpts{Ascending: ascending})
	}

	// Apply pagination
	query = query.Range((opts.Page-1)*opts.PerPage, opts.Page*opts.PerPage-1, "")

	_, err := query.ExecuteTo(&userBadges)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to list user badges")
	}

	return userBadges, nil
}

func (r *userBadgeRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	query := r.supabaseClient.
		From(r.table).
		Select("id", "exact", false)

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
		return 0, errors.Wrap(err, errors.ErrDatabase, "failed to count user badges")
	}

	return int(count), nil
}

func (r *userBadgeRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *userBadgeRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.UserBadgeDTO, error) {
	var userBadges []models.UserBadgeDTO
	_, err := r.supabaseClient.
		From(r.table).
		Select("*, badge:badges(*)", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&userBadges)

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find user badges by field")
	}

	if len(userBadges) == 0 {
		return nil, sql.ErrNoRows
	}

	return userBadges, nil
}

func (r *userBadgeRepository) BulkCreate(ctx context.Context, values []*models.UserBadge) ([]models.UserBadge, error) {
	var results []models.UserBadge
	for _, badge := range values {
		createdBadge, err := r.Create(ctx, badge)
		if err != nil {
			return nil, err
		}
		results = append(results, *createdBadge)
	}
	return results, nil
}

func (r *userBadgeRepository) BulkUpdate(ctx context.Context, values []*models.UserBadge) ([]models.UserBadge, error) {
	var results []models.UserBadge
	for _, badge := range values {
		updatedBadge, err := r.Update(ctx, badge)
		if err != nil {
			return nil, err
		}
		results = append(results, *updatedBadge)
	}
	return results, nil
}

func (r *userBadgeRepository) BulkDelete(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *userBadgeRepository) FindUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadgeDTO, error) {
	return r.FindByField(ctx, "user_id", userID)
}

func (r *userBadgeRepository) FindUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadgeDTO, error) {
	return r.FindByField(ctx, "badge_id", badgeID)
}

func (r *userBadgeRepository) CountUserBadges(ctx context.Context, userID string) (int, error) {
	return r.Count(ctx, []base.FilterOption{
		{
			Field:    "user_id",
			Operator: base.OperatorEqual,
			Value:    userID,
		},
	})
}

func (r *userBadgeRepository) RemoveBadgeFromUser(ctx context.Context, payload *models.UserBadgePayload) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Eq("user_id", payload.UserID.String()).
		Eq("badge_id", payload.BadgeID.String()).
		Execute()

	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to remove badge from user")
	}

	return nil
}

func (r *userBadgeRepository) AddBadgeToUser(ctx context.Context, payload *models.UserBadge) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(payload, false, "", "minimal", "").
		Execute()

	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to add badge to user")
	}

	return nil
}

func (r *userBadgeRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.UserBadgeDTO, int, error) {
	var userBadges []models.UserBadgeDTO
	query := r.supabaseClient.
		From(r.table).
		Select("*, badge:badges(*)", "", false)

	// Apply filters
	for _, filter := range opts.Filters {
		switch filter.Operator {
		case base.OperatorEqual:
			query = query.Eq(filter.Field, fmt.Sprintf("%v", filter.Value))
		case base.OperatorLike:
			query = query.Like(filter.Field, fmt.Sprintf("%%%v%%", filter.Value))
		}
	}

	// Count total matching records
	_, totalCount, err := query.Execute()
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to count user badges")
	}

	// Apply pagination
	offset := (opts.Page - 1) * opts.PerPage
	query = query.Range(offset, offset+opts.PerPage-1, "")

	// Execute query to get results
	_, err = query.ExecuteTo(&userBadges)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to search user badges")
	}

	return userBadges, int(totalCount), nil
}
