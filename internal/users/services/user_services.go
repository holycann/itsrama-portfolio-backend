package services

import (
	"context"
	"database/sql"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/internal/users/repositories"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
)

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.UserCreate) (*models.UserDTO, error) {
	// Validate input
	if err := base.ValidateModel(user); err != nil {
		return nil, err
	}

	// Check if user with email already exists
	exists, err := s.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrDatabase, "error checking user existence")
	}
	if exists {
		return nil, errors.New(errors.ErrDatabase, "user with email already exists", err)
	}

	// Convert UserCreate to User
	userModel := &models.User{
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}

	// Create user
	createdUser, err := s.repo.Create(ctx, userModel)
	if err != nil {
		return nil, err
	}

	dto := createdUser.ToDTO()

	// Convert to DTO
	return &dto, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.UserDTO, error) {
	// Retrieve user
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.ErrNotFound, "user not found", err)
		}
		return nil, errors.Wrap(err, errors.ErrDatabase, "error retrieving user")
	}

	// Convert to DTO
	dto := user.ToDTO()
	return &dto, nil
}

func (s *userService) ListUsers(ctx context.Context, opts base.ListOptions) ([]models.UserDTO, int, error) {
	// Validate ListOptions
	if err := opts.Validate(); err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrValidation, "invalid list options")
	}

	users, count, err := s.repo.Search(ctx, opts)
	if err != nil {
		return nil, 0, err
	}

	// Convert to DTOs
	var userDTOs []models.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, user.ToDTO())
	}

	return userDTOs, count, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *models.UserUpdate) (*models.UserDTO, error) {
	// Validate input
	if err := base.ValidateModel(user); err != nil {
		return nil, err
	}

	// Convert UserUpdate to User
	userModel := &models.User{
		ID:    user.ID,
		Email: user.Email,
		Phone: user.Phone,
	}

	// Perform update
	updatedUser, err := s.repo.Update(ctx, userModel)
	if err != nil {
		return nil, err
	}

	// Convert to DTO
	dto := updatedUser.ToDTO()
	return &dto, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	_, err := s.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete the user
	return s.repo.Delete(ctx, id)
}

func (s *userService) CountUsers(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.repo.Count(ctx, filters)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.UserDTO, error) {
	// Retrieve user
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.ErrNotFound, "user not found", err)
		}
		return nil, errors.Wrap(err, errors.ErrDatabase, "error retrieving user")
	}

	// Convert to DTO
	dto := user.ToDTO()
	return &dto, nil
}

func (s *userService) SearchUsers(ctx context.Context, opts base.ListOptions) ([]models.UserDTO, int, error) {
	// Validate ListOptions
	if err := opts.Validate(); err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrValidation, "invalid list options")
	}

	return s.ListUsers(ctx, opts)
}
