package server

import (
	"context"
	"net/http"
)

type AuthService interface {
	Middleware(handler func(context.Context, http.ResponseWriter, *http.Request)) http.HandlerFunc
	Login(w http.ResponseWriter, r *http.Request)
	Refresh(w http.ResponseWriter, r *http.Request)
	Logout(ctx context.Context, w http.ResponseWriter, r *http.Request)
}
