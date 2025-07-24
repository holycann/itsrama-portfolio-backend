package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
)

type ThreadService interface {
	CreateThread(ctx context.Context, thread *models.Thread) error
	GetThreads(ctx context.Context, limit, offset int) ([]*models.Thread, error)
	GetThreadByID(ctx context.Context, id string) (*models.Thread, error)
	GetThreadByTitle(ctx context.Context, title string) (*models.Thread, error)
	UpdateThread(ctx context.Context, thread *models.Thread) error
	DeleteThread(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}

type MessageService interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessages(ctx context.Context, limit, offset int) ([]*models.Message, error)
	GetMessageByID(ctx context.Context, id string) (*models.Message, error)
	GetMessagesByThreadID(ctx context.Context, id string, limit, offset int) ([]*models.Message, error)
	GetMessageByContent(ctx context.Context, content string) ([]*models.Message, error)
	UpdateMessage(ctx context.Context, message *models.Message) error
	DeleteMessage(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
	CountByThreadID(ctx context.Context, id string) (int, error)
}
