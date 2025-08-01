package services

import (
	"context"
	"mime/multipart"

	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type UserBadgeService interface {
	CreateUserBadge(ctx context.Context, userBadgeCreate *models.UserBadgeCreate) error
	GetUserBadgeByID(ctx context.Context, id string) (*models.UserBadge, error)
	ListUserBadges(ctx context.Context, opts repository.ListOptions) ([]models.UserBadge, error)
	DeleteUserBadge(ctx context.Context, id string) error
	CountUserBadges(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetUserBadgesByUser(ctx context.Context, userID string) ([]models.UserBadge, error)
	GetUserBadgesByBadge(ctx context.Context, badgeID string) ([]models.UserBadge, error)
	SearchUserBadges(ctx context.Context, query string, opts repository.ListOptions) ([]models.UserBadge, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context, opts repository.ListOptions) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
	CountUsers(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	SearchUsers(ctx context.Context, query string, opts repository.ListOptions) ([]models.User, error)
}

type UserProfileService interface {
	CreateProfile(ctx context.Context, userProfile *models.UserProfile) error
	GetProfileByID(ctx context.Context, id string) (*models.UserProfile, error)
	ListProfiles(ctx context.Context, opts repository.ListOptions) ([]models.UserProfile, error)
	UpdateProfile(ctx context.Context, userProfile *models.UserProfile, avatar *multipart.FileHeader, identity *multipart.FileHeader) error
	DeleteProfile(ctx context.Context, id string) error
	CountProfiles(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetProfileByUserID(ctx context.Context, userID string) (*models.UserProfile, error)
	SearchProfiles(ctx context.Context, query string, opts repository.ListOptions) ([]models.UserProfile, error)
}
