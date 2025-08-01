package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/repository"
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

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.supabaseAuth.AdminCreateUser(types.AdminCreateUserRequest{
		Email:    user.Email,
		Password: &user.Password,
		Phone:    user.Phone,
		Role:     user.Role,
	})

	return err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := r.supabaseAuth.AdminGetUser(types.AdminGetUserRequest{
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by id: %w", err)
	}

	return &models.User{
		ID:           user.ID.String(),
		Email:        user.Email,
		Phone:        user.Phone,
		Role:         user.Role,
		LastSignInAt: user.LastSignInAt,
		CreatedAt:    &user.CreatedAt,
		UpdatedAt:    &user.UpdatedAt,
	}, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	UserID, err := uuid.Parse(user.ID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	_, err = r.supabaseAuth.AdminUpdateUser(types.AdminUpdateUserRequest{
		UserID:   UserID,
		Email:    user.Email,
		Password: user.Password,
		Phone:    user.Phone,
		Role:     user.Role,
	})

	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	UserID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	return r.supabaseAuth.AdminDeleteUser(types.AdminDeleteUserRequest{
		UserID: UserID,
	})
}

func (r *userRepository) List(ctx context.Context, opts repository.ListOptions) ([]models.User, error) {
	// Default pagination
	page := opts.Offset / opts.Limit
	if page < 1 {
		page = 1
	}

	listOpts := types.AdminListUsersRequest{
		Page:    &page,
		PerPage: &opts.Limit,
	}

	list, err := r.supabaseAuth.AdminListUsers(listOpts)
	if err != nil {
		return nil, err
	}

	var mappedUsers []models.User
	for _, user := range list.Users {
		// Apply filtering
		if !r.matchesFilters(user, opts.Filters) {
			continue
		}

		// Adjust time handling to use pointers
		mappedUsers = append(mappedUsers, models.User{
			ID:           user.ID.String(),
			Email:        user.Email,
			Phone:        user.Phone,
			Role:         user.Role,
			LastSignInAt: user.LastSignInAt,
			CreatedAt:    &user.CreatedAt,
			UpdatedAt:    &user.UpdatedAt,
		})
	}

	// Apply sorting if needed
	r.sortUsers(&mappedUsers, opts.SortBy, opts.SortOrder)

	return mappedUsers, nil
}

func (r *userRepository) Count(ctx context.Context, filters []repository.FilterOption) (int, error) {
	list, err := r.supabaseAuth.AdminListUsers(types.AdminListUsersRequest{})
	if err != nil {
		return 0, err
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
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *userRepository) FindByField(ctx context.Context, field string, value interface{}) ([]models.User, error) {
	list, err := r.supabaseAuth.AdminListUsers(types.AdminListUsersRequest{})
	if err != nil {
		return nil, err
	}

	var matchedUsers []models.User
	for _, user := range list.Users {
		var fieldValue interface{}
		switch strings.ToLower(field) {
		case "id":
			fieldValue = user.ID.String()
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
				ID:           user.ID.String(),
				Email:        user.Email,
				Phone:        user.Phone,
				Role:         user.Role,
				LastSignInAt: user.LastSignInAt,
				CreatedAt:    &user.CreatedAt,
				UpdatedAt:    &user.UpdatedAt,
			})
		}
	}

	return matchedUsers, nil
}

func (r *userRepository) Search(ctx context.Context, opts repository.ListOptions) ([]models.User, int, error) {
	// Perform list operation first
	users, err := r.List(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	// Count total matching users
	totalCount, err := r.Count(ctx, opts.Filters)
	if err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// Helper methods for filtering and sorting
func (r *userRepository) matchesFilters(user types.User, filters []repository.FilterOption) bool {
	if len(filters) == 0 {
		return true
	}

	for _, filter := range filters {
		var value interface{}
		switch strings.ToLower(filter.Field) {
		case "id":
			value = user.ID.String()
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
		case "=":
			if fmt.Sprintf("%v", value) != fmt.Sprintf("%v", filter.Value) {
				return false
			}
		case "like":
			if !strings.Contains(
				strings.ToLower(fmt.Sprintf("%v", value)),
				strings.ToLower(fmt.Sprintf("%v", filter.Value)),
			) {
				return false
			}
		}
	}
	return true
}

func (r *userRepository) sortUsers(users *[]models.User, sortBy string, sortOrder repository.SortOrder) {
	if sortBy == "" {
		return
	}

	sort.Slice(*users, func(i, j int) bool {
		var less bool
		switch strings.ToLower(sortBy) {
		case "id":
			less = (*users)[i].ID < (*users)[j].ID
		case "email":
			less = (*users)[i].Email < (*users)[j].Email
		case "phone":
			less = (*users)[i].Phone < (*users)[j].Phone
		case "role":
			less = (*users)[i].Role < (*users)[j].Role
		case "last_sign_in_at":
			if (*users)[i].LastSignInAt == nil || (*users)[j].LastSignInAt == nil {
				less = false
			} else {
				less = (*users)[i].LastSignInAt.Before(*(*users)[j].LastSignInAt)
			}
		case "created_at":
			if (*users)[i].CreatedAt == nil || (*users)[j].CreatedAt == nil {
				less = false
			} else {
				less = (*users)[i].CreatedAt.Before(*(*users)[j].CreatedAt)
			}
		default:
			return false
		}

		if sortOrder == repository.SortDescending {
			less = !less
		}
		return less
	})
}

// Additional methods specific to users
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
