package scheduler

//nolint:depguard
import (
	"context"
	"encoding/json"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/mb"
)

type Scheduler struct {
	logger    Logger
	storage   Storage
	producer  mb.ProducerMB
	routeKey  string
	frequency time.Duration
}

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

type Storage interface {
	DeleteOldEvent(ctx context.Context) error
	GetNotifications(ctx context.Context) ([]models.Notification, error)
	Close()
}

func New(logger Logger, storage Storage, producer mb.ProducerMB, frequency time.Duration, routeKey string) *Scheduler {
	return &Scheduler{logger: logger, storage: storage, producer: producer, frequency: frequency, routeKey: routeKey}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("stopping scheduler...", nil)
			return
		case <-ticker.C:
			s.logger.Info("deleting old events...", map[string]interface{}{"time": time.Now()})
			err := s.storage.DeleteOldEvent(ctx)
			if err != nil {
				s.logger.Error("error while deleting old events", map[string]interface{}{"error": err})
			}

			s.logger.Info("getting new notifications...", nil)
			notifications, err := s.storage.GetNotifications(ctx)
			if err != nil {
				s.logger.Error("error while getting new notifications", map[string]interface{}{"error": err})
			}

			if len(notifications) != 0 {
				for _, notification := range notifications {
					data, err := json.Marshal(notification)
					if err != nil {
						s.logger.Error("error while marshaling new notification", map[string]interface{}{"error": err})
					}

					err = s.producer.Publish(s.routeKey, data)
					if err != nil {
						s.logger.Error("error while publishing new notification", map[string]interface{}{"error": err})
					}
				}
			}
		}
	}
}
