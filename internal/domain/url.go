package domain

import "time"

// URL represents the core data model
type URL struct {
    ID          int64     `json:"id"`
    OriginalURL string    `json:"original_url"`
    ShortCode   string    `json:"short_code"`
    CreatedAt   time.Time `json:"created_at"`
    ExpiresAt   time.Time `json:"expires_at"`
};

// URLRepository defines how we talk to databases
type URLRepository interface {
    Save(url *URL) error
    GetByCode(code string) (*URL, error)
}

// ShortenerService defines the business logic
type ShortenerService interface {
    Shorten(originalURL string) (string, error)
    Resolve(code string) (string, error)
}