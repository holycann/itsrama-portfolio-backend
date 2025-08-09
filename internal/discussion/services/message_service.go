package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
)

type messageService struct {
	messageRepo repositories.MessageRepository
}

func NewMessageService(messageRepo repositories.MessageRepository) MessageService {
	return &messageService{
		messageRepo: messageRepo,
	}
}

func (s *messageService) CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	// Validate message object
	if message == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Message cannot be nil",
			nil,
		)
	}

	// Validate model
	if err := base.ValidateModel(message); err != nil {
		return nil, err
	}

	// Additional message content validation
	if len(message.Content) == 0 {
		return nil, errors.New(
			errors.ErrValidation,
			"Message content cannot be empty",
			nil,
		)
	}

	if len(message.Content) > 1000 {
		return nil, errors.New(
			errors.ErrValidation,
			"Message content exceeds maximum length of 1000 characters",
			nil,
		)
	}

	// Sanitize message content (remove excessive whitespace)
	message.Content = strings.TrimSpace(message.Content)

	// Validate thread and sender
	if message.ThreadID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread ID is required",
			nil,
		)
	}

	if message.SenderID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Sender ID is required",
			nil,
		)
	}

	// Set default values
	if message.ID == uuid.Nil {
		message.ID = uuid.New()
	}

	now := time.Now().UTC()
	message.CreatedAt = &now
	message.UpdatedAt = &now

	// Set default type if not provided
	if message.Type == "" {
		message.Type = models.DiscussionMessageType
	}

	// Call repository to create message
	return s.messageRepo.Create(ctx, message)
}

func (s *messageService) GetMessageByID(ctx context.Context, id string) (*models.MessageDTO, error) {
	// Validate ID
	if id == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Message ID cannot be empty",
			nil,
		)
	}

	// Retrieve message from repository
	return s.messageRepo.FindByID(ctx, id)
}

func (s *messageService) ListMessages(ctx context.Context, opts base.ListOptions) ([]models.MessageDTO, error) {
	// Set default pagination
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 10
	}

	return s.messageRepo.List(ctx, opts)
}

func (s *messageService) UpdateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	// Validate message object
	if message == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Message cannot be nil",
			nil,
		)
	}

	// Validate model
	if err := base.ValidateModel(message); err != nil {
		return nil, err
	}

	// Validate required fields
	if message.ID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Message ID is required for update",
			nil,
		)
	}

	// Update timestamp
	now := time.Now()
	message.UpdatedAt = &now

	// Call repository to update message
	return s.messageRepo.Update(ctx, message)
}

func (s *messageService) DeleteMessage(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"Message ID cannot be empty",
			nil,
		)
	}

	// Call repository to delete message
	return s.messageRepo.Delete(ctx, id)
}

func (s *messageService) CountMessages(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.messageRepo.Count(ctx, filters)
}

func (s *messageService) GetMessagesByThread(ctx context.Context, threadID string) ([]models.MessageDTO, error) {
	// Validate thread ID
	if threadID == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Thread ID cannot be empty",
			nil,
		)
	}

	return s.messageRepo.FindMessagesByThread(ctx, threadID)
}

func (s *messageService) GetMessagesByUser(ctx context.Context, userID string) ([]models.MessageDTO, error) {
	// Validate user ID
	if userID == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"User ID cannot be empty",
			nil,
		)
	}

	// Retrieve messages by sender ID
	return s.messageRepo.FindMessagesByUser(ctx, userID)
}

func (s *messageService) GetRecentMessages(ctx context.Context, limit int) ([]models.MessageDTO, error) {
	// Set default limit
	if limit <= 0 {
		limit = 10
	}
	return s.messageRepo.FindRecentMessages(ctx, limit)
}

func (s *messageService) SearchMessages(ctx context.Context, query string, opts base.ListOptions) ([]models.MessageDTO, int, error) {
	// Set default pagination
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 10
	}

	// Search messages
	messages, count, err := s.messageRepo.Search(ctx, opts)
	if err != nil {
		return nil, 0, errors.Wrap(
			err,
			errors.ErrDatabase,
			"Failed to search messages",
		)
	}
	return messages, count, nil
}
