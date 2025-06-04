package models

// ShortenURLRequest represents the request payload for shortening a single URL
type ShortenURLRequest struct {
    URL string `json:"url"`
}

// ShortenURLResponse represents the response payload containing the shortened URL
type ShortenURLResponse struct {
    Result string `json:"result"`
}

// BatchUnitURLRequest represents a single URL shortening request in a batch operation
type BatchUnitURLRequest struct {
    ID     string `json:"correlation_id"`
    URL    string `json:"original_url"`
    UserID string `json:"user_id"`
}

// BatchUnitURLResponse represents a single URL shortening response in a batch operation
type BatchUnitURLResponse struct {
    ID    string `json:"correlation_id"`
    Short string `json:"short_url"`
}

// UserURLResponse represents a user's URL mapping containing both short and original URLs
type UserURLResponse struct {
    ShortURL    string `json:"short_url"`
    OriginalURL string `json:"original_url"`
}

// URLRecord represents a complete URL record stored in the system
type URLRecord struct {
    UserID   string `json:"user_id"`
    ShortURL string `json:"short_url"`
    URL      string `json:"original_url"`
    Deleted  bool   `json:"deleted"`
}

// DeleteRecord represents a record for URL deletion
type DeleteRecord struct {
    UserID   string `json:"user_id"`
    ShortURL string `json:"short_url"`
}