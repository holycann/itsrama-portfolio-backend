package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type ThreadService interface {
	CreateThread(ctx context.Context, thread *models.Thread) error
	GetThreadByID(ctx context.Context, id string) (*models.ResponseThread, error)
	ListThreads(ctx context.Context, opts repository.ListOptions) ([]models.ResponseThread, error)
	UpdateThread(ctx context.Context, thread *models.Thread) error
	DeleteThread(ctx context.Context, id string) error
	CountThreads(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetThreadByEvent(ctx context.Context, eventID string) (*models.ResponseThread, error)
	GetActiveThreads(ctx context.Context, limit int) ([]models.ResponseThread, error)
	SearchThreads(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseThread, error)
	JoinThread(ctx context.Context, threadID, userID string) error
}

type MessageService interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessageByID(ctx context.Context, id string) (*models.ResponseMessage, error)
	ListMessages(ctx context.Context, opts repository.ListOptions) ([]models.ResponseMessage, error)
	UpdateMessage(ctx context.Context, message *models.Message) error
	DeleteMessage(ctx context.Context, id string) error
	CountMessages(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetMessagesByThread(ctx context.Context, threadID string) ([]models.ResponseMessage, error)
	GetMessagesByUser(ctx context.Context, userID string) ([]models.ResponseMessage, error)
	GetRecentMessages(ctx context.Context, limit int) ([]models.ResponseMessage, error)
	SearchMessages(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseMessage, error)
}

type ParticipantService interface {
	CreateParticipant(ctx context.Context, participant *models.Participant) error
	GetParticipantByID(ctx context.Context, id string) (*models.ResponseParticipant, error)
	ListParticipants(ctx context.Context, opts repository.ListOptions) ([]models.ResponseParticipant, error)
	UpdateParticipant(ctx context.Context, participant *models.Participant) error
	CountParticipants(ctx context.Context, filters []repository.FilterOption) (int, error)
	GetParticipantsByThread(ctx context.Context, threadID string) ([]models.ResponseParticipant, error)
	GetThreadParticipants(ctx context.Context, threadID string) ([]models.ResponseParticipant, error)
	RemoveParticipant(ctx context.Context, threadID, userID string) error
	SearchParticipants(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseParticipant, error)
}
