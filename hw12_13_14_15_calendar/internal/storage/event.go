package storage

type Event struct {
	ID            string `db:"id" json:"id"`
	Title         string `db:"title" json:"title"`
	Date          int64  `db:"date" json:"date"`
	DurationUntil int64  `db:"duration_until" json:"duration_until"`
	Description   string `db:"description" json:"description"`
	OwnerID       string `db:"owner_id" json:"owner_id"`
	NoticeBefore  int64  `db:"notice_before" json:"notice_before"`
}
