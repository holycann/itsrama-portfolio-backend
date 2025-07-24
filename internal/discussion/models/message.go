package models

import (
	"time"
)

// Message merepresentasikan entitas pesan dalam sistem diskusi
type Message struct {
	ID        string    `json:"id" db:"id" example:"msg_12345"`                  // ID unik untuk pesan, contoh: "msg_12345"
	ThreadID  string    `json:"thread_id" db:"thread_id" example:"thread_12345"` // ID thread tempat pesan berada
	Content   string    `json:"content" db:"content" example:"Ini isi pesan"`    // Isi pesan
	CreatedAt time.Time `json:"created_at" db:"created_at"`                      // Waktu pembuatan pesan
}

// RequestMessage digunakan untuk permintaan pembuatan atau pembaruan data pesan
type RequestMessage struct {
	Message
}

// ResponseMessage digunakan untuk mengembalikan data pesan ke client
type ResponseMessage struct {
	Message
}

// UserID    string    `json:"user_id" db:"user_id" example:"user_67890"`       // ID pengguna yang mengirim pesan
