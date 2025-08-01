package models

import (
	"time"

	"github.com/google/uuid"
	placeModel "github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/users/models"
)

// Event represents a cultural event with location and city references
type Event struct {
	ID            uuid.UUID `json:"id" db:"id"`                           // Unique identifier for the event
	UserID        uuid.UUID `json:"user_id" db:"user_id"`                 // ID of the user who created the event
	LocationID    uuid.UUID `json:"location_id" db:"location_id"`         // Reference to location ID
	CityID        uuid.UUID `json:"city_id" db:"city_id"`                 // Reference to city ID
	ProvinceID    uuid.UUID `json:"province_id" db:"province_id"`         // Reference to province ID
	Name          string    `json:"name" db:"name"`                       // Event name
	Description   string    `json:"description" db:"description"`         // Event description
	ImageURL      string    `json:"image_url" db:"image_url"`             // Event image URL
	StartDate     time.Time `json:"start_date" db:"start_date"`           // Event start date
	EndDate       time.Time `json:"end_date" db:"end_date"`               // Event end date
	IsKidFriendly bool      `json:"is_kid_friendly" db:"is_kid_friendly"` // Whether the event is kid-friendly
	CreatedAt     time.Time `json:"created_at" db:"created_at"`           // Event creation time
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`           // Event last update time
}

// RequestEvent is used for creating or updating an event
type RequestEvent struct {
	Event
	Location *placeModel.Location `json:"location"`
}

// ResponseEvent is used for returning event data to the client
type ResponseEvent struct {
	Event
	User     *models.User         `json:"user,omitempty"`
	Location *placeModel.Location `json:"location,omitempty"`
	City     *placeModel.City     `json:"city,omitempty"`
	Province *placeModel.Province `json:"province,omitempty"`
	Views    int                  `json:"views,omitempty"`
}
