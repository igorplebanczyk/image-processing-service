package util

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	type response struct {
		Error string `json:"error"`
	}

	payload := response{Error: msg}
	resp, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	applySecurityHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func RespondWithoutContent(w http.ResponseWriter, code int) {
	applySecurityHeaders(w)
	w.WriteHeader(code)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error marshalling response")
		return
	}

	applySecurityHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error writing response")
		return
	}
}

func RespondWithText(w http.ResponseWriter, code int, text string) {
	applySecurityHeaders(w)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)

	_, err := w.Write([]byte(text))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error writing response")
		return
	}
}

func RespondWithImage(w http.ResponseWriter, code int, imageBytes []byte, imageName string) {
	contentType := http.DetectContentType(imageBytes)

	applySecurityHeaders(w)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", imageName))
	w.WriteHeader(code)

	_, err := w.Write(imageBytes)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to send image")
		return
	}
}

func applySecurityHeaders(w http.ResponseWriter) {
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "deny")
    w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self';")
}
