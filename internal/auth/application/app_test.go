package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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
	type fields struct {
		userRepo         domain.UserRepository
		refreshTokenRepo domain.RefreshTokenRepository
		secret           string
		issuer           string
		accessExpiry     time.Duration
		refreshExpiry    time.Duration
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "invalid username",
			fields: fields{
				userRepo: &MockUserRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
						return nil, fmt.Errorf("user not found")
					},
				},
				refreshTokenRepo: &MockRefreshTokenRepository{},
				secret:           "secret",
				issuer:           "testIssuer",
				accessExpiry:     time.Hour,
				refreshExpiry:    time.Hour * 24,
			},
			args: args{
				username: "invaliduser",
				password: "password123",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "invalid password",
			fields: fields{
				userRepo: &MockUserRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
						return &domain.User{
							ID:       uuid.New(),
							Username: "testuser",
							Password: "$2a$10$KIXZnTTK1AtfQ5teOfRo3ePZ6VuMKhjFtsKxfwHMfNzm2ZFi0IdFS",
						}, nil
					},
				},
				refreshTokenRepo: &MockRefreshTokenRepository{},
				secret:           "secret",
				issuer:           "testIssuer",
				accessExpiry:     time.Hour,
				refreshExpiry:    time.Hour * 24,
			},
			args: args{
				username: "testuser",
				password: "wrongpassword",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "error generating token",
			fields: fields{
				userRepo: &MockUserRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
						return &domain.User{
							ID:       uuid.New(),
							Username: "testuser",
							Password: "$2a$10$KIXZnTTK1AtfQ5teOfRo3ePZ6VuMKhjFtsKxfwHMfNzm2ZFi0IdFS",
						}, nil
					},
				},
				refreshTokenRepo: &MockRefreshTokenRepository{
					CreateRefreshTokenFunc: func(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
						return nil
					},
				},
				secret:        "secret",
				issuer:        "testIssuer",
				accessExpiry:  time.Hour,
				refreshExpiry: time.Hour * 24,
			},
			args: args{
				username: "testuser",
				password: "password123",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name: "error adding refresh token to db",
			fields: fields{
				userRepo: &MockUserRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
						return &domain.User{
							ID:       uuid.New(),
							Username: "testuser",
							Password: "$2a$10$KIXZnTTK1AtfQ5teOfRo3ePZ6VuMKhjFtsKxfwHMfNzm2ZFi0IdFS",
						}, nil
					},
				},
				refreshTokenRepo: &MockRefreshTokenRepository{
					CreateRefreshTokenFunc: func(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
						return fmt.Errorf("db error")
					},
				},
				secret:        "secret",
				issuer:        "testIssuer",
				accessExpiry:  time.Hour,
				refreshExpiry: time.Hour * 24,
			},
			args: args{
				username: "testuser",
				password: "password123",
			},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthService{
				userRepo:         tt.fields.userRepo,
				refreshTokenRepo: tt.fields.refreshTokenRepo,
				secret:           tt.fields.secret,
				issuer:           tt.fields.issuer,
				accessExpiry:     tt.fields.accessExpiry,
				refreshExpiry:    tt.fields.refreshExpiry,
			}
			got, got1, err := s.Login(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Login() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	type fields struct {
		userRepo         domain.UserRepository
		refreshTokenRepo domain.RefreshTokenRepository
		secret           string
		issuer           string
		accessExpiry     time.Duration
		refreshExpiry    time.Duration
	}
	type args struct {
		rawRefreshToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "invalid refresh token format",
			fields: fields{
				userRepo: &MockUserRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
						return &domain.User{
							ID:       uuid.New(),
							Username: username,
							Password: "hashedPassword",
						}, nil
					},
				},
				refreshTokenRepo: &MockRefreshTokenRepository{},
				secret:           "secret",
				issuer:           "testIssuer",
				accessExpiry:     time.Hour,
				refreshExpiry:    time.Hour * 24,
			},
			args: args{
				rawRefreshToken: "invalid-refresh-token",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "token not found",
			fields: fields{
				userRepo: &MockUserRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
						return &domain.User{
							ID:       uuid.New(),
							Username: username,
							Password: "hashedPassword",
						}, nil
					},
				},
				refreshTokenRepo: &MockRefreshTokenRepository{
					GetRefreshTokensByUserIDFunc: func(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error) {
						return []*domain.RefreshToken{
							{Token: "stored-refresh-token", UserID: userID},
						}, nil
					},
				},
				secret:        "secret",
				issuer:        "testIssuer",
				accessExpiry:  time.Hour,
				refreshExpiry: time.Hour * 24,
			},
			args: args{
				rawRefreshToken: "non-matching-refresh-token",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "error fetching refresh tokens",
			fields: fields{
				userRepo: &MockUserRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
						return &domain.User{
							ID:       uuid.New(),
							Username: username,
							Password: "hashedPassword",
						}, nil
					},
				},
				refreshTokenRepo: &MockRefreshTokenRepository{
					GetRefreshTokensByUserIDFunc: func(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error) {
						return nil, fmt.Errorf("error fetching tokens")
					},
				},
				secret:        "secret",
				issuer:        "testIssuer",
				accessExpiry:  time.Hour,
				refreshExpiry: time.Hour * 24,
			},
			args: args{
				rawRefreshToken: "valid-refresh-token",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "error generating access token",
			fields: fields{
				userRepo: &MockUserRepository{
					GetUserByUsernameFunc: func(ctx context.Context, username string) (*domain.User, error) {
						return &domain.User{
							ID:       uuid.New(),
							Username: username,
							Password: "hashedPassword",
						}, nil
					},
				},
				refreshTokenRepo: &MockRefreshTokenRepository{
					GetRefreshTokensByUserIDFunc: func(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error) {
						return []*domain.RefreshToken{
							{Token: "valid-refresh-token", UserID: userID},
						}, nil
					},
				},
				secret:        "secret",
				issuer:        "testIssuer",
				accessExpiry:  time.Hour,
				refreshExpiry: time.Hour * 24,
			},
			args: args{
				rawRefreshToken: "valid-refresh-token",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthService{
				userRepo:         tt.fields.userRepo,
				refreshTokenRepo: tt.fields.refreshTokenRepo,
				secret:           tt.fields.secret,
				issuer:           tt.fields.issuer,
				accessExpiry:     tt.fields.accessExpiry,
				refreshExpiry:    tt.fields.refreshExpiry,
			}
			got, err := s.Refresh(tt.args.rawRefreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Refresh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Refresh() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Authenticate(t *testing.T) {
	type fields struct {
		userRepo         domain.UserRepository
		refreshTokenRepo domain.RefreshTokenRepository
		secret           string
		issuer           string
		accessExpiry     time.Duration
		refreshExpiry    time.Duration
	}
	type args struct {
		accessToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		{
			name: "invalid access token",
			fields: fields{
				userRepo:         &MockUserRepository{},
				refreshTokenRepo: &MockRefreshTokenRepository{},
				secret:           "secret",
				issuer:           "testIssuer",
				accessExpiry:     time.Hour,
				refreshExpiry:    time.Hour * 24,
			},
			args: args{
				accessToken: "invalid-access-token",
			},
			want:    uuid.Nil,
			wantErr: true,
		},
		{
			name: "error in token verification",
			fields: fields{
				userRepo:         &MockUserRepository{},
				refreshTokenRepo: &MockRefreshTokenRepository{},
				secret:           "secret",
				issuer:           "testIssuer",
				accessExpiry:     time.Hour,
				refreshExpiry:    time.Hour * 24,
			},
			args: args{
				accessToken: "malformed-token",
			},
			want:    uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthService{
				userRepo:         tt.fields.userRepo,
				refreshTokenRepo: tt.fields.refreshTokenRepo,
				secret:           tt.fields.secret,
				issuer:           tt.fields.issuer,
				accessExpiry:     tt.fields.accessExpiry,
				refreshExpiry:    tt.fields.refreshExpiry,
			}
			got, err := s.Authenticate(tt.args.accessToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Authenticate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	type fields struct {
		userRepo         domain.UserRepository
		refreshTokenRepo domain.RefreshTokenRepository
		secret           string
		issuer           string
		accessExpiry     time.Duration
		refreshExpiry    time.Duration
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful logout",
			fields: fields{
				userRepo: &MockUserRepository{},
				refreshTokenRepo: &MockRefreshTokenRepository{
					RevokeRefreshTokenFunc: func(ctx context.Context, userID uuid.UUID) error {
						return nil
					},
				},
				secret:        "secret",
				issuer:        "testIssuer",
				accessExpiry:  time.Hour,
				refreshExpiry: time.Hour * 24,
			},
			args: args{
				userID: uuid.New(),
			},
			wantErr: false,
		},
		{
			name: "error during logout",
			fields: fields{
				userRepo: &MockUserRepository{},
				refreshTokenRepo: &MockRefreshTokenRepository{
					RevokeRefreshTokenFunc: func(ctx context.Context, userID uuid.UUID) error {
						return fmt.Errorf("error revoking refresh token") // Simulating an error
					},
				},
				secret:        "secret",
				issuer:        "testIssuer",
				accessExpiry:  time.Hour,
				refreshExpiry: time.Hour * 24,
			},
			args: args{
				userID: uuid.New(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthService{
				userRepo:         tt.fields.userRepo,
				refreshTokenRepo: tt.fields.refreshTokenRepo,
				secret:           tt.fields.secret,
				issuer:           tt.fields.issuer,
				accessExpiry:     tt.fields.accessExpiry,
				refreshExpiry:    tt.fields.refreshExpiry,
			}
			err := s.Logout(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
