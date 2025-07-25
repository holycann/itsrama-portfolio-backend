package models

// Province represents a province entity in the system
type Province struct {
	ID   string `json:"id" db:"id" example:"province_12345"` // Unique ID for the province, example: "province_12345"
	Name string `json:"name" db:"name" example:"West Java"`  // Province name, example: "West Java"
}

// RequestProvince is used for province data creation or update requests
type RequestProvince struct {
	Province
}

// ResponseProvince is used for returning province data to the client
type ResponseProvince struct {
	Province
}
