package services

import (
	"context"
	"fmt"

	// "time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/discussion/models"
	"github.com/holycann/cultour-backend/internal/discussion/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type participantService struct {
	participantRepo repositories.ParticipantRepository
}

func NewParticipantService(participantRepo repositories.ParticipantRepository) ParticipantService {
	return &participantService{
		participantRepo: participantRepo,
	}
}

func (s *participantService) CreateParticipant(ctx context.Context, participant *models.Participant) error {
	// Validasi objek participant
	if participant == nil {
		return fmt.Errorf("participant tidak boleh nil")
	}

	// participant.JoinedAt = time.Now()

	// Panggil repository untuk membuat participant
	return s.participantRepo.Create(ctx, participant)
}

func (s *participantService) GetParticipantByID(ctx context.Context, id string) (*models.ResponseParticipant, error) {
	// Validasi ID
	if id == "" {
		return nil, fmt.Errorf("participant ID tidak boleh kosong")
	}

	// Ambil participant dari repository
	return s.participantRepo.FindByID(ctx, id)
}

func (s *participantService) ListParticipants(ctx context.Context, opts repository.ListOptions) ([]models.ResponseParticipant, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.participantRepo.List(ctx, opts)
}

func (s *participantService) UpdateParticipant(ctx context.Context, participant *models.Participant) error {
	// Validasi objek participant
	if participant == nil {
		return fmt.Errorf("participant tidak boleh nil")
	}

	// Validasi field yang wajib diisi
	if participant.ThreadID == uuid.Nil || participant.UserID == uuid.Nil {
		return fmt.Errorf("thread ID dan user ID wajib diisi untuk update")
	}

	// Update timestamp
	// participant.UpdatedAt = time.Now()

	// Panggil repository untuk update participant
	return s.participantRepo.Update(ctx, participant)
}

func (s *participantService) CountParticipants(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.participantRepo.Count(ctx, filters)
}

func (s *participantService) GetParticipantsByThread(ctx context.Context, threadID string) ([]models.ResponseParticipant, error) {
	// Validasi thread ID
	if threadID == "" {
		return nil, fmt.Errorf("thread ID tidak boleh kosong")
	}

	return s.participantRepo.FindParticipantsByThread(ctx, threadID)
}

func (s *participantService) GetThreadParticipants(ctx context.Context, threadID string) ([]models.ResponseParticipant, error) {
	// Validasi thread ID
	if threadID == "" {
		return nil, fmt.Errorf("thread ID tidak boleh kosong")
	}

	return s.participantRepo.FindThreadParticipants(ctx, threadID)
}

func (s *participantService) RemoveParticipant(ctx context.Context, threadID, userID string) error {
	// Validasi input parameters
	if threadID == "" {
		return fmt.Errorf("thread ID tidak boleh kosong")
	}
	if userID == "" {
		return fmt.Errorf("user ID tidak boleh kosong")
	}

	// Panggil repository untuk menghapus participant
	return s.participantRepo.RemoveParticipant(ctx, threadID, userID)
}

func (s *participantService) SearchParticipants(ctx context.Context, query string, opts repository.ListOptions) ([]models.ResponseParticipant, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	// Set search query
	opts.SearchQuery = query

	// Panggil repository untuk mencari participants
	results, _, err := s.participantRepo.Search(ctx, opts)
	return results, err
}
