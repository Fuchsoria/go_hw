package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(ctx context.Context, connectionString string) (*Storage, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("cannot open db, %w", err)
	}

	return &Storage{db}, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("cannot connect to db, %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddEvent(event storage.Event) error {
	_, err := s.db.NamedExec("INSERT INTO events (id, title, date, duration_until, description, owner_id, notice_before) VALUES (:id, :title, :date, :duration_until, :description, :owner_id, :notice_before)", &event)

	return err
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	_, err := s.db.NamedExec("UPDATE events SET title=:title, date=:date, duration_until=:duration_until, description=:description, notice_before=:notice_before WHERE id = :id", &event)

	return err
}

func (s *Storage) RemoveEvent(eventID string) error {
	result, err := s.db.Exec("DELETE FROM events WHERE id=$1", eventID)

	if count, _ := result.RowsAffected(); count == 0 {
		return fmt.Errorf("events were not found, %w", storage.ErrNotFound)
	}

	return err
}

func (s *Storage) RemoveOldEvents(ts int64) error {
	_, err := s.db.Exec("DELETE FROM events WHERE date<$1", ts)

	return err
}

func (s *Storage) eventsFromTo(fromTS int64, toTS int64) ([]storage.Event, error) {
	var result []storage.Event

	err := s.db.Select(&result, "SELECT * FROM events WHERE date BETWEEN $1 AND $2", fromTS, toTS)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) DailyEvents(date time.Time) ([]storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	from := start.Unix()
	to := start.AddDate(0, 0, 1).Unix()

	result, err := s.eventsFromTo(from, to)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) WeeklyEvents(date time.Time) ([]storage.Event, error) {
	offset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	week := date.Add(time.Duration(offset*24) * time.Hour)
	start := time.Date(week.Year(), week.Month(), week.Day(), 0, 0, 0, 0, week.Location())
	from := start.Unix()
	to := start.AddDate(0, 0, 7).Unix()

	result, err := s.eventsFromTo(from, to)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) MonthEvents(date time.Time) ([]storage.Event, error) {
	current := time.Date(date.Year(), date.Month(), 0, 0, 0, 0, 0, date.Location())
	from := current.Unix()
	to := current.AddDate(0, 1, 0).Unix()

	result, err := s.eventsFromTo(from, to)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Storage) GetEventsForNotification() ([]storage.Event, error) {
	var result []storage.Event

	current := time.Now()
	fromTS := current.Add(-(time.Second * 30)).Unix()
	toTS := current.Add(time.Second * 30).Unix()

	err := s.db.Select(&result, "SELECT * FROM events WHERE notice_before != -1 AND date-notice_before BETWEEN $1 AND $2", fromTS, toTS)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, err = s.db.Exec("UPDATE events SET notice_before=-1 WHERE notice_before != -1 AND date-notice_before BETWEEN $1 AND $2", fromTS, toTS)
	if err != nil {
		return nil, err
	}

	return result, nil
}
