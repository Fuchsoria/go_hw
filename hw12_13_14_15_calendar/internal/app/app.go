package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

const ownerID string = "main_owner_id"

var ErrCantFindID = errors.New("cannot find event id")

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type Storage interface {
	AddEvent(event storage.Event) error
	UpdateEvent(event storage.Event) error
	RemoveEvent(eventID string) error
	DailyEvents(date time.Time) ([]storage.Event, error)
	WeeklyEvents(date time.Time) ([]storage.Event, error)
	MonthEvents(date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(title string, date int64) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("cannot create unique id, %w", err)
	}

	event := storage.Event{ID: id.String(), Title: title, Date: date, OwnerID: ownerID}

	return a.storage.AddEvent(event)
}

func (a *App) UpdateEvent(event storage.Event) error {
	if event.ID == "" {
		return fmt.Errorf("cannot update event, %w", ErrCantFindID)
	}

	return a.storage.UpdateEvent(event)
}

func (a *App) RemoveEvent(id string) error {
	return a.storage.RemoveEvent(id)
}

func (a *App) DailyEvents(date time.Time) ([]storage.Event, error) {
	return a.storage.DailyEvents(date)
}

func (a *App) WeeklyEvents(date time.Time) ([]storage.Event, error) {
	return a.storage.WeeklyEvents(date)
}

func (a *App) MonthEvents(date time.Time) ([]storage.Event, error) {
	return a.storage.MonthEvents(date)
}

func (a *App) GetLogger() Logger {
	return a.logger
}
