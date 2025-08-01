package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type ThreadRepository interface {
	repository.BaseRepository[models.Thread, models.ResponseThread]
	FindThreadByEvent(ctx context.Context, eventID string) (*models.ResponseThread, error)
	FindActiveThreads(ctx context.Context, limit int) ([]models.ResponseThread, error)
	JoinThread(ctx context.Context, threadID, userID string) error
}

type MessageRepository interface {
	repository.BaseRepository[models.Message, models.ResponseMessage]
	FindMessagesByThread(ctx context.Context, threadID string) ([]models.ResponseMessage, error)
	FindMessagesByUser(ctx context.Context, userID string) ([]models.ResponseMessage, error)
	FindRecentMessages(ctx context.Context, limit int) ([]models.ResponseMessage, error)
}
type ParticipantRepository interface {
	repository.BaseRepository[models.Participant, models.ResponseParticipant]
	FindParticipantsByThread(ctx context.Context, threadID string) ([]models.ResponseParticipant, error)
	FindThreadParticipants(ctx context.Context, threadID string) ([]models.ResponseParticipant, error)
	RemoveParticipant(ctx context.Context, threadID, userID string) error
}
