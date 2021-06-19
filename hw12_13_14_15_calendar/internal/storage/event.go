package storage

type Event struct {
	ID            string `db:"id"`
	Title         string `db:"title"`
	Date          int64  `db:"date"`
	DurationUntil int64  `db:"duration_until"`
	Description   string `db:"description"`
	OwnerID       string `db:"owner_id"`
	NoticeBefore  int64  `db:"notice_before"`
}
