package server

import (
	"image-processing-service/src/internal/common/server/respond"
	"net/http"
)

// A simple health check endpoint. If it can be reached, it will always return 200 OK. Meaning that it only fails if
// the server is down.

func health(w http.ResponseWriter, _ *http.Request) {
	respond.WithoutContent(w, http.StatusOK)
}
