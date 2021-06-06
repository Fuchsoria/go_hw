package storage

type Event struct {
	ID            string
	Title         string
	Date          int64
	DurationUntil int64
	Description   string
	OwnerID       string
	NoticeBefore  int64
}
