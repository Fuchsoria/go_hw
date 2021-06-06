package sqlstorage

import (
	"context"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct { // TODO
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", "user=root dbname=calendar sslmode=disable")
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddEvent(event storage.Event) error {

	return nil
}

func (s *Storage) UpdateEvent(eventId string, event storage.Event) error {
	return nil
}

func (s *Storage) RemoveEvent(eventId string) error {

	return nil
}

func (s *Storage) DailyEvents(date time.Time) []storage.Event {
	result := []storage.Event{}

	return result
}

func (s *Storage) WeeklyEvents(date time.Time) []storage.Event {
	result := []storage.Event{}

	return result
}

func (s *Storage) MonthEvents(date time.Time) []storage.Event {
	result := []storage.Event{}

	return result
}
