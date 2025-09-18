package experience

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/internal/tech_stack"
	"github.com/holycann/itsrama-portfolio-backend/internal/utils"
)

// Experience represents a professional work experience
//
// @Description Detailed information about a professional work experience
// @Name Experience
type Experience struct {
	// Identification
	// @Description Unique identifier for the experience
	ID uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Job Details
	// @Description Job role and company information
	Role    string `json:"role" db:"role" validate:"required" example:"Senior Software Engineer"`
	Company string `json:"company" db:"company" validate:"required" example:"Tech Innovations Inc."`
	LogoUrl string `json:"logo_url" db:"logo_url" example:"https://example.com/company-logo.png"`
	JobType string `json:"job_type" db:"job_type" example:"Full-time"`

	// Timing and Location
	// @Description Job timing and location details
	StartDate   utils.CustomDate  `json:"start_date" db:"start_date" validate:"required" example:"2020-01-15" swaggertype:"string"`
	EndDate     *utils.CustomDate `json:"end_date" db:"end_date" example:"2023-06-30" swaggertype:"string"`
	Location    string            `json:"location" db:"location" example:"San Francisco, CA"`
	Arrangement string            `json:"arrangement" db:"arrangement" example:"Remote"`

	// Job Description
	// @Description Detailed description of work and achievements
	WorkDescription string   `json:"work_description" db:"work_description" example:"Led development of scalable web applications"`
	Impact          []string `json:"impact" db:"impact" pg:"array" example:"Increased system performance by 40%"`
	ImagesUrl       []string `json:"images_url" db:"images_url" pg:"array" example:"https://example.com/project1.png"`

	// Metadata
	// @Description Additional metadata for the experience
	IsFeatured bool       `json:"is_featured" db:"is_featured" example:"true"`
	CreatedAt  *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// ExperienceTechStack represents the association between an experience and tech stack
//
// @Description Association between an experience and its related technologies
// @Name ExperienceTechStack
type ExperienceTechStack struct {
	// @Description ID of the associated experience
	ExperienceID uuid.UUID `json:"experience_id" db:"experience_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`

	// @Description ID of the associated tech stack
	TechStackID uuid.UUID `json:"tech_stack_id" db:"tech_stack_id" validate:"required" example:"650f9500-f39c-52d5-b827-557766550001"`
}

// ExperienceTechStackDTO represents a data transfer object for experience tech stack
//
// @Description Data transfer object for experience and tech stack association
// @Name ExperienceTechStackDTO
type ExperienceTechStackDTO struct {
	// @Description ID of the associated experience
	ExperienceID uuid.UUID `json:"experience_id" db:"experience_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`

	// @Description ID of the associated tech stack
	TechStackID uuid.UUID `json:"tech_stack_id" db:"tech_stack_id" validate:"required" example:"650f9500-f39c-52d5-b827-557766550001"`

	// @Description Detailed information about the tech stack
	TechStack tech_stack.TechStack `json:"tech_stack" db:"tech_stack" pg:"array"`
}

// ExperienceDTO represents a data transfer object for experience
//
// @Description Data transfer object for experience with additional details
// @Name ExperienceDTO
type ExperienceDTO struct {
	// Identification
	// @Description Unique identifier for the experience
	ID uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// Job Details
	// @Description Job role and company information
	Role    string `json:"role" db:"role" validate:"required" example:"Senior Software Engineer"`
	Company string `json:"company" db:"company" validate:"required" example:"Tech Innovations Inc."`
	LogoUrl string `json:"logo_url" db:"logo_url" example:"https://example.com/company-logo.png"`
	JobType string `json:"job_type" db:"job_type" example:"Full-time"`

	// Timing and Location
	// @Description Job timing and location details
	StartDate   utils.CustomDate  `json:"start_date" db:"start_date" validate:"required" example:"2020-01-15" swaggertype:"string"`
	EndDate     *utils.CustomDate `json:"end_date" db:"end_date" example:"2023-06-30" swaggertype:"string"`
	Location    string            `json:"location" db:"location" example:"San Francisco, CA"`
	Arrangement string            `json:"arrangement" db:"arrangement" example:"Remote"`

	// Job Description
	// @Description Detailed description of work and achievements
	WorkDescription string   `json:"work_description" db:"work_description" example:"Led development of scalable web applications"`
	Impact          []string `json:"impact" db:"impact" pg:"array" example:"Increased system performance by 40%"`
	ImagesUrl       []string `json:"images_url" db:"images_url" pg:"array" example:"https://example.com/project1.png"`

	// Metadata
	// @Description Additional metadata for the experience
	IsFeatured bool       `json:"is_featured" db:"is_featured" example:"true"`
	CreatedAt  *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty" db:"updated_at"`

	// Relationships
	// @Description Associated tech stacks for the experience
	ExperienceTechStack []ExperienceTechStackDTO `json:"experience_tech_stack" db:"experience_tech_stack" pg:"array"`
}

// ExperienceCreate represents the input for creating a new experience
//
// @Description Input model for creating a new experience
// @Name ExperienceCreate
type ExperienceCreate struct {
	// @Description IDs of associated tech stacks
	// @Enums ["550e8400-e29b-41d4-a716-446655440000", "650f9500-f39c-52d5-b827-557766550001"]
	TechStackIds []uuid.UUID `json:"tech_stack_ids"`

	// @Description Job role
	// @Format string
	Role string `json:"role" example:"Senior Software Engineer"`

	// @Description Company name
	// @Format string
	Company string `json:"company" example:"Tech Innovations Inc."`

	// @Description Logo image file
	LogoImage *multipart.FileHeader `json:"logo_image" swaggerignore:"true"`

	// @Description Job type
	// @Enums ["Full-time", "Part-time", "Contract", "Freelance"]
	JobType string `json:"job_type" example:"Full-time"`

	// @Description Start date of the job
	// @Format date
	StartDate utils.CustomDate `json:"start_date" example:"2020-01-15" swaggertype:"string"`

	// @Description End date of the job (optional)
	// @Format date
	EndDate *utils.CustomDate `json:"end_date" example:"2023-06-30" swaggertype:"string"`

	// @Description Job location
	// @Format string
	Location string `json:"location" example:"San Francisco, CA"`

	// @Description Work arrangement
	// @Enums ["Remote", "Hybrid", "On-site"]
	Arrangement string `json:"arrangement" example:"Remote"`

	// @Description Work description
	// @Format string
	WorkDescription string `json:"work_description" example:"Led development of scalable web applications"`

	// @Description Key impacts or achievements
	Impact []string `json:"impact" example:"Increased system performance by 40%"`

	// @Description Job-related images
	Images []*multipart.FileHeader `json:"images" swaggerignore:"true"`

	// @Description Flag to mark as featured experience
	IsFeatured bool `json:"is_featured" example:"true"`
}

// ExperienceUpdate represents the input for updating an existing experience
//
// @Description Input model for updating an existing experience
// @Name ExperienceUpdate
type ExperienceUpdate struct {
	// @Description Unique identifier of the experience to update
	ID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`

	// @Description IDs of associated tech stacks
	// @Enums ["550e8400-e29b-41d4-a716-446655440000", "650f9500-f39c-52d5-b827-557766550001"]
	TechStackIds []uuid.UUID `json:"tech_stack_ids"`

	// @Description Job role
	// @Format string
	Role string `json:"role" example:"Senior Software Engineer"`

	// @Description Company name
	// @Format string
	Company string `json:"company" example:"Tech Innovations Inc."`

	// @Description Logo image file
	LogoImage *multipart.FileHeader `json:"logo_image" swaggerignore:"true"`

	// @Description Job type
	// @Enums ["Full-time", "Part-time", "Contract", "Freelance"]
	JobType string `json:"job_type" example:"Full-time"`

	// @Description Start date of the job
	// @Format date
	StartDate utils.CustomDate `json:"start_date" example:"2020-01-15" swaggertype:"string"`

	// @Description End date of the job (optional)
	// @Format date
	EndDate *utils.CustomDate `json:"end_date" example:"2023-06-30" swaggertype:"string"`

	// @Description Job location
	// @Format string
	Location string `json:"location" example:"San Francisco, CA"`

	// @Description Work arrangement
	// @Enums ["Remote", "Hybrid", "On-site"]
	Arrangement string `json:"arrangement" example:"Remote"`

	// @Description Work description
	// @Format string
	WorkDescription string `json:"work_description" example:"Led development of scalable web applications"`

	// @Description Key impacts or achievements
	Impact []string `json:"impact" example:"Increased system performance by 40%"`

	// @Description Job-related images
	Images []*multipart.FileHeader `json:"images" swaggerignore:"true"`

	// @Description Flag to mark as featured experience
	IsFeatured bool `json:"is_featured" example:"true"`
}

// ToDTO converts an Experience to an ExperienceDTO
//
// @Description Converts Experience model to ExperienceDTO
func (e *Experience) ToDTO(experienceTechStack []ExperienceTechStackDTO) ExperienceDTO {
	return ExperienceDTO{
		ID:                  e.ID,
		Role:                e.Role,
		Company:             e.Company,
		LogoUrl:             e.LogoUrl,
		JobType:             e.JobType,
		StartDate:           e.StartDate,
		EndDate:             e.EndDate,
		Location:            e.Location,
		Arrangement:         e.Arrangement,
		WorkDescription:     e.WorkDescription,
		Impact:              e.Impact,
		ImagesUrl:           e.ImagesUrl,
		IsFeatured:          e.IsFeatured,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
		ExperienceTechStack: experienceTechStack,
	}
}

// ToExperience converts ExperienceCreate to Experience
//
// @Description Converts ExperienceCreate input to Experience model
func (ec *ExperienceCreate) ToExperience() Experience {
	return Experience{
		Role:            ec.Role,
		Company:         ec.Company,
		LogoUrl:         "", // Will be set during file upload
		JobType:         ec.JobType,
		StartDate:       ec.StartDate,
		EndDate:         ec.EndDate,
		Location:        ec.Location,
		Arrangement:     ec.Arrangement,
		WorkDescription: ec.WorkDescription,
		Impact:          ec.Impact,
		ImagesUrl:       nil, // Will be set during file upload
		IsFeatured:      ec.IsFeatured,
	}
}

// ToExperience converts ExperienceUpdate to Experience
//
// @Description Converts ExperienceUpdate input to Experience model
func (eu *ExperienceUpdate) ToExperience() Experience {
	return Experience{
		ID:              eu.ID,
		Role:            eu.Role,
		Company:         eu.Company,
		LogoUrl:         "", // Will be set during file upload
		JobType:         eu.JobType,
		StartDate:       eu.StartDate,
		EndDate:         eu.EndDate,
		Location:        eu.Location,
		Arrangement:     eu.Arrangement,
		WorkDescription: eu.WorkDescription,
		Impact:          eu.Impact,
		ImagesUrl:       nil, // Will be set during file upload
		IsFeatured:      eu.IsFeatured,
	}
}
