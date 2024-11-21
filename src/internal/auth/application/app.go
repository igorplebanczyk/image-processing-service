package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"image-processing-service/src/internal/auth/domain"
	commonerrors "image-processing-service/src/internal/common/errors"
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
		return "", "", commonerrors.NewUnauthorized("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", commonerrors.NewUnauthorized("invalid username or password")
	}

	accessToken, err := generateAccessToken(s.secret, s.issuer, user.ID.String(), s.accessExpiry)
	if err != nil {
		return "", "", commonerrors.NewInternal(fmt.Sprintf("error generating access token: %v", err))
	}

	refreshToken, err := generateRefreshToken(s.secret, s.issuer, user.ID.String(), s.refreshExpiry)
	if err != nil {
		return "", "", commonerrors.NewInternal(fmt.Sprintf("error generating refresh token: %v", err))
	}

	err = s.refreshTokenRepo.CreateRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(s.refreshExpiry))
	if err != nil {
		return "", "", commonerrors.NewInternal(fmt.Sprintf("error saving refresh token: %v", err))
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(rawRefreshToken string) (string, error) {
	id, err := verifyAndParseToken(s.secret, s.issuer, rawRefreshToken)
	if err != nil {
		return "", commonerrors.NewUnauthorized("invalid refresh token")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storedTokens, err := s.refreshTokenRepo.GetRefreshTokensByUserID(ctx, id)
	if err != nil {
		return "", commonerrors.NewInternal(fmt.Sprintf("error reading refresh tokens from database: %v", err))
	}

	found := false
	for _, t := range storedTokens {
		if t.Token == rawRefreshToken {
			found = true
			break
		}
	}
	if !found {
		return "", commonerrors.NewUnauthorized("invalid refresh token")
	}

	accessToken, err := generateAccessToken(s.secret, s.issuer, id.String(), s.accessExpiry)
	if err != nil {
		return "", commonerrors.NewInternal(fmt.Sprintf("error generating access token: %v", err))
	}

	return accessToken, nil
}

func (s *AuthService) Authenticate(accessToken string) (uuid.UUID, error) {
	id, err := verifyAndParseToken(s.secret, s.issuer, accessToken)
	if err != nil {
		return uuid.Nil, commonerrors.NewUnauthorized("invalid access token")
	}

	return id, nil
}

func (s *AuthService) AuthenticateAdmin(accessToken string) (uuid.UUID, error) {
	id, err := s.Authenticate(accessToken)
	if err != nil {
		return uuid.Nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	role, err := s.userRepo.GetUserRoleByID(ctx, id)
	if err != nil {
		return uuid.Nil, commonerrors.NewInternal(fmt.Sprintf("error reading user role from database: %v", err))
	}

	if role != domain.AdminRole {
		return uuid.Nil, commonerrors.NewForbidden("permission denied")
	}

	return id, nil
}

func (s *AuthService) Logout(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.refreshTokenRepo.RevokeRefreshToken(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to revoke refresh token: %v", err))
	}

	return nil
}
