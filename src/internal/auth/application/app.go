package application

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"image-processing-service/src/internal/auth/domain"
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
		return "", "", domain.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := generateAccessToken(s.secret, s.issuer, user.ID.String(), s.accessExpiry)
	if err != nil {
		return "", "", errors.Join(domain.ErrInternal, err)
	}

	refreshToken, err := generateRefreshToken(s.secret, s.issuer, user.ID.String(), s.refreshExpiry)
	if err != nil {
		return "", "", errors.Join(domain.ErrInternal, err)
	}

	err = s.refreshTokenRepo.CreateRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(s.refreshExpiry))
	if err != nil {
		return "", "", errors.Join(domain.ErrInternal, err)
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(rawRefreshToken string) (string, error) {
	id, err := verifyAndParseToken(s.secret, s.issuer, rawRefreshToken)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storedTokens, err := s.refreshTokenRepo.GetRefreshTokensByUserID(ctx, id)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	found := false
	for _, t := range storedTokens {
		if t.Token == rawRefreshToken {
			found = true
			break
		}
	}
	if !found {
		return "", domain.ErrInvalidToken
	}

	accessToken, err := generateAccessToken(s.secret, s.issuer, id.String(), s.accessExpiry)
	if err != nil {
		return "", errors.Join(domain.ErrInternal, err)
	}

	return accessToken, nil
}

func (s *AuthService) Authenticate(accessToken string) (uuid.UUID, error) {
	id, err := verifyAndParseToken(s.secret, s.issuer, accessToken)
	if err != nil {
		return uuid.Nil, domain.ErrInvalidToken
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
		return uuid.Nil, errors.Join(domain.ErrInternal, err)
	}

	if role != domain.AdminRole {
		return uuid.Nil, domain.ErrPermissionDenied
	}

	return id, nil
}

func (s *AuthService) Logout(userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.refreshTokenRepo.RevokeRefreshToken(ctx, userID)
	if err != nil {
		return errors.Join(domain.ErrInternal, err)
	}

	return nil
}
