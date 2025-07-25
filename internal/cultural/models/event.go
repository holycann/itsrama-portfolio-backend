package models

import (
	"time"

	"github.com/holycann/cultour-backend/internal/users/models"
)

// Event represents a place with coordinates and city reference
type Event struct {
	ID            string    `json:"id" db:"id" example:"loc_12345"`                                   // Unique identifier for the event
	UserID        string    `json:"user_id" db:"user_id" example:"user_67890"`                        // ID of the user who created the event
	LocationID    string    `json:"location_id" db:"location_id" example:"location_67890"`            // Reference to location ID
	CityID        string    `json:"city_id" db:"city_id" example:"city_12345"`                        // Reference to city ID
	ProvinceID    string    `json:"province_id" db:"province_id" example:"province_67890"`            // Reference to province ID
	Name          string    `json:"name" db:"name" example:"Monas"`                                   // Event name
	Description   string    `json:"description" db:"description" example:"Monas"`                     // Event description
	ImageURL      string    `json:"image_url" db:"image_url" example:"https://example.com/image.jpg"` // Event image URL
	StartDate     time.Time `json:"start_date" db:"start_date" example:"2024-06-01T08:00:00+07:00"`   // Event start date (format: YYYY-MM-DDTHH:MM:SS±HH:MM)
	EndDate       time.Time `json:"end_date" db:"end_date" example:"2024-06-01T09:00:00+07:00"`       // Event end date (format: YYYY-MM-DDTHH:MM:SS±HH:MM)
	IsKidFriendly bool      `json:"is_kid_friendly" db:"is_kid_friendly" example:"true"`              // Whether the event is kid-friendly
	Views         int8      `json:"views" db:"views" example:"10"`
}

// RequestEvent is used for creating or updating a location
type RequestEvent struct {
	Event
}

// ResponseEvent is used for returning location data to the client
type ResponseEvent struct {
	Event
	User *models.User `json:"user,omitempty"`
}
