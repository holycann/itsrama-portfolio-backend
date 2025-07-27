package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type ThreadRepository interface {
	repository.BaseRepository[models.Thread]

	// Specialized methods for threads
	FindByTitle(ctx context.Context, title string) (*models.Thread, error)
	FindThreadsByEvent(ctx context.Context, eventID string) ([]models.Thread, error)
	FindActiveThreads(ctx context.Context, limit int) ([]models.Thread, error)
}

type MessageRepository interface {
	repository.BaseRepository[models.Message]

	// Specialized methods for messages
	FindMessagesByThread(ctx context.Context, threadID string) ([]models.Message, error)
	FindMessagesByUser(ctx context.Context, userID string) ([]models.Message, error)
	FindRecentMessages(ctx context.Context, limit int) ([]models.Message, error)
}
