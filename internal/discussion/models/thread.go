package models

// Thread merepresentasikan entitas thread diskusi dalam sistem
type Thread struct {
	ID    string `json:"id" db:"id" example:"thread_12345"`                // ID unik untuk thread, contoh: "thread_12345"
	Title string `json:"title" db:"title" example:"Diskusi Sejarah Monas"` // Judul thread, contoh: "Diskusi Sejarah Monas"
	// Content string `json:"content" db:"content" example:"Ayo diskusi tentang sejarah Monas!"` // Konten thread
	// UserID    string `json:"user_id" db:"user_id" example:"user_67890"`                         // ID pengguna yang membuat thread
	EventID   string `json:"event_id" db:"event_id" example:"event-xx"`                 // Referensi ke event terkait (jika ada)
	Status    string `json:"status" db:"status" example:"open"`                         // Status thread, contoh: "open"
	CreatedAt string `json:"created_at" db:"created_at" example:"2024-06-01T15:04:05Z"` // Waktu pembuatan thread
}

// RequestThread digunakan untuk permintaan pembuatan atau pembaruan thread
type RequestThread struct {
	Thread
}

// ResponseThread digunakan untuk mengembalikan data thread ke client
type ResponseThread struct {
	Thread
}
