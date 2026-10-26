package models

import "time"

// URL represents the core URL mapping entity stored in the database.
type URL struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
}

// ShortenRequest represents the JSON payload submitted by the client.
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse represents the JSON returned upon successfully shortening a URL.
type ShortenResponse struct {
	ShortCode   string    `json:"short_code"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

// ErrorResponse represents a standardized JSON error message.
type ErrorResponse struct {
	Error string `json:"error"`
}
