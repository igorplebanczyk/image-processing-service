package interfaces

import (
	"encoding/json"
	"github.com/google/uuid"
	"image-processing-service/internal/auth/application"
	"image-processing-service/internal/common/server/respond"
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
			respond.WithError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := s.service.Authenticate(bearerToken)
		if err != nil {
			respond.WithError(w, http.StatusUnauthorized, "token expired or invalid")
			return
		}

		handler(userID, w, r)
	}
}

func (s *AuthServer) Login(w http.ResponseWriter, r *http.Request) {
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
		respond.WithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	accessToken, refreshToken, err := s.service.Login(p.Username, p.Password)
	if err != nil {
		respond.WithError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *AuthServer) Refresh(w http.ResponseWriter, r *http.Request) {
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
		respond.WithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	accessToken, err := s.service.Refresh(p.RefreshToken)
	if err != nil {
		respond.WithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		AccessToken: accessToken,
	})
}

func (s *AuthServer) Logout(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	err := s.service.Logout(userID)
	if err != nil {
		respond.WithError(w, http.StatusInternalServerError, "error logging out")
		return
	}

	respond.WithoutContent(w, http.StatusOK)
}
