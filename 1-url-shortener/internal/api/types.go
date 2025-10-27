package api

type URLRequest struct {
	URL string `json:"url"`
}

type URLResponse struct {
	Shortcode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
