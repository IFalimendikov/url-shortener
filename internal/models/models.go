package models

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}

type BatchUnitURLRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

type BatchUnitURLResponse struct {
	ID    string `json:"correlation_id"`
	Short string `json:"short_url"`
}
