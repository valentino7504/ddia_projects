package db

import "time"

type ShortURL struct {
	ID          int64
	ShortCode   string
	OriginalURL string
	CreatedAt   *time.Time
}
