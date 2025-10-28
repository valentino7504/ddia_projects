package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"valentino7504/1-url-shortener/internal/db"
)

type ShortenService struct {
	sqliteDB *sql.DB
	Logger   *slog.Logger
}

var MalformattedURLErr error = errors.New("invalid url provided")

func NewShortenService(sqliteDB *sql.DB, logger *slog.Logger) *ShortenService {
	svc := ShortenService{sqliteDB: sqliteDB, Logger: logger}
	err := db.InitDB(sqliteDB)
	if err != nil {
		svc.Logger.Error("unable to initialize db - quitting", "error", err)
		os.Exit(1)
	}
	return &svc
}

func (s *ShortenService) CreateShortURL(longURL string) (string, error) {
	u, err := url.Parse(longURL)
	if err != nil {
		return "", MalformattedURLErr
	}
	if u.Scheme == "" {
		longURL = "https://" + longURL
		if err != nil {
			return "", MalformattedURLErr
		}
	}
	if u.Host == "" {
		return "", MalformattedURLErr
	}
	url, err := db.InsertShortURL(s.sqliteDB, longURL)
	if err != nil {
		s.Logger.Error("unable to create new short URL", "error", err)
		return "", fmt.Errorf("error creating new short URL - %w", err)
	}
	s.Logger.Info("new shorturl added",
		"shortcode", url.ShortCode,
		"original_url", url.OriginalURL,
		"id", url.ID,
	)
	return url.ShortCode, nil
}

func (s *ShortenService) ResolveShortURL(code string) (string, error) {
	url, err := db.GetShortURL(s.sqliteDB, code)
	if err != nil {
		return "", fmt.Errorf("error fetching url destination - %w", err)
	}
	if url == nil {
		s.Logger.Error("unable to redirect user", "error", "short code does not exist")
		return "", db.ErrNotFound
	}
	s.Logger.Info("redirected user",
		"shortcode", url.ShortCode,
		"original_url", url.OriginalURL,
	)
	return url.OriginalURL, nil
}

func (s *ShortenService) GetURLDetails(code string) (*db.ShortURL, error) {
	url, err := db.GetShortURL(s.sqliteDB, code)
	if err != nil {
		return nil, fmt.Errorf("error fetching url details - %w", err)
	}
	if url == nil {
		s.Logger.Error("unable to get URL details", "error", "short code does not exist")
		return nil, db.ErrNotFound
	}
	s.Logger.Info("url details fetched",
		"id", url.ID,
		"shortcode", url.ShortCode,
		"original_url", url.OriginalURL,
		"created_at", url.CreatedAt,
	)
	return url, nil
}

func (s *ShortenService) DeleteShortURL(code string) error {
	err := db.DeleteShortURL(s.sqliteDB, code)
	if err != nil {
		s.Logger.Error("unable to delete short url", "error", err)
		return fmt.Errorf("error deleting shortURL - %w", err)
	}
	s.Logger.Info("delete operation completed", "shortcode", code)
	return nil
}
