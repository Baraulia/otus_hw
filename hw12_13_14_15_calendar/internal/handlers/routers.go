package handlers

import (
	"context"
	"net/http"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/gorilla/mux"
)

type ApplicationInterface interface {
	CreateEvent(ctx context.Context, id, title string) error
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

	return r
}
