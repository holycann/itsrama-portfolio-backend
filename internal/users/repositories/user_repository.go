package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
	"github.com/supabase-community/auth-go"
	"github.com/supabase-community/auth-go/types"
)

type userRepository struct {
	supabaseAuth auth.Client
}

func NewUserRepository(client auth.Client) UserRepository {
	return &userRepository{
		supabaseAuth: client,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	// Create user
	createdUser, err := r.supabaseAuth.AdminCreateUser(types.AdminCreateUserRequest{
		Email:    user.Email,
		Password: &user.Password,
		Phone:    user.Phone,
		Role:     user.Role,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to create user")
	}

	return &models.User{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		Phone:     createdUser.Phone,
		Role:      createdUser.Role,
		CreatedAt: &createdUser.CreatedAt,
		UpdatedAt: &createdUser.UpdatedAt,
	}, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New(errors.ErrValidation, "invalid user ID", err)
	}

	user, err := r.supabaseAuth.AdminGetUser(types.AdminGetUserRequest{
		UserID: userID,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to fetch user")
	}

	return &models.User{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	updatedUser, err := r.supabaseAuth.AdminUpdateUser(types.AdminUpdateUserRequest{
		UserID:   user.ID,
		Email:    user.Email,
		Password: user.Password,
		Phone:    user.Phone,
		Role:     user.Role,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to update user")
	}

	return &models.User{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Phone:     updatedUser.Phone,
		Role:      updatedUser.Role,
		CreatedAt: &updatedUser.CreatedAt,
		UpdatedAt: &updatedUser.UpdatedAt,
	}, nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	UserID, err := uuid.Parse(id)
	if err != nil {
		return errors.New(errors.ErrValidation, "invalid user ID", err)
	}

	return r.supabaseAuth.AdminDeleteUser(types.AdminDeleteUserRequest{
		UserID: UserID,
	})
}

func (r *userRepository) List(ctx context.Context, opts base.ListOptions) ([]models.User, error) {
	listOpts := types.AdminListUsersRequest{
		Page:    &opts.Page,
		PerPage: &opts.PerPage,
	}

	list, err := r.supabaseAuth.AdminListUsers(listOpts)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to list users")
	}

	var mappedUsers []models.User
	for _, user := range list.Users {
		// Apply filtering
		if !r.matchesFilters(user, opts.Filters) {
			continue
		}

		mappedUsers = append(mappedUsers, models.User{
			ID:        user.ID,
			Email:     user.Email,
			Phone:     user.Phone,
			Role:      user.Role,
			CreatedAt: &user.CreatedAt,
			UpdatedAt: &user.UpdatedAt,
		})
	}

	return mappedUsers, nil
}

func (r *userRepository) Count(ctx context.Context, filters []base.FilterOption) (int, error) {
	list, err := r.supabaseAuth.AdminListUsers(types.AdminListUsersRequest{})
	if err != nil {
		return 0, errors.Wrap(err, errors.ErrDatabase, "failed to count users")
	}

	count := 0
	for _, user := range list.Users {
		if r.matchesFilters(user, filters) {
			count++
		}
	}

	return count, nil
}

func (r *userRepository) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *userRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.User, error) {
	list, err := r.supabaseAuth.AdminListUsers(types.AdminListUsersRequest{})
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "failed to find users by field")
	}

	var matchedUsers []models.User
	for _, user := range list.Users {
		var fieldValue interface{}
		switch field {
		case "id":
			fieldValue = user.ID
		case "email":
			fieldValue = user.Email
		case "phone":
			fieldValue = user.Phone
		case "role":
			fieldValue = user.Role
		default:
			continue
		}

		if fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", value) {
			matchedUsers = append(matchedUsers, models.User{
				ID:        user.ID,
				Email:     user.Email,
				Phone:     user.Phone,
				Role:      user.Role,
				CreatedAt: &user.CreatedAt,
				UpdatedAt: &user.UpdatedAt,
			})
		}
	}

	return matchedUsers, nil
}

func (r *userRepository) Search(ctx context.Context, opts base.ListOptions) ([]models.User, int, error) {
	users, err := r.List(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := r.Count(ctx, opts.Filters)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	users, err := r.FindByField(ctx, "email", email)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}
	return &users[0], nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	users, err := r.FindByField(ctx, "email", email)
	if err != nil {
		return false, err
	}
	return len(users) > 0, nil
}

func (r *userRepository) ChangeUserRole(ctx context.Context, payload *models.UserRoleUpdate) error {
	_, err := r.supabaseAuth.AdminUpdateUser(types.AdminUpdateUserRequest{
		UserID: payload.ID,
		Role:   payload.NewRole,
	})

	return err
}

func (r *userRepository) BulkCreate(ctx context.Context, values []*models.User) ([]models.User, error) {
	var results []models.User
	for _, user := range values {
		createdUser, err := r.Create(ctx, user)
		if err != nil {
			return nil, err
		}
		results = append(results, *createdUser)
	}
	return results, nil
}

func (r *userRepository) BulkUpdate(ctx context.Context, values []*models.User) ([]models.User, error) {
	var results []models.User
	for _, user := range values {
		updatedUser, err := r.Update(ctx, user)
		if err != nil {
			return nil, err
		}
		results = append(results, *updatedUser)
	}
	return results, nil
}

func (r *userRepository) BulkDelete(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *userRepository) SoftDelete(ctx context.Context, id string) error {
	parsedUserID, err := uuid.Parse(id)
	if err != nil {
		return errors.New(errors.ErrValidation, "invalid user ID", err)
	}

	_, err = r.supabaseAuth.AdminUpdateUser(types.AdminUpdateUserRequest{
		UserID: parsedUserID,
	})

	return err
}

func (r *userRepository) BulkUpsert(ctx context.Context, values []*models.User) ([]models.User, error) {
	var results []models.User
	for _, user := range values {
		updatedUser, err := r.Update(ctx, user)
		if err != nil {
			createdUser, createErr := r.Create(ctx, user)
			if createErr != nil {
				return nil, createErr
			}
			results = append(results, *createdUser)
		} else {
			results = append(results, *updatedUser)
		}
	}
	return results, nil
}

// Helper methods for filtering
func (r *userRepository) matchesFilters(user types.User, filters []base.FilterOption) bool {
	if len(filters) == 0 {
		return true
	}

	for _, filter := range filters {
		var value interface{}
		switch filter.Field {
		case "id":
			value = user.ID
		case "email":
			value = user.Email
		case "phone":
			value = user.Phone
		case "role":
			value = user.Role
		default:
			continue
		}

		switch filter.Operator {
		case base.OperatorEqual:
			if fmt.Sprintf("%v", value) != fmt.Sprintf("%v", filter.Value) {
				return false
			}
		case base.OperatorLike:
			if !containsIgnoreCase(fmt.Sprintf("%v", value), fmt.Sprintf("%v", filter.Value)) {
				return false
			}
		}
	}
	return true
}

// Helper function for case-insensitive contains
func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(
		strings.ToLower(str),
		strings.ToLower(substr),
	)
}
