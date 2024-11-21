package telemetry

import (
	"log/slog"
	"net/http"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func TelemetryMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rr := &responseRecorder{w, http.StatusOK}

		slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path)

		next.ServeHTTP(rr, r)

		if rr.statusCode >= 400 {
			slog.Error(
				"HTTP request failed",
				"method", r.Method,
				"path", r.URL.Path,
				"status_code", rr.statusCode,
			)
		}
	}
}
