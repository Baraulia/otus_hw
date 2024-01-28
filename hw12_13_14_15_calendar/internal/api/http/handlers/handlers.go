package handlers

//nolint:depguard
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/google/uuid"
)

func (h *Handler) helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello"))
	if err != nil {
		h.logger.Error("error while writing response", map[string]interface{}{"handler": "helloHandler", "error": err})
		http.Error(w, fmt.Sprintf("helloHandler: error while writing response:%s", err), 500)
		return
	}
}

func (h *Handler) createEvent(w http.ResponseWriter, req *http.Request) {
	var input models.Event
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Error("Error while decoding request", map[string]interface{}{"error": err})
		http.Error(w, err.Error(), 400)
		return
	}

	id, err := h.app.CreateEvent(req.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("id", id)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) updateEvent(w http.ResponseWriter, req *http.Request) {
	var input models.Event
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&input); err != nil {
		h.logger.Error("Error while decoding request", map[string]interface{}{"error": err})
		http.Error(w, err.Error(), 400)
		return
	}

	id := strings.TrimPrefix(req.URL.Path, "/event/")
	err := uuid.Validate(id)
	if err != nil {
		h.logger.Error("invalid id", map[string]interface{}{"error": err})
		http.Error(w, err.Error(), 400)
		return
	}

	input.ID = id
	err = h.app.UpdateEvent(req.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteEvent(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/event/")
	err := uuid.Validate(id)
	if err != nil {
		h.logger.Error("invalid id", map[string]interface{}{"error": err})
		http.Error(w, err.Error(), 400)
		return
	}

	err = h.app.DeleteEvent(req.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getListEvents(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	startParam := query.Get("start")
	if startParam == "" {
		h.logger.Error("start is required parameter", nil)
		http.Error(w, "start is required parameter", 400)
		return
	}

	start, err := time.Parse(time.DateOnly, startParam)
	if err != nil {
		h.logger.Error("Invalid start parameter", map[string]interface{}{"error": err})
		http.Error(w, fmt.Sprintf("invalid start parameter: %s", err), http.StatusBadRequest)
		return
	}

	amountDaysParam := query.Get("amount_days")
	var events []models.Event
	switch amountDaysParam {
	case "":
		events, err = h.app.GetListEventsDuringDay(req.Context(), start)
		if err != nil {
			http.Error(w, fmt.Sprintf("server error: %s", err), 500)
			return
		}
	default:
		amountDays, err := strconv.Atoi(amountDaysParam)
		if err != nil {
			h.logger.Error("Invalid amount_days parameter", map[string]interface{}{"error": err})
			http.Error(w, "Invalid amount_days parameter", http.StatusBadRequest)
			return
		}

		events, err = h.app.GetListEventsDuringFewDays(req.Context(), start, amountDays)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	output, err := json.Marshal(events)
	if err != nil {
		h.logger.Error("getEventsDuringDay: error while marshaling list of events", map[string]interface{}{"error": err})
		http.Error(w, fmt.Sprintf("getEventsDuringDay: error while marshaling list of events: %s", err), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		h.logger.Error("getEventsDuringDay: error while writing response", map[string]interface{}{"error": err})
		http.Error(w, fmt.Sprintf("getEventsDuringDay: error while writing response:%s", err), 500)
		return
	}
}
