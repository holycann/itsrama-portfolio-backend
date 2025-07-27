package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/repositories"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type provinceService struct {
	provinceRepo repositories.ProvinceRepository
}

func NewProvinceService(provinceRepo repositories.ProvinceRepository) ProvinceService {
	return &provinceService{
		provinceRepo: provinceRepo,
	}
}

func (s *provinceService) CreateProvince(ctx context.Context, province *models.Province) error {
	// Validate province object
	if province == nil {
		return fmt.Errorf("province cannot be nil")
	}

	// Validate required fields
	if province.Name == "" {
		return fmt.Errorf("province name is required")
	}

	// Set default values
	province.ID = uuid.New()
	now := time.Now()
	province.CreatedAt = now

	// Call repository to create province
	return s.provinceRepo.Create(ctx, province)
}

func (s *provinceService) GetProvinceByID(ctx context.Context, id string) (*models.Province, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("province ID cannot be empty")
	}

	// Retrieve province from repository
	return s.provinceRepo.FindByID(ctx, id)
}

func (s *provinceService) ListProvinces(ctx context.Context, opts repository.ListOptions) ([]models.Province, error) {
	// Set default values if not provided
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}

	return s.provinceRepo.List(ctx, opts)
}

func (s *provinceService) UpdateProvince(ctx context.Context, province *models.Province) error {
	// Validate province object
	if province == nil {
		return fmt.Errorf("province cannot be nil")
	}

	// Validate required fields
	if province.ID == uuid.Nil {
		return fmt.Errorf("province ID is required for update")
	}

	// Update timestamp
	province.UpdatedAt = time.Now()

	// Call repository to update province
	return s.provinceRepo.Update(ctx, province)
}

func (s *provinceService) DeleteProvince(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("province ID cannot be empty")
	}

	// Call repository to delete province
	return s.provinceRepo.Delete(ctx, id)
}

func (s *provinceService) CountProvinces(ctx context.Context, filters []repository.FilterOption) (int, error) {
	return s.provinceRepo.Count(ctx, filters)
}

func (s *provinceService) GetProvinceByName(ctx context.Context, name string) (*models.Province, error) {
	return s.provinceRepo.FindProvinceByName(ctx, name)
}

func (s *provinceService) SearchProvinces(ctx context.Context, query string, opts repository.ListOptions) ([]models.Province, error) {
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
			Field:    "name",
			Operator: "like",
			Value:    query,
		},
		repository.FilterOption{
			Field:    "description",
			Operator: "like",
			Value:    query,
		},
	)

	return s.provinceRepo.List(ctx, opts)
}
