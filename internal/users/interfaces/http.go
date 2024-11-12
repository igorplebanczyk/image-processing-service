package interfaces

import (
	"encoding/json"
	"github.com/google/uuid"
	"image-processing-service/internal/common/server/respond"
	"image-processing-service/internal/users/application"
	"net/http"
)

type UserServer struct {
	service *application.UserService
}

func NewServer(service *application.UserService) *UserServer {
	return &UserServer{service: service}
}

func (s *UserServer) Register(w http.ResponseWriter, r *http.Request) {
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
		respond.WithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := s.service.Register(p.Username, p.Email, p.Password)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "error registering user")
		return
	}

	respond.WithJSON(w, http.StatusCreated, response{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	})
}

func (s *UserServer) Info(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
	type response struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	user, err := s.service.GetUser(userID)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "error getting user")
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
		respond.WithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = s.service.UpdateUser(userID, p.Username, p.Email)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "error updating user")
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (s *UserServer) Delete(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
	err := s.service.DeleteUser(userID)
	if err != nil {
		respond.WithError(w, http.StatusBadRequest, "error deleting user")
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}
