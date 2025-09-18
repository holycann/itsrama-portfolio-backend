package project

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/internal/tech_stack"
)

// DevelopmentStatus represents the development stage of a project
// @Description Development status of a project
// @Name DevelopmentStatus
type DevelopmentStatus string

// ProgressStatus represents the current progress of a project
// @Description Current progress status of a project
// @Name ProgressStatus
type ProgressStatus string

// ProjectCategory represents the type of project
// @Description Category of the project
// @Name ProjectCategory
type ProjectCategory string

const (
	Alpha DevelopmentStatus = "Alpha"
	Beta  DevelopmentStatus = "Beta"
	MVP   DevelopmentStatus = "MVP"

	InProgress ProgressStatus = "In Progress"
	InRevision ProgressStatus = "In Revision"
	OnHold     ProgressStatus = "On Hold"
	Completed  ProgressStatus = "Completed"

	WebDevelopment ProjectCategory = "Web Development"
	ApiDevelopment ProjectCategory = "API Development"
	BotDevelopment ProjectCategory = "Bot Development"
	MobileApp      ProjectCategory = "Mobile App"
	DesktopApp     ProjectCategory = "Desktop App"
	UIUX           ProjectCategory = "UI/UX Design"
	Other          ProjectCategory = "Other"
)

// ProjectImage represents an image associated with a project
// @Description Image details for a project
// @Name ProjectImage
type ProjectImage struct {
	Src         string `json:"src" example:"https://example.com/image.jpg"`
	Alt         string `json:"alt" example:"Project screenshot"`
	IsThumbnail bool   `json:"is_thumbnail" example:"false"`
}

// Project represents the main project model
// @Description Detailed information about a project
// @Name Project
type Project struct {
	// Identification
	ID   uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Slug string    `json:"slug" db:"slug" validate:"required" example:"portfolio-website"`

	// Project Details
	Title       string          `json:"title" db:"title" validate:"required" example:"Portfolio Website"`
	Subtitle    string          `json:"subtitle" db:"subtitle" example:"Personal portfolio showcasing projects"`
	Description string          `json:"description" db:"description" validate:"required" example:"A responsive website to display my professional projects and skills"`
	MyRole      []string        `json:"my_role" db:"my_role" example:"Full-stack Developer,UI/UX Designer"`
	Category    ProjectCategory `json:"category" db:"category" example:"Web Development"`

	// URLs
	GithubUrl string `json:"github_url,omitempty" db:"github_url" example:"https://github.com/username/project"`
	WebUrl    string `json:"web_url,omitempty" db:"web_url" example:"https://myportfolio.com"`

	// Project Content
	Images   []ProjectImage `json:"images" db:"images" pg:"array"`
	Features []string       `json:"features" db:"features" pg:"array" example:"Responsive Design,Dark Mode"`

	// Status
	DevelopmentStatus  DevelopmentStatus `json:"development_status" db:"development_status" example:"Beta"`
	ProgressStatus     ProgressStatus    `json:"progress_status" db:"progress_status" example:"In Progress"`
	ProgressPercentage int               `json:"progress_percentage" db:"progress_percentage" example:"75"`
	IsFeatured         bool              `json:"is_featured" db:"is_featured" example:"true"`

	// Metadata
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// ProjectTechStack represents the relationship between a project and its tech stack
// @Description Tech stack associated with a project
// @Name ProjectTechStack
type ProjectTechStack struct {
	ProjectID   uuid.UUID `json:"project_id" db:"project_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	TechStackID uuid.UUID `json:"tech_stack_id" db:"tech_stack_id" validate:"required" example:"650f9500-f39c-52d5-b827-557766550001"`
}

// ProjectTechStackDTO is a Data Transfer Object for ProjectTechStack
// @Description Data transfer object for project tech stack with additional tech stack details
// @Name ProjectTechStackDTO
type ProjectTechStackDTO struct {
	ProjectID   uuid.UUID            `json:"project_id" db:"project_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	TechStackID uuid.UUID            `json:"tech_stack_id" db:"tech_stack_id" validate:"required" example:"650f9500-f39c-52d5-b827-557766550001"`
	TechStack   tech_stack.TechStack `json:"tech_stack" db:"tech_stack" pg:"array"`
}

// ProjectDTO is a Data Transfer Object for Project
// @Description Data transfer object for project with additional tech stack information
// @Name ProjectDTO
type ProjectDTO struct {
	// Identification
	ID   uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Slug string    `json:"slug" db:"slug" validate:"required" example:"portfolio-website"`

	// Project Details
	Title       string          `json:"title" db:"title" validate:"required" example:"Portfolio Website"`
	Subtitle    string          `json:"subtitle" db:"subtitle" example:"Personal portfolio showcasing projects"`
	Description string          `json:"description" db:"description" validate:"required" example:"A responsive website to display my professional projects and skills"`
	MyRole      []string        `json:"my_role" db:"my_role" example:"Full-stack Developer,UI/UX Designer"`
	Category    ProjectCategory `json:"category" db:"category" example:"Web Development"`

	// URLs
	GithubUrl string `json:"github_url,omitempty" db:"github_url" example:"https://github.com/username/project"`
	WebUrl    string `json:"web_url,omitempty" db:"web_url" example:"https://myportfolio.com"`

	// Project Content
	Images   []ProjectImage `json:"images" db:"images" pg:"array"`
	Features []string       `json:"features" db:"features" pg:"array" example:"Responsive Design,Dark Mode"`

	// Status
	DevelopmentStatus  DevelopmentStatus `json:"development_status" db:"development_status" example:"Beta"`
	ProgressStatus     ProgressStatus    `json:"progress_status" db:"progress_status" example:"In Progress"`
	ProgressPercentage int               `json:"progress_percentage" db:"progress_percentage" example:"75"`
	IsFeatured         bool              `json:"is_featured" db:"is_featured" example:"true"`

	// Metadata
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`

	// Relationships
	ProjectTechStack []ProjectTechStackDTO `json:"project_tech_stack" db:"project_tech_stack" pg:"array"`
}

// ProjectCreate represents the input for creating a new project
// @Description Input model for creating a new project
// @Name ProjectCreate
type ProjectCreate struct {
	Slug         string          `json:"slug" db:"slug" validate:"required" example:"portfolio-website"`
	TechStackIds []uuid.UUID     `json:"tech_stack_ids" validate:"required" example:"[\"550e8400-e29b-41d4-a716-446655440000\"]"`
	Title        string          `json:"title" validate:"required" example:"Portfolio Website"`
	Subtitle     string          `json:"subtitle" example:"Personal portfolio showcasing projects"`
	Description  string          `json:"description" validate:"required" example:"A responsive website to display my professional projects and skills"`
	MyRole       []string        `json:"my_role" example:"Full-stack Developer,UI/UX Designer"`
	Category     ProjectCategory `json:"category" example:"Web Development"`

	GithubUrl string `json:"github_url,omitempty" example:"https://github.com/username/project"`
	WebUrl    string `json:"web_url,omitempty" example:"https://myportfolio.com"`

	Features []string `json:"features" example:"Responsive Design,Dark Mode"`

	DevelopmentStatus  DevelopmentStatus `json:"development_status" example:"Beta"`
	ProgressStatus     ProgressStatus    `json:"progress_status" example:"In Progress"`
	ProgressPercentage int               `json:"progress_percentage" example:"75"`
	IsFeatured         bool              `json:"is_featured" example:"true"`

	UploadedImages []*multipart.FileHeader `json:"uploaded_images" swaggerignore:"true"`
}

// ProjectUpdate represents the input for updating an existing project
// @Description Input model for updating an existing project
// @Name ProjectUpdate
type ProjectUpdate struct {
	ID           uuid.UUID   `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Slug         string      `json:"slug" db:"slug" validate:"required" example:"portfolio-website"`
	TechStackIds []uuid.UUID `json:"tech_stack_ids" validate:"required" example:"[\"550e8400-e29b-41d4-a716-446655440000\"]"`

	Title       string          `json:"title" example:"Updated Portfolio Website"`
	Subtitle    string          `json:"subtitle" example:"Updated personal portfolio showcasing projects"`
	Description string          `json:"description" example:"An improved responsive website to display my professional projects and skills"`
	MyRole      []string        `json:"my_role" example:"Full-stack Developer,DevOps Engineer"`
	Category    ProjectCategory `json:"category" example:"Web Development"`

	GithubUrl string `json:"github_url,omitempty" example:"https://github.com/username/updated-project"`
	WebUrl    string `json:"web_url,omitempty" example:"https://updated-myportfolio.com"`

	Features []string `json:"features" example:"Responsive Design,Dark Mode,Performance Optimization"`

	DevelopmentStatus  DevelopmentStatus `json:"development_status" example:"Beta"`
	ProgressStatus     ProgressStatus    `json:"progress_status" example:"Completed"`
	ProgressPercentage int               `json:"progress_percentage" example:"100"`
	IsFeatured         bool              `json:"is_featured" example:"true"`

	UploadedImages []*multipart.FileHeader `json:"uploaded_images" swaggerignore:"true"`
}

// ToProject converts ProjectCreate to Project
func (pc *ProjectCreate) ToProject() Project {
	now := time.Now().UTC()
	return Project{
		ID:                 uuid.New(),
		Slug:               pc.Slug,
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
		IsFeatured:         pc.IsFeatured,
		Images:             nil, // Will be set during file upload
		CreatedAt:          &now,
		UpdatedAt:          &now,
	}
}

// ToProject converts ProjectUpdate to Project
func (pu *ProjectUpdate) ToProject() Project {
	now := time.Now().UTC()
	return Project{
		ID:                 pu.ID,
		Slug:               pu.Slug,
		Title:              pu.Title,
		Subtitle:           pu.Subtitle,
		Description:        pu.Description,
		MyRole:             pu.MyRole,
		Category:           pu.Category,
		GithubUrl:          pu.GithubUrl,
		WebUrl:             pu.WebUrl,
		IsFeatured:         pu.IsFeatured,
		Images:             nil, // Will be set during file upload
		Features:           pu.Features,
		DevelopmentStatus:  pu.DevelopmentStatus,
		ProgressStatus:     pu.ProgressStatus,
		ProgressPercentage: pu.ProgressPercentage,
		UpdatedAt:          &now,
	}
}

// ToDTO converts a Project to a ProjectDTO
func (p *Project) ToDTO(projectTechStack []ProjectTechStackDTO) ProjectDTO {
	return ProjectDTO{
		ID:                 p.ID,
		Slug:               p.Slug,
		Title:              p.Title,
		Subtitle:           p.Subtitle,
		Description:        p.Description,
		MyRole:             p.MyRole,
		Category:           p.Category,
		GithubUrl:          p.GithubUrl,
		WebUrl:             p.WebUrl,
		IsFeatured:         p.IsFeatured,
		Images:             p.Images,
		Features:           p.Features,
		DevelopmentStatus:  p.DevelopmentStatus,
		ProgressStatus:     p.ProgressStatus,
		ProgressPercentage: p.ProgressPercentage,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
		ProjectTechStack:   projectTechStack,
	}
}
