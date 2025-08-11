package services

import (
	"context"

	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/pkg/base"
)

// Thread-related service methods for managing discussion threads
type ThreadService interface {
	// Thread creation and retrieval methods
	CreateThread(ctx context.Context, thread *models.CreateThread) (*models.Thread, error)
	GetThreadByID(ctx context.Context, id string) (*models.ThreadDTO, error)
	GetThreadByEvent(ctx context.Context, eventID string) (*models.ThreadDTO, error)
	GetActiveThreads(ctx context.Context, limit int) ([]models.ThreadDTO, error)

	// Thread listing and search methods
	ListThreads(ctx context.Context, opts base.ListOptions) ([]models.ThreadDTO, error)
	SearchThreads(ctx context.Context, query string, opts base.ListOptions) ([]models.ThreadDTO, int, error)
	CountThreads(ctx context.Context, filters []base.FilterOption) (int, error)

	// Thread modification methods
	UpdateThread(ctx context.Context, thread *models.Thread) (*models.Thread, error)
	DeleteThread(ctx context.Context, id string) error
	JoinThread(ctx context.Context, threadID, userID string) error
}

// Message-related service methods for managing discussion messages
type MessageService interface {
	// Message creation and retrieval methods
	CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	GetMessageByID(ctx context.Context, id string) (*models.MessageDTO, error)
	GetMessagesByThread(ctx context.Context, threadID string) ([]models.MessageDTO, error)
	GetMessagesByUser(ctx context.Context, userID string) ([]models.MessageDTO, error)
	GetRecentMessages(ctx context.Context, limit int) ([]models.MessageDTO, error)

	// Message listing and search methods
	ListMessages(ctx context.Context, opts base.ListOptions) ([]models.MessageDTO, error)
	SearchMessages(ctx context.Context, query string, opts base.ListOptions) ([]models.MessageDTO, int, error)
	CountMessages(ctx context.Context, filters []base.FilterOption) (int, error)

	// Message modification methods
	UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
	DeleteMessage(ctx context.Context, id string) error
}

// Participant-related service methods for managing discussion participants
type ParticipantService interface {
	// Participant creation and retrieval methods
	CreateParticipant(ctx context.Context, participant *models.Participant) (*models.Participant, error)
	GetParticipantByID(ctx context.Context, id string) (*models.ParticipantDTO, error)
	GetParticipantsByThread(ctx context.Context, threadID string) ([]models.ParticipantDTO, error)
	GetParticipantByThread(ctx context.Context, userID, threadID string) (*models.ParticipantDTO, error)
	GetThreadParticipants(ctx context.Context, threadID string) ([]models.ParticipantDTO, error)

	// Participant listing and search methods
	ListParticipants(ctx context.Context, opts base.ListOptions) ([]models.ParticipantDTO, error)
	SearchParticipants(ctx context.Context, query string, opts base.ListOptions) ([]models.ParticipantDTO, int, error)
	CountParticipants(ctx context.Context, filters []base.FilterOption) (int, error)

	// Participant modification methods
	UpdateParticipant(ctx context.Context, participant *models.Participant) (*models.Participant, error)
	RemoveParticipant(ctx context.Context, threadID, userID string) error
}
