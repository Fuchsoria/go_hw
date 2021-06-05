package storage

import "time"

type Event struct {
	ID            string
	Title         string
	Date          time.Time
	DurationUntil time.Time
	Description   string
	OwnerID       string
	NoticeBefore  time.Duration
}
