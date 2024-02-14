package internalhttp

import (
	"net/http"
)

func loggingMiddleware(_ http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		// TODO
	})
}
