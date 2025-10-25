package db

import "time"

func ParseDateTime(strTime string) *time.Time {
	t, _ := time.Parse(time.DateTime, strTime)
	return &t
}

func FormatDateTime(t *time.Time) string {
	return t.Format(time.DateTime)
}
