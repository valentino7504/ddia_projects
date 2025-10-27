package api

import (
	"net/http"

	"valentino7504/1-url-shortener/internal/service"
)

func Routes(s *service.ShortenService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{short_code}", redirect(s))
	mux.HandleFunc("GET /api/urls/{short_code}", getShortURLDetails(s))
	mux.HandleFunc("DELETE /api/urls/{short_code}", deleteShortURL(s))
	mux.HandleFunc("POST /api/shorten", createShortURL(s))
	return mux
}
