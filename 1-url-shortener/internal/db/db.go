package db

import (
	"database/sql"
	"errors"
	"fmt"

	"valentino7504/1-url-shortener/internal/base62"

	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("short URL not found")

func GetConnection(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("cannot open db: %w", err)
	}
	return db, nil
}

func InitDB(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("no db pointer passed to InitDB")
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("could not establish connection to the db")
	}

	const query = `CREATE TABLE IF NOT EXISTS short_urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			short_code TEXT UNIQUE NOT NULL,
			original_url TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (CURRENT_TIMESTAMP)
		);`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("could not initialize database - %w", err)
	}
	return nil
}

func GetAllURLs(sqliteDb *sql.DB) ([]*ShortURL, error) {
	const query = `SELECT id, short_code, original_url, created_at
		FROM short_urls;`
	rows, err := sqliteDb.Query(query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching jobs - %w", err)
	}
	var shortURLs []*ShortURL
	for rows.Next() {
		var url ShortURL
		var createdAt string
		err := rows.Scan(
			&url.ID,
			&url.ShortCode,
			&url.OriginalURL,
			&createdAt,
		)
		if err != nil {
			return shortURLs, fmt.Errorf("error parsing url %s", url.OriginalURL)
		}
		url.CreatedAt = ParseDateTime(createdAt)
		shortURLs = append(shortURLs, &url)
	}
	return shortURLs, nil
}

func GetShortURL(sqliteDB *sql.DB, shortCode string) (*ShortURL, error) {
	const query = `SELECT
		id, short_code, original_url, created_at
		FROM short_urls WHERE short_code = ?;
	`
	var shortURL ShortURL
	var createdAt string
	err := sqliteDB.QueryRow(query, shortCode).Scan(
		&shortURL.ID,
		&shortURL.ShortCode,
		&shortURL.OriginalURL,
		&createdAt,
	)
	shortURL.CreatedAt = ParseDateTime(createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unable to fetch shorturl - %w", err)
	}
	return &shortURL, nil
}

func InsertShortURL(sqliteDB *sql.DB, url string) (*ShortURL, error) {
	tx, err := sqliteDB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db transaction for createioN")
	}
	defer func() { _ = tx.Rollback() }()
	const query = `INSERT INTO short_urls
		(short_code, original_url)
		VALUES ("", ?) RETURNING id, created_at;`
	var createdAt string
	var urlID int64
	err = tx.QueryRow(query, url).Scan(&urlID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("error in adding url - %w", err)
	}
	shortCode := base62.EncodeBase62(urlID)
	_, err = tx.Exec("UPDATE short_urls SET short_code = ? WHERE id = ?", shortCode, urlID)
	if err != nil {
		return nil, fmt.Errorf("error in adding url - %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error in committing transaction")
	}
	return &ShortURL{ID: urlID, OriginalURL: url, ShortCode: shortCode, CreatedAt: ParseDateTime(createdAt)}, nil
}

func DeleteShortURL(sqliteDB *sql.DB, shortCode string) error {
	const query = `DELETE FROM short_urls
		WHERE short_code = ?
		RETURNING id;`
	var delID int64
	if err := sqliteDB.QueryRow(query, shortCode).Scan(&delID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("Error deleting record with provided short code")
	}
	return nil
}
