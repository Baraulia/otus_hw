package app

//nolint:depguard
import (
	"context"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/models"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
	Fatal(msg string, fields map[string]interface{})
}

type Storage interface {
	CreateEvent(ctx context.Context, eventDTO models.Event) (string, error)
	UpdateEvent(ctx context.Context, eventDTO models.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetListEventsDuringDay(ctx context.Context, day time.Time) ([]models.Event, error)
	GetListEventsDuringFewDays(ctx context.Context, start time.Time, amountDays int) ([]models.Event, error)
	Close()
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

//nolint:revive
func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
