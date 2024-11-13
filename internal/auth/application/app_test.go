package application

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func parseClaims(secret, token string) (*jwt.RegisteredClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsedToken.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func TestGenerateAccessToken(t *testing.T) {
	secret := "test_secret"
	issuer := "test_issuer"
	userID := uuid.New().String()
	expiry := 15 * time.Minute

	token, err := generateAccessToken(secret, issuer, userID, expiry)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected a token, got an empty string")
	}

	if !strings.HasPrefix(token, "eyJ") {
		t.Errorf("expected token to start with 'eyJ', got %s", token[:3])
	}

	claims, err := parseClaims(secret, token)
	if err != nil {
		t.Fatalf("expected valid token claims, got error: %v", err)
	}

	if claims.Subject != userID {
		t.Errorf("expected Subject %s, got %s", userID, claims.Subject)
	}
	if claims.Issuer != issuer {
		t.Errorf("expected Issuer %s, got %s", issuer, claims.Issuer)
	}

	expectedExpiry := time.Now().Add(expiry).Unix()
	if claims.ExpiresAt.Unix() > expectedExpiry || claims.ExpiresAt.Unix() < expectedExpiry-5 {
		t.Errorf("expected ExpiresAt close to %v, got %v", expectedExpiry, claims.ExpiresAt.Unix())
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	secret := "test_secret"
	issuer := "test_issuer"
	userID := uuid.New().String()
	expiry := 24 * time.Hour

	token, err := generateRefreshToken(secret, issuer, userID, expiry)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected a token, got an empty string")
	}

	if !strings.HasPrefix(token, "eyJ") {
		t.Errorf("expected token to start with 'eyJ', got %s", token[:3])
	}

	claims, err := parseClaims(secret, token)
	if err != nil {
		t.Fatalf("expected valid token claims, got error: %v", err)
	}

	if claims.Subject != userID {
		t.Errorf("expected Subject %s, got %s", userID, claims.Subject)
	}
	if claims.Issuer != issuer {
		t.Errorf("expected Issuer %s, got %s", issuer, claims.Issuer)
	}

	expectedExpiry := time.Now().Add(expiry).Unix()
	if claims.ExpiresAt.Unix() > expectedExpiry || claims.ExpiresAt.Unix() < expectedExpiry-5 {
		t.Errorf("expected ExpiresAt close to %v, got %v", expectedExpiry, claims.ExpiresAt.Unix())
	}
}

func TestVerifyAndParseToken(t *testing.T) {
	secret := "test_secret"
	issuer := "test_issuer"
	userID := uuid.New().String()
	expiry := 15 * time.Minute

	token, err := generateAccessToken(secret, issuer, userID, expiry)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	id, err := verifyAndParseToken(secret, issuer, token)
	if err != nil {
		t.Errorf("expected valid token, got error: %v", err)
	}
	if id.String() != userID {
		t.Errorf("expected userID %s, got %s", userID, id)
	}

	invalidToken := token + "tampered"
	_, err = verifyAndParseToken(secret, issuer, invalidToken)
	if err == nil {
		t.Errorf("expected error for tampered token, got none")
	}

	_, err = verifyAndParseToken(secret, "wrong_issuer", token)
	if err == nil {
		t.Errorf("expected error for wrong issuer, got none")
	}

	expiredToken, _ := generateAccessToken(secret, issuer, userID, -1*time.Minute) // Expired token
	_, err = verifyAndParseToken(secret, issuer, expiredToken)
	if err == nil {
		t.Errorf("expected error for expired token, got none")
	}
}
