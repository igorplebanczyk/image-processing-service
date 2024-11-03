package auth

import (
	"encoding/json"
	"image-processing-service/internal/server/util"
	"net/http"
)

func (s *Service) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var p parameters
	err := decoder.Decode(&p)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "error decoding request")
		return
	}

	user, err := s.repo.GetUserByValue("username", p.Username)
	if err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if !VerifyPassword(p.Password, user.Password) {
		util.RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	accessToken, err := s.generateAccessToken(*user)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "error generating access token")
		return
	}

	refreshToken, err := s.generateRefreshToken(*user)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, "error generating refresh token")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, response{AccessToken: accessToken, RefreshToken: refreshToken})
}

func (s *Service) refresh(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	var p parameters
	err := decoder.Decode(&p)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "error decoding request")
		return
	}

	resp, err := s.Refresh(p.RefreshToken)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "error refreshing token")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, resp)
}
