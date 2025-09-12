package project

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/internal/tech_stack"
)

type DevelopmentStatus string
type ProgressStatus string

const (
	Alpha DevelopmentStatus = "Alpha"
	Beta  DevelopmentStatus = "Beta"
	MVP   DevelopmentStatus = "MVP"

	InProgress ProgressStatus = "In Progress"
	InRevision ProgressStatus = "In Revision"
	OnHold     ProgressStatus = "On Hold"
	Completed  ProgressStatus = "Completed"
)

type ProjectImage struct {
	Src         string `json:"src"`
	Alt         string `json:"alt"`
	IsThumbnail bool   `json:"is_thumbnail"`
}

type Project struct {
	// Identification
	ID uuid.UUID `json:"id" db:"id"`

	// Project Details
	Title       string   `json:"title" db:"title" validate:"required"`
	Subtitle    string   `json:"subtitle" db:"subtitle"`
	Description string   `json:"description" db:"description" validate:"required"`
	MyRole      []string `json:"my_role" db:"my_role"`
	Category    string   `json:"category" db:"category"`

	// URLs
	GithubUrl string `json:"github_url,omitempty" db:"github_url"`
	WebUrl    string `json:"web_url,omitempty" db:"web_url"`

	// Project Content
	ImagesSrc []ProjectImage `json:"images" db:"images" pg:"array"`
	Features  []string       `json:"features" db:"features" pg:"array"`

	// Status
	DevelopmentStatus  DevelopmentStatus `json:"development_status" db:"development_status"`
	ProgressStatus     ProgressStatus    `json:"progress_status" db:"progress_status"`
	ProgressPercentage int               `json:"progress_percentage" db:"progress_percentage"`

	// Metadata
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type ProjectTechStack struct {
	ProjectID   uuid.UUID `json:"project_id" db:"project_id" validate:"required"`
	TechStackID uuid.UUID `json:"tech_stack_id" db:"tech_stack_id" validate:"required"`
}

type ProjectDTO struct {
	// Identification
	ID uuid.UUID `json:"id" db:"id"`

	// Project Details
	Title       string   `json:"title" db:"title" validate:"required"`
	Subtitle    string   `json:"subtitle" db:"subtitle"`
	Description string   `json:"description" db:"description" validate:"required"`
	MyRole      []string `json:"my_role" db:"my_role"`
	Category    string   `json:"category" db:"category"`

	// URLs
	GithubUrl string `json:"github_url,omitempty" db:"github_url"`
	WebUrl    string `json:"web_url,omitempty" db:"web_url"`

	// Project Content
	ImagesSrc []ProjectImage `json:"images" db:"images" pg:"array"`
	Features  []string       `json:"features" db:"features" pg:"array"`

	// Status
	DevelopmentStatus  DevelopmentStatus `json:"development_status" db:"development_status"`
	ProgressStatus     ProgressStatus    `json:"progress_status" db:"progress_status"`
	ProgressPercentage int               `json:"progress_percentage" db:"progress_percentage"`

	// Metadata
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`

	// Relationships
	TechStack []tech_stack.TechStack `json:"tech_stack" db:"tech_stack" pg:"array"`
}

type ProjectCreate struct {
	TechStackIds []uuid.UUID `json:"tech_stack_ids" validate:"required"`
	Title        string      `json:"title" validate:"required"`
	Subtitle     string      `json:"subtitle"`
	Description  string      `json:"description" validate:"required"`
	MyRole       []string    `json:"my_role"`
	Category     string      `json:"category"`

	GithubUrl string `json:"github_url,omitempty"`
	WebUrl    string `json:"web_url,omitempty"`

	Features []string `json:"features"`

	DevelopmentStatus  DevelopmentStatus `json:"development_status"`
	ProgressStatus     ProgressStatus    `json:"progress_status"`
	ProgressPercentage int               `json:"progress_percentage"`

	UploadedImages []*multipart.FileHeader `json:"uploaded_images"`
}

type ProjectUpdate struct {
	ID           uuid.UUID   `json:"id" validate:"required"`
	TechStackIds []uuid.UUID `json:"tech_stack_ids" validate:"required"`

	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle"`
	Description string   `json:"description"`
	MyRole      []string `json:"my_role"`
	Category    string   `json:"category"`

	GithubUrl string `json:"github_url,omitempty"`
	WebUrl    string `json:"web_url,omitempty"`

	Features []string `json:"features"`

	DevelopmentStatus  DevelopmentStatus `json:"development_status"`
	ProgressStatus     ProgressStatus    `json:"progress_status"`
	ProgressPercentage int               `json:"progress_percentage"`

	UploadedImages []*multipart.FileHeader `json:"uploaded_images"`
}

// ToProject converts ProjectCreate to Project
func (pc *ProjectCreate) ToProject() Project {
	now := time.Now().UTC()
	return Project{
		ID:                 uuid.New(),
		Title:              pc.Title,
		Subtitle:           pc.Subtitle,
		Description:        pc.Description,
		MyRole:             pc.MyRole,
		Category:           pc.Category,
		GithubUrl:          pc.GithubUrl,
		WebUrl:             pc.WebUrl,
		Features:           pc.Features,
		DevelopmentStatus:  pc.DevelopmentStatus,
		ProgressStatus:     pc.ProgressStatus,
		ProgressPercentage: pc.ProgressPercentage,
		ImagesSrc:          nil, // Will be set during file upload
		CreatedAt:          &now,
		UpdatedAt:          &now,
	}
}

// ToProject converts ProjectUpdate to Project
func (pu *ProjectUpdate) ToProject() Project {
	now := time.Now().UTC()
	return Project{
		ID:                 pu.ID,
		Title:              pu.Title,
		Subtitle:           pu.Subtitle,
		Description:        pu.Description,
		MyRole:             pu.MyRole,
		Category:           pu.Category,
		GithubUrl:          pu.GithubUrl,
		WebUrl:             pu.WebUrl,
		ImagesSrc:          nil, // Will be set during file upload
		Features:           pu.Features,
		DevelopmentStatus:  pu.DevelopmentStatus,
		ProgressStatus:     pu.ProgressStatus,
		ProgressPercentage: pu.ProgressPercentage,
		UpdatedAt:          &now,
	}
}

// ToDTO converts a Project to a ProjectDTO
func (p *Project) ToDTO(techStacks []tech_stack.TechStack) ProjectDTO {
	return ProjectDTO{
		ID:                 p.ID,
		Title:              p.Title,
		Subtitle:           p.Subtitle,
		Description:        p.Description,
		MyRole:             p.MyRole,
		Category:           p.Category,
		GithubUrl:          p.GithubUrl,
		WebUrl:             p.WebUrl,
		ImagesSrc:          p.ImagesSrc,
		Features:           p.Features,
		DevelopmentStatus:  p.DevelopmentStatus,
		ProgressStatus:     p.ProgressStatus,
		ProgressPercentage: p.ProgressPercentage,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
		TechStack:          techStacks,
	}
}
