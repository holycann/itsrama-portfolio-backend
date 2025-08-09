package repositories

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

type ThreadRepository interface {
	base.BaseRepository[models.Thread, models.ThreadDTO]
	FindThreadByEvent(ctx context.Context, eventID string) (*models.ThreadDTO, error)
	FindActiveThreads(ctx context.Context, limit int) ([]models.ThreadDTO, error)
	JoinThread(ctx context.Context, threadID, userID string) error
}

type MessageRepository interface {
	base.BaseRepository[models.Message, models.MessageDTO]
	FindMessagesByThread(ctx context.Context, threadID string) ([]models.MessageDTO, error)
	FindMessagesByUser(ctx context.Context, userID string) ([]models.MessageDTO, error)
	FindRecentMessages(ctx context.Context, limit int) ([]models.MessageDTO, error)
}

type ParticipantRepository interface {
	base.BaseRepository[models.Participant, models.ParticipantDTO]
	FindParticipantsByThread(ctx context.Context, threadID string) ([]models.ParticipantDTO, error)
	FindThreadParticipants(ctx context.Context, threadID string) ([]models.ParticipantDTO, error)
	RemoveParticipant(ctx context.Context, threadID, userID string) error
}
