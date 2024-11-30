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

type Service struct {
	port     int
	server   *http.Server
	authAPI  *authInterface.AuthAPI
	userAPI  *userInterface.UserAPI
	imageAPI *imageInterface.ImageAPI
}

func NewService(port int, authService *authInterface.AuthAPI, usersCfg *userInterface.UserAPI, imagesCfg *imageInterface.ImageAPI) *Service {
	service := &Service{
		port:     port,
		authAPI:  authService,
		userAPI:  usersCfg,
		imageAPI: imagesCfg,
	}
	service.setup()
	return service
}

func (s *Service) setup() {
	mux := http.NewServeMux()
	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           telemetry.Middleware(mux),
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
	}

	mux.Handle("/metrics", metrics.Handler())
	mux.HandleFunc("/health", health)

	mux.HandleFunc("POST /auth/login/one", s.authAPI.LoginOne)
	mux.HandleFunc("POST /auth/login/two", s.authAPI.LoginTwo)
	mux.HandleFunc("POST /auth/refresh", s.authAPI.Refresh)
	mux.HandleFunc("DELETE /auth/logout", s.authAPI.UserMiddleware(s.authAPI.Logout))

	mux.HandleFunc("POST /users", s.userAPI.Register)
	mux.HandleFunc("GET /users", s.authAPI.UserMiddleware(s.userAPI.GetDetails))
	mux.HandleFunc("PUT /users", s.authAPI.UserMiddleware(s.userAPI.UpdateDetails))
	mux.HandleFunc("DELETE /users", s.authAPI.UserMiddleware(s.userAPI.Delete))
	mux.HandleFunc("POST /users/verify", s.authAPI.UserMiddleware(s.userAPI.ResendVerificationCode))
	mux.HandleFunc("PATCH  /users/verify", s.authAPI.UserMiddleware(s.userAPI.Verify))
	mux.HandleFunc("POST /users/reset-password", s.userAPI.SendForgotPasswordCode)
	mux.HandleFunc("PATCH /users/reset-password", s.userAPI.ResetPassword)

	mux.HandleFunc("POST /images", s.authAPI.UserMiddleware(s.imageAPI.Upload))
	mux.HandleFunc("GET /images", s.authAPI.UserMiddleware(s.imageAPI.Get))
	mux.HandleFunc("GET /images/all", s.authAPI.UserMiddleware(s.imageAPI.GetAll))
	mux.HandleFunc("PUT /images", s.authAPI.UserMiddleware(s.imageAPI.UpdateDetails))
	mux.HandleFunc("PATCH /images", s.authAPI.UserMiddleware(s.imageAPI.Transform))
	mux.HandleFunc("DELETE /images", s.authAPI.UserMiddleware(s.imageAPI.Delete))

	mux.HandleFunc("POST /admin/broadcast", s.authAPI.AdminMiddleware(s.userAPI.AdminBroadcast))
	mux.HandleFunc("GET /admin/auth", s.authAPI.AdminMiddleware(s.authAPI.AdminAccess))
	mux.HandleFunc("DELETE /admin/auth/{id}", s.authAPI.AdminMiddleware(s.authAPI.AdminLogoutUser))
	mux.HandleFunc("GET /admin/users/{id}", s.authAPI.AdminMiddleware(s.userAPI.AdminGetUserDetails))
	mux.HandleFunc("GET /admin/users", s.authAPI.AdminMiddleware(s.userAPI.AdminGetAllUsersDetails))
	mux.HandleFunc("PATCH /admin/users/{id}", s.authAPI.AdminMiddleware(s.userAPI.AdminUpdateRole))
	mux.HandleFunc("DELETE /admin/users/{id}", s.authAPI.AdminMiddleware(s.userAPI.AdminDeleteUser))
	mux.HandleFunc("GET /admin/images", s.authAPI.AdminMiddleware(s.imageAPI.AdminListAllImages))
	mux.HandleFunc("DELETE /admin/images/{id}", s.authAPI.AdminMiddleware(s.imageAPI.AdminDeleteImage))
}

func (s *Service) Start() {
	slog.Info("Starting server", "port", s.port)
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
	slog.Info("Server stopped")
}
