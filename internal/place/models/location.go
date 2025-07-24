package models

// Location represents a place with coordinates and city reference
type Location struct {
	ID        string  `json:"id" db:"id" example:"loc_12345"`                // Unique identifier for the location
	CityID    string  `json:"city_id" db:"city_id" example:"city_67890"`     // Reference to the city ID
	Name      string  `json:"name" db:"name" example:"Monas"`                // Name of the location
	Latitude  float64 `json:"latitude" db:"latitude" example:"-6.175392"`    // Latitude in decimal degrees
	Longitude float64 `json:"longitude" db:"longitude" example:"106.827153"` // Longitude in decimal degrees
}

// RequestLocation is used for creating or updating a location
type RequestLocation struct {
	Location
}

// ResponseLocation is used for returning location data to the client
type ResponseLocation struct {
	Location
}
