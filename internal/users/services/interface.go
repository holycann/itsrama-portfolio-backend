package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/users/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
}

type UserProfileService interface {
	CreateProfile(ctx context.Context, userProfile *models.UserProfile) error
	GetProfiles(ctx context.Context, limit, offset int) ([]*models.UserProfile, error)
	GetProfileByID(ctx context.Context, id string) (*models.UserProfile, error)
	GetProfileByUserID(ctx context.Context, userID string) (*models.UserProfile, error)
	UpdateProfile(ctx context.Context, userProfile *models.UserProfile) error
	DeleteProfile(ctx context.Context, id string) error
}
