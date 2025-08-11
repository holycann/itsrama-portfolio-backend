package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	placeModel "github.com/holycann/cultour-backend/internal/place/models"
	"github.com/holycann/cultour-backend/internal/users/models"
	"github.com/lib/pq"
)

// LocalStory represents a local story with detailed information
type LocalStory struct {
	// Unique identifier for the local story
	// @example "story_123"
	ID uuid.UUID `json:"id" db:"id"`

	// ID of the user who created the local story
	// @example "user_456"
	CreatorID uuid.UUID `json:"creator_id" db:"creator_id" validate:"required"`

	// Reference to location ID
	// @example "location_789"
	LocationID uuid.UUID `json:"location_id" db:"location_id" validate:"required"`

	// Local story title
	// @example "Legend of Sangkuriang"
	Title string `json:"title" db:"title" validate:"required,min=2,max=100"`

	// Local story summary
	// @example "Story about the origin of Tangkuban Perahu"
	Summary string `json:"summary" db:"summary" validate:"required,max=500"`

	// Full story text
	StoryText string `json:"story_text" db:"story_text" validate:"required"`

	// Whether the story is suitable for kids
	IsForKids bool `json:"is_for_kids" db:"is_for_kids"`

	// Audio URL for the story
	// @example "https://example.com/story_audio.mp3"
	AudioURL string `json:"audio_url,omitempty" db:"audio_url" validate:"omitempty,url"`

	// Image URL for the story
	// @example "https://example.com/story_image.jpg"
	ImageURL string `json:"image_url,omitempty" db:"image_url" validate:"omitempty,url"`

	// Origin culture of the story
	// @example "Sunda"
	OriginCulture string `json:"origin_culture" db:"origin_culture" validate:"required"`

	// Language of the story
	// @example "EN"
	Language string `json:"language" db:"language" validate:"required"`

	// Tags associated with the story
	Tags pq.StringArray `json:"tags" db:"tags"`

	// Story creation time
	CreatedAt *time.Time `json:"created_at" db:"created_at"`

	// Story last update time
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

// LocalStoryDTO represents the data transfer object for returning local story details
type LocalStoryDTO struct {
	// Unique identifier for the local story
	// @example "story_123"
	ID uuid.UUID `json:"id"`

	// Local story title
	// @example "Legend of Sangkuriang"
	Title string `json:"title"`

	// Local story summary
	// @example "Story about the origin of Tangkuban Perahu"
	Summary string `json:"summary"`

	// Full story text
	StoryText string `json:"story_text"`

	// Whether the story is suitable for kids
	IsForKids bool `json:"is_for_kids"`

	// Audio URL for the story
	// @example "https://example.com/story_audio.mp3"
	AudioURL string `json:"audio_url,omitempty"`

	// Image URL for the story
	// @example "https://example.com/story_image.jpg"
	ImageURL string `json:"image_url,omitempty"`

	// Location details
	Location *placeModel.Location `json:"location,omitempty"`

	// City details
	City *placeModel.City `json:"city,omitempty"`

	// Origin culture of the story
	// @example "Sunda"
	OriginCulture string `json:"origin_culture"`

	// Language of the story
	// @example "EN"
	Language string `json:"language"`

	// Tags associated with the story
	Tags pq.StringArray `json:"tags"`

	// Story creator details
	Creator *models.User `json:"creator,omitempty"`

	// Story creation time
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// ToDTO converts a LocalStoryDTO to a map representation
func (r *LocalStoryDTO) ToDTO() map[string]interface{} {
	dto := map[string]interface{}{
		"id":             r.ID,
		"title":          r.Title,
		"summary":        r.Summary,
		"story_text":     r.StoryText,
		"is_for_kids":    r.IsForKids,
		"origin_culture": r.OriginCulture,
		"language":       r.Language,
		"tags":           r.Tags,
	}

	// Optional fields
	if r.AudioURL != "" {
		dto["audio_url"] = r.AudioURL
	}

	if r.ImageURL != "" {
		dto["image_url"] = r.ImageURL
	}

	if r.Location != nil {
		dto["location"] = r.Location
	}

	if r.City != nil {
		dto["city"] = r.City
	}

	if r.Creator != nil {
		dto["creator"] = r.Creator
	}

	if r.CreatedAt != nil {
		dto["created_at"] = r.CreatedAt
	}

	return dto
}

// LocalStoryPayload represents the payload for creating or updating a local story
type LocalStoryPayload struct {
	// Local story title
	Title string `form:"title" json:"title" validate:"required,min=2,max=100"`

	// Local story summary
	Summary string `form:"summary" json:"summary" validate:"required,max=500"`

	// Full story text
	StoryText string `form:"story_text" json:"story_text" validate:"required"`

	// City ID
	CityID uuid.UUID `form:"city_id" json:"city_id" validate:"required"`

	// Whether the story is suitable for kids
	IsForKids bool `form:"is_for_kids" json:"is_for_kids"`

	// Origin culture of the story
	OriginCulture string `form:"origin_culture" json:"origin_culture" validate:"required"`

	// Language of the story
	Language string `form:"language" json:"language" validate:"required"`

	// Tags associated with the story
	Tags []string `form:"tags" json:"tags"`

	// Story image
	Image *multipart.FileHeader `form:"image" json:"-" validate:"omitempty"`

	// Story audio
	Audio *multipart.FileHeader `form:"audio" json:"-" validate:"omitempty"`
}
