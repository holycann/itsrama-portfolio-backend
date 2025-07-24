package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/users/models"
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
	})

	return err
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
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
		return err
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
		return err
	}

	fmt.Println("User ID:", UserID)

	err = r.supabaseAuth.AdminDeleteUser(types.AdminDeleteUserRequest{
		UserID: UserID,
	})

	fmt.Println("Error:", err)

	return err
}

func (r *userRepository) SoftDelete(ctx context.Context, id string) error {
	return nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]models.User, error) {
	list, err := r.supabaseAuth.AdminListUsers(types.AdminListUsersRequest{
		Page:    &offset,
		PerPage: &limit,
	})
	if err != nil {
		return nil, err
	}

	var mappedUsers []models.User
	for _, user := range list.Users {
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

	return mappedUsers, nil
}

func (r *userRepository) Count(ctx context.Context) (int, error) {
	list, err := r.supabaseAuth.AdminListUsers(types.AdminListUsersRequest{})
	if err != nil {
		return 0, err
	}

	return len(list.Users), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	list, err := r.supabaseAuth.AdminListUsers(types.AdminListUsersRequest{})
	if err != nil {
		return nil, err
	}

	if len(list.Users) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, user := range list.Users {
		if user.Email == email {
			users := []models.User{{
				ID:           user.ID.String(),
				Email:        user.Email,
				Phone:        user.Phone,
				Role:         user.Role,
				LastSignInAt: user.LastSignInAt,
				CreatedAt:    &user.CreatedAt,
				UpdatedAt:    &user.UpdatedAt,
			}}
			return &users[0], nil
		}
	}

	return nil, sql.ErrNoRows
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	list, err := r.supabaseAuth.AdminListUsers(types.AdminListUsersRequest{})

	if err != nil {
		return false, err
	}

	for _, user := range list.Users {
		if user.Email == email {
			return true, nil
		}
	}

	return false, nil
}
