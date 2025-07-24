package models

import (
	"time"

	"github.com/lib/pq"
)

// LocalStory represents a local story with detailed information
type LocalStory struct {
	ID            string         `json:"id" db:"id" example:"story_12345"`
	Title         string         `json:"title" db:"title" example:"Legenda Sangkuriang"`
	Summary       string         `json:"summary" db:"summary" example:"Cerita tentang asal-usul Tangkuban Perahu"`
	StoryText     string         `json:"story_text" db:"story_text"`
	IsForKids     bool           `json:"is_for_kids" db:"is_for_kids" example:"true"`
	AudioURL      string         `json:"audio_url,omitempty" db:"audio_url"`
	ImageURL      string         `json:"image_url,omitempty" db:"image_url"`
	LocationID    string         `json:"location_id,omitempty" db:"location_id"`
	CityID        string         `json:"city_id,omitempty" db:"city_id"`
	OriginCulture string         `json:"origin_culture" db:"origin_culture" example:"Sunda"`
	Language      string         `json:"language" db:"language" example:"ID"`
	Tags          pq.StringArray `json:"tags" db:"tags"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

// RequestLocalStory is used for creating or updating a local story
type RequestLocalStory struct {
	LocalStory
}

// ResponseLocalStory is used for returning local story data to the client
type ResponseLocalStory struct {
	LocalStory
}
