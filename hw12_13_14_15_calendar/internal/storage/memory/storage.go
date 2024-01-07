package memorystorage

//nolint:depguard
import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/app"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/models"
	"github.com/google/uuid"
)

type Storage struct {
	repository map[string]models.Event
	logger     app.Logger
	mu         sync.RWMutex
}

func New(logger app.Logger) *Storage {
	repo := make(map[string]models.Event)
	return &Storage{repository: repo, logger: logger, mu: sync.RWMutex{}}
}

func (s *Storage) Close() {
}

func (s *Storage) CreateEvent(_ context.Context, eventDTO models.Event) (string, error) {
	newUUID := uuid.New()
	s.mu.Lock()
	s.repository[newUUID.String()] = eventDTO
	s.mu.Unlock()
	s.logger.Info("event was created", map[string]interface{}{"id": newUUID})

	return newUUID.String(), nil
}

func (s *Storage) UpdateEvent(_ context.Context, eventDTO models.Event) error {
	if eventDTO.ID == "" {
		s.logger.Error("event id is required parameter", nil)
		return fmt.Errorf("event id is required parameter")
	}
	_, ok := s.repository[eventDTO.ID]
	if !ok {
		s.logger.Error("event with such an id is does not exist", map[string]interface{}{"id": eventDTO.ID})
		return fmt.Errorf("event with such an id is does not exist")
	}

	s.mu.Lock()
	s.repository[eventDTO.ID] = eventDTO
	s.mu.Unlock()
	s.logger.Info("event was updated", map[string]interface{}{"id": eventDTO.ID})

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	_, ok := s.repository[id]
	if !ok {
		s.logger.Error("event with such an id is does not exist", map[string]interface{}{"id": id})
		return fmt.Errorf("event with such an id is does not exist")
	}

	s.mu.Lock()
	delete(s.repository, id)
	s.mu.Unlock()
	s.logger.Info("event was deleted", map[string]interface{}{"id": id})

	return nil
}

func (s *Storage) GetListEventsDuringDay(_ context.Context, targetDay time.Time) ([]models.Event, error) {
	events := make([]models.Event, 0)
	s.mu.RLock()
	for id, event := range s.repository {
		if event.EventTime.Before(targetDay.Add(24*time.Hour).Truncate(24*time.Hour)) &&
			targetDay.Before(event.EventTime.Add(24*time.Hour).Truncate(24*time.Hour)) {
			event.ID = id
			events = append(events, event)
		}
	}
	s.mu.RUnlock()

	return events, nil
}

func (s *Storage) GetListEventsDuringFewDays(
	_ context.Context, start time.Time, amountDays int,
) ([]models.Event, error) {
	events := make([]models.Event, 0)
	s.mu.RLock()
	for id, event := range s.repository {
		startDay := start.Truncate(24 * time.Hour)
		finishDay := startDay.AddDate(0, 0, amountDays).Truncate(24 * time.Hour)

		if event.EventTime.After(startDay) && event.EventTime.Before(finishDay) {
			event.ID = id
			events = append(events, event)
		}
	}
	s.mu.RUnlock()

	return events, nil
}
