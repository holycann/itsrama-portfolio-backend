package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
)

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	// Validate input
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	fmt.Println(user.Email)

	// Check if user with email already exists
	exists, err := s.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}
	if exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Set default values
	now := time.Now().UTC()
	user.CreatedAt = &now
	user.UpdatedAt = &now

	// Create user
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve users
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	// Convert to pointer slice
	userPtrs := make([]*models.User, len(users))
	for i := range users {
		userPtrs[i] = &users[i]
	}

	return userPtrs, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	// Validate input
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	// Retrieve user
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", id)
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Validate input
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	// Retrieve user
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	// Validate input
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}
	if user.ID == "" {
		return fmt.Errorf("user ID is required for update")
	}

	// Check if user exists
	existingUser, err := s.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}

	// Update timestamps
	now := time.Now().UTC()
	user.UpdatedAt = &now

	// Preserve certain fields
	user.CreatedAt = existingUser.CreatedAt

	// Perform update
	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	// Validate input
	if id == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Check if user exists
	_, err := s.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	// Soft delete the user
	return s.repo.Delete(ctx, id)
}
