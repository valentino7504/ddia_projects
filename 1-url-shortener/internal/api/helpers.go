package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func buildShortenedURL(r *http.Request, shortcode string) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s", scheme, r.Host, shortcode)
}

func sendJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendJSONError(w http.ResponseWriter, msg string, status int) {
	sendJSON(w, ErrorResponse{Error: msg}, status)
}
