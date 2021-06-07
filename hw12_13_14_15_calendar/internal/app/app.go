package app

import (
	"context"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
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
	UpdateEvent(eventId string, event storage.Event) error
	RemoveEvent(eventId string) error
	DailyEvents(date time.Time) []storage.Event
	WeeklyEvents(date time.Time) []storage.Event
	MonthEvents(date time.Time) []storage.Event
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

func (a *App) GetLogger() Logger {
	return a.logger
}

// TODO
