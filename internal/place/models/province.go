package models

import (
	"time"

	"github.com/google/uuid"
)

// Province represents a province entity in the system
type Province struct {
	ID          uuid.UUID `json:"id" db:"id"`                         // Unique ID for the province
	Name        string    `json:"name" db:"name" example:"West Java"` // Province name, example: "West Java"
	Description string    `json:"description" db:"description"`       // Province description
	CreatedAt   time.Time `json:"created_at" db:"created_at"`         // Province creation time
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`         // Province last update time
}

// RequestProvince is used for province data creation or update requests
type RequestProvince struct {
	Province
}

// ResponseProvince is used for returning province data to the client
type ResponseProvince struct {
	Province
}
