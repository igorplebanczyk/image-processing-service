package auth

import (
	"encoding/json"
	"image-processing-service/internal/server/util"
	"net/http"
)

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
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

	resp, err := s.authenticate(p.Username, p.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "error authenticating user")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, resp)
}

func (s *Service) Refresh(w http.ResponseWriter, r *http.Request) {
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

	resp, err := s.refresh(p.RefreshToken)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "error refreshing token")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, resp)
}
