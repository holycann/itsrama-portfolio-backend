package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/repository"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type userProfileRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewUserProfileRepository(client *supabase.Client) UserProfileRepository {
	return &userProfileRepository{
		supabaseClient: client,
		table:          "users_profile",
	}
}

func (r *userProfileRepository) Create(ctx context.Context, value *models.UserProfile) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(value, false, "", "minimal", "").
		Execute()

	return err
}

func (r *userProfileRepository) FindByID(ctx context.Context, id string) (*models.UserProfile, error) {
	var userProfile []models.UserProfile

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("id", id).
		Limit(1, "").
		ExecuteTo(&userProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile by id: %w", err)
	}

	if len(userProfile) == 0 {
		return nil, sql.ErrNoRows
	}

	return &userProfile[0], nil
}

func (r *userProfileRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.UserProfile, error) {
	var userProfiles []models.UserProfile

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&userProfiles)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile by field %s: %w", field, err)
	}

	if len(userProfiles) == 0 {
		return nil, sql.ErrNoRows
	}

	return userProfiles, nil
}

func (r *userProfileRepository) Update(ctx context.Context, value *models.UserProfile) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(value, "minimal", "").
		Eq("id", value.ID.String()).
		Execute()

	return err
}

func (r *userProfileRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Single().
		Eq("id", id).
		Execute()

	return err
}

func (r *userProfileRepository) SoftDelete(ctx context.Context, id string) error {
	updateData := map[string]interface{}{
		"deleted_at": time.Now().UTC(),
	}

	_, _, err := r.supabaseClient.
		From(r.table).
		Update(updateData, "minimal", "").
		Eq("id", id).
		Execute()

	return err
}

func (r *userProfileRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.UserProfile, error) {
	var userProfiles []models.UserProfile

	query := r.supabaseClient.
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

	_, err := query.ExecuteTo(&userProfiles)
	if err != nil {
		return nil, err
	}

	return userProfiles, nil
}

func (r *userProfileRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
	query := r.supabaseClient.
		From(r.table).
		Select("id", "exact", true)

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

func (r *userProfileRepository) FindByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
	var users []models.UserProfile

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Eq("user_id", userID).
		Single().
		ExecuteTo(&users)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return &users[0], nil
}

func (r *userProfileRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", true).
		Eq("id", id).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userProfileRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", true).
		Eq("user_id", userID).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userProfileRepository) Search(ctx context.Context, option repository.ListOptions) ([]models.UserProfile, int, error) {
	var userProfiles []models.UserProfile

	query := option.SearchQuery

	_, count, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Or(
			fmt.Sprintf("username.ilike.%%%s%%", query),
			fmt.Sprintf("email.ilike.%%%s%%", query),
		).
		Execute()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search user profiles: %w", err)
	}

	_, err = r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Or(
			fmt.Sprintf("username.ilike.%%%s%%", query),
			fmt.Sprintf("email.ilike.%%%s%%", query),
		).
		ExecuteTo(&userProfiles)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute search user profiles: %w", err)
	}

	return userProfiles, int(count), nil
}
