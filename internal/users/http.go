package users

import (
	"encoding/json"
	"fmt"
	"image-processing-service/internal/services/server/util"
	"net/http"
)

type Config struct {
	UserRepo         UserRepository
	RefreshTokenRepo RefreshTokenRepository
}

func (cfg *Config) RegisterUser(w http.ResponseWriter, r *http.Request) {
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

	err = validate(cfg.UserRepo, p.Username, p.Email, p.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("error validating users: %v", err))
		return
	}

	user, err := cfg.UserRepo.CreateUser(p.Username, p.Email, p.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating users: %v", err))
		return
	}

	fmt.Printf("fine")
	util.RespondWithJSON(w, http.StatusCreated, response{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	})
}
