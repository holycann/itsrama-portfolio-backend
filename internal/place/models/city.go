package models

import (
	"time"

	"github.com/google/uuid"
)

// City represents a city entity in the system
type City struct {
	ID         uuid.UUID `json:"id" db:"id"`                       // Unique ID for the city
	Name       string    `json:"name" db:"name" example:"Jakarta"` // City name, example: "Jakarta"
	ProvinceID uuid.UUID `json:"province_id" db:"province_id"`     // ID of the province where the city is located
	ImageURL   string    `json:"image_url" db:"image_url"`         // URL of the city's image
	CreatedAt  time.Time `json:"created_at" db:"created_at"`       // City creation time
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`       // City last update time
}

// RequestCity is used for city data creation or update requests
type RequestCity struct {
	City
}

// ResponseCity is used for returning city data to the client
type ResponseCity struct {
	City
}
