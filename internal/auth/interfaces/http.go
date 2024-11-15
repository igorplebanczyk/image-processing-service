package interfaces

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"image-processing-service/internal/auth/application"
	"image-processing-service/internal/auth/domain"
	"image-processing-service/internal/common/server/respond"
	"log/slog"
	"net/http"
	"strings"
)

type AuthServer struct {
	service *application.AuthService
}

func NewServer(authService *application.AuthService) *AuthServer {
	return &AuthServer{
		service: authService,
	}
}

func (s *AuthServer) Middleware(handler func(uuid.UUID, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			slog.Error("HTTP request error", "error", domain.ErrInvalidRequest.Error())
			respond.WithError(w, http.StatusUnauthorized, domain.ErrInvalidRequest.Error())
			return
		}
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := s.service.Authenticate(bearerToken)
		if err != nil {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, http.StatusUnauthorized, domain.ErrInvalidToken.Error())
			return
		}

		slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path)
		handler(userID, w, r)
	}
}

func (s *AuthServer) Login(w http.ResponseWriter, r *http.Request) {
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

func (s *AuthServer) Refresh(w http.ResponseWriter, r *http.Request) {
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

func (s *AuthServer) Logout(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	err := s.service.Logout(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
		return
	}

	respond.WithoutContent(w, http.StatusOK)
}
