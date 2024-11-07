package server

import (
	"image-processing-service/internal/services/server/util"
	"net/http"
)

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	util.RespondWithText(w, http.StatusOK, "OK")
}
