package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const issuer string = "image-processing-service"

type Service struct {
	repo          UserRepository
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewService(repo UserRepository, secret string, accessExpiry time.Duration, refreshExpiry time.Duration) *Service {
	return &Service{
		repo:          repo,
		jwtSecret:     secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

type response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Authenticate(username string, password string) (response, error) {
	user, err := s.repo.GetUserByValue("username", username)
	if err != nil {
		return response{}, fmt.Errorf("error getting user by username: %w", err)
	}

	if !VerifyPassword(password, user.Password) {
		return response{}, fmt.Errorf("invalid password")
	}

	accessToken, err := s.generateAccessToken(*user)
	if err != nil {
		return response{}, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(*user)
	if err != nil {
		return response{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	return response{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *Service) Refresh(refreshToken string) (response, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return response{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid || claims.Issuer != issuer {
		return response{}, fmt.Errorf("invalid or expired refresh token")
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return response{}, fmt.Errorf("invalid user id in refresh token: %w", err)
	}

	user, err := s.repo.GetUserByValue("id", id.String())
	if err != nil {
		return response{}, fmt.Errorf("error fetching user: %w", err)
	}

	accessToken, err := s.generateAccessToken(*user)
	if err != nil {
		return response{}, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err = s.generateRefreshToken(*user)
	if err != nil {
		return response{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	return response{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
