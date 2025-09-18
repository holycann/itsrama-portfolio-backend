package project

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/internal/base"
	"github.com/holycann/itsrama-portfolio-backend/internal/tech_stack"
	"github.com/holycann/itsrama-portfolio-backend/internal/validator"
	"github.com/holycann/itsrama-portfolio-backend/pkg/errors"
	"github.com/holycann/itsrama-portfolio-backend/pkg/supabase"
	storage_go "github.com/supabase-community/storage-go"
)

type ProjectService interface {
	CreateProject(ctx context.Context, projectCreate *ProjectCreate) (*ProjectDTO, error)
	GetProjectByID(ctx context.Context, id string) (*ProjectDTO, error)
	UpdateProject(ctx context.Context, projectUpdate *ProjectUpdate) (*ProjectDTO, error)
	DeleteProject(ctx context.Context, id string) error
	ListProjects(ctx context.Context, opts base.ListOptions) ([]ProjectDTO, error)
	CountProjects(ctx context.Context, filters []base.FilterOption) (int, error)
	SearchProjects(ctx context.Context, opts base.ListOptions) ([]ProjectDTO, int, error)
	BulkCreateProjects(ctx context.Context, projectsCreate []*ProjectCreate) ([]ProjectDTO, error)
	BulkUpdateProjects(ctx context.Context, projectsUpdate []*ProjectUpdate) ([]ProjectDTO, error)
	BulkDeleteProjects(ctx context.Context, ids []string) error
	uploadProjectImages(ctx context.Context, projectID string, files []*multipart.FileHeader) ([]string, error)
}

type projectService struct {
	projectRepo      ProjectRepository
	techStackService tech_stack.TechStackService
	storage          supabase.SupabaseStorage
}

func NewProjectService(projectRepo ProjectRepository, techStackService tech_stack.TechStackService, storage supabase.SupabaseStorage) ProjectService {
	return &projectService{
		projectRepo:      projectRepo,
		techStackService: techStackService,
		storage:          storage,
	}
}

func (s *projectService) CreateProject(ctx context.Context, projectCreate *ProjectCreate) (*ProjectDTO, error) {
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

		project.Images = imageData
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

	projectTechStack := make([]ProjectTechStackDTO, len(projectCreate.TechStackIds))
	for _, techStackID := range projectCreate.TechStackIds {
		techStackDTO, err := s.techStackService.GetTechStackByID(ctx, techStackID.String())
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to find tech stack",
				errors.WithContext("tech_stack_id", techStackID),
			)
		}
		projectTechStack = append(projectTechStack, ProjectTechStackDTO{
			ProjectID:   createdProject.ID,
			TechStackID: techStackID,
			TechStack:   *techStackDTO,
		})
	}

	createdProjectDTO := createdProject.ToDTO(projectTechStack)

	return &createdProjectDTO, nil
}

func (s *projectService) GetProjectByID(ctx context.Context, id string) (*ProjectDTO, error) {
	if id == "" {
		return nil, fmt.Errorf("project ID cannot be empty")
	}

	projects, err := s.projectRepo.FindByField(ctx, "id", id)
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("project not found")
	}

	return &projects[0], nil
}

func (s *projectService) UpdateProject(ctx context.Context, projectUpdate *ProjectUpdate) (*ProjectDTO, error) {
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

		project.Images = imageData
	} else {
		project.Images = existingProject.Images
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

	projectTechStack := make([]ProjectTechStackDTO, len(projectUpdate.TechStackIds))
	for _, techStackID := range projectUpdate.TechStackIds {
		techStackDTO, err := s.techStackService.GetTechStackByID(ctx, techStackID.String())
		if err != nil {
			return nil, errors.Wrap(err,
				errors.ErrDatabase,
				"Failed to find tech stack",
				errors.WithContext("tech_stack_id", techStackID),
			)
		}
		projectTechStack = append(projectTechStack, ProjectTechStackDTO{
			ProjectID:   updatedProject.ID,
			TechStackID: techStackID,
			TechStack:   *techStackDTO,
		})
	}

	updatedProjectDTO := updatedProject.ToDTO(projectTechStack)

	return &updatedProjectDTO, nil
}

func (s *projectService) DeleteProject(ctx context.Context, id string) error {
	if id == "" {
		return errors.New(
			errors.ErrValidation,
			"Project ID cannot be empty",
			nil,
		)
	}

	// Retrieve existing project to get image URL
	existingProject, err := s.GetProjectByID(ctx, id)
	if err != nil {
		return errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to retrieve existing project",
			errors.WithContext("project_id", id),
		)
	}

	// Delete project from repository
	err = s.projectRepo.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err,
			errors.ErrDatabase,
			"Failed to delete project",
			errors.WithContext("project_id", id),
		)
	}

	// Optional: Delete associated tech stack associations
	err = s.projectRepo.DeleteProjectTechStack(ctx, id)
	if err != nil {
		// Log the error but don't return it to avoid blocking the deletion
		fmt.Printf("Failed to delete project tech stack associations: %v\n", err)
	}

	// Delete associated images if exists
	for _, image := range existingProject.Images {
		if image.Src != "" {
			imagePath := filepath.Join(fmt.Sprintf("itsrama/images/project/%s/", existingProject.ID), filepath.Base(image.Src))
			_, err = s.storage.Delete(ctx, imagePath)
			if err != nil {
				// Log the error but don't return it to avoid blocking the deletion
				fmt.Printf("Failed to delete project image: %v\n", err)
			}
		}
	}

	return nil
}

func (s *projectService) ListProjects(ctx context.Context, opts base.ListOptions) ([]ProjectDTO, error) {
	// Validate list options
	if err := opts.Validate(); err != nil {
		return nil, errors.Wrap(err,
			errors.ErrValidation,
			"Invalid list options",
			errors.WithContext("options", opts),
		)
	}

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

func (s *projectService) SearchProjects(ctx context.Context, opts base.ListOptions) ([]ProjectDTO, int, error) {
	// Validate list options
	if err := opts.Validate(); err != nil {
		return nil, 0, errors.Wrap(err,
			errors.ErrValidation,
			"Invalid list options",
			errors.WithContext("options", opts),
		)
	}

	return s.projectRepo.Search(ctx, opts)
}

func (s *projectService) BulkCreateProjects(ctx context.Context, projectsCreate []*ProjectCreate) ([]ProjectDTO, error) {
	projects := make([]ProjectDTO, len(projectsCreate))

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

func (s *projectService) BulkUpdateProjects(ctx context.Context, projectsUpdate []*ProjectUpdate) ([]ProjectDTO, error) {
	projects := make([]ProjectDTO, len(projectsUpdate))

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
