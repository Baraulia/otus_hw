package handlers

import (
	"fmt"
	"net/http"
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
