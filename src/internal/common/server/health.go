package server

import (
	"image-processing-service/src/internal/common/server/respond"
	"log/slog"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path)
	respond.WithoutContent(w, http.StatusOK)
}
