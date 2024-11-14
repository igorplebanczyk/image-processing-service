package application

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"testing"
	"time"
)

func Test_generateAccessToken(t *testing.T) {
	type args struct {
		secret string
		issuer string
		userID string
		expiry time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful token generation",
			args: args{
				secret: "testsecret",
				issuer: "testissuer",
				userID: "12345",
				expiry: time.Hour,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateAccessToken(tt.args.secret, tt.args.issuer, tt.args.userID, tt.args.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("generateAccessToken() got = %v, want a non-empty string", got)
			}
		})
	}
}

func Test_generateRefreshToken(t *testing.T) {
	type args struct {
		secret string
		issuer string
		userID string
		expiry time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful refresh token generation",
			args: args{
				secret: "testsecret",
				issuer: "testissuer",
				userID: "12345",
				expiry: time.Hour * 24,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateRefreshToken(tt.args.secret, tt.args.issuer, tt.args.userID, tt.args.expiry)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("generateRefreshToken() got = %v, want a non-empty string", got)
			}
		})
	}
}

func Test_verifyAndParseToken(t *testing.T) {
	type args struct {
		secret   string
		issuer   string
		rawToken string
	}
	tests := []struct {
		name    string
		args    args
		want    uuid.UUID
		wantErr bool
	}{
		{
			name: "successful token verification and parsing",
			args: args{
				secret:   "testsecret",
				issuer:   "testissuer",
				rawToken: generateValidToken("testsecret", "testissuer", "89fe715c-8cf7-423d-9f50-d02c41a589b4", time.Hour),
			},
			want:    uuid.MustParse("89fe715c-8cf7-423d-9f50-d02c41a589b4"),
			wantErr: false,
		},
		{
			name: "invalid token",
			args: args{
				secret:   "testsecret",
				issuer:   "testissuer",
				rawToken: "invalidtoken",
			},
			want:    uuid.Nil,
			wantErr: true,
		},
		{
			name: "invalid token issuer",
			args: args{
				secret:   "testsecret",
				issuer:   "wrongissuer",
				rawToken: generateValidToken("testsecret", "testissuer", "12345", time.Hour),
			},
			want:    uuid.Nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := verifyAndParseToken(tt.args.secret, tt.args.issuer, tt.args.rawToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyAndParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("verifyAndParseToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func generateValidToken(secret, issuer, userID string, expiry time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
	})
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}
