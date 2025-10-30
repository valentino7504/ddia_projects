package service

import (
	"database/sql"
	"errors"
	"log/slog"
	"testing"
	"time"

	"valentino7504/1-url-shortener/internal/base62"
	"valentino7504/1-url-shortener/internal/db"

	_ "modernc.org/sqlite"
)

type TestCase struct {
	name  string
	input string
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	sqliteDB, err := db.GetConnection(":memory:")
	if err != nil {
		t.Fatalf("failed to open test database - %v", err)
	}
	t.Cleanup(func() {
		_ = sqliteDB.Close()
	})
	return sqliteDB
}

func newLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}

func setupTestService(t *testing.T) *ShortenService {
	t.Helper()
	db := setupTestDB(t)
	svc, err := NewShortenService(db, newLogger())
	if err != nil {
		t.Fatalf("failed to initialize service: %v", err)
	}
	if svc == nil {
		t.Fatal("expected valid pointer to service, got nil")
	}
	return svc
}

func createDummyURLEntries(t *testing.T, s *ShortenService, urls []TestCase) []string {
	t.Helper()
	codes := make([]string, len(urls))
	for i, url := range urls {
		code, err := s.CreateShortURL(url.input)
		if err != nil {
			t.Fatal("error initializing dummy urls")
		}
		codes[i] = code
	}
	return codes
}

func TestCreateShortenService_Success(t *testing.T) {
	sqliteDB := setupTestDB(t)
	svc, err := NewShortenService(sqliteDB, newLogger())
	if err != nil {
		t.Fatalf("error initializing db - %v", err)
	}
	if svc.sqliteDB == nil {
		t.Fatal("failed to initialise db")
	}
	if svc.Logger == nil {
		t.Fatal("failed to initialise logger")
	}
}

func TestCreateShortenService_ClosedDB(t *testing.T) {
	sqliteDB := setupTestDB(t)
	_ = sqliteDB.Close()
	svc, err := NewShortenService(sqliteDB, newLogger())
	if err == nil {
		t.Fatal("expected error for closed DB, got nil")
	}
	if svc != nil {
		t.Fatal("expected nil service for closed DB")
	}
}

func TestCreateNewShortURL_Success(t *testing.T) {
	svc := setupTestService(t)
	testCases := []TestCase{
		{"valid_https", "https://example.com"},
		{"valid_https_with_space", " https://example.com "},
		{"valid_domain_no_scheme", "example.com"},
		{"valid_ipv4", "https://192.168.0.1"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(*testing.T) {
			code, err := svc.CreateShortURL(testCase.input)
			if err != nil {
				t.Fatalf("expected no errors for %s - val %v, got %v", testCase.name, testCase.input, err)
			}
			if code == "" || !base62.IsValidBase62(code) {
				t.Fatalf("expected valid base 62 shortcode, got %v", code)
			}
		})
	}
}

func TestCreateNewShortURL_InvalidURL(t *testing.T) {
	svc := setupTestService(t)
	testCases := []TestCase{
		{"missing_host", "https://"},
		{"bad_scheme", "://example.com"},
		{"no_dot", "https://example"},
		{"garbage", "someRandomString"},
		{"invalid_ip", "http://999.999.999.999"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			code, err := svc.CreateShortURL(testCase.input)
			if err == nil {
				t.Fatalf("expected error for %v - val: %v, got nil", testCase.name, testCase.input)
			}
			if code != "" {
				t.Fatalf("expected no short code for invalid input, got %v", code)
			}
		})
	}
}

func TestResolveShortURL_Success(t *testing.T) {
	s := setupTestService(t)
	testCases := []TestCase{
		{"example.com", "https://www.example.com"},
		{"google.com", "https://www.google.com"},
		{"auxtoria.com", "https://auxtoria.com"},
	}
	shortCodes := createDummyURLEntries(t, s, testCases)
	for i, code := range shortCodes {
		t.Run(testCases[i].name, func(t *testing.T) {
			url, err := s.ResolveShortURL(code)
			if err != nil {
				t.Fatal("error when resolving short codes")
			}
			if url != testCases[i].input {
				t.Fatalf("expected %v from URL resolution, got %v", testCases[i], url)
			}
		})
	}
}

func TestResolveShortURL_NonExistent(t *testing.T) {
	s := setupTestService(t)
	testCases := []TestCase{
		{"abc123", "abc123"},
		{"12345", "12345"},
		{"nonexistent", "nonexistent"},
		{"0", "0"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := s.ResolveShortURL(tc.input)
			if !errors.Is(err, db.ErrNotFound) {
				t.Fatalf("expected not found error, got %v error", err)
			}
			if url != "" {
				t.Fatalf("expected no url, got %v", url)
			}
		})
	}
}

func TestDeleteShortURL(t *testing.T) {
	s := setupTestService(t)
	testCases := []TestCase{
		{"example.com", "https://www.example.com"},
		{"google.com", "https://www.google.com"},
		{"auxtoria.com", "https://auxtoria.com"},
	}
	shortCodes := createDummyURLEntries(t, s, testCases)
	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.DeleteShortURL(shortCodes[i])
			if err != nil {
				t.Fatalf("unexpected error while deleting - %v", err)
			}
			url, err := s.ResolveShortURL(shortCodes[i])
			if !errors.Is(err, db.ErrNotFound) {
				t.Fatalf("expected not found error, got %v error", err)
			}
			if url != "" {
				t.Fatalf("expected no url, got %v", url)
			}
		})
	}
	t.Run("delete_nonexistent_succeeds", func(t *testing.T) {
		err := s.DeleteShortURL("abcdefghij1234")
		if err != nil {
			t.Fatalf("unexpected error while deleting - %v", err)
		}
	})
}

func TestGetURLDetails_Success(t *testing.T) {
	tolerance := time.Second
	s := setupTestService(t)
	testCases := []TestCase{
		{"fedora43", "https://news.itsfoss.com/fedora-43-release/"},
		{"google.com", "https://www.google.com"},
		{"auxtoria.com", "https://auxtoria.com"},
	}
	before := time.Now().Add(-tolerance)
	shortCodes := createDummyURLEntries(t, s, testCases)
	after := time.Now().Add(tolerance)

	for i, code := range shortCodes {
		url, err := s.GetURLDetails(code)
		if errors.Is(err, db.ErrNotFound) {
			t.Fatal("expected entry to exist but got NotFound error")
		}
		if err != nil {
			t.Fatal("unexpected DB error")
		}
		if url == nil {
			t.Fatalf("expected non nil result=%v, got nil", testCases[i].input)
		}
		if url.ID != int64(i+1) {
			t.Fatal("ID missing in ShortURL struct")
		}
		if url.ShortCode != code {
			t.Fatalf("expected shortcode %v, got %v", code, url.ShortCode)
		}
		if url.OriginalURL != testCases[i].input {
			t.Fatalf("expected originalurl %v, got %v", testCases[i].input, url.OriginalURL)
		}
		if url.CreatedAt.Before(before) || url.CreatedAt.After(after) {
			t.Fatalf("CreatedAt %v outside expected range %v to %v", url.CreatedAt, before, after)
		}
	}
}

func TestGetURLDetails_Nonexistent(t *testing.T) {
	s := setupTestService(t)
	testCases := []TestCase{
		{"abc123", "abc123"},
		{"12345", "12345"},
		{"nonexistent", "nonexistent"},
		{"0", "0"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := s.GetURLDetails(tc.input)
			if !errors.Is(err, db.ErrNotFound) {
				t.Fatalf("expected not found error, got %v error", err)
			}
			if url != nil {
				t.Fatalf("expected no url, got %v", url)
			}
		})
	}
}
