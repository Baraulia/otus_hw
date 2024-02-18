package app

//nolint:depguard
import (
	"context"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
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
	CheckReadness(ctx context.Context) (bool, error)
	Close()
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, dto models.Event) (string, error) {
	return a.storage.CreateEvent(ctx, dto)
}

func (a *App) UpdateEvent(ctx context.Context, eventDTO models.Event) error {
	return a.storage.UpdateEvent(ctx, eventDTO)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetListEventsDuringDay(ctx context.Context, day time.Time) ([]models.Event, error) {
	return a.storage.GetListEventsDuringDay(ctx, day)
}

func (a *App) GetListEventsDuringFewDays(ctx context.Context, start time.Time, amountDays int) ([]models.Event, error) {
	return a.storage.GetListEventsDuringFewDays(ctx, start, amountDays)
}

func (a *App) CheckReadness(ctx context.Context) (bool, error) {
	return a.storage.CheckReadness(ctx)
}
