package handlers

import (
	"context"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
	"net/http"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/gorilla/mux"
)

//go:generate mockgen -source=routers.go -destination=mocks/service_mock.go
type ApplicationInterface interface {
	CreateEvent(ctx context.Context, eventDTO models.Event) (string, error)
	UpdateEvent(ctx context.Context, eventDTO models.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetListEventsDuringDay(ctx context.Context, day time.Time) ([]models.Event, error)
	GetListEventsDuringFewDays(ctx context.Context, start time.Time, amountDays int) ([]models.Event, error)
}

type Handler struct {
	logger app.Logger
	app    ApplicationInterface
}

func NewHandler(logger app.Logger, app ApplicationInterface) *Handler {
	return &Handler{logger: logger, app: app}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware(h.logger))

	r.HandleFunc("/hello", h.helloHandler).Methods(http.MethodGet)

	r.HandleFunc("/event", h.createEvent).Methods(http.MethodPost)
	r.HandleFunc("/event/{id}", h.updateEvent).Methods(http.MethodPut)
	r.HandleFunc("/event/{id}", h.deleteEvent).Methods(http.MethodDelete)
	r.HandleFunc("/event/list", h.getListEvents).Methods(http.MethodGet)

	return r
}
