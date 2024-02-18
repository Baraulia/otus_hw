package api

//nolint:depguard
import (
	"context"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
)

//go:generate mockgen -source=serviceInterface.go -destination=mocks/service_mock.go -package=mockservice
type ApplicationInterface interface {
	CreateEvent(ctx context.Context, eventDTO models.Event) (string, error)
	UpdateEvent(ctx context.Context, eventDTO models.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetListEventsDuringDay(ctx context.Context, day time.Time) ([]models.Event, error)
	GetListEventsDuringFewDays(ctx context.Context, start time.Time, amountDays int) ([]models.Event, error)
	CheckReadness(ctx context.Context) (bool, error)
}
