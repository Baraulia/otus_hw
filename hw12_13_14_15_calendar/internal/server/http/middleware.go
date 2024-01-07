package internalhttp

import (
	"net/http"
)

func loggingMiddleware(_ http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO
	})
}
