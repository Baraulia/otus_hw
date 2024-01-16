package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/stretchr/testify/require"
)

var testEvents = []models.Event{
	{
		Header:    "testHeader1",
		EventTime: time.Now(),
	},
	{
		Header:    "testHeader2",
		EventTime: time.Now().Add(time.Hour * 24),
	},
	{
		Header:    "testHeader3",
		EventTime: time.Now().Add(time.Hour * 24 * 2),
	},
	{
		Header:    "testHeader4",
		EventTime: time.Now().Add(time.Hour * 24 * 3),
	},
	{
		Header:    "testHeader5",
		EventTime: time.Now().Add(time.Hour * 24 * 4),
	},
	{
		Header:    "testHeader6",
		EventTime: time.Now().Add(time.Hour * 24 * 5),
	},
	{
		Header:    "testHeader7",
		EventTime: time.Now().Add(time.Hour * 24 * 6),
	},
	{
		Header:    "testHeader8",
		EventTime: time.Now().Add(time.Hour * 24 * 7),
	},
	{
		Header:    "testHeader9",
		EventTime: time.Now().Add(time.Hour * 24 * 8),
	},
	{
		Header:    "testHeader10",
		EventTime: time.Now().Add(time.Hour * 24 * 9),
	},
}

func TestCreateEvent(t *testing.T) {
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	storage := New(logg)

	id, err := storage.CreateEvent(context.Background(), testEvents[0])
	require.NoError(t, err)

	eventFromStorage, ok := storage.repository[id]
	if !ok {
		t.Error("Event not added to the repository")
	}

	require.Equal(t, testEvents[0].Header, eventFromStorage.Header)
}

func TestUpdateEvent(t *testing.T) {
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	storage := New(logg)

	id, err := storage.CreateEvent(context.Background(), testEvents[0])
	require.NoError(t, err)

	eventFromStorageOld, ok := storage.repository[id]
	if !ok {
		t.Error("Event not added to the repository")
	}

	newHeader := "newHeader"
	err = storage.UpdateEvent(context.Background(), models.Event{
		ID:     id,
		Header: newHeader,
	})
	require.NoError(t, err)

	eventFromStorageNew, ok := storage.repository[id]
	if !ok {
		t.Error("Event not added to the repository")
	}

	require.NotEqual(t, eventFromStorageOld.Header, eventFromStorageNew.Header)
	require.Equal(t, eventFromStorageNew.Header, newHeader)
}

func TestGetEvents(t *testing.T) {
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	storage := New(logg)
	ctx := context.Background()

	for _, event := range testEvents {
		_, err := storage.CreateEvent(ctx, event)
		require.NoError(t, err)
	}

	eventsPerDay, err := storage.GetListEventsDuringDay(ctx, time.Now())
	require.NoError(t, err)

	eventsPerWeek, err := storage.GetListEventsDuringFewDays(ctx, time.Now(), 7)
	require.NoError(t, err)

	eventsPerMonth, err := storage.GetListEventsDuringFewDays(ctx, time.Now(), 30)
	require.NoError(t, err)

	require.Equal(t, len(eventsPerDay), 1)
	require.Equal(t, len(eventsPerWeek), 7)
	require.Equal(t, len(eventsPerMonth), 10)
}
