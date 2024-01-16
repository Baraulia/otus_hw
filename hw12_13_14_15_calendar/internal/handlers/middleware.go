package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/app"
)

func loggingMiddleware(logger app.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			lrw := NewLoggingResponseWriter(w)
			next.ServeHTTP(lrw, r)
			timeFormatted := startTime.Format("[02/Jan/2006:15:04:05 -0700]")

			logger.Info(fmt.Sprintf("%s %s %s %s %s %d %s %s",
				r.RemoteAddr, timeFormatted, r.Method,
				r.URL.Path, r.Proto, lrw.statusCode,
				time.Since(startTime).String(), r.UserAgent()),
				nil)
		})
	}
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
