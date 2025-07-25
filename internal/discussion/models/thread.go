package models

// Thread represents a discussion thread entity in the system
type Thread struct {
	ID    string `json:"id" db:"id" example:"thread_12345"`                         // Unique ID for the thread, example: "thread_12345"
	Title string `json:"title" db:"title" example:"Discussion about Monas History"` // Thread title, example: "Discussion about Monas History"
	// Content string `json:"content" db:"content" example:"Let's discuss about Monas history!"` // Thread content
	// UserID    string `json:"user_id" db:"user_id" example:"user_67890"`                         // ID of the user who created the thread
	EventID   string `json:"event_id" db:"event_id" example:"event-xx"`                 // Reference to related event (if any)
	Status    string `json:"status" db:"status" example:"open"`                         // Thread status, example: "open"
	CreatedAt string `json:"created_at" db:"created_at" example:"2024-06-01T15:04:05Z"` // Thread creation time
}

// RequestThread is used for thread creation or update requests
type RequestThread struct {
	Thread
}

// ResponseThread is used for returning thread data to the client
type ResponseThread struct {
	Thread
}
