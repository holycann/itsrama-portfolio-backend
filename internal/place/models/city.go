package models

// City represents a city entity in the system
type City struct {
	ID       string `json:"id" db:"id" example:"city_12345"`              // Unique ID for the city, example: "city_12345"
	Name     string `json:"name" db:"name" example:"Jakarta"`             // City name, example: "Jakarta"
	Province string `json:"province" db:"province" example:"DKI Jakarta"` // Name of the province where the city is located, example: "DKI Jakarta"
}

// RequestCity is used for city data creation or update requests
type RequestCity struct {
	City
}

// ResponseCity is used for returning city data to the client
type ResponseCity struct {
	City
}
