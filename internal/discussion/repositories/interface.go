package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

// ThreadRepository mendefinisikan antarmuka untuk operasi CRUD pada entitas Thread
type ThreadRepository interface {
	repository.BaseRepository[models.Thread]

	FindByTitle(ctx context.Context, id string) (*models.Thread, error)
}

// MessageRepository mendefinisikan antarmuka untuk operasi CRUD pada entitas Message
type MessageRepository interface {
	repository.BaseRepository[models.Message]

	ListByThreadID(ctx context.Context, threadID string, limit, offset int) ([]models.Message, error)
	CountByThreadID(ctx context.Context, threadID string) (int, error)
}
