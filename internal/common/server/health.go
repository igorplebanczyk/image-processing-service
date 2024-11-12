package server

import (
	"image-processing-service/internal/common/server/respond"
	"net/http"
)

func health(w http.ResponseWriter, _ *http.Request) {
	respond.WithoutContent(w, http.StatusOK)
}
