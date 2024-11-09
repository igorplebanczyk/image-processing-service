package users

import (
	"context"
	"encoding/json"
	"fmt"
	"image-processing-service/internal/services/server/util"
	"net/http"
)

type Config struct {
	userRepo         UserRepository
	refreshTokenRepo RefreshTokenRepository
}

func NewConfig(userRepo UserRepository, refreshTokenRepo RefreshTokenRepository) *Config {
	return &Config{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}

func (cfg *Config) Register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = validateUsername(cfg.userRepo, p.Username)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = validateEmail(cfg.userRepo, p.Email)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = validatePassword(p.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user, err := cfg.userRepo.CreateUser(r.Context(), p.Username, p.Email, p.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating users: %v", err))
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, response{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}

func (cfg *Config) Info(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	type response struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	user := ctx.Value("user").(*User)

	util.RespondWithJSON(w, http.StatusOK, response{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}

func (cfg *Config) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value("user").(*User)

	err := cfg.userRepo.DeleteUser(r.Context(), user.ID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error deleting users: %v", err))
		return
	}

	util.RespondWithoutContent(w, http.StatusNoContent)
}

func (cfg *Config) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	type response struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	var p parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&p)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	user := ctx.Value("user").(*User)

	if p.Username == "" && p.Email == "" {
		util.RespondWithError(w, http.StatusBadRequest, "invalid request")
		return
	} else if p.Username == "" {
		err = validateEmail(cfg.userRepo, p.Email)
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
			return
		}
		p.Username = user.Username
	} else if p.Email == "" {
		err = validateUsername(cfg.userRepo, p.Username)
		if err != nil {
			util.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid request: %v", err))
			return
		}
		p.Email = user.Email
	}

	err = cfg.userRepo.UpdateUser(r.Context(), user.ID, p.Username, p.Email)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error updating users: %v", err))
		return
	}

	util.RespondWithJSON(w, http.StatusOK, response{
		Username:  p.Username,
		Email:     p.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}
