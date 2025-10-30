package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"regexp"
	"strings"

	"valentino7504/1-url-shortener/internal/db"
)

type ShortenService struct {
	sqliteDB *sql.DB
	Logger   *slog.Logger
}

var ErrMalformattedURL error = errors.New("invalid url provided")

var (
	domainRe = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	ipv4Re   = regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])$`)
)

func NewShortenService(sqliteDB *sql.DB, logger *slog.Logger) (*ShortenService, error) {
	svc := ShortenService{sqliteDB: sqliteDB, Logger: logger}
	err := db.InitDB(sqliteDB)
	if err != nil {
		svc.Logger.Error("unable to initialize db", "error", err)
		return nil, err
	}
	return &svc, nil
}

func (s *ShortenService) CreateShortURL(longURL string) (string, error) {
	longURL = strings.TrimSpace(longURL)
	u, err := url.Parse(longURL)
	if err != nil {
		return "", ErrMalformattedURL
	}
	if u.Scheme == "" {
		longURL = "https://" + longURL
		u, err = url.Parse(longURL)
		if err != nil {
			return "", ErrMalformattedURL
		}
	}
	if !domainRe.MatchString(u.Host) && !ipv4Re.MatchString(u.Host) {
		return "", ErrMalformattedURL
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
		if url == nil {
			return "", db.ErrNotFound
		}
		return "", fmt.Errorf("error fetching url destination - %w", err)
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
		if url == nil {
			return nil, db.ErrNotFound
		}
		return nil, fmt.Errorf("error fetching url details - %w", err)
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
