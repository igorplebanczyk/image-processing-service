package server

import (
	"context"
	"fmt"
	authInterface "image-processing-service/src/internal/auth/interfaces"
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
		Handler:           mux,
		ReadTimeout:       time.Minute,
		ReadHeaderTimeout: 5 * time.Second,
	}

	mux.HandleFunc("/health", health)

	mux.HandleFunc("POST /users", s.userAPI.Register)
	mux.HandleFunc("GET /users", s.authAPI.UserMiddleware(s.userAPI.Info))
	mux.HandleFunc("PUT /users", s.authAPI.UserMiddleware(s.userAPI.Update))
	mux.HandleFunc("DELETE /users", s.authAPI.UserMiddleware(s.userAPI.Delete))

	mux.HandleFunc("POST /login", s.authAPI.Login)
	mux.HandleFunc("POST /refresh", s.authAPI.Refresh)
	mux.HandleFunc("DELETE /logout", s.authAPI.UserMiddleware(s.authAPI.Logout))

	mux.HandleFunc("POST /images", s.authAPI.UserMiddleware(s.imageAPI.Upload))
	mux.HandleFunc("GET /images/file", s.authAPI.UserMiddleware(s.imageAPI.Download))
	mux.HandleFunc("GET /images/info", s.authAPI.UserMiddleware(s.imageAPI.Info))
	mux.HandleFunc("GET /images/list", s.authAPI.UserMiddleware(s.imageAPI.List))
	mux.HandleFunc("PUT /images", s.authAPI.UserMiddleware(s.imageAPI.Transform))
	mux.HandleFunc("DELETE /images", s.authAPI.UserMiddleware(s.imageAPI.Delete))

	mux.HandleFunc("GET /admin/access", s.authAPI.AdminMiddleware(s.authAPI.AdminAccess))
	mux.HandleFunc("GET /admin/users", s.authAPI.AdminMiddleware(s.userAPI.AdminListAllUsers))
	mux.HandleFunc("PATCH /admin/users", s.authAPI.AdminMiddleware(s.userAPI.AdminUpdateRole))
	mux.HandleFunc("DELETE /admin/users", s.authAPI.AdminMiddleware(s.userAPI.AdminDeleteUser))
	mux.HandleFunc("DELETE /admin/logout", s.authAPI.AdminMiddleware(s.authAPI.AdminLogoutUser))
	mux.HandleFunc("GET /admin/images", s.authAPI.AdminMiddleware(s.imageAPI.AdminListAllImages))
	mux.HandleFunc("DELETE /admin/images", s.authAPI.AdminMiddleware(s.imageAPI.AdminDeleteImage))
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
