package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/supabase-community/supabase-go"
)

type userProfileRepository struct {
	supabaseClient *supabase.Client
	table          string
}

func NewUserProfileRepository(client *supabase.Client, table string) UserProfileRepository {

	return &userProfileRepository{
		supabaseClient: client,
		table:          table,
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

func (r *userProfileRepository) Update(ctx context.Context, value *models.UserProfile) error {
	_, _, err := r.supabaseClient.
		From(r.table).
		Update(value, "minimal", "").
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

func (r *userProfileRepository) List(ctx context.Context, limit, offset int) ([]models.UserProfile, error) {
	var userProfiles []models.UserProfile

	_, err := r.supabaseClient.
		From(r.table).
		Select("*", "", false).
		Range(offset, offset+limit-1, "").
		Filter("deleted_at", "is", "null").
		ExecuteTo(&userProfiles)

	if err != nil {
		return nil, err
	}

	return userProfiles, nil
}

func (r *userProfileRepository) Count(ctx context.Context) (int, error) {
	_, count, err := r.supabaseClient.
		From(r.table).
		Select("id", "exact", true).
		Execute()
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
