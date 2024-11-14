package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"image-processing-service/internal/auth/domain"
	"testing"
	"time"
)

// Mocks

type MockUserRepository struct {
	GetUserByUsernameFunc func(ctx context.Context, username string) (*domain.User, error)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return m.GetUserByUsernameFunc(ctx, username)
}

type MockRefreshTokenRepository struct {
	CreateRefreshTokenFunc       func(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	GetRefreshTokensByUserIDFunc func(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error)
	RevokeRefreshTokenFunc       func(ctx context.Context, userID uuid.UUID) error
}

func (m *MockRefreshTokenRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	return m.CreateRefreshTokenFunc(ctx, userID, token, expiresAt)
}

func (m *MockRefreshTokenRepository) GetRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error) {
	return m.GetRefreshTokensByUserIDFunc(ctx, userID)
}

func (m *MockRefreshTokenRepository) RevokeRefreshToken(ctx context.Context, userID uuid.UUID) error {
	return m.RevokeRefreshTokenFunc(ctx, userID)
}

// Tests

func TestAuthService_Login(t *testing.T) {
	mockUserRepo := &MockUserRepository{
		GetUserByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
			if username == "valid_user" {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("valid_password"), bcrypt.DefaultCost)
				return &domain.User{ID: uuid.New(), Username: "valid_user", Password: string(hashedPassword)}, nil
			}
			return nil, fmt.Errorf("invalid username")
		},
	}

	mockRefreshTokenRepo := &MockRefreshTokenRepository{
		CreateRefreshTokenFunc: func(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
			return nil
		},
	}

	authService := AuthService{
		userRepo:         mockUserRepo,
		refreshTokenRepo: mockRefreshTokenRepo,
		secret:           "secret_key",
		issuer:           "test_issuer",
		accessExpiry:     15 * time.Minute,
		refreshExpiry:    7 * 24 * time.Hour,
	}

	t.Run("successful login", func(t *testing.T) {
		accessToken, refreshToken, err := authService.Login("valid_user", "valid_password")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if accessToken == "" || refreshToken == "" {
			t.Errorf("expected tokens, got empty strings")
		}
	})

	t.Run("invalid username", func(t *testing.T) {
		_, _, err := authService.Login("invalid_user", "any_password")
		if err == nil || err.Error() != "invalid username" {
			t.Errorf("expected invalid username error, got %v", err)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		_, _, err := authService.Login("valid_user", "invalid_password")
		if err == nil || err.Error() != "invalid password" {
			t.Errorf("expected invalid password error, got %v", err)
		}
	})
}

func TestAuthService_Refresh(t *testing.T) {
	secret := "secret_key"
	issuer := "test_issuer"
	validUserID := uuid.New()
	refreshExpiry := 24 * time.Hour

	validRefreshToken, err := generateRefreshToken(secret, issuer, validUserID.String(), refreshExpiry)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRefreshTokenRepo := &MockRefreshTokenRepository{
		GetRefreshTokensByUserIDFunc: func(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error) {
			if userID == validUserID {
				return []*domain.RefreshToken{
					{Token: validRefreshToken},
				}, nil
			}
			return nil, fmt.Errorf("user not found")
		},
	}

	authService := AuthService{
		refreshTokenRepo: mockRefreshTokenRepo,
		secret:           secret,
		issuer:           issuer,
		accessExpiry:     15 * time.Minute,
	}

	t.Run("successful refresh", func(t *testing.T) {
		accessToken, err := authService.Refresh(validRefreshToken)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if accessToken == "" {
			t.Errorf("expected a valid access token, got an empty string")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := authService.Refresh("invalid_refresh_token")
		if err == nil || err.Error() != "invalid token" {
			t.Errorf("expected invalid token error, got %v", err)
		}
	})

	t.Run("refresh token not found in database", func(t *testing.T) {
		mockRefreshTokenRepo.GetRefreshTokensByUserIDFunc = func(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error) {
			return []*domain.RefreshToken{}, nil
		}

		_, err := authService.Refresh(validRefreshToken)
		if err == nil || err.Error() != "invalid token" {
			t.Errorf("expected invalid token error, got %v", err)
		}
	})
}

func TestAuthService_Authenticate(t *testing.T) {
	validUserID := uuid.New()
	secret := "secret_key"
	issuer := "test_issuer"
	accessExpiry := 15 * time.Minute

	validAccessToken, err := generateAccessToken(secret, issuer, validUserID.String(), accessExpiry)
	if err != nil {
		t.Fatalf("failed to generate a valid access token: %v", err)
	}

	authService := AuthService{
		secret:       secret,
		issuer:       issuer,
		accessExpiry: accessExpiry,
	}

	t.Run("successful authentication", func(t *testing.T) {
		userID, err := authService.Authenticate(validAccessToken)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if userID != validUserID {
			t.Errorf("expected user ID %v, got %v", validUserID, userID)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := authService.Authenticate("invalid_access_token")
		if err == nil || err.Error() != "invalid token" {
			t.Errorf("expected invalid token error, got %v", err)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		expiredAccessToken, err := generateAccessToken(secret, issuer, validUserID.String(), -1*time.Minute)
		if err != nil {
			t.Fatalf("failed to generate an expired access token: %v", err)
		}

		_, err = authService.Authenticate(expiredAccessToken)
		if err == nil || err.Error() != "invalid token" {
			t.Errorf("expected invalid token error due to expiration, got %v", err)
		}
	})
}

func TestAuthService_Logout(t *testing.T) {
	validUserID := uuid.New()

	mockRefreshTokenRepo := &MockRefreshTokenRepository{
		RevokeRefreshTokenFunc: func(ctx context.Context, userID uuid.UUID) error {
			if userID == validUserID {
				return nil
			}
			return fmt.Errorf("user not found")
		},
	}

	authService := AuthService{
		refreshTokenRepo: mockRefreshTokenRepo,
	}

	t.Run("successful logout", func(t *testing.T) {
		err := authService.Logout(validUserID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		nonExistentUserID := uuid.New() // Simulate a different user ID
		err := authService.Logout(nonExistentUserID)
		if err == nil || err.Error() != "error deleting refresh tokens: user not found" {
			t.Errorf("expected user not found error, got %v", err)
		}
	})
}
