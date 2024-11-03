package auth

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"image-processing-service/internal/server/util"
	"net/http"
	"strings"
	"time"
)

func (s *Service) Middleware(handler func(*User, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			util.RespondWithError(w, http.StatusUnauthorized, "missing or invalid Authorization header")
			return
		}
		bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(bearerToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.jwtSecret), nil
		})
		if err != nil || !token.Valid {
			util.RespondWithError(w, http.StatusUnauthorized, "token expired or invalid")
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok || claims.ExpiresAt.Time.Before(time.Now()) {
			util.RespondWithError(w, http.StatusUnauthorized, "token expired or invalid")
			return
		}

		if claims.Issuer != issuer {
			util.RespondWithError(w, http.StatusUnauthorized, "invalid token issuer")
			return
		}

		user, err := s.userRepo.GetUserByValue("id", claims.Subject)
		if err != nil {
			util.RespondWithError(w, http.StatusInternalServerError, "error fetching user")
			return
		}

		handler(user, w, r)
	}
}

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
