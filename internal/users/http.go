package users

import (
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
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	decoder := json.NewDecoder(r.Body)
	var p parameters
	err := decoder.Decode(&p)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error decoding request: %v", err))
		return
	}

	err = validate(cfg.userRepo, p.Username, p.Email, p.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error validating users: %v", err))
		return
	}

	user, err := cfg.userRepo.CreateUser(p.Username, p.Email, p.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating users: %v", err))
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, response{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}

func (cfg *Config) Delete(user *User, w http.ResponseWriter, _ *http.Request) {
	err := cfg.userRepo.DeleteUser(user.ID)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error deleting users: %v", err))
		return
	}

	util.RespondWithText(w, http.StatusOK, "user deleted successfully")
}
