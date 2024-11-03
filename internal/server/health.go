package server

import (
	"image-processing-service/internal/server/util"
	"net/http"
)

func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("OK"))
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "error writing response")
	}
}
