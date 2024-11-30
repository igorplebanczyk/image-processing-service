package interfaces

import (
	"encoding/json"
	"github.com/google/uuid"
	commonerrors "image-processing-service/src/internal/common/errors"
	"image-processing-service/src/internal/common/server/respond"
	"image-processing-service/src/internal/users/application"
	"image-processing-service/src/internal/users/domain"
	"log/slog"
	"net/http"
	"strconv"
)

type UserAPI struct {
	service *application.UserService
}

func NewAPI(service *application.UserService) *UserAPI {
	return &UserAPI{service: service}
}

func (a *UserAPI) Register(w http.ResponseWriter, r *http.Request) {
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
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	user, err := a.service.Register(p.Username, p.Email, p.Password)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithJSON(w, http.StatusCreated, response{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
	})
}

func (a *UserAPI) GetDetails(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
	type response struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		Verified  string `json:"verified"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	user, err := a.service.GetDetails(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.String(),
		Verified:  strconv.FormatBool(user.Verified),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}

func (a *UserAPI) UpdateDetails(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = a.service.UpdateDetails(userID, p.Username, p.Email)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *UserAPI) Delete(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
	err := a.service.Delete(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *UserAPI) ResendVerificationCode(userID uuid.UUID, w http.ResponseWriter, _ *http.Request) {
	err := a.service.ResendVerificationCode(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *UserAPI) Verify(userID uuid.UUID, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		OTP string `json:"otp"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = a.service.Verify(userID, p.OTP)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *UserAPI) SendForgotPasswordCode(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = a.service.SendForgotPasswordCode(p.Email)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *UserAPI) ResetPassword(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = a.service.ResetPassword(p.Email, p.OTP, p.NewPassword)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *UserAPI) AdminGetUserDetails(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		Verified  string `json:"verified"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid user ID"))
		return
	}

	user, err := a.service.GetDetails(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithJSON(w, http.StatusOK, response{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.String(),
		Verified:  strconv.FormatBool(user.Verified),
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}

func (a *UserAPI) AdminGetAllUsersDetails(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		Verified  string `json:"verified"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid page"))
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid limit"))
		return
	}

	users, err := a.service.AdminGetAllUsers(page, limit)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	var resp []response
	for _, user := range users {
		resp = append(resp, response{
			ID:        user.ID.String(),
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role.String(),
			Verified:  strconv.FormatBool(user.Verified),
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		})
	}

	respond.WithJSON(w, http.StatusOK, resp)
}

func (a *UserAPI) AdminUpdateRole(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid user ID"))
		return
	}

	roleStr := r.URL.Query().Get("role")
	role := domain.Role(roleStr)

	err = a.service.AdminUpdateUserRole(userID, role)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *UserAPI) AdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid user ID"))
		return
	}

	err = a.service.Delete(userID)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}

func (a *UserAPI) AdminBroadcast(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, commonerrors.NewInvalidInput("invalid body"))
		return
	}

	err = a.service.AdminBroadcast(p.Subject, p.Body)
	if err != nil {
		slog.Error("HTTP request error", "error", err)
		respond.WithError(w, err)
		return
	}

	respond.WithoutContent(w, http.StatusNoContent)
}
