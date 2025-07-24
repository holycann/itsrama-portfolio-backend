package models

import "time"

// Event represents a place with coordinates and city reference
type Event struct {
	ID            string     `json:"id" db:"id" example:"loc_12345"`                                 // Unique identifier untuk event
	LocationID    string     `json:"location_id" db:"location_id" example:"location_67890"`          // Referensi ke ID lokasi
	Name          string     `json:"name" db:"name" example:"Monas"`                                 // Nama event
	Description   string     `json:"description" db:"description" example:"Monas"`                   // Deskripsi event
	StartDate     *time.Time `json:"start_date" db:"start_date" example:"2024-06-01T08:00:00+07:00"` // Tanggal mulai event (format: YYYY-MM-DDTHH:MM:SS±HH:MM)
	EndDate       *time.Time `json:"end_date" db:"end_date" example:"2024-06-01T09:00:00+07:00"`     // Tanggal selesai event (format: YYYY-MM-DDTHH:MM:SS±HH:MM)
	IsKidFriendly bool       `json:"is_kid_friendly" db:"is_kid_friendly" example:"true"`            // Apakah event ramah anak
	Views         int8       `json:"views" db:"views" example:"10"`
}

// RequestEvent is used for creating or updating a location
type RequestEvent struct {
	Event
}

// ResponseEvent is used for returning location data to the client
type ResponseEvent struct {
	Event
}
