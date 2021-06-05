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

func (s *Storage) AddEvent(event storage.Event) {
	s.store[event.ID] = event
}

func (s *Storage) UpdateEvent(eventId string, event storage.Event) {
	s.store[eventId] = event
}

func (s *Storage) RemoveEvent(eventId string) {
	delete(s.store, eventId)
}

func (s *Storage) DailyEvents(date time.Time) []storage.Event {
	result := []storage.Event{}

	for _, event := range s.store {
		if event.Date.Format("02-01-2006") == date.Format("02-01-2006") {
			result = append(result, event)
		}
	}

	return result
}

func (s *Storage) WeeklyEvents(date time.Time) []storage.Event {
	result := []storage.Event{}

	for _, event := range s.store {
		eYear, eWeek := event.Date.ISOWeek()
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
		if event.Date.Format("01-2006") == date.Format("01-2006") {
			result = append(result, event)
		}
	}

	return result
}

func New() *Storage {
	return &Storage{}
}

// TODO
