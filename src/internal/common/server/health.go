package server

import (
	"image-processing-service/src/internal/common/server/respond"
	"net/http"
)

func health(w http.ResponseWriter, _ *http.Request) {
	respond.WithoutContent(w, http.StatusOK)
}
