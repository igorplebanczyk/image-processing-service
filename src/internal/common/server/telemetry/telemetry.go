package telemetry

import (
	"image-processing-service/src/internal/common/metrics"
	"log/slog"
	"net/http"
	"time"
)

// The Middleware uses responseRecorder to record the status code of the response, which is then logged and added
// to the metrics.
// The Middleware wraps the ServeMux directly, thus wrapping all requests.

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

		slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path)

		start := time.Now()

		next.ServeHTTP(rr, r)

		duration := time.Since(start).Seconds()
		metrics.HttpRequestsTotal.WithLabelValues(r.Method, http.StatusText(rr.statusCode)).Inc()
		metrics.HttpDurationSeconds.WithLabelValues(r.Method, http.StatusText(rr.statusCode)).Observe(duration)

		if rr.statusCode >= 400 {
			metrics.HttpErrorsTotal.WithLabelValues(http.StatusText(rr.statusCode)).Inc()
			slog.Error(
				"HTTP request failed",
				"method", r.Method,
				"path", r.URL.Path,
				"status_code", rr.statusCode,
			)
		}
	}
}
