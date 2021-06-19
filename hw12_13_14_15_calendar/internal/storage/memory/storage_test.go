package memorystorage

import (
	"testing"
	"time"

	"github.com/Fuchsoria/go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	const ownerID = "test_owner_id"
	event := storage.Event{
		ID:            "id",
		Title:         "test title",
		Date:          time.Now().Unix(),
		DurationUntil: time.Now().AddDate(0, 1, 0).Unix(),
		Description:   "test desc",
		OwnerID:       ownerID,
		NoticeBefore:  0,
	}

	t.Run("test add event", func(t *testing.T) {
		id, err := uuid.NewRandom()
		require.NoError(t, err, "should not have error in uuid")

		event := event
		event.ID = id.String()

		s := New()
		require.NotNil(t, s, "storage should not be nil")

		err = s.AddEvent(event)
		require.NoError(t, err, "should not have error in add event method")
		require.Len(t, s.store, 1, "length should be 1, because of added event")
	})

	t.Run("test update event", func(t *testing.T) {
		s := New()

		err := s.AddEvent(event)
		require.NoError(t, err, "should not have error in add event method")

		newEvent := event
		newEvent.Title = "NEW TITLE"

		err = s.UpdateEvent(newEvent)
		require.NoError(t, err, "should not have error in update event method")

		require.Equal(t, "NEW TITLE", s.store["id"].Title, "title should be updated")
	})

	t.Run("test delete event", func(t *testing.T) {
		s := New()

		err := s.AddEvent(event)
		require.NoError(t, err, "should not have error in add event method")

		require.NotEmpty(t, s.store["id"], "event should exist")

		err = s.RemoveEvent(event.ID)
		require.NoError(t, err, "should not have error in remove event method")

		require.Empty(t, s.store["id"], "event should be removed")
	})

	t.Run("test lists", func(t *testing.T) {
		s := New()

		err := s.AddEvent(event)
		require.NoError(t, err, "should not have error in add event method")

		now := time.Now()
		oneYearLater := now.AddDate(1, 0, 0)

		// now.
		daily, err := s.DailyEvents(now)
		require.NoError(t, err, "should not have error in DailyEvents method")
		require.Len(t, daily, 1, "DailyEvents should not be empty")

		weekly, err := s.WeeklyEvents(now)
		require.NoError(t, err, "should not have error in WeeklyEvents method")
		require.Len(t, weekly, 1, "WeeklyEvents should not be empty")

		month, err := s.MonthEvents(now)
		require.NoError(t, err, "should not have error in MonthEvents method")
		require.Len(t, month, 1, "MonthEvents should not be empty")

		// one year later.
		daily, err = s.DailyEvents(oneYearLater)
		require.NoError(t, err, "should not have error in DailyEvents method")
		require.Empty(t, daily, "DailyEvents should be empty")

		weekly, err = s.WeeklyEvents(oneYearLater)
		require.NoError(t, err, "should not have error in WeeklyEvents method")
		require.Empty(t, weekly, "WeeklyEvents should be empty")

		month, err = s.MonthEvents(oneYearLater)
		require.NoError(t, err, "should not have error in MonthEvents method")
		require.Empty(t, month, "MonthEvents should be empty")
	})
}
