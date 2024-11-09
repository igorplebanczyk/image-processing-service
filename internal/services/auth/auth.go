package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"image-processing-service/internal/services/auth/util"
	"image-processing-service/internal/users"
	"time"
)

const issuer string = "image-processing-service"

type Service struct {
	userRepo         users.UserRepository
	refreshTokenRepo users.RefreshTokenRepository
	jwtSecret        string
	accessExpiry     time.Duration
	refreshExpiry    time.Duration
}

func New(userRepo users.UserRepository, refreshTokenRepo users.RefreshTokenRepository, secret string, accessExpiry time.Duration, refreshExpiry time.Duration) *Service {
	return &Service{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        secret,
		accessExpiry:     accessExpiry,
		refreshExpiry:    refreshExpiry,
	}
}

type response struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) authenticate(username string, password string) (response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return response{}, fmt.Errorf("error getting users by username: %w", err)
	}

	if !util.VerifyPassword(password, user.Password) {
		return response{}, fmt.Errorf("invalid password")
	}

	accessToken, err := s.generateAccessToken(*user)
	if err != nil {
		return response{}, fmt.Errorf("error generating access token: %w", err)
	}

	rawRefreshToken, err := s.generateRefreshToken(*user)
	if err != nil {
		return response{}, fmt.Errorf("error generating refresh token: %w", err)
	}
	_, err = s.refreshTokenRepo.CreateRefreshToken(ctx, user.ID, rawRefreshToken, time.Now().Add(s.refreshExpiry))
	if err != nil {
		return response{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	return response{AccessToken: accessToken, RefreshToken: rawRefreshToken}, nil
}

func (s *Service) refresh(refreshToken string) (response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		return response{}, fmt.Errorf("invalid users id in refresh token: %w", err)
	}

	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return response{}, fmt.Errorf("error fetching users: %w", err)
	}

	storedRefreshTokens, err := s.refreshTokenRepo.GetRefreshTokensByUserID(ctx, user.ID)
	if err != nil {
		return response{}, fmt.Errorf("error fetching refresh token: %w", err)
	}

	valid := false
	for _, storedRefreshToken := range storedRefreshTokens {
		if storedRefreshToken.Token == refreshToken {
			valid = true
			break
		}
	}
	if !valid {
		return response{}, fmt.Errorf("invalid refresh token")
	}

	accessToken, err := s.generateAccessToken(*user)
	if err != nil {
		return response{}, fmt.Errorf("error generating access token: %w", err)
	}

	return response{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
