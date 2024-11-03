package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Service struct {
	repo      UserRepository
	jwtSecret string
	jwtExpiry time.Duration
}

func NewService(repo UserRepository, jwtSecret string, jwtExpiry time.Duration) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *Service) GenerateJWT(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		Issuer:    "image-processing-service",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiry)),
	})

	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return signedToken, nil
}

func (s *Service) Authenticate(username string, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", fmt.Errorf("error getting user by username: %w", err)
	}

	if !CheckPasswordHash(password, user.Password) {
		return "", fmt.Errorf("invalid password")
	}

	token, err := s.GenerateJWT(*user)
	if err != nil {
		return "", fmt.Errorf("error generating JWT: %w", err)
	}

	return token, nil
}
