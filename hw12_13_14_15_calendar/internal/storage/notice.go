package storage

import "time"

type Notice struct {
	EventID string
	Title   string
	Date    time.Time
	UserID  string
}
