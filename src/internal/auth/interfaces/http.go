package interfaces

import (
	"encoding/json"
	"github.com/google/uuid"
	"image-processing-service/src/internal/auth/application"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/common/logs"
	"image-processing-service/src/internal/common/server/respond"
	"log/slog"
	"net/http"
	"strings"
)

type AuthAPI struct {
	service *application.AuthService
}

func NewServer(authService *application.AuthService) *AuthAPI {
	return &AuthAPI{
		service: authService,
	}
}

func extractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", commonerrors.NewInvalidInput("missing or invalid Authorization header")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func (s *AuthAPI) UserMiddleware(handler func(uuid.UUID, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := extractTokenFromHeader(r)
		if err != nil {
			slog.Error("HTTP request error", "type", logs.Error, "error", err)
			respond.WithError(w, err)
			return
		}

		userID, err := s.service.Authenticate(token)
		if err != nil {
			slog.Error("HTTP request error", "type", logs.Error, "error", err)
			respond.WithError(w, err)
			return
		}

		handler(userID, w, r)
	}
}

func (s *AuthAPI) AdminMiddleware(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := extractTokenFromHeader(r)
		if err != nil {
			slog.Error("HTTP request error", "type", logs.Error, "error", err)
			respond.WithError(w, err)
			return
		}

		_, err = s.service.AuthenticateAdmin(token)
		if err != nil {
			slog.Error("HTTP request error", "type", logs.Error, "error", err)
			respond.WithError(w, err)
			return
		}

		handler(w, r)
	}
}

func (s *AuthAPI) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "type", logs.Error, "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	accessToken, refreshToken, err := s.service.Login(p.Username, p.Password)
	if err != nil {
		slog.Error("HTTP request error", "type", logs.Error, "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *AuthAPI) Refresh(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RefreshToken string `json:"refresh_token"`
	}

	type response struct {
		AccessToken string `json:"access_token"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "type", logs.Error, "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	accessToken, err := s.service.Refresh(p.RefreshToken)
	if err != nil {
		slog.Error("HTTP request error", "type", logs.Error, "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		AccessToken: accessToken,
	})
}

func (s *AuthAPI) Logout(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
	err := s.service.Logout(userID)
	if err != nil {
		slog.Error("HTTP request error", "type", logs.Error, "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusOK)
}

func (s *AuthAPI) AdminLogoutUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID uuid.UUID `json:"user_id"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "type", logs.Error, "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = s.service.AdminLogoutUser(p.UserID)
	if err != nil {
		slog.Error("HTTP request error", "type", logs.Error, "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusOK)
}

func (s *AuthAPI) AdminAccess(w http.ResponseWriter, _ *http.Request) {
	respond.WithoutContent(w, http.StatusOK)
}
