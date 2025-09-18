package tech_stack

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/internal/base"
	"github.com/holycann/itsrama-portfolio-backend/internal/validator"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/supabase"
	storage_go "github.com/supabase-community/storage-go"
)

type TechStackService interface {
	CreateTechStack(ctx context.Context, techStackCreate *TechStackCreate) (*TechStack, error)
	GetTechStackByID(ctx context.Context, id string) (*TechStack, error)
	UpdateTechStack(ctx context.Context, techStackUpdate *TechStackUpdate) (*TechStack, error)
	DeleteTechStack(ctx context.Context, id string) error
	ListTechStacks(ctx context.Context, opts base.ListOptions) ([]TechStack, error)
	CountTechStacks(ctx context.Context, filters []base.FilterOption) (int, error)
	SearchTechStacks(ctx context.Context, opts base.ListOptions) ([]TechStack, int, error)
	BulkCreateTechStacks(ctx context.Context, techStacksCreate []*TechStackCreate) ([]TechStack, error)
	BulkUpdateTechStacks(ctx context.Context, techStacksUpdate []*TechStackUpdate) ([]TechStack, error)
	BulkDeleteTechStacks(ctx context.Context, ids []string) error
}

type techStackService struct {
	techStackRepo TechStackRepository
	storage       supabase.SupabaseStorage
}

func NewTechStackService(techStackRepo TechStackRepository, storage supabase.SupabaseStorage) TechStackService {
	return &techStackService{
		techStackRepo: techStackRepo,
		storage:       storage,
	}
}

func (s *techStackService) CreateTechStack(ctx context.Context, techStackCreate *TechStackCreate) (*TechStack, error) {
	// Validate input
	if err := validator.ValidateModel(techStackCreate); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	techStack := techStackCreate.ToTechStack()
	techStack.ID = uuid.New()
	techStack.CreatedAt = &now
	techStack.UpdatedAt = &now

	// Upload image if provided
	if techStackCreate.Image != nil {
		imageURL, err := s.uploadTechStackImage(ctx, techStack.ID.String(), techStackCreate.Image)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload tech stack image",
				errors.WithContext("tech_stack_id", techStack.ID),
			)
		}
		techStack.ImageUrl = imageURL
	}

	// Create tech stack in repository
	createdTechStack, err := s.techStackRepo.Create(ctx, &techStack)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to create tech stack",
			errors.WithContext("tech_stack_name", techStack.Name),
		)
	}

	return createdTechStack, nil
}

func (s *techStackService) GetTechStackByID(ctx context.Context, id string) (*TechStack, error) {
	if id == "" {
		return nil, fmt.Errorf("tech stack ID cannot be empty")
	}

	techStacks, err := s.techStackRepo.FindByField(ctx, "id", id)
	if err != nil {
		return nil, err
	}

	if len(techStacks) == 0 {
		return nil, fmt.Errorf("tech stack with ID %s not found", id)
	}

	return &techStacks[0], nil
}

func (s *techStackService) UpdateTechStack(ctx context.Context, techStackUpdate *TechStackUpdate) (*TechStack, error) {
	// Validate input
	if err := validator.ValidateModel(techStackUpdate); err != nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Invalid tech stack payload",
			err,
			errors.WithContext("payload", techStackUpdate),
		)
	}

	// Retrieve existing tech stack
	existingTechStack, err := s.GetTechStackByID(ctx, techStackUpdate.ID.String())
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to retrieve existing tech stack",
			errors.WithContext("tech_stack_id", techStackUpdate.ID),
		)
	}

	now := time.Now().UTC()
	techStack := techStackUpdate.ToTechStack()
	techStack.CreatedAt = existingTechStack.CreatedAt
	techStack.UpdatedAt = &now

	// Conditionally update fields
	if !validator.IsValueChanged(&existingTechStack.Name, &techStack.Name) {
		techStack.Name = existingTechStack.Name
	}
	if !validator.IsValueChanged(&existingTechStack.Category, &techStack.Category) {
		techStack.Category = existingTechStack.Category
	}
	if !validator.IsValueChanged(&existingTechStack.Version, &techStack.Version) {
		techStack.Version = existingTechStack.Version
	}
	if !validator.IsValueChanged(&existingTechStack.Role, &techStack.Role) {
		techStack.Role = existingTechStack.Role
	}

	// Upload image if provided
	if techStackUpdate.Image != nil {
		imageURL, err := s.uploadTechStackImage(ctx, techStack.ID.String(), techStackUpdate.Image)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload tech stack image",
				errors.WithContext("tech_stack_id", techStack.ID),
			)
		}
		techStack.ImageUrl = imageURL
	} else {
		techStack.ImageUrl = existingTechStack.ImageUrl
	}

	// Update tech stack in repository
	updatedTechStack, err := s.techStackRepo.Update(ctx, &techStack)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to update tech stack",
			errors.WithContext("tech_stack_id", techStack.ID),
		)
	}

	return updatedTechStack, nil
}

func (s *techStackService) DeleteTechStack(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"Tech stack ID cannot be empty",
			nil,
		)
	}

	// Retrieve existing tech stack to get image URL
	existingTechStack, err := s.GetTechStackByID(ctx, id)
	if err != nil {
		return errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to retrieve existing tech stack",
			errors.WithContext("tech_stack_id", id),
		)
	}

	// Delete tech stack from repository
	err = s.techStackRepo.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to delete tech stack",
			errors.WithContext("tech_stack_id", id),
		)
	}

	// Delete associated image if exists
	if existingTechStack.ImageUrl != "" {
		imagePath := filepath.Join("itsrama/images/tech_stack/", filepath.Base(existingTechStack.ImageUrl))
		_, err = s.storage.Delete(ctx, imagePath)
		if err != nil {
			// Log the error but don't return it to avoid blocking the deletion
			fmt.Printf("Failed to delete tech stack image: %v\n", err)
		}
	}

	return nil
}

func (s *techStackService) ListTechStacks(ctx context.Context, opts base.ListOptions) ([]TechStack, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err,
			errors.ErrValidation,
			"Invalid list options",
			errors.WithContext("options", opts),
		)
	}

	techStacks, err := s.techStackRepo.List(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to list tech stacks",
			errors.WithContext("options", opts),
		)
	}

	return techStacks, nil
}

func (s *techStackService) CountTechStacks(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.techStackRepo.Count(ctx, filters)
}

func (s *techStackService) SearchTechStacks(ctx context.Context, opts base.ListOptions) ([]TechStack, int, error) {
	if err := opts.Validate(); err != nil {
		return nil, 0, errors.Wrap(err,
			errors.ErrValidation,
			"Invalid list options",
			errors.WithContext("options", opts),
		)
	}

	return s.techStackRepo.Search(ctx, opts)
}

func (s *techStackService) BulkCreateTechStacks(ctx context.Context, techStacksCreate []*TechStackCreate) ([]TechStack, error) {
	techStacks := make([]TechStack, len(techStacksCreate))

	for i, techStackCreate := range techStacksCreate {
		createdTechStack, err := s.CreateTechStack(ctx, techStackCreate)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to create tech stack",
				errors.WithContext("tech_stack_name", techStackCreate.Name),
			)
		}
		techStacks[i] = *createdTechStack
	}

	return techStacks, nil
}

func (s *techStackService) BulkUpdateTechStacks(ctx context.Context, techStacksUpdate []*TechStackUpdate) ([]TechStack, error) {
	techStacks := make([]TechStack, len(techStacksUpdate))

	for i, techStackUpdate := range techStacksUpdate {
		updatedTechStack, err := s.UpdateTechStack(ctx, techStackUpdate)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to update tech stack",
				errors.WithContext("tech_stack_id", techStackUpdate.ID),
			)
		}
		techStacks[i] = *updatedTechStack
	}

	return techStacks, nil
}

func (s *techStackService) BulkDeleteTechStacks(ctx context.Context, ids []string) error {
	for _, id := range ids {
		err := s.DeleteTechStack(ctx, id)
		if err != nil {
			return errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to delete tech stack",
				errors.WithContext("tech_stack_id", id),
			)
		}
	}
	return nil
}

func (s *techStackService) uploadTechStackImage(ctx context.Context, techStackID string, file *multipart.FileHeader) (string, error) {
	if techStackID == "" {
		return "", fmt.Errorf("tech stack ID cannot be empty")
	}
	if file == nil {
		return "", fmt.Errorf("file is required")
	}

	destPath := fmt.Sprintf("images/tech_stack/%s%s", techStackID, filepath.Ext(file.Filename))

	_, err := s.storage.Upload(ctx, file, destPath, storage_go.FileOptions{
		ContentType: func(s string) *string { return &s }("image"),
		Upsert:      func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		return "", errors.Wrap(err,
			errors.ErrInternal,
			"Failed to upload tech stack image",
			errors.WithContext("tech_stack_id", techStackID),
		)
	}

	signedURL, err := s.storage.GetPublicURL(destPath)
	if err != nil {
		return "", errors.Wrap(err,
			errors.ErrInternal,
			"Failed to get public URL for tech stack image",
			errors.WithContext("dest_path", destPath),
		)
	}

	return signedURL, nil
}
