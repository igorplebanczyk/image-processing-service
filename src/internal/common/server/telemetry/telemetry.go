package telemetry

import (
	"image-processing-service/src/internal/common/logs"
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

func Middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rr := &responseRecorder{w, http.StatusOK}

		slog.Info("HTTP request", "type", logs.HTTP, "method", r.Method, "path", r.URL.Path)

		next.ServeHTTP(rr, r)

		if rr.statusCode >= 400 {
			slog.Error(
				"HTTP request failed",
				"type", logs.HTTPErr,
				"method", r.Method,
				"path", r.URL.Path,
				"status_code", rr.statusCode,
			)
		}
	}
}
