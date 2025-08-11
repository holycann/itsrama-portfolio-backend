package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/repositories"
	"github.com/holycann/cultour-backend/pkg/base"
	"github.com/holycann/cultour-backend/pkg/errors"
)

type provinceService struct {
	provinceRepo repositories.ProvinceRepository
}

func NewProvinceService(provinceRepo repositories.ProvinceRepository) ProvinceService {
	return &provinceService{
		provinceRepo: provinceRepo,
	}
}

func (s *provinceService) CreateProvince(ctx context.Context, province *models.Province) (*models.Province, error) {
	// Validate province object
	if province == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Province cannot be nil",
			nil,
			errors.WithContext("input", "nil province"),
		)
	}

	// Validate required fields
	if err := base.ValidateModel(province); err != nil {
		return nil, err
	}

	// Set default values
	province.ID = uuid.New()
	now := time.Now()
	province.CreatedAt = &now
	province.UpdatedAt = &now

	// Call repository to create province
	return s.provinceRepo.Create(ctx, province)
}

func (s *provinceService) GetProvinceByID(ctx context.Context, id string) (*models.ProvinceDTO, error) {
	// Validate ID
	if id == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Province ID cannot be empty",
			nil,
			errors.WithContext("input", "empty ID"),
		)
	}

	// Retrieve province from repository
	return s.provinceRepo.FindByID(ctx, id)
}

func (s *provinceService) ListProvinces(ctx context.Context, opts base.ListOptions) ([]models.ProvinceDTO, error) {
	return s.provinceRepo.List(ctx, opts)
}

func (s *provinceService) UpdateProvince(ctx context.Context, province *models.Province) (*models.Province, error) {
	// Validate province object
	if province == nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Province cannot be nil",
			nil,
			errors.WithContext("input", "nil province"),
		)
	}

	// Validate required fields
	if province.ID == uuid.Nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Province ID is required for update",
			nil,
			errors.WithContext("input", "missing ID"),
		)
	}

	// Validate model
	if err := base.ValidateModel(province); err != nil {
		return nil, err
	}

	// Update timestamp
	now := time.Now()
	province.UpdatedAt = &now

	return s.provinceRepo.Update(ctx, province)
}

func (s *provinceService) DeleteProvince(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"Province ID cannot be empty",
			nil,
			errors.WithContext("input", "empty ID"),
		)
	}

	// Call repository to delete province
	return s.provinceRepo.Delete(ctx, id)
}

func (s *provinceService) CountProvinces(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.provinceRepo.Count(ctx, filters)
}

func (s *provinceService) GetProvinceByName(ctx context.Context, name string) (*models.ProvinceDTO, error) {
	// Validate name
	if name == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Province name cannot be empty",
			nil,
			errors.WithContext("input", "empty name"),
		)
	}

	return s.provinceRepo.FindProvinceByName(ctx, name)
}

func (s *provinceService) GetProvinceByCode(ctx context.Context, code string) (*models.ProvinceDTO, error) {
	// Validate code
	if code == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Province code cannot be empty",
			nil,
			errors.WithContext("input", "empty code"),
		)
	}

	return s.provinceRepo.FindProvinceByCode(ctx, code)
}

func (s *provinceService) ListProvincesByRegion(ctx context.Context, region string) ([]models.ProvinceDTO, error) {
	// Validate region
	if region == "" {
		return nil, errors.New(
			errors.ErrValidation,
			"Region cannot be empty",
			nil,
			errors.WithContext("input", "empty region"),
		)
	}

	return s.provinceRepo.ListProvincesByRegion(ctx, region)
}

func (s *provinceService) SearchProvinces(ctx context.Context, opts base.ListOptions) ([]models.ProvinceDTO, int, error) {
	return s.provinceRepo.Search(ctx, opts)
}
