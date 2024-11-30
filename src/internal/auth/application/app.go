package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"image-processing-service/src/internal/auth/domain"
	"image-processing-service/src/internal/common/emails"
	commonerrors "image-processing-service/src/internal/common/errors"
	"time"
)

type AuthService struct {
	userDBRepo         domain.UserDBRepository
	refreshTokenDBRepo domain.RefreshTokenDBRepository
	mailService        *emails.Service
	secret             string
	issuer             string
	accessExpiry       time.Duration
	refreshExpiry      time.Duration
	otpExpiry          uint
}

func NewService(
	userRepo domain.UserDBRepository,
	refreshTokenRepo domain.RefreshTokenDBRepository,
	mailService *emails.Service,
	secret string,
	issuer string,
	accessExpiry time.Duration,
	refreshExpiry time.Duration,
	otpExpiry uint,
) *AuthService {
	return &AuthService{
		userDBRepo:         userRepo,
		refreshTokenDBRepo: refreshTokenRepo,
		mailService:        mailService,
		secret:             secret,
		issuer:             issuer,
		accessExpiry:       accessExpiry,
		refreshExpiry:      refreshExpiry,
		otpExpiry:          otpExpiry,
	}
}

func (s *AuthService) LoginOne(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.userDBRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return commonerrors.NewUnauthorized("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return commonerrors.NewUnauthorized("invalid username or password")
	}

	otp, err := totp.GenerateCode(user.OTPSecret, time.Now())
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error generating OTP code: %v", err))
	}

	err = s.mailService.SendText(user.Email, "OTP code", fmt.Sprintf("Your 2FA code is: %s", otp))
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("error sending OTP code: %v", err))
	}

	return nil
}

func (s *AuthService) LoginTwo(username, otp string) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.userDBRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", commonerrors.NewUnauthorized("invalid username or password")
	}

	ok := totp.Validate(otp, user.OTPSecret)
	if !ok {
		return "", "", commonerrors.NewUnauthorized("invalid OTP code")
	}

	accessToken, err := generateAccessToken(s.secret, s.issuer, user.ID.String(), s.accessExpiry)
	if err != nil {
		return "", "", commonerrors.NewInternal(fmt.Sprintf("error generating access token: %v", err))
	}

	refreshToken, err := generateRefreshToken(s.secret, s.issuer, user.ID.String(), s.refreshExpiry)
	if err != nil {
		return "", "", commonerrors.NewInternal(fmt.Sprintf("error generating refresh token: %v", err))
	}

	err = s.refreshTokenDBRepo.CreateRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(s.refreshExpiry))
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

	storedToken, err := s.refreshTokenDBRepo.GetRefreshTokenByUserIDandToken(ctx, id, rawRefreshToken)
	if err != nil {
		return "", commonerrors.NewUnauthorized("invalid refresh token")
	}
	if storedToken.ExpiresAt.Before(time.Now()) {
		_ = s.refreshTokenDBRepo.RevokeAllUserRefreshTokens(ctx, id) // Optionally log or track revoked tokens
		return "", commonerrors.NewUnauthorized("refresh token expired or revoked")
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

	role, err := s.userDBRepo.GetUserRoleByID(ctx, id)
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

	err := s.refreshTokenDBRepo.RevokeAllUserRefreshTokens(ctx, userID)
	if err != nil {
		return commonerrors.NewInternal(fmt.Sprintf("failed to revoke refresh token: %v", err))
	}

	return nil
}
