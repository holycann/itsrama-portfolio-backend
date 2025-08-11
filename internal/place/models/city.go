package models

import (
	"time"

	"github.com/google/uuid"
)

// City represents a geographical city in the system
// @Description Comprehensive model for tracking and managing city information
// @Description Provides a structured representation of cities with detailed metadata and geographical context
// @Tags Geographical Locations
type City struct {
	// Unique identifier for the city
	// @Description Globally unique UUID for the city, generated automatically
	// @Description Serves as the primary key and reference for the city
	// @Example "city_123"
	// @Format uuid
	ID uuid.UUID `json:"id" db:"id"`

	// City name
	// @Description Official name of the city
	// @Description Provides a clear, distinctive identifier for the location
	// @Example "Jakarta"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name" db:"name" validate:"required,min=2,max=100"`

	// City description
	// @Description Detailed explanation of the city's characteristics, history, or significance
	// @Description Provides context and additional information about the city
	// @Example "The capital city of Indonesia"
	// @MaxLength 500
	Description string `json:"description" db:"description" validate:"omitempty,max=500"`

	// ID of the province where the city is located
	// @Description Unique identifier linking the city to its parent province
	// @Description Enables hierarchical geographical organization
	// @Format uuid
	ProvinceID uuid.UUID `json:"province_id" db:"province_id"`

	// URL of the city's image
	// @Description Public URL pointing to a representative image of the city
	// @Description Serves as a visual representation for the city
	// @Format uri
	ImageURL string `json:"image_url" db:"image_url"`

	// Timestamp when the city was created
	// @Description Precise timestamp of city record creation in UTC
	// @Description Helps track city information lifecycle and origin
	// @Format date-time
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Timestamp when the city was last updated
	// @Description Precise timestamp of the last modification to the city details in UTC
	// @Description Indicates when city information was last changed
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// CityDTO represents the data transfer object for city information
// @Description Comprehensive data transfer object for city details
// @Description Used for API responses to provide rich city information with related entities
// @Tags Geographical Locations
type CityDTO struct {
	// Unique identifier for the city
	// @Description Globally unique UUID for the city
	// @Example "city_123"
	// @Format uuid
	ID uuid.UUID `json:"id"`

	// City name
	// @Description Official name of the city
	// @Example "Jakarta"
	Name string `json:"name"`

	// City description
	// @Description Detailed explanation of the city's characteristics
	// @Example "The capital city of Indonesia"
	Description string `json:"description"`

	// ID of the province where the city is located
	// @Description Unique identifier of the parent province
	// @Format uuid
	ProvinceID uuid.UUID `json:"province_id"`

	// Province details
	// @Description Comprehensive information about the province containing the city
	// @Description Provides broader geographical context
	Province *Province `json:"province,omitempty"`

	// URL of the city's image
	// @Description Public URL pointing to a representative image of the city
	// @Format uri
	ImageURL string `json:"image_url"`

	// Timestamp when the city was created
	// @Description Precise timestamp of city record creation in UTC
	// @Format date-time
	CreatedAt *time.Time `json:"created_at"`

	// Timestamp when the city was last updated
	// @Description Precise timestamp of the last modification to the city details in UTC
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at"`
}

// ToDTO converts a City to a CityDTO
// @Description Transforms a City model into a lightweight data transfer object
// @Description Useful for API responses and data serialization
// @Return CityDTO Converted city data transfer object
func (c *City) ToDTO() CityDTO {
	return CityDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		ProvinceID:  c.ProvinceID,
		ImageURL:    c.ImageURL,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// CityCreate represents the payload for creating a city
// @Description Structured payload for city creation operations
// @Description Supports input for initializing new city records
// @Tags Geographical Locations
type CityCreate struct {
	// City name
	// @Description Official name of the new city
	// @Description Must be unique and descriptive
	// @Example "Jakarta"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name" validate:"required,min=2,max=100"`

	// City description
	// @Description Detailed explanation of the city's characteristics
	// @Description Provides context and additional information
	// @Example "The capital city of Indonesia"
	// @MaxLength 500
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`

	// ID of the province where the city is located
	// @Description Unique identifier of the parent province
	// @Description Required for geographical organization
	// @Format uuid
	ProvinceID uuid.UUID `json:"province_id" validate:"required"`

	// URL of the city's image
	// @Description Optional public URL pointing to a representative image
	// @Description Serves as a visual representation for the city
	// @Format uri
	ImageURL string `json:"image_url,omitempty"`
}

// CityUpdate represents the payload for updating a city
// @Description Structured payload for city update operations
// @Description Supports partial updates with optional fields
// @Tags Geographical Locations
type CityUpdate struct {
	// Unique identifier for the city
	// @Description Globally unique UUID of the city to be updated
	// @Description Must match an existing city in the system
	// @Example "city_123"
	// @Format uuid
	ID uuid.UUID `json:"id" validate:"required"`

	// City name
	// @Description Updated official name for the city
	// @Description Optional field for renaming the city
	// @Example "Jakarta Raya"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`

	// City description
	// @Description Updated description of the city's characteristics
	// @Description Optional field for refining city information
	// @Example "The expanded capital region of Indonesia"
	// @MaxLength 500
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`

	// ID of the province where the city is located
	// @Description Optional update for the city's parent province
	// @Description Allows changing the geographical context
	// @Format uuid
	ProvinceID uuid.UUID `json:"province_id,omitempty"`

	// URL of the city's image
	// @Description Optional update for the city's representative image
	// @Description Allows changing the visual representation
	// @Format uri
	ImageURL string `json:"image_url,omitempty"`
}
