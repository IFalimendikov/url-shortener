package types

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}

type BatchUnitURLRequest struct {
	URL    string `json:"original_url"`
	UserID string `json:"user_id"`
}

type BatchUnitURLResponse struct {
	Short  string `json:"short_url"`
	UserID string `json:"user_id"`
}

type UserURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type URLRecord struct {
	ID       int   `json:"uuid"`
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
	URL      string `json:"original_url"`
	Deleted  bool   `json:"deleted"`
}

type DeleteRecord struct {
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
}
