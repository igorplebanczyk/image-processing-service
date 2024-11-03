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

func (s *Service) Refresh(refreshToken string) (Response, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return Response{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid || claims.Issuer != issuer {
		return Response{}, fmt.Errorf("invalid or expired refresh token")
	}

	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return Response{}, fmt.Errorf("invalid user id in refresh token: %w", err)
	}

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return Response{}, fmt.Errorf("error fetching user: %w", err)
	}

	accessToken, err := s.GenerateAccessToken(*user)
	if err != nil {
		return Response{}, fmt.Errorf("error generating access token: %w", err)
	}

	refreshToken, err = s.GenerateRefreshToken(*user)
	if err != nil {
		return Response{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	return Response{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
