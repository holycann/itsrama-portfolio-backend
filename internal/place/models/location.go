package models

import (
	"time"

	"github.com/google/uuid"
)

// Location represents a place with coordinates and city reference
type Location struct {
	ID        uuid.UUID `json:"id" db:"id"`                                    // Unique identifier for the location
	CityID    uuid.UUID `json:"city_id" db:"city_id"`                          // Reference to the city ID
	Name      string    `json:"name" db:"name" example:"Monas"`                // Name of the location
	Latitude  float64   `json:"latitude" db:"latitude" example:"-6.175392"`    // Latitude in decimal degrees
	Longitude float64   `json:"longitude" db:"longitude" example:"106.827153"` // Longitude in decimal degrees
	CreatedAt time.Time `json:"created_at" db:"created_at"`                    // Location creation time
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`                    // Location last update time
}

// RequestLocation is used for creating or updating a location
type RequestLocation struct {
	Location
}

// ResponseLocation is used for returning location data to the client
type ResponseLocation struct {
	Location
}
