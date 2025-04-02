package models

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}

type ShortenURLBatchRequest struct {
	URLs []BatchUnitURLRequest `json:"urls"`
}

type BatchUnitURLRequest struct {
	ID  string `json:"correlation_id"`
	URL string `json:"original_url"`
}

type ShortenURLBatchResponse struct {
	URLs []BatchUnitURLResponse `json:"urls"`
}

type BatchUnitURLResponse struct {
	ID  string `json:"correlation_id"`
	Short string `json:"short_url"`
}
