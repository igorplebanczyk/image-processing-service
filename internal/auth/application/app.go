package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"image-processing-service/internal/auth/domain"
	"time"
)

type AuthService struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	secret           string
	issuer           string
	accessExpiry     time.Duration
	refreshExpiry    time.Duration
}

func NewService(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	secret string,
	issuer string,
	accessExpiry time.Duration,
	refreshExpiry time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		secret:           secret,
		issuer:           issuer,
		accessExpiry:     accessExpiry,
		refreshExpiry:    refreshExpiry,
	}
}

func (s *AuthService) Login(username, password string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", fmt.Errorf("invalid username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("invalid password")
	}

	accessToken, err := generateAccessToken(s.secret, s.issuer, user.ID.String(), s.accessExpiry)
	refreshToken, err := generateRefreshToken(s.secret, s.issuer, user.ID.String(), s.refreshExpiry)
	if err != nil {
		return "", "", fmt.Errorf("error generating token: %w", err)
	}

	err = s.refreshTokenRepo.CreateRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(s.refreshExpiry))
	if err != nil {
		return "", "", fmt.Errorf("error adding refresh token to database: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(rawRefreshToken string) (string, error) {
	id, err := verifyAndParseToken(s.secret, s.issuer, rawRefreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storedTokens, err := s.refreshTokenRepo.GetRefreshTokensByUserID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("error getting refresh tokens: %w", err)
	}

	found := false
	for _, t := range storedTokens {
		if t.Token == rawRefreshToken {
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("invalid token")
	}

	accessToken, err := generateAccessToken(s.secret, s.issuer, id.String(), s.accessExpiry)
	if err != nil {
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return accessToken, nil
}

func (s *AuthService) Authenticate(accessToken string) (uuid.UUID, error) {
	id, err := verifyAndParseToken(s.secret, s.issuer, accessToken)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	return id, nil
}

func (s *AuthService) Logout(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.refreshTokenRepo.RevokeRefreshToken(ctx, userID)
	if err != nil {
		return fmt.Errorf("error deleting refresh tokens: %w", err)
	}

	return nil
}
