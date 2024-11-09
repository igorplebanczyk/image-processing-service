package server

import (
	"context"
	"fmt"
	"image-processing-service/internal/images"
	"image-processing-service/internal/users"
	"log/slog"
	"net/http"
	"time"
)

type Service struct {
	port        int
	server      *http.Server
	authService AuthService
	usersCfg    *users.Config
	imagesCfg   *images.Config
}

func New(port int, authService AuthService, usersCfg *users.Config, imagesCfg *images.Config) *Service {
	service := &Service{
		port:        port,
		authService: authService,
		usersCfg:    usersCfg,
		imagesCfg:   imagesCfg,
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

	mux.HandleFunc("POST /users", s.usersCfg.Register)
	mux.HandleFunc("GET /users", s.authService.Middleware(s.usersCfg.Info))
	mux.HandleFunc("PUT /users", s.authService.Middleware(s.usersCfg.Update))
	mux.HandleFunc("DELETE /users", s.authService.Middleware(s.usersCfg.Delete))

	mux.HandleFunc("POST /login", s.authService.Login)
	mux.HandleFunc("POST /refresh", s.authService.Refresh)
	mux.HandleFunc("DELETE /logout", s.authService.Middleware(s.authService.Logout))

	mux.HandleFunc("POST /images", s.authService.Middleware(s.imagesCfg.Upload))
	mux.HandleFunc("GET /images/file", s.authService.Middleware(s.imagesCfg.Download))
	mux.HandleFunc("GET /images/info", s.authService.Middleware(s.imagesCfg.Info))
	mux.HandleFunc("GET /images/list", s.authService.Middleware(s.imagesCfg.List))
	mux.HandleFunc("PUT /images", s.authService.Middleware(s.imagesCfg.Transform))
	mux.HandleFunc("DELETE /images", s.authService.Middleware(s.imagesCfg.Delete))
}

func (s *Service) Start() {
	slog.Info("Server starting", "port", s.port)
	err := s.server.ListenAndServe()
	if err != nil {
		slog.Error("Error starting server", "error", err)
		return
	}
}

func (s *Service) Stop() {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		slog.Error("Error shutting down server", "error", err)
	}
	slog.Info("Server stopped")
}
