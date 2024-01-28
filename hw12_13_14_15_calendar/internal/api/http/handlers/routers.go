package handlers

//nolint:depguard
import (
	"net/http"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/api"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/gorilla/mux"
)

type Handler struct {
	logger app.Logger
	app    api.ApplicationInterface
}

func NewHandler(logger app.Logger, app api.ApplicationInterface) *Handler {
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
