package service

import (
	"database/sql"
	"log/slog"
	"testing"

	"valentino7504/1-url-shortener/internal/base62"
	"valentino7504/1-url-shortener/internal/db"

	_ "modernc.org/sqlite"
)

type TestCase struct {
	name    string
	input   string
	wantErr bool
}

func setupTestDB(t *testing.T) *sql.DB {
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
	sqliteDB := setupTestDB(t)
	svc, err := NewShortenService(sqliteDB, newLogger())
	if err != nil {
		t.Fatal("error initializing db")
	}
	if svc == nil {
		t.Fatal("expected valid pointer to service, got nil")
	}
	testCases := []TestCase{
		{"valid_https", "https://example.com", false},
		{"valid_domain_no_scheme", "example.com", false},
		{"valid_ipv4", "https://192.168.0.1", false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(*testing.T) {
			code, err := svc.CreateShortURL(testCase.input)
			if err != nil && !testCase.wantErr {
				t.Fatalf("expected no errors for %s - val %v, got %v", testCase.name, testCase.input, err)
			}
			if code == "" || !base62.IsValidBase62(code) {
				t.Fatalf("expected valid base 62 shortcode, got %v", code)
			}
		})
	}
}

func TestCreateNewShortURL_InvalidURL(t *testing.T) {
	sqliteDB := setupTestDB(t)
	svc, err := NewShortenService(sqliteDB, newLogger())
	if err != nil {
		t.Fatal("error initializing db")
	}
	if svc == nil {
		t.Fatal("expected valid pointer to service, got nil")
	}
	testCases := []TestCase{
		{"missing_host", "https://", true},
		{"bad_scheme", "://example.com", true},
		{"no_dot", "https://example", true},
		{"garbage", "someRandomString", true},
		{"invalid_ip", "http://999.999.999.999", true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			code, err := svc.CreateShortURL(testCase.input)
			if err == nil && testCase.wantErr {
				t.Fatalf("expected error for %v - val: %v, got nil", testCase.name, testCase.input)
			}
			if code != "" {
				t.Fatalf("expected no short code for invalid input, got %v", code)
			}
		})
	}
}
