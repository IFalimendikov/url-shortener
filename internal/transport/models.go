package transport

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortneURLResponse struct {
	Result string `json:"result"`
}