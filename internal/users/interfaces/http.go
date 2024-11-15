package interfaces

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"image-processing-service/internal/common/server/respond"
	"image-processing-service/internal/users/application"
	"image-processing-service/internal/users/domain"
	"log/slog"
	"net/http"
)

type UserServer struct {
	service *application.UserService
}

func NewServer(service *application.UserService) *UserServer {
	return &UserServer{service: service}
}

func (s *UserServer) Register(w http.ResponseWriter, r *http.Request) {
	slog.Info("HTTP request", "method", r.Method, "path", r.URL.Path)
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	user, err := s.service.Register(p.Username, p.Email, p.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, http.StatusBadRequest, domain.ErrInternal.Error())
			return
		}
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
	}

	respond.WithJSON(w, http.StatusCreated, response{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	})
}

func (s *UserServer) Info(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type response struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	user, err := s.service.GetUser(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}

func (s *UserServer) Update(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	err = s.service.UpdateUser(userID, p.Username, p.Email)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
			return
		}
		if errors.Is(err, domain.ErrInvalidRequest) {
			slog.Error("HTTP request error", "error", err)
			respond.WithError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
			return
		}
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusBadRequest, domain.ErrValidationFailed.Error())
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (s *UserServer) Delete(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	err := s.service.DeleteUser(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}
