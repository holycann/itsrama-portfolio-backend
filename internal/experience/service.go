package experience

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/pkg/base"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/supabase"
	"github.com/holycann/itsrama-portfolio-backend/pkg/validator"
	storage_go "github.com/supabase-community/storage-go"
)

type ExperienceService interface {
	CreateExperience(ctx context.Context, experienceCreate *ExperienceCreate) (*Experience, error)
	GetExperienceByID(ctx context.Context, id string) (*Experience, error)
	UpdateExperience(ctx context.Context, experienceUpdate *ExperienceUpdate) (*Experience, error)
	DeleteExperience(ctx context.Context, id string) error
	ListExperiences(ctx context.Context, opts base.ListOptions) ([]Experience, error)
	CountExperiences(ctx context.Context, filters []base.FilterOption) (int, error)
	GetExperiencesByCompany(ctx context.Context, company string) ([]Experience, error)
	SearchExperiences(ctx context.Context, opts base.ListOptions) ([]Experience, int, error)
	BulkCreateExperiences(ctx context.Context, experiencesCreate []*ExperienceCreate) ([]Experience, error)
	BulkUpdateExperiences(ctx context.Context, experiencesUpdate []*ExperienceUpdate) ([]Experience, error)
	BulkDeleteExperiences(ctx context.Context, ids []string) error
}

type experienceService struct {
	experienceRepo ExperienceRepository
	storage        supabase.SupabaseStorage
}

func NewExperienceService(experienceRepo ExperienceRepository, storage supabase.SupabaseStorage) ExperienceService {
	return &experienceService{
		experienceRepo: experienceRepo,
		storage:        storage,
	}
}

func (s *experienceService) CreateExperience(ctx context.Context, experienceCreate *ExperienceCreate) (*Experience, error) {
	// Validate input
	if err := validator.ValidateModel(experienceCreate); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	experience := experienceCreate.ToExperience()
	experience.ID = uuid.New()
	experience.CreatedAt = &now
	experience.UpdatedAt = &now

	// Upload logo if provided
	if experienceCreate.LogoImage != nil {
		logoURL, err := s.uploadExperienceLogo(ctx, experience.ID.String(), experienceCreate.LogoImage)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload experience logo",
				errors.WithContext("experience_id", experience.ID),
			)
		}
		experience.LogoUrl = logoURL
	}

	// Upload images if provided
	if len(experienceCreate.Images) > 0 {
		imageURLs, err := s.uploadExperienceImages(ctx, experience.ID.String(), experienceCreate.Images)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload experience images",
				errors.WithContext("experience_id", experience.ID),
			)
		}

		experience.ImagesUrl = imageURLs
	}

	// Create experience in repository
	createdExperience, err := s.experienceRepo.Create(ctx, &experience)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to create experience",
			errors.WithContext("experience_company", experience.Company),
		)
	}

	// Create tech stack associations if provided
	if len(experienceCreate.TechStackIds) > 0 {
		for _, techStackID := range experienceCreate.TechStackIds {
			experienceTechStack := &ExperienceTechStack{
				ExperienceID: createdExperience.ID,
				TechStackID:  techStackID,
			}

			_, err := s.experienceRepo.CreateExperienceTechStack(ctx, experienceTechStack)
			if err != nil {
				return nil, errors.Wrap(err,
					errors.ErrDatabase,
					"Failed to create experience tech stack association",
					errors.WithContext("experience_id", createdExperience.ID),
					errors.WithContext("tech_stack_id", techStackID),
				)
			}
		}
	}

	return createdExperience, nil
}

func (s *experienceService) GetExperienceByID(ctx context.Context, id string) (*Experience, error) {
	if id == "" {
		return nil, fmt.Errorf("experience ID cannot be empty")
	}

	return s.experienceRepo.FindByID(ctx, id)
}

func (s *experienceService) UpdateExperience(ctx context.Context, experienceUpdate *ExperienceUpdate) (*Experience, error) {
	// Validate input
	if err := validator.ValidateModel(experienceUpdate); err != nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Invalid experience payload",
			err,
			errors.WithContext("payload", experienceUpdate),
		)
	}

	// Retrieve existing experience
	existingExperience, err := s.GetExperienceByID(ctx, experienceUpdate.ID.String())
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to retrieve existing experience",
			errors.WithContext("experience_id", experienceUpdate.ID),
		)
	}

	now := time.Now().UTC()
	experience := experienceUpdate.ToExperience()
	experience.CreatedAt = existingExperience.CreatedAt
	experience.UpdatedAt = &now

	// Conditionally update fields
	if !validator.IsValueChanged(&existingExperience.Role, &experience.Role) {
		experience.Role = existingExperience.Role
	}
	if !validator.IsValueChanged(&existingExperience.Company, &experience.Company) {
		experience.Company = existingExperience.Company
	}
	if !validator.IsValueChanged(&existingExperience.JobType, &experience.JobType) {
		experience.JobType = existingExperience.JobType
	}
	if !validator.IsValueChanged(&existingExperience.StartDate, &experience.StartDate) {
		experience.StartDate = existingExperience.StartDate
	}
	if !validator.IsValueChanged(&existingExperience.EndDate, &experience.EndDate) {
		experience.EndDate = existingExperience.EndDate
	}
	if !validator.IsValueChanged(&existingExperience.Location, &experience.Location) {
		experience.Location = existingExperience.Location
	}
	if !validator.IsValueChanged(&existingExperience.Arrangement, &experience.Arrangement) {
		experience.Arrangement = existingExperience.Arrangement
	}
	if !validator.IsValueChanged(&existingExperience.WorkDescription, &experience.WorkDescription) {
		experience.WorkDescription = existingExperience.WorkDescription
	}
	if len(experience.Impact) == 0 {
		experience.Impact = existingExperience.Impact
	}
	experience.LogoUrl = existingExperience.LogoUrl
	experience.ImagesUrl = existingExperience.ImagesUrl

	// Upload logo if provided
	if experienceUpdate.LogoImage != nil {
		logoURL, err := s.uploadExperienceLogo(ctx, experience.ID.String(), experienceUpdate.LogoImage)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload experience logo",
				errors.WithContext("experience_id", experience.ID),
			)
		}
		experience.LogoUrl = logoURL
	}

	// Upload images if provided
	if len(experienceUpdate.Images) > 0 {
		imageURLs, err := s.uploadExperienceImages(ctx, experience.ID.String(), experienceUpdate.Images)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload experience images",
				errors.WithContext("experience_id", experience.ID),
			)
		}

		experience.ImagesUrl = imageURLs
	}

	// Update experience in repository
	updatedExperience, err := s.experienceRepo.Update(ctx, &experience)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to update experience",
			errors.WithContext("experience_id", experience.ID),
		)
	}

	// Create tech stack associations if provided
	if len(experienceUpdate.TechStackIds) > 0 {
		// First, delete existing tech stack associations
		err := s.experienceRepo.DeleteExperienceTechStack(ctx, updatedExperience.ID.String())
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to delete existing experience tech stack associations",
			)
		}

		// Then create new tech stack associations
		for _, techStackID := range experienceUpdate.TechStackIds {
			experienceTechStack := &ExperienceTechStack{
				ExperienceID: updatedExperience.ID,
				TechStackID:  techStackID,
			}

			_, err := s.experienceRepo.CreateExperienceTechStack(ctx, experienceTechStack)
			if err != nil {
				return nil, errors.Wrap(err,
					errors.ErrDatabase,
					"Failed to create experience tech stack association",
					errors.WithContext("experience_id", updatedExperience.ID),
					errors.WithContext("tech_stack_id", techStackID),
				)
			}
		}
	}

	return updatedExperience, nil
}

func (s *experienceService) DeleteExperience(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"Experience ID cannot be empty",
			nil,
		)
	}

	return s.experienceRepo.Delete(ctx, id)
}

func (s *experienceService) ListExperiences(ctx context.Context, opts base.ListOptions) ([]Experience, error) {
	experience, err := s.experienceRepo.List(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to list experience",
			errors.WithContext("options", opts),
		)
	}

	return experience, nil
}

func (s *experienceService) CountExperiences(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.experienceRepo.Count(ctx, filters)
}

func (s *experienceService) GetExperiencesByCompany(ctx context.Context, company string) ([]Experience, error) {
	return s.experienceRepo.FindByCompany(ctx, company)
}

func (s *experienceService) SearchExperiences(ctx context.Context, opts base.ListOptions) ([]Experience, int, error) {
	return s.experienceRepo.Search(ctx, opts)
}

func (s *experienceService) BulkCreateExperiences(ctx context.Context, experiencesCreate []*ExperienceCreate) ([]Experience, error) {
	experiences := make([]Experience, len(experiencesCreate))

	for i, experienceCreate := range experiencesCreate {
		createdExperience, err := s.CreateExperience(ctx, experienceCreate)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to create experience",
				errors.WithContext("experience_id", createdExperience.ID),
			)
		}
		experiences[i] = *createdExperience
	}

	return experiences, nil
}

func (s *experienceService) BulkUpdateExperiences(ctx context.Context, experiencesUpdate []*ExperienceUpdate) ([]Experience, error) {
	experiences := make([]Experience, len(experiencesUpdate))

	for i, experienceUpdate := range experiencesUpdate {

		updatedExperience, err := s.UpdateExperience(ctx, experienceUpdate)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to update experience",
				errors.WithContext("experience_id", experienceUpdate.ID),
			)
		}
		experiences[i] = *updatedExperience
	}

	return experiences, nil
}

func (s *experienceService) BulkDeleteExperiences(ctx context.Context, ids []string) error {
	for _, id := range ids {
		err := s.DeleteExperience(ctx, id)
		if err != nil {
			return errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to delete experience",
				errors.WithContext("experience_id", id),
			)
		}
	}
	return nil
}

func (s *experienceService) uploadExperienceLogo(ctx context.Context, experienceID string, file *multipart.FileHeader) (string, error) {
	if experienceID == "" {
		return "", fmt.Errorf("experience ID cannot be empty")
	}
	if file == nil {
		return "", fmt.Errorf("file data is required")
	}

	destPath := "images/experience/logos/" + experienceID + filepath.Ext(file.Filename)

	_, err := s.storage.Upload(ctx, file, destPath, storage_go.FileOptions{
		ContentType: func(s string) *string { return &s }("image"),
		Upsert:      func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		return "", errors.Wrap(err,
			errors.ErrInternal,
			"Failed to upload experience logo",
			errors.WithContext("experience_id", experienceID),
		)
	}

	signedURL, err := s.storage.GetPublicURL(destPath)
	if err != nil {
		return "", errors.Wrap(err,
			errors.ErrInternal,
			"Failed to get public URL for experience logo",
			errors.WithContext("dest_path", destPath),
		)
	}

	return signedURL, nil
}

func (s *experienceService) uploadExperienceImages(ctx context.Context, experienceID string, files []*multipart.FileHeader) ([]string, error) {
	if experienceID == "" {
		return nil, fmt.Errorf("experience ID cannot be empty")
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one file is required")
	}

	imageURLs := make([]string, len(files))
	for i, file := range files {
		destPath := fmt.Sprintf("images/experience/%s/%d%s", experienceID, i, filepath.Ext(file.Filename))

		_, err := s.storage.Upload(ctx, file, destPath, storage_go.FileOptions{
			ContentType: func(s string) *string { return &s }("image"),
			Upsert:      func(b bool) *bool { return &b }(true),
		})
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload experience image",
				errors.WithContext("experience_id", experienceID),
			)
		}

		signedURL, err := s.storage.GetPublicURL(destPath)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to get public URL for experience image",
				errors.WithContext("dest_path", destPath),
			)
		}

		imageURLs[i] = signedURL
	}

	return imageURLs, nil
}
