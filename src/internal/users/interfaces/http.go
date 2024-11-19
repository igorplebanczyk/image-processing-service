package interfaces

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"image-processing-service/src/internal/common/server/respond"
	"image-processing-service/src/internal/users/application"
	"image-processing-service/src/internal/users/domain"
	"log/slog"
	"net/http"
)

type UserAPI struct {
	service *application.UserService
}

func NewServer(service *application.UserService) *UserAPI {
	return &UserAPI{service: service}
}

func (s *UserAPI) Register(w http.ResponseWriter, r *http.Request) {
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

func (s *UserAPI) GetData(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
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

func (s *UserAPI) Update(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
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

func (s *UserAPI) Delete(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
	err := s.service.DeleteUser(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (s *UserAPI) AdminListAllUsers(w http.ResponseWriter, _ *http.Request) {
	type response struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	users, err := s.service.AdminGetAllUsers()
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
		return
	}

	var resp []response
	for _, user := range users {
		resp = append(resp, response{
			ID:        user.ID.String(),
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role.String(),
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		})
	}

	respond.WithJSON(w, http.StatusOK, resp)
}

func (s *UserAPI) AdminUpdateRole(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID uuid.UUID   `json:"user_id"`
		Role   domain.Role `json:"role"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusBadRequest, domain.ErrInvalidRequest.Error())
		return
	}

	err = s.service.AdminUpdateUserRole(p.UserID, p.Role)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (s *UserAPI) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
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

	err = s.service.DeleteUser(p.UserID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, http.StatusInternalServerError, domain.ErrInternal.Error())
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}
