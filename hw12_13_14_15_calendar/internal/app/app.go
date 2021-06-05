package app

import (
	"context"
)

type App struct { // TODO
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

type Storage interface {
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
