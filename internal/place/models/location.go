package models

import (
	"time"

	"github.com/google/uuid"
)

// Location represents a specific geographical location in the system
// @Description Comprehensive model for tracking and managing precise geographical points
// @Description Provides a structured representation of locations with detailed geospatial information
// @Tags Geographical Locations
type Location struct {
	// Unique identifier for the location
	// @Description Globally unique UUID for the location, generated automatically
	// @Description Serves as the primary key and reference for the location
	// @Example "location_123"
	// @Format uuid
	ID uuid.UUID `json:"id" db:"id"`

	// Name of the location
	// @Description Official or descriptive name of the geographical point
	// @Description Provides a clear, distinctive identifier for the location
	// @Example "Monas"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name" db:"name" validate:"required,min=2,max=100"`

	// ID of the city where the location is located
	// @Description Unique identifier linking the location to its parent city
	// @Description Enables hierarchical geographical organization
	// @Format uuid
	CityID uuid.UUID `json:"city_id" db:"city_id" validate:"required"`

	// Latitude in decimal degrees
	// @Description Geographic latitude coordinate for precise location mapping
	// @Description Represents the north-south position on the Earth's surface
	// @Example -6.175392
	Latitude float64 `json:"latitude" db:"latitude" validate:"required"`

	// Longitude in decimal degrees
	// @Description Geographic longitude coordinate for precise location mapping
	// @Description Represents the east-west position on the Earth's surface
	// @Example 106.827153
	Longitude float64 `json:"longitude" db:"longitude" validate:"required"`

	// PostGIS location point (automatically set by database trigger)
	// @Description Native PostGIS spatial data point for advanced geospatial queries
	// @Description Automatically generated from latitude and longitude
	Location interface{} `json:"-" db:"location"`

	// Timestamp when the location was created
	// @Description Precise timestamp of location record creation in UTC
	// @Description Helps track location information lifecycle and origin
	// @Format date-time
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Timestamp when the location was last updated
	// @Description Precise timestamp of the last modification to the location details in UTC
	// @Description Indicates when location information was last changed
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// LocationDTO represents the data transfer object for location information
// @Description Comprehensive data transfer object for location details
// @Description Used for API responses to provide rich location information with related entities
// @Tags Geographical Locations
type LocationDTO struct {
	// Unique identifier for the location
	// @Description Globally unique UUID for the location
	// @Example "location_123"
	// @Format uuid
	ID uuid.UUID `json:"id"`

	// Name of the location
	// @Description Official or descriptive name of the geographical point
	// @Example "Monas"
	Name string `json:"name"`

	// ID of the city where the location is located
	// @Description Unique identifier of the parent city
	// @Format uuid
	CityID uuid.UUID `json:"city_id"`

	// City details
	// @Description Comprehensive information about the city containing the location
	// @Description Provides broader geographical context
	City *City `json:"city,omitempty"`

	// Latitude in decimal degrees
	// @Description Geographic latitude coordinate for precise location mapping
	// @Example -6.175392
	Latitude float64 `json:"latitude"`

	// Longitude in decimal degrees
	// @Description Geographic longitude coordinate for precise location mapping
	// @Example 106.827153
	Longitude float64 `json:"longitude"`

	// Timestamp when the location was created
	// @Description Precise timestamp of location record creation in UTC
	// @Format date-time
	CreatedAt *time.Time `json:"created_at"`

	// Timestamp when the location was last updated
	// @Description Precise timestamp of the last modification to the location details in UTC
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at"`
}

// ToDTO converts a Location to a LocationDTO
// @Description Transforms a Location model into a lightweight data transfer object
// @Description Useful for API responses and data serialization
// @Return LocationDTO Converted location data transfer object
func (l *Location) ToDTO() LocationDTO {
	return LocationDTO{
		ID:        l.ID,
		Name:      l.Name,
		CityID:    l.CityID,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
	}
}

// LocationCreate represents the payload for creating a location
// @Description Structured payload for location creation operations
// @Description Supports input for initializing new location records with geospatial data
// @Tags Geographical Locations
type LocationCreate struct {
	// Location name
	// @Description Official or descriptive name of the new location
	// @Description Must be unique and descriptive
	// @Example "Monas"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name" validate:"required,min=2,max=100"`

	// ID of the city where the location is located
	// @Description Unique identifier of the parent city
	// @Description Required for geographical organization
	// @Format uuid
	CityID uuid.UUID `json:"city_id" validate:"required"`

	// Latitude in decimal degrees
	// @Description Geographic latitude coordinate for precise location mapping
	// @Description Represents the north-south position on the Earth's surface
	// @Example -6.175392
	Latitude float64 `json:"latitude" validate:"required"`

	// Longitude in decimal degrees
	// @Description Geographic longitude coordinate for precise location mapping
	// @Description Represents the east-west position on the Earth's surface
	// @Example 106.827153
	Longitude float64 `json:"longitude" validate:"required"`
}

// LocationUpdate represents the payload for updating a location
// @Description Structured payload for location update operations
// @Description Supports partial updates with optional fields
// @Tags Geographical Locations
type LocationUpdate struct {
	// Unique identifier for the location
	// @Description Globally unique UUID of the location to be updated
	// @Description Must match an existing location in the system
	// @Example "location_123"
	// @Format uuid
	ID uuid.UUID `json:"id" validate:"required"`

	// Location name
	// @Description Updated name for the location
	// @Description Optional field for renaming the location
	// @Example "National Monument"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`

	// ID of the city where the location is located
	// @Description Optional update for the location's parent city
	// @Description Allows changing the geographical context
	// @Format uuid
	CityID uuid.UUID `json:"city_id,omitempty"`

	// Latitude in decimal degrees
	// @Description Optional update for the location's latitude coordinate
	// @Description Allows precise geospatial repositioning
	// @Example -6.175392
	Latitude float64 `json:"latitude,omitempty"`

	// Longitude in decimal degrees
	// @Description Optional update for the location's longitude coordinate
	// @Description Allows precise geospatial repositioning
	// @Example 106.827153
	Longitude float64 `json:"longitude,omitempty"`
}
