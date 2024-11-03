package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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

func (s *Service) GenerateAccessToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessExpiry)),
	})

	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error signing access token: %w", err)
	}

	return signedToken, nil
}

func (s *Service) GenerateRefreshToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpiry)),
	})

	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error signing refresh token: %w", err)
	}

	return signedToken, nil
}

type Response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Authenticate(username string, password string) (Response, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return Response{}, fmt.Errorf("error getting user by username: %w", err)
	}

	if !CheckPasswordHash(password, user.Password) {
		return Response{}, fmt.Errorf("invalid password")
	}

	accessToken, err := s.GenerateAccessToken(*user)
	if err != nil {
		return Response{}, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err := s.GenerateRefreshToken(*user)
	if err != nil {
		return Response{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	return Response{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
