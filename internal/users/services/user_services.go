package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
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

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	// Create user
	return s.repo.Create(ctx, user)
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

func (s *userService) ListUsers(ctx context.Context, opts repository.ListOptions) ([]models.User, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.repo.List(ctx, opts)
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

func (s *userService) CountUsers(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.repo.Count(ctx, filters)
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

func (s *userService) SearchUsers(ctx context.Context, query string, opts repository.ListOptions) ([]models.User, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	// Add search query to filters
	opts.Filters = append(opts.Filters,
		repository.FilterOption{
			Field:    "email",
			Operator: "like",
			Value:    query,
		},
		repository.FilterOption{
			Field:    "phone",
			Operator: "like",
			Value:    query,
		},
	)

	return s.repo.List(ctx, opts)
}
