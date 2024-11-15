package server

import (
	"image-processing-service/internal/common/server/respond"
	"log/slog"
	"net/http"
)

func health(w http.ResponseWriter, _ *http.Request) {
	respond.WithoutContent(w, http.StatusOK)
	slog.Info("Server health check OK")
}
