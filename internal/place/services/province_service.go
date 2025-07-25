package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/place/repositories"
)

type provinceService struct {
	provinceRepo repositories.ProvinceRepository
}

// NewProvinceService creates a new instance of the province service
// with the given province repository.
func NewProvinceService(provinceRepo repositories.ProvinceRepository) ProvinceService {
	return &provinceService{
		provinceRepo: provinceRepo,
	}
}

// CreateProvince adds a new province to the database
// Validates the province object before creating
func (s *provinceService) CreateProvince(ctx context.Context, province *models.Province) error {
	// Validate province object
	if province == nil {
		return fmt.Errorf("province cannot be nil")
	}

	// Validate required fields (example validation)
	if province.Name == "" {
		return fmt.Errorf("province name is required")
	}

	province.ID = uuid.NewString()

	// Call repository to create province
	return s.provinceRepo.Create(ctx, province)
}

// GetProvinces retrieves a list of provinces with pagination
func (s *provinceService) GetProvinces(ctx context.Context, limit, offset int) ([]*models.Province, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	// Retrieve provinces from repository
	provinces, err := s.provinceRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert []models.Province to []*models.Province
	provincePtrs := make([]*models.Province, len(provinces))
	for i := range provinces {
		provincePtrs[i] = &provinces[i]
	}

	return provincePtrs, nil
}

// GetProvinceByID retrieves a single province by its unique identifier
func (s *provinceService) GetProvinceByID(ctx context.Context, id string) (*models.Province, error) {
	// Validate ID
	if id == "" {
		return nil, fmt.Errorf("province ID cannot be empty")
	}

	// Retrieve province from repository
	return s.provinceRepo.FindByID(ctx, id)
}

// GetProvinceByName retrieves a province by its name
// Note: This method is not directly supported by the current repository implementation
// You might need to add a custom method in the repository or implement filtering
func (s *provinceService) GetProvinceByName(ctx context.Context, name string) (*models.Province, error) {
	// Validate name
	if name == "" {
		return nil, fmt.Errorf("province name cannot be empty")
	}

	// Since the current repository doesn't have a direct method for this,
	// we'll use a workaround by listing all provinces and finding by name
	provinces, err := s.provinceRepo.List(ctx, 1, 0)
	if err != nil {
		return nil, err
	}

	// Find province by name (linear search)
	for _, province := range provinces {
		if province.Name == name {
			return &province, nil
		}
	}

	return nil, fmt.Errorf("province with name %s not found", name)
}

// UpdateProvince updates an existing province in the database
func (s *provinceService) UpdateProvince(ctx context.Context, province *models.Province) error {
	// Validate province object
	if province == nil {
		return fmt.Errorf("province cannot be nil")
	}

	// Validate required fields
	if province.ID == "" {
		return fmt.Errorf("province ID is required for update")
	}

	// Call repository to update province
	return s.provinceRepo.Update(ctx, province)
}

// DeleteProvince removes a province from the database by its ID
func (s *provinceService) DeleteProvince(ctx context.Context, id string) error {
	// Validate ID
	if id == "" {
		return fmt.Errorf("province ID cannot be empty")
	}

	// Call repository to delete province
	return s.provinceRepo.Delete(ctx, id)
}

// Count calculates the total number of stored provinces
func (s *provinceService) Count(ctx context.Context) (int, error) {
	return s.provinceRepo.Count(ctx)
}
