package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	placeModel "github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/users/models"
)

// Event represents a cultural event with detailed location and user information
// @Description Comprehensive model for tracking and managing cultural events
// @Description Provides a structured representation of events with rich metadata, location details, and user context
// @Tags Cultural Events
type Event struct {
	// Unique identifier for the event
	// @Description Globally unique UUID for the event, generated automatically
	// @Description Serves as the primary key and reference for the event
	// @Example "event_123"
	// @Format uuid
	ID uuid.UUID `json:"id" db:"id"`

	// ID of the user who created the event
	// @Description References the user responsible for creating the event
	// @Description Helps track event ownership and attribution
	// @Example "user_456"
	// @Format uuid
	UserID uuid.UUID `json:"user_id" db:"user_id" validate:"required"`

	// Reference to location ID
	// @Description Unique identifier linking the event to a specific location
	// @Description Enables geospatial context and discovery
	// @Example "location_789"
	// @Format uuid
	LocationID uuid.UUID `json:"location_id" db:"location_id" validate:"required"`

	// Event name
	// @Description Official name of the cultural event
	// @Description Provides a concise, descriptive title for the event
	// @Example "Summer Music Festival"
	// @MinLength 2
	// @MaxLength 100
	Name string `json:"name" db:"name" validate:"required,min=2,max=100"`

	// Event description
	// @Description Detailed explanation of the event's purpose, activities, and highlights
	// @Description Provides comprehensive context for potential attendees
	// @Example "A vibrant music festival featuring local and international artists"
	// @MaxLength 500
	Description string `json:"description" db:"description" validate:"required,max=500"`

	// Event image URL
	// @Description Public URL pointing to the event's primary image
	// @Description Serves as a visual representation and promotional material
	// @Example "https://example.com/event_image.jpg"
	// @Format uri
	ImageURL string `json:"image_url" db:"image_url" validate:"omitempty,url" format:"uri"`

	// Event start date
	// @Description Official start date and time of the event
	// @Description Indicates when the event begins
	// @Format date-time
	StartDate time.Time `json:"start_date" db:"start_date" validate:"required"`

	// Event end date
	// @Description Official end date and time of the event
	// @Description Indicates when the event concludes
	// @Format date-time
	EndDate time.Time `json:"end_date" db:"end_date" validate:"required,gtfield=StartDate"`

	// Whether the event is kid-friendly
	// @Description Indicates if the event is suitable for children
	// @Description Helps families and parents make informed attendance decisions
	// @Example true
	IsKidFriendly bool `json:"is_kid_friendly" db:"is_kid_friendly"`

	// Event creation time
	// @Description Timestamp of when the event was first created in the system
	// @Description Helps track event lifecycle and origin
	// @Format date-time
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Event last update time
	// @Description Timestamp of the most recent update to the event details
	// @Description Indicates when event information was last modified
	// @Format date-time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// EventDTO represents the data transfer object for returning comprehensive event details
// @Description Structured representation of event data for API responses
// @Description Includes related entities and provides a clean, comprehensive view of event information
// @Tags Cultural Events
type EventDTO struct {
	// Unique identifier for the event
	// @Description Globally unique UUID for the event
	// @Example "event_123"
	// @Format uuid
	ID uuid.UUID `json:"id"`

	// Event name
	// @Description Official name of the cultural event
	// @Example "Summer Music Festival"
	Name string `json:"name"`

	// Event description
	// @Description Detailed explanation of the event's purpose, activities, and highlights
	// @Example "A vibrant music festival featuring local and international artists"
	Description string `json:"description"`

	// Event image URL
	// @Description Public URL pointing to the event's primary image
	// @Example "https://example.com/event_image.jpg"
	// @Format uri
	ImageURL string `json:"image_url,omitempty"`

	// Event start date
	// @Description Official start date and time of the event
	// @Format date-time
	StartDate time.Time `json:"start_date"`

	// Event end date
	// @Description Official end date and time of the event
	// @Format date-time
	EndDate time.Time `json:"end_date"`

	// Whether the event is kid-friendly
	// @Description Indicates if the event is suitable for children
	// @Example true
	IsKidFriendly bool `json:"is_kid_friendly"`

	// Number of views for the event
	// @Description Total number of times the event has been viewed
	// @Description Helps track event popularity and engagement
	// @Example 13
	Views map[string]int `json:"views,omitempty"`

	// Location details
	// @Description Detailed information about the event's location
	// @Description Provides geographical context and navigation information
	Location *placeModel.Location `json:"location,omitempty"`

	// City details
	// @Description City where the event is taking place
	// @Description Helps users understand the event's urban context
	City *placeModel.City `json:"city,omitempty"`

	// Province details
	// @Description Province or region hosting the event
	// @Description Provides broader geographical information
	Province *placeModel.Province `json:"province,omitempty"`

	// Event creator details
	// @Description Information about the user who created the event
	// @Description Enables attribution and trust
	Creator *models.User `json:"creator,omitempty"`

	// Event creation time
	// @Description Timestamp of when the event was first created in the system
	// @Format date-time
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// ToDTO converts an EventDTO to a map representation
// @Description Transforms the EventDTO into a flexible map for API responses
// @Description Allows dynamic serialization with optional fields
// @Return map[string]interface{} Structured event data
func (r *EventDTO) ToDTO() map[string]interface{} {
	dto := map[string]interface{}{
		"id":              r.ID,
		"name":            r.Name,
		"description":     r.Description,
		"start_date":      r.StartDate,
		"end_date":        r.EndDate,
		"is_kid_friendly": r.IsKidFriendly,
	}

	// Optional fields
	if r.ImageURL != "" {
		dto["image_url"] = r.ImageURL
	}

	if r.Views != nil {
		dto["views"] = r.Views["views"]
	}

	if r.Location != nil {
		dto["location"] = r.Location
	}

	if r.City != nil {
		dto["city"] = r.City
	}

	if r.Province != nil {
		dto["province"] = r.Province
	}

	if r.Creator != nil {
		dto["creator"] = r.Creator
	}

	if r.CreatedAt != nil {
		dto["created_at"] = r.CreatedAt
	}

	return dto
}

// EventPayload represents the input data for creating or updating an event
// @Description Structured payload for event creation and update operations
// @Description Supports multipart form data with flexible input options
// @Tags Cultural Events
type EventPayload struct {
	// Event ID
	// @Description Unique identifier for the event (optional during creation)
	// @Description Used for identifying the event during updates
	// @Format uuid
	ID uuid.UUID `form:"id" json:"id"`

	// User ID of the event creator
	// @Description References the user creating the event
	// @Description Automatically set during event creation
	// @Example "user_456"
	// @Format uuid
	UserID uuid.UUID `form:"user_id" json:"user_id" validate:"required"`

	// Event name
	// @Description Official name of the cultural event
	// @Description Must be unique and descriptive
	// @Example "Summer Music Festival"
	// @MinLength 2
	// @MaxLength 100
	Name string `form:"name" json:"name" validate:"required,min=2,max=100"`

	// Event description
	// @Description Detailed explanation of the event's purpose, activities, and highlights
	// @Description Provides comprehensive context for potential attendees
	// @Example "A vibrant music festival featuring local and international artists"
	// @MaxLength 500
	Description string `form:"description" json:"description" validate:"required,max=500"`

	// Location details
	// @Description Comprehensive location information for the event
	// @Description Enables precise geographical context and discovery
	Location *placeModel.LocationCreate `form:"location" json:"location" validate:"required"`

	// Event start date
	// @Description Official start date and time of the event
	// @Description Indicates when the event begins
	// @Format date-time
	StartDate time.Time `form:"start_date" json:"start_date" validate:"required"`

	// Event end date
	// @Description Official end date and time of the event
	// @Description Indicates when the event concludes
	// @Format date-time
	EndDate time.Time `form:"end_date" json:"end_date" validate:"required"`

	// Whether the event is kid-friendly
	// @Description Indicates if the event is suitable for children
	// @Description Helps families and parents make informed attendance decisions
	// @Example true
	IsKidFriendly bool `form:"is_kid_friendly" json:"is_kid_friendly"`

	// Event image URL
	// @Description Public URL pointing to the event's primary image
	// @Description Optional field for pre-existing image URLs
	// @Example "https://example.com/event_image.jpg"
	// @Format uri
	ImageURL string `json:"image_url,omitempty"`

	// Event image file
	// @Description Multipart file upload for the event image
	// @Description Allows direct image upload during event creation/update
	Image *multipart.FileHeader `form:"image" json:"-" validate:"omitempty"`
}
