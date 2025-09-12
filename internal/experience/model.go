package experience

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/holycann/itsrama-portfolio-backend/internal/tech_stack"
	"github.com/holycann/itsrama-portfolio-backend/pkg/utils"
)

type Experience struct {
	// Identification
	ID uuid.UUID `json:"id" db:"id"`

	// Job Details
	Role    string `json:"role" db:"role" validate:"required"`
	Company string `json:"company" db:"company" validate:"required"`
	LogoUrl string `json:"logo_url" db:"logo_url"`
	JobType string `json:"job_type" db:"job_type"`

	// Timing and Location
	StartDate   utils.CustomDate  `json:"start_date" db:"start_date" validate:"required"`
	EndDate     *utils.CustomDate `json:"end_date" db:"end_date"`
	Location    string            `json:"location" db:"location"`
	Arrangement string            `json:"arrangement" db:"arrangement"`

	// Job Description
	WorkDescription string   `json:"work_description" db:"work_description"`
	Impact          []string `json:"impact" db:"impact" pg:"array"`
	ImagesUrl       []string `json:"images_url" db:"images_url" pg:"array"`

	// Metadata
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type ExperienceTechStack struct {
	ExperienceID uuid.UUID `json:"experience_id" db:"experience_id" validate:"required"`
	TechStackID  uuid.UUID `json:"tech_stack_id" db:"tech_stack_id" validate:"required"`
}

type ExperienceDTO struct {
	// Identification
	ID uuid.UUID `json:"id" db:"id"`

	// Job Details
	Role    string `json:"role" db:"role" validate:"required"`
	Company string `json:"company" db:"company" validate:"required"`
	LogoUrl string `json:"logo_url" db:"logo_url"`
	JobType string `json:"job_type" db:"job_type"`

	// Timing and Location
	StartDate   utils.CustomDate  `json:"start_date" db:"start_date" validate:"required"`
	EndDate     *utils.CustomDate `json:"end_date" db:"end_date"`
	Location    string            `json:"location" db:"location"`
	Arrangement string            `json:"arrangement" db:"arrangement"`

	// Job Description
	WorkDescription string   `json:"work_description" db:"work_description"`
	Impact          []string `json:"impact" db:"impact" pg:"array"`
	ImagesUrl       []string `json:"images_url" db:"images_url" pg:"array"`

	// Metadata
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`

	// Relationships
	TechStack []tech_stack.TechStack `json:"tech_stack" db:"tech_stack" pg:"array"`
}

type ExperienceCreate struct {
	TechStackIds    []uuid.UUID             `json:"tech_stack_ids"`
	Role            string                  `json:"role"`
	Company         string                  `json:"company"`
	LogoImage       *multipart.FileHeader   `json:"logo_image"`
	JobType         string                  `json:"job_type"`
	StartDate       utils.CustomDate        `json:"start_date"`
	EndDate         *utils.CustomDate       `json:"end_date"`
	Location        string                  `json:"location"`
	Arrangement     string                  `json:"arrangement"`
	WorkDescription string                  `json:"work_description"`
	Impact          []string                `json:"impact"`
	Images          []*multipart.FileHeader `json:"images"`
}

type ExperienceUpdate struct {
	ID              uuid.UUID               `json:"id"`
	TechStackIds    []uuid.UUID             `json:"tech_stack_ids"`
	Role            string                  `json:"role"`
	Company         string                  `json:"company"`
	LogoImage       *multipart.FileHeader   `json:"logo_image"`
	JobType         string                  `json:"job_type"`
	StartDate       utils.CustomDate        `json:"start_date"`
	EndDate         *utils.CustomDate       `json:"end_date"`
	Location        string                  `json:"location"`
	Arrangement     string                  `json:"arrangement"`
	WorkDescription string                  `json:"work_description"`
	Impact          []string                `json:"impact"`
	Images          []*multipart.FileHeader `json:"images"`
}

// ToDTO converts an Experience to an ExperienceDTO
func (e *Experience) ToDTO(techStacks []tech_stack.TechStack) ExperienceDTO {
	return ExperienceDTO{
		ID:              e.ID,
		Role:            e.Role,
		Company:         e.Company,
		LogoUrl:         e.LogoUrl,
		JobType:         e.JobType,
		StartDate:       e.StartDate,
		EndDate:         e.EndDate,
		Location:        e.Location,
		Arrangement:     e.Arrangement,
		WorkDescription: e.WorkDescription,
		Impact:          e.Impact,
		ImagesUrl:       e.ImagesUrl,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		TechStack:       techStacks,
	}
}

// ToExperience converts ExperienceCreate to Experience
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
	}
}

// ToExperience converts ExperienceUpdate to Experience
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
	}
}
