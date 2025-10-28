package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"valentino7504/1-url-shortener/internal/base62"
	"valentino7504/1-url-shortener/internal/db"
	"valentino7504/1-url-shortener/internal/service"
)

func redirect(s *service.ShortenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("short_code")
		if !base62.IsValidBase62(code) {
			http.NotFound(w, r)
			return
		}
		url, err := s.ResolveShortURL(code)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				http.NotFound(w, r)
			} else {
				http.Error(
					w,
					"failed to resolve short code - internal server error",
					http.StatusInternalServerError,
				)
			}
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func deleteShortURL(s *service.ShortenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("short_code")
		if !base62.IsValidBase62(code) {
			sendJSONError(w, "invalid short code", http.StatusBadRequest)
			return
		}
		if err := s.DeleteShortURL(code); err != nil {
			sendJSONError(
				w, "unable to delete short URL - internal error",
				http.StatusInternalServerError,
			)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func getShortURLDetails(s *service.ShortenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("short_code")
		if !base62.IsValidBase62(code) {
			sendJSONError(w, "invalid short code", http.StatusBadRequest)
			return
		}
		url, err := s.GetURLDetails(code)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				sendJSONError(w, "short code does not exist", http.StatusNotFound)
			} else {
				sendJSONError(
					w,
					"failed to resolve short code - internal server error",
					http.StatusInternalServerError,
				)
			}
			return
		}
		sendJSON(w, *url, http.StatusOK)
	}
}

func createShortURL(s *service.ShortenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var req URLRequest
		if err := dec.Decode(&req); err != nil {
			sendJSONError(w, "invalid json request", http.StatusBadRequest)
			return
		}
		if req.URL == "" {
			sendJSONError(w, "URL field is required", http.StatusBadRequest)
			return
		}
		code, err := s.CreateShortURL(req.URL)
		if err != nil {
			if errors.Is(err, service.MalformattedURLErr) {
				sendJSONError(w, "invalid URL", http.StatusBadRequest)
			} else {
				s.Logger.Error("unable to create shortcode", "error", err)
				sendJSONError(w, "failed to create short URL", http.StatusInternalServerError)
			}
			return
		}
		sendJSON(w, URLResponse{
			Shortcode:   code,
			ShortURL:    buildShortenedURL(r, code),
			OriginalURL: req.URL,
		}, http.StatusCreated)
	}
}
