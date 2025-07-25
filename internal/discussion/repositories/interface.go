package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

// ThreadRepository defines the interface for CRUD operations on Thread entities
type ThreadRepository interface {
	repository.BaseRepository[models.Thread]

	FindByTitle(ctx context.Context, id string) (*models.Thread, error)
}

// MessageRepository defines the interface for CRUD operations on Message entities
type MessageRepository interface {
	repository.BaseRepository[models.Message]

	ListByThreadID(ctx context.Context, threadID string, limit, offset int) ([]models.Message, error)
	CountByThreadID(ctx context.Context, threadID string) (int, error)
}
