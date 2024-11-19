package interfaces

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"image-processing-service/src/internal/auth/application"
	"image-processing-service/src/internal/auth/domain"
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
		return "", domain.ErrInvalidRequest
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func (s *AuthAPI) UserMiddleware(handler func(uuid.UUID, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := extractTokenFromHeader(r)
		if err != nil {
			slog.Error("HTTP request error", "error", domain.ErrInvalidRequest.Error())
			respond.WithError(w, http.StatusUnauthorized, domain.ErrInvalidRequest.Error())
			return
		}

		userID, err := s.service.Authenticate(token)
		if err != nil {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, http.StatusUnauthorized, domain.ErrInvalidToken.Error())
			return
		}

		slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path, "user_id", userID)
		handler(userID, w, r)
	}
}

func (s *AuthAPI) AdminMiddleware(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := extractTokenFromHeader(r)
		if err != nil {
			slog.Error("HTTP request error", "error", domain.ErrInvalidRequest.Error())
			respond.WithError(w, http.StatusUnauthorized, domain.ErrInvalidRequest.Error())
			return
		}

		userID, err := s.service.AuthenticateAdmin(token)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidToken) {
				slog.Error("HTTP request error", "error", err)
				respond.WithError(w, http.StatusUnauthorized, domain.ErrInvalidToken.Error())
				return
			}
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, http.StatusForbidden, domain.ErrPermissionDenied.Error())
			return
		}

		slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path, "user_id", userID)
		handler(w, r)
	}
}

func (s *AuthAPI) Login(w http.ResponseWriter, r *http.Request) {
	slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path)

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
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	accessToken, refreshToken, err := s.service.Login(p.Username, p.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
			return
		}
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusUnauthorized, domain.ErrInvalidCredentials.Error())
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *AuthAPI) Refresh(w http.ResponseWriter, r *http.Request) {
	slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path)

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
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	accessToken, err := s.service.Refresh(p.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
			return
		}
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusUnauthorized, domain.ErrInvalidToken.Error())
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		AccessToken: accessToken,
	})
}

func (s *AuthAPI) Logout(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
	err := s.service.Logout(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
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
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	err = s.service.AdminLogoutUser(p.UserID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
		return
	}

	respond.WithoutContent(w, http.StatusOK)
}

func (s *AuthAPI) AdminAccess(w http.ResponseWriter, _ *http.Request) {
	respond.WithoutContent(w, http.StatusOK)
}
