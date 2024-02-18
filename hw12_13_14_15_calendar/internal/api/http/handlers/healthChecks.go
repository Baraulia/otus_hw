package handlers

import (
	"net/http"
)

func (h *Handler) livenessProbe(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) readinessProbe(w http.ResponseWriter, req *http.Request) {
	result, err := h.app.CheckReadness(req.Context())

	if err == nil && result {
		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusServiceUnavailable)
}
