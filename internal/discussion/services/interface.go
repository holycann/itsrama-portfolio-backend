package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type ThreadService interface {
	CreateThread(ctx context.Context, thread *models.Thread) error
	GetThreadByID(ctx context.Context, id string) (*models.Thread, error)
	ListThreads(ctx context.Context, opts repository.ListOptions) ([]models.Thread, error)
	UpdateThread(ctx context.Context, thread *models.Thread) error
	DeleteThread(ctx context.Context, id string) error
	CountThreads(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetThreadByTitle(ctx context.Context, title string) (*models.Thread, error)
	GetThreadsByEvent(ctx context.Context, eventID string) ([]models.Thread, error)
	GetActiveThreads(ctx context.Context, limit int) ([]models.Thread, error)
	SearchThreads(ctx context.Context, query string, opts repository.ListOptions) ([]models.Thread, error)
}

type MessageService interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessageByID(ctx context.Context, id string) (*models.Message, error)
	ListMessages(ctx context.Context, opts repository.ListOptions) ([]models.Message, error)
	UpdateMessage(ctx context.Context, message *models.Message) error
	DeleteMessage(ctx context.Context, id string) error
	CountMessages(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetMessagesByThread(ctx context.Context, threadID string) ([]models.Message, error)
	GetMessagesByUser(ctx context.Context, userID string) ([]models.Message, error)
	GetRecentMessages(ctx context.Context, limit int) ([]models.Message, error)
	SearchMessages(ctx context.Context, query string, opts repository.ListOptions) ([]models.Message, error)
}
