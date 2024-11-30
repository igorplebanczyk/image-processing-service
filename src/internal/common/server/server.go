package server

import (
	"context"
	"fmt"
	authInterface "image-processing-service/src/internal/auth/interfaces"
	"image-processing-service/src/internal/common/metrics"
	"image-processing-service/src/internal/common/server/telemetry"
	imageInterface "image-processing-service/src/internal/images/interfaces"
	userInterface "image-processing-service/src/internal/users/interfaces"
	"log/slog"
	"net/http"
	"time"
)

const (
	readTimeout       = 1 * time.Minute
	readHeaderTimeout = 5 * time.Second
	writeTimeout      = 1 * time.Minute
	idleTimeout       = 1 * time.Minute
)

type Service struct {
	port      int
	server    *http.Server
	authAPI   *authInterface.AuthAPI
	usersAPI  *userInterface.UserAPI
	imagesAPI *imageInterface.ImageAPI
}

func NewService(
	port int,
	authAPI *authInterface.AuthAPI,
	usersAPI *userInterface.UserAPI,
	imagesAPI *imageInterface.ImageAPI,
) *Service {
	service := &Service{
		port:      port,
		authAPI:   authAPI,
		usersAPI:  usersAPI,
		imagesAPI: imagesAPI,
	}
	service.setup()

	slog.Info("Init step 16: server set up")

	return service
}

func (s *Service) setup() {
	mux := http.NewServeMux()
	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           telemetry.Middleware(mux),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	mux.Handle("/metrics", metrics.Handler())
	mux.HandleFunc("/health", health)

	mux.HandleFunc("POST /auth/login/one", s.authAPI.LoginOne)
	mux.HandleFunc("POST /auth/login/two", s.authAPI.LoginTwo)
	mux.HandleFunc("POST /auth/refresh", s.authAPI.Refresh)
	mux.HandleFunc("DELETE /auth/logout", s.authAPI.UserMiddleware(s.authAPI.Logout))

	mux.HandleFunc("POST /users", s.usersAPI.Register)
	mux.HandleFunc("GET /users", s.authAPI.UserMiddleware(s.usersAPI.GetDetails))
	mux.HandleFunc("PUT /users", s.authAPI.UserMiddleware(s.usersAPI.UpdateDetails))
	mux.HandleFunc("DELETE /users", s.authAPI.UserMiddleware(s.usersAPI.Delete))
	mux.HandleFunc("POST /users/verify", s.authAPI.UserMiddleware(s.usersAPI.ResendVerificationCode))
	mux.HandleFunc("PATCH  /users/verify", s.authAPI.UserMiddleware(s.usersAPI.Verify))
	mux.HandleFunc("POST /users/reset-password", s.usersAPI.SendForgotPasswordCode)
	mux.HandleFunc("PATCH /users/reset-password", s.usersAPI.ResetPassword)

	mux.HandleFunc("POST /images", s.authAPI.UserMiddleware(s.imagesAPI.Upload))
	mux.HandleFunc("GET /images", s.authAPI.UserMiddleware(s.imagesAPI.Get))
	mux.HandleFunc("GET /images/all", s.authAPI.UserMiddleware(s.imagesAPI.GetAll))
	mux.HandleFunc("PUT /images", s.authAPI.UserMiddleware(s.imagesAPI.UpdateDetails))
	mux.HandleFunc("PATCH /images", s.authAPI.UserMiddleware(s.imagesAPI.Transform))
	mux.HandleFunc("DELETE /images", s.authAPI.UserMiddleware(s.imagesAPI.Delete))

	mux.HandleFunc("POST /admin/broadcast", s.authAPI.AdminMiddleware(s.usersAPI.AdminBroadcast))
	mux.HandleFunc("GET /admin/auth", s.authAPI.AdminMiddleware(s.authAPI.AdminAccess))
	mux.HandleFunc("DELETE /admin/auth/{id}", s.authAPI.AdminMiddleware(s.authAPI.AdminLogoutUser))
	mux.HandleFunc("GET /admin/users/{id}", s.authAPI.AdminMiddleware(s.usersAPI.AdminGetUserDetails))
	mux.HandleFunc("GET /admin/users", s.authAPI.AdminMiddleware(s.usersAPI.AdminGetAllUsersDetails))
	mux.HandleFunc("PATCH /admin/users/{id}", s.authAPI.AdminMiddleware(s.usersAPI.AdminUpdateRole))
	mux.HandleFunc("DELETE /admin/users/{id}", s.authAPI.AdminMiddleware(s.usersAPI.AdminDeleteUser))
	mux.HandleFunc("GET /admin/images", s.authAPI.AdminMiddleware(s.imagesAPI.AdminListAllImages))
	mux.HandleFunc("DELETE /admin/images/{id}", s.authAPI.AdminMiddleware(s.imagesAPI.AdminDeleteImage))
}

func (s *Service) Start() {
	slog.Info("Init step 19: starting server")

	err := s.server.ListenAndServe()
	if err != nil {
		slog.Error("Init error: error starting server", "error", err)
		return
	}
}

func (s *Service) Stop() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		slog.Error("Shutdown error: error shutting down server", "error", err)
	}
	slog.Info("Shutdown step 5: server stopped")
}
