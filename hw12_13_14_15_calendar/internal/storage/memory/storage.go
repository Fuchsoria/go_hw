package memorystorage

import (
	"sync"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	store map[string]storage.Event
	mu    sync.RWMutex
}

func (s *Storage) AddEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[event.ID] = event

	return nil
}

func (s *Storage) UpdateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[event.ID] = event

	return nil
}

func (s *Storage) RemoveEvent(eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, eventID)

	return nil
}

func (s *Storage) DailyEvents(date time.Time) ([]storage.Event, error) {
	var result []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.store {
		eventDate := time.Unix(event.Date, 0)

		if eventDate.Year() == date.Year() && eventDate.Month() == date.Month() && eventDate.Day() == date.Day() {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) WeeklyEvents(date time.Time) ([]storage.Event, error) {
	var result []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.store {
		eventDate := time.Unix(event.Date, 0)
		eYear, eWeek := eventDate.ISOWeek()
		cYear, cWeek := date.ISOWeek()

		if eYear == cYear && eWeek == cWeek {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) MonthEvents(date time.Time) ([]storage.Event, error) {
	var result []storage.Event

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, event := range s.store {
		eventDate := time.Unix(event.Date, 0)
		if eventDate.Year() == date.Year() && eventDate.Month() == date.Month() {
			result = append(result, event)
		}
	}

	return result, nil
}

func New() *Storage {
	return &Storage{store: map[string]storage.Event{}}
}
