package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type messageService struct {
	messageRepo repositories.MessageRepository
}

func NewMessageService(messageRepo repositories.MessageRepository) MessageService {
	return &messageService{
		messageRepo: messageRepo,
	}
}

func (s *messageService) CreateMessage(ctx context.Context, message *models.Message) error {
	// Validate message object
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	// Validate required fields
	if message.Content == "" {
		return fmt.Errorf("message content is required")
	}

	// Set default values
	message.ID = uuid.New()
	message.CreatedAt = time.Now()

	// Set default type if not provided
	if message.Type == "" {
		message.Type = models.DiscussionMessageType
	}

	// Call repository to create message
	return s.messageRepo.Create(ctx, message)
}

func (s *messageService) GetMessageByID(ctx context.Context, id string) (*models.ResponseMessage, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("message ID cannot be empty")
	}

	// Retrieve message from repository
	return s.messageRepo.FindByID(ctx, id)
}

func (s *messageService) ListMessages(ctx context.Context, opts repository.ListOptions) ([]models.ResponseMessage, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.messageRepo.List(ctx, opts)
}

func (s *messageService) UpdateMessage(ctx context.Context, message *models.Message) error {
	// Validate message object
	if message == nil {
		return fmt.Errorf("message cannot be nil")
	}

	// Validate required fields
	if message.ID == uuid.Nil {
		return fmt.Errorf("message ID is required for update")
	}

	// Update timestamp
	message.UpdatedAt = time.Now()

	// Call repository to update message
	return s.messageRepo.Update(ctx, message)
}

func (s *messageService) DeleteMessage(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("message ID cannot be empty")
	}

	// Call repository to delete message
	return s.messageRepo.Delete(ctx, id)
}

func (s *messageService) CountMessages(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.messageRepo.Count(ctx, filters)
}

func (s *messageService) GetMessagesByThread(ctx context.Context, threadID string) ([]models.ResponseMessage, error) {
	return s.messageRepo.FindMessagesByThread(ctx, threadID)
}

func (s *messageService) GetMessagesByUser(ctx context.Context, userID string) ([]models.ResponseMessage, error) {
	return s.messageRepo.FindMessagesByUser(ctx, userID)
}

func (s *messageService) GetRecentMessages(ctx context.Context, limit int) ([]models.ResponseMessage, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.messageRepo.FindRecentMessages(ctx, limit)
}

func (s *messageService) SearchMessages(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseMessage, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	// Add search query to filters
	opts.Filters = append(opts.Filters,
		repository.FilterOption{
			Field:    "content",
			Operator: "like",
			Value:    query,
		},
	)

	return s.messageRepo.List(ctx, opts)
}
