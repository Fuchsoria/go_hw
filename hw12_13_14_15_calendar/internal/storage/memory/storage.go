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
	s.store[event.ID] = event

	return nil
}

func (s *Storage) UpdateEvent(eventId string, event storage.Event) error {
	s.store[eventId] = event

	return nil
}

func (s *Storage) RemoveEvent(eventId string) error {
	delete(s.store, eventId)

	return nil
}

func (s *Storage) DailyEvents(date time.Time) []storage.Event {
	result := []storage.Event{}

	for _, event := range s.store {
		eventDate := time.Unix(event.Date, 0)

		if eventDate.Year() == date.Year() && eventDate.Month() == date.Month() && eventDate.Day() == date.Day() {
			result = append(result, event)
		}
	}

	return result
}

func (s *Storage) WeeklyEvents(date time.Time) []storage.Event {
	result := []storage.Event{}

	for _, event := range s.store {
		eventDate := time.Unix(event.Date, 0)
		eYear, eWeek := eventDate.ISOWeek()
		cYear, cWeek := date.ISOWeek()

		if eYear == cYear && eWeek == cWeek {
			result = append(result, event)
		}
	}

	return result
}

func (s *Storage) MonthEvents(date time.Time) []storage.Event {
	result := []storage.Event{}

	for _, event := range s.store {
		eventDate := time.Unix(event.Date, 0)
		if eventDate.Year() == date.Year() && eventDate.Month() == date.Month() {
			result = append(result, event)
		}
	}

	return result
}

func New() *Storage {
	return &Storage{}
}
