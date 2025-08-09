package models

import (
	"time"

	"github.com/google/uuid"
)

// Province represents an administrative province in the system
// @Description Comprehensive model for tracking and managing administrative regions
// @Description Provides a structured representation of provinces with detailed geographical context
// @Tags Geographical Locations
type Province struct {
	// Unique identifier for the province
	// @Description Globally unique UUID for the province, generated automatically
	// @Description Serves as the primary key and reference for the province
	// @Example "province_123"
	// @Format uuid
	ID uuid.UUID `json:"id" db:"id"`

	// Province name
	// @Description Official name of the administrative region
	// @Description Provides a clear, distinctive identifier for the province
	// @Example "West Java"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name" db:"name" validate:"required,min=2,max=100"`

	// Province description
	// @Description Detailed explanation of the province's characteristics, history, or significance
	// @Description Provides context and additional information about the administrative region
	// @Example "A beautiful province with rich cultural heritage"
	// @MaxLength 500
	Description string `json:"description" db:"description" validate:"omitempty,max=500"`

	// Timestamp when the province was created
	// @Description Precise timestamp of province record creation in UTC
	// @Description Helps track province information lifecycle and origin
	// @Format date-time
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Timestamp when the province was last updated
	// @Description Precise timestamp of the last modification to the province details in UTC
	// @Description Indicates when province information was last changed
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// ProvinceDTO represents the data transfer object for province information
// @Description Comprehensive data transfer object for province details
// @Description Used for API responses to provide rich province information with related entities
// @Tags Geographical Locations
type ProvinceDTO struct {
	// Unique identifier for the province
	// @Description Globally unique UUID for the province
	// @Example "province_123"
	// @Format uuid
	ID uuid.UUID `json:"id"`

	// Province name
	// @Description Official name of the administrative region
	// @Example "West Java"
	Name string `json:"name"`

	// Province description
	// @Description Detailed explanation of the province's characteristics
	// @Example "A beautiful province with rich cultural heritage"
	Description string `json:"description"`

	// Timestamp when the province was created
	// @Description Precise timestamp of province record creation in UTC
	// @Format date-time
	CreatedAt *time.Time `json:"created_at"`

	// Timestamp when the province was last updated
	// @Description Precise timestamp of the last modification to the province details in UTC
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at"`
}

// ToDTO converts a Province to a ProvinceDTO
// @Description Transforms a Province model into a lightweight data transfer object
// @Description Useful for API responses and data serialization
// @Return ProvinceDTO Converted province data transfer object
func (p *Province) ToDTO() ProvinceDTO {
	return ProvinceDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// ProvinceCreate represents the payload for creating a province
// @Description Structured payload for province creation operations
// @Description Supports input for initializing new province records
// @Tags Geographical Locations
type ProvinceCreate struct {
	// Province name
	// @Description Official name of the new administrative region
	// @Description Must be unique and descriptive
	// @Example "West Java"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name" validate:"required,min=2,max=100"`

	// Province description
	// @Description Detailed explanation of the province's characteristics
	// @Description Provides context and additional information
	// @Example "A beautiful province with rich cultural heritage"
	// @MaxLength 500
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// ProvinceUpdate represents the payload for updating a province
// @Description Structured payload for province update operations
// @Description Supports partial updates with optional fields
// @Tags Geographical Locations
type ProvinceUpdate struct {
	// Unique identifier for the province
	// @Description Globally unique UUID of the province to be updated
	// @Description Must match an existing province in the system
	// @Example "province_123"
	// @Format uuid
	ID uuid.UUID `json:"id" validate:"required"`

	// Province name
	// @Description Updated name for the administrative region
	// @Description Optional field for renaming the province
	// @Example "Greater West Java"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`

	// Province description
	// @Description Updated description of the province's characteristics
	// @Description Optional field for refining province information
	// @Example "An expanded province with diverse cultural landscapes"
	// @MaxLength 500
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
}
