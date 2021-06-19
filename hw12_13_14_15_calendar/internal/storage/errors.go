package storage

import "errors"

var (
	ErrDateBusy             = errors.New("date is already busy")
	ErrCantCreateEvent      = errors.New("cannot create event")
	ErrCantUpdateEvent      = errors.New("cannot update event")
	ErrCantRemoveEvent      = errors.New("cannot remove event")
	ErrCantConnectToStorage = errors.New("cannot connect to storage")
)
