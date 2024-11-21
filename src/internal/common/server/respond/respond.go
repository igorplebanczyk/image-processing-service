package respond

import (
	"encoding/json"
	"errors"
	"fmt"
	commonerrors "image-processing-service/src/internal/common/errors"
	"net/http"
)

const Version = "1.0.0"

func WithError(w http.ResponseWriter, error error) {
	applyCommonHeaders(w)

	var commonError commonerrors.Error
	ok := errors.As(error, &commonError)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	switch commonError.Type() {
	case commonerrors.InvalidInput:
		http.Error(w, commonError.Error(), http.StatusBadRequest)
	case commonerrors.Unauthorized:
		http.Error(w, commonError.Error(), http.StatusUnauthorized)
	case commonerrors.Forbidden:
		http.Error(w, commonError.Error(), http.StatusForbidden)
	case commonerrors.Internal:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	case commonerrors.Unknown:
		http.Error(w, commonError.Error(), http.StatusInternalServerError)
	}
}

func WithoutContent(w http.ResponseWriter, code int) {
	applyCommonHeaders(w)
	w.WriteHeader(code)
}

func WithJSON(w http.ResponseWriter, code int, payload any) {
	response, err := json.Marshal(payload)
	if err != nil {
		WithError(w, commonerrors.NewInternal("error marshalling response"))
		return
	}

	applyCommonHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		WithError(w, commonerrors.NewInternal("error sending response"))
		return
	}
}

func WithImage(w http.ResponseWriter, code int, imageBytes []byte, imageName string) {
	contentType := http.DetectContentType(imageBytes)

	applyCommonHeaders(w)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", imageName))
	w.WriteHeader(code)

	_, err := w.Write(imageBytes)
	if err != nil {
		WithError(w, commonerrors.NewInternal("error sending image"))
		return
	}
}

func applyCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("API-Version", Version)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "deny")
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self';")
}
