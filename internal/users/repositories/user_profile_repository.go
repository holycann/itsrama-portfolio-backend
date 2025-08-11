package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
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

func (r *userProfileRepository) Create(ctx context.Context, value *models.UserProfile) (*models.UserProfile, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Insert(value, false, "", "minimal", "").
		Execute()

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create user profile")
	}

	return value, nil
}

func (r *userProfileRepository) FindByID(ctx context.Context, id string) (*models.UserProfileDTO, error) {
	var userProfile models.UserProfileDTO

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, user:users_view!users_profile_user_id_fkey(*)", "", false).
		Eq("id", id).
		Single().
		ExecuteTo(&userProfile)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to fetch user profile")
	}

	return &userProfile, nil
}

func (r *userProfileRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.UserProfileDTO, error) {
	var userProfiles []models.UserProfileDTO

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, user:users_view!users_profile_user_id_fkey(*)", "", false).
		Eq(field, fmt.Sprintf("%v", value)).
		ExecuteTo(&userProfiles)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to fetch user profiles")
	}

	if len(userProfiles) == 0 {
		return nil, sql.ErrNoRows
	}

	return userProfiles, nil
}

func (r *userProfileRepository) Update(ctx context.Context, value *models.UserProfile) (*models.UserProfile, error) {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(value, "minimal", "").
		Eq("id", value.ID.String()).
		Execute()

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update user profile")
	}

	return value, nil
}

func (r *userProfileRepository) Delete(ctx context.Context, id string) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Delete("minimal", "").
		Single().
		Eq("id", id).
		Execute()

	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to delete user profile")
	}

	return nil
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

	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to soft delete user profile")
	}

	return nil
}

func (r *userProfileRepository) List(ctx context.Context, opts base.ListOptions) ([]models.UserProfileDTO, error) {
	var userProfiles []models.UserProfileDTO

	query := r.supabaseClient.
		From(r.table).
		Select("*, user:users_view!users_profile_user_id_fkey(*)", "", false)

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

	_, err := query.ExecuteTo(&userProfiles)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to list user profiles")
	}

	return userProfiles, nil
}

func (r *userProfileRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	query := r.supabaseClient.
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
		return 0, errors.Wrap(err, errors.ErrDatabase, "failed to count user profiles")
	}

	return int(count), nil
}

func (r *userProfileRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", true).
		Eq("id", id).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, errors.Wrap(err, errors.ErrDatabase, "failed to check user profile existence")
	}

	return count > 0, nil
}

func (r *userProfileRepository) FindByUserID(ctx context.Context, userID string) (*models.UserProfileDTO, error) {
	var userProfile models.UserProfileDTO

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, user:users_view!users_profile_user_id_fkey(*)", "", false).
		Eq("user_id", userID).
		Single().
		ExecuteTo(&userProfile)

	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to fetch user profile by user ID")
	}

	return &userProfile, nil
}

func (r *userProfileRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", true).
		Eq("user_id", userID).
		Limit(1, "").
		Execute()

	if err != nil {
		return false, errors.Wrap(err, errors.ErrDatabase, "failed to check user profile existence by user ID")
	}

	return count > 0, nil
}

func (r *userProfileRepository) FindByFullname(ctx context.Context, fullname string) ([]models.UserProfileDTO, error) {
	var userProfiles []models.UserProfileDTO

	_, err := r.supabaseClient.
		From(r.table).
		Select("*, user:users_view!users_profile_user_id_fkey(*)", "", false).
		Like("fullname", fmt.Sprintf("%%%s%%", fullname)).
		ExecuteTo(&userProfiles)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find user profiles by fullname")
	}

	return userProfiles, nil
}

func (r *userProfileRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.UserProfileDTO, int, error) {
	var userProfiles []models.UserProfileDTO

	query := r.supabaseClient.
		From(r.table).
		Select("*, user:users_view!users_profile_user_id_fkey(*)", "", false)

	// Apply search if provided
	if opts.Search != "" {
		query = query.Or(
			fmt.Sprintf("fullname.ilike.%%%s%%", opts.Search),
			fmt.Sprintf("bio.ilike.%%%s%%", opts.Search),
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

	// Apply pagination
	offset := (opts.Page - 1) * opts.PerPage
	query = query.Range(offset, offset+opts.PerPage-1, "")

	_, err = query.ExecuteTo(&userProfiles)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrDatabase, "failed to execute search")
	}

	return userProfiles, int(totalCount), nil
}

func (r *userProfileRepository) UpdateAvatarImage(ctx context.Context, payload *models.UserProfileAvatarUpdate) error {
	updateData := map[string]interface{}{
		"avatar_url": payload.AvatarUrl,
	}

	_, _, err := r.supabaseClient.
		From(r.table).
		Update(updateData, "minimal", "").
		Eq("id", payload.ID.String()).
		Execute()

	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to update avatar image")
	}

	return nil
}

func (r *userProfileRepository) UpdatePersonalInfo(ctx context.Context, payload *models.UserProfileUpdate) error {
	updateData := map[string]interface{}{
		"fullname": payload.Fullname,
		"bio":      payload.Bio,
	}

	_, _, err := r.supabaseClient.
		From(r.table).
		Update(updateData, "minimal", "").
		Eq("id", payload.ID.String()).
		Execute()

	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to update personal info")
	}

	return nil
}

func (r *userProfileRepository) VerifyIdentity(ctx context.Context, payload *models.UserProfileIdentityUpdate) error {
	updateData := map[string]interface{}{
		"identity_image_url": payload.IdentityImageUrl,
	}

	_, _, err := r.supabaseClient.
		From(r.table).
		Update(updateData, "minimal", "").
		Eq("id", payload.ID.String()).
		Execute()

	if err != nil {
		return errors.Wrap(err, errors.ErrDatabase, "failed to verify identity")
	}

	return nil
}

func (r *userProfileRepository) BulkCreate(ctx context.Context, values []*models.UserProfile) ([]models.UserProfile, error) {
	var results []models.UserProfile
	for _, profile := range values {
		createdProfile, err := r.Create(ctx, profile)
		if err != nil {
			return nil, err
		}
		results = append(results, *createdProfile)
	}
	return results, nil
}

func (r *userProfileRepository) BulkUpdate(ctx context.Context, values []*models.UserProfile) ([]models.UserProfile, error) {
	var results []models.UserProfile
	for _, profile := range values {
		updatedProfile, err := r.Update(ctx, profile)
		if err != nil {
			return nil, err
		}
		results = append(results, *updatedProfile)
	}
	return results, nil
}

func (r *userProfileRepository) BulkDelete(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *userProfileRepository) BulkUpsert(ctx context.Context, values []*models.UserProfile) ([]models.UserProfile, error) {
	var results []models.UserProfile
	for _, profile := range values {
		updatedProfile, err := r.Update(ctx, profile)
		if err != nil {
			createdProfile, createErr := r.Create(ctx, profile)
			if createErr != nil {
				return nil, createErr
			}
			results = append(results, *createdProfile)
		} else {
			results = append(results, *updatedProfile)
		}
	}
	return results, nil
}
