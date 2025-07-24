package models

// City merepresentasikan entitas kota dalam sistem
type City struct {
	ID       string `json:"id" db:"id" example:"city_12345"`              // ID unik untuk kota, contoh: "city_12345"
	Name     string `json:"name" db:"name" example:"Jakarta"`             // Nama kota, contoh: "Jakarta"
	Province string `json:"province" db:"province" example:"DKI Jakarta"` // Nama provinsi tempat kota berada, contoh: "DKI Jakarta"
}

// RequestCity digunakan untuk permintaan pembuatan atau pembaruan data kota
type RequestCity struct {
	City
}

// ResponseCity digunakan untuk mengembalikan data kota ke client
type ResponseCity struct {
	City
}
