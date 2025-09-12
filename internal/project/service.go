package project

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

type ProjectService interface {
	CreateProject(ctx context.Context, projectCreate *ProjectCreate) (*Project, error)
	GetProjectByID(ctx context.Context, id string) (*Project, error)
	UpdateProject(ctx context.Context, projectUpdate *ProjectUpdate) (*Project, error)
	DeleteProject(ctx context.Context, id string) error
	ListProjects(ctx context.Context, opts base.ListOptions) ([]Project, error)
	CountProjects(ctx context.Context, filters []base.FilterOption) (int, error)
	GetProjectsByCategory(ctx context.Context, category string) ([]Project, error)
	SearchProjects(ctx context.Context, opts base.ListOptions) ([]Project, int, error)
	BulkCreateProjects(ctx context.Context, projectsCreate []*ProjectCreate) ([]Project, error)
	BulkUpdateProjects(ctx context.Context, projectsUpdate []*ProjectUpdate) ([]Project, error)
	BulkDeleteProjects(ctx context.Context, ids []string) error
	uploadProjectImages(ctx context.Context, projectID string, files []*multipart.FileHeader) ([]string, error)
}

type projectService struct {
	projectRepo ProjectRepository
	storage     supabase.SupabaseStorage
}

func NewProjectService(projectRepo ProjectRepository, storage supabase.SupabaseStorage) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		storage:     storage,
	}
}

func (s *projectService) CreateProject(ctx context.Context, projectCreate *ProjectCreate) (*Project, error) {
	// Validate input
	if err := validator.ValidateModel(projectCreate); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	project := projectCreate.ToProject()
	project.ID = uuid.New()
	project.CreatedAt = &now
	project.UpdatedAt = &now

	// Upload images if provided
	if len(projectCreate.UploadedImages) > 0 {
		imageURLs, err := s.uploadProjectImages(ctx, project.ID.String(), projectCreate.UploadedImages)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload project images",
				errors.WithContext("project_id", project.ID),
			)
		}

		imageData := make([]ProjectImage, len(imageURLs))
		for i, url := range imageURLs {
			imageData[i] = ProjectImage{
				Src:         url,
				Alt:         filepath.Base(url),
				IsThumbnail: i == 0,
			}
		}

		project.ImagesSrc = imageData
	}

	// Create project in repository
	createdProject, err := s.projectRepo.Create(ctx, &project)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to create project",
			errors.WithContext("project_title", project.Title),
		)
	}

	// Create project tech stack if provided
	if len(projectCreate.TechStackIds) > 0 {
		for _, techStackID := range projectCreate.TechStackIds {
			projectTechStack := &ProjectTechStack{
				ProjectID:   project.ID,
				TechStackID: techStackID,
			}

			_, err := s.projectRepo.CreateProjectTechStack(ctx, projectTechStack)
			if err != nil {
				return nil, errors.Wrap(err,
					errors.ErrDatabase,
					"Failed to create project tech stack",
					errors.WithContext("project_id", project.ID),
					errors.WithContext("tech_stack_id", techStackID),
				)
			}
		}
	}

	return createdProject, nil
}

func (s *projectService) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	if id == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}

	return s.projectRepo.FindByID(ctx, id)
}

func (s *projectService) UpdateProject(ctx context.Context, projectUpdate *ProjectUpdate) (*Project, error) {
	// Validate input
	if err := validator.ValidateModel(projectUpdate); err != nil {
		return nil, errors.New(
			errors.ErrValidation,
			"Invalid project payload",
			err,
			errors.WithContext("payload", projectUpdate),
		)
	}

	// Retrieve existing project
	existingProject, err := s.GetProjectByID(ctx, projectUpdate.ID.String())
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to retrieve existing project",
			errors.WithContext("project_id", projectUpdate.ID),
		)
	}

	now := time.Now().UTC()
	project := projectUpdate.ToProject()
	project.CreatedAt = existingProject.CreatedAt
	project.UpdatedAt = &now

	// Upload images if provided
	if len(projectUpdate.UploadedImages) > 0 {
		imageURLs, err := s.uploadProjectImages(ctx, project.ID.String(), projectUpdate.UploadedImages)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload project images",
				errors.WithContext("project_id", project.ID),
			)
		}

		imageData := make([]ProjectImage, len(imageURLs))
		for i, url := range imageURLs {
			imageData[i] = ProjectImage{
				Src:         url,
				Alt:         filepath.Base(url),
				IsThumbnail: i == 0,
			}
		}

		project.ImagesSrc = imageData
	} else {
		project.ImagesSrc = existingProject.ImagesSrc
	}

	// Conditionally update fields
	if !validator.IsValueChanged(&existingProject.Title, &project.Title) {
		project.Title = existingProject.Title
	}
	if !validator.IsValueChanged(&existingProject.Subtitle, &project.Subtitle) {
		project.Subtitle = existingProject.Subtitle
	}
	if !validator.IsValueChanged(&existingProject.Description, &project.Description) {
		project.Description = existingProject.Description
	}
	if !validator.IsValueChanged(&existingProject.MyRole, &project.MyRole) {
		project.MyRole = existingProject.MyRole
	}
	if !validator.IsValueChanged(&existingProject.Category, &project.Category) {
		project.Category = existingProject.Category
	}
	if !validator.IsValueChanged(&existingProject.GithubUrl, &project.GithubUrl) {
		project.GithubUrl = existingProject.GithubUrl
	}
	if !validator.IsValueChanged(&existingProject.WebUrl, &project.WebUrl) {
		project.WebUrl = existingProject.WebUrl
	}
	if !validator.IsValueChanged(&existingProject.Features, &project.Features) {
		project.Features = existingProject.Features
	}
	if !validator.IsValueChanged(&existingProject.DevelopmentStatus, &project.DevelopmentStatus) {
		project.DevelopmentStatus = existingProject.DevelopmentStatus
	}
	if !validator.IsValueChanged(&existingProject.ProgressStatus, &project.ProgressStatus) {
		project.ProgressStatus = existingProject.ProgressStatus
	}
	if !validator.IsValueChanged(&existingProject.ProgressPercentage, &project.ProgressPercentage) {
		project.ProgressPercentage = existingProject.ProgressPercentage
	}

	// Update project in repository
	updatedProject, err := s.projectRepo.Update(ctx, &project)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to update project",
			errors.WithContext("project_id", project.ID),
		)
	}

	// Create project tech stack if provided
	if len(projectUpdate.TechStackIds) > 0 {
		// Delete existing project tech stack
		err = s.projectRepo.DeleteProjectTechStack(ctx, updatedProject.ID.String())
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to delete existing project tech stack",
				errors.WithContext("project_id", updatedProject.ID),
			)
		}

		// Create new project tech stack entries
		for _, techStackID := range projectUpdate.TechStackIds {
			projectTechStack := &ProjectTechStack{
				ProjectID:   updatedProject.ID,
				TechStackID: techStackID,
			}
			_, err := s.projectRepo.CreateProjectTechStack(ctx, projectTechStack)
			if err != nil {
				return nil, errors.Wrap(err,
					errors.ErrDatabase,
					"Failed to create project tech stack",
					errors.WithContext("project_id", updatedProject.ID),
					errors.WithContext("tech_stack_id", techStackID),
				)
			}
		}
	}

	return updatedProject, nil
}

func (s *projectService) DeleteProject(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"Project ID cannot be empty",
			nil,
		)
	}

	return s.projectRepo.Delete(ctx, id)
}

func (s *projectService) ListProjects(ctx context.Context, opts base.ListOptions) ([]Project, error) {
	projects, err := s.projectRepo.List(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to list projects",
			errors.WithContext("options", opts),
		)
	}

	return projects, nil
}

func (s *projectService) CountProjects(ctx context.Context, filters []base.FilterOption) (int, error) {
	return s.projectRepo.Count(ctx, filters)
}

func (s *projectService) GetProjectsByCategory(ctx context.Context, category string) ([]Project, error) {
	if category == "" {
		return nil, fmt.Errorf("category cannot be empty")
	}

	return s.projectRepo.FindByCategory(ctx, category)
}

func (s *projectService) SearchProjects(ctx context.Context, opts base.ListOptions) ([]Project, int, error) {
	return s.projectRepo.Search(ctx, opts)
}

func (s *projectService) BulkCreateProjects(ctx context.Context, projectsCreate []*ProjectCreate) ([]Project, error) {
	projects := make([]Project, len(projectsCreate))

	for i, projectCreate := range projectsCreate {
		createdProject, err := s.CreateProject(ctx, projectCreate)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to create project",
				errors.WithContext("project_id", createdProject.ID),
			)
		}
		projects[i] = *createdProject
	}

	return projects, nil
}

func (s *projectService) BulkUpdateProjects(ctx context.Context, projectsUpdate []*ProjectUpdate) ([]Project, error) {
	projects := make([]Project, len(projectsUpdate))

	for i, projectUpdate := range projectsUpdate {
		updatedProject, err := s.UpdateProject(ctx, projectUpdate)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to update project",
				errors.WithContext("project_id", projectUpdate.ID),
			)
		}
		projects[i] = *updatedProject
	}

	return projects, nil
}

func (s *projectService) BulkDeleteProjects(ctx context.Context, ids []string) error {
	for _, id := range ids {
		err := s.DeleteProject(ctx, id)
		if err != nil {
			return errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to delete project",
				errors.WithContext("project_id", id),
			)
		}
	}
	return nil
}

func (s *projectService) uploadProjectImages(ctx context.Context, projectID string, files []*multipart.FileHeader) ([]string, error) {
	if projectID == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one file is required")
	}

	imageURLs := make([]string, len(files))
	for i, file := range files {
		destPath := fmt.Sprintf("images/project/%s/%d%s", projectID, i, filepath.Ext(file.Filename))

		_, err := s.storage.Upload(ctx, file, destPath, storage_go.FileOptions{
			ContentType: func(s string) *string { return &s }("image"),
			Upsert:      func(b bool) *bool { return &b }(true),
		})
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to upload project image",
				errors.WithContext("project_id", projectID),
			)
		}

		signedURL, err := s.storage.GetPublicURL(destPath)
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrInternal,
				"Failed to get public URL for project image",
				errors.WithContext("dest_path", destPath),
			)
		}

		imageURLs[i] = signedURL
	}

	return imageURLs, nil
}
