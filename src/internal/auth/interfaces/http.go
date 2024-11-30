package interfaces

import (
	"encoding/json"
	"github.com/google/uuid"
	"image-processing-service/src/internal/auth/application"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/common/server/respond"
	"log/slog"
	"net/http"
	"time"
)

type AuthAPI struct {
	service            *application.AuthService
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewAPI(authService *application.AuthService, accessTokenExpiry, refreshTokenExpiry time.Duration) *AuthAPI {
	return &AuthAPI{
		service:            authService,
		accessTokenExpiry:  accessTokenExpiry,
		refreshTokenExpiry: refreshTokenExpiry,
	}
}

func (a *AuthAPI) UserMiddleware(handler func(uuid.UUID, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getAccessTokenFromCookie(r)
		if err != nil {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, err)
			return
		}

		userID, err := a.service.Authenticate(token)
		if err != nil {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, err)
			return
		}

		handler(userID, w, r)
	}
}

func (a *AuthAPI) AdminMiddleware(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getAccessTokenFromCookie(r)
		if err != nil {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, err)
			return
		}

		_, err = a.service.AuthenticateAdmin(token)
		if err != nil {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, err)
			return
		}

		handler(w, r)
	}
}

func (a *AuthAPI) LoginOne(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = a.service.LoginOne(p.Username, p.Password)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusOK)
}

func (a *AuthAPI) LoginTwo(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		OTP      string `json:"otp"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	accessToken, refreshToken, err := a.service.LoginTwo(p.Username, p.OTP)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	setAccessTokenInCookie(w, accessToken, a.accessTokenExpiry)
	setRefreshTokenInCookie(w, refreshToken, a.refreshTokenExpiry)

	respond.WithoutContent(w, http.StatusOK)
}

func (a *AuthAPI) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := getRefreshTokenFromCookie(r)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	accessToken, err := a.service.Refresh(refreshToken)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	setAccessTokenInCookie(w, accessToken, a.accessTokenExpiry)

	respond.WithoutContent(w, http.StatusOK)
}

func (a *AuthAPI) Logout(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
	err := a.service.Logout(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	setAccessTokenInCookie(w, "", -time.Hour)
	setRefreshTokenInCookie(w, "", -time.Hour)

	respond.WithoutContent(w, http.StatusOK)
}

func (a *AuthAPI) AdminAccess(w http.ResponseWriter, _ *http.Request) {
	respond.WithoutContent(w, http.StatusOK)
}

func (a *AuthAPI) AdminLogoutUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid user ID"))
		return
	}

	err = a.service.AdminLogoutUser(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusOK)
}
