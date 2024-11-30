package interfaces

import (
	"errors"
	commonerrors "image-processing-service/src/internal/common/errors"
	"net/http"
	"time"
)

func setAccessTokenInCookie(w http.ResponseWriter, accessToken string, accessTokenExpiry time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(accessTokenExpiry),
	})
}

func setRefreshTokenInCookie(w http.ResponseWriter, refreshToken string, refreshTokenExpiry time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(refreshTokenExpiry),
	})
}

func getAccessTokenFromCookie(r *http.Request) (string, error) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", commonerrors.NewUnauthorized("no access token cookie")
		}
		return "", commonerrors.NewInternal("failed to get access token cookie")
	}
	return accessCookie.Value, nil
}

func getRefreshTokenFromCookie(r *http.Request) (string, error) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", commonerrors.NewUnauthorized("no refresh token cookie")
		}
		return "", commonerrors.NewInternal("failed to get refresh token cookie")
	}
	return refreshCookie.Value, nil
}
