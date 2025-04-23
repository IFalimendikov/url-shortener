package types

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}

type BatchUnitURLRequest struct {
	ID     string `json:"correlation_id"`
	URL    string `json:"original_url"`
	UserID string `json:"user_id"`
}

type BatchUnitURLResponse struct {
	ID     string `json:"correlation_id"`
	Short  string `json:"short_url"`
}

type UserURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLRecord struct {
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
	URL      string `json:"original_url"`
	Deleted  bool   `json:"deleted"`
}

type DeleteRecord struct {
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
}
