package server

import (
	"fmt"
	"image-processing-service/internal/images"
	"image-processing-service/internal/services/auth"
	"image-processing-service/internal/users"
	"log/slog"
	"net/http"
)

type Service struct {
	mux         *http.ServeMux
	server      *http.Server
	port        int
	authService *auth.Service
	usersCfg    *users.Config
	imagesCfg   *images.Config
}

func NewService(port int, authService *auth.Service, usersCfg *users.Config, imagesCfg *images.Config) *Service {
	service := &Service{
		port:        port,
		authService: authService,
		usersCfg:    usersCfg,
		imagesCfg:   imagesCfg,
	}
	service.init()
	return service
}

func (s *Service) Start() error {
	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	slog.Info("Server starting", "port", s.server.Addr)
	return nil
}

func (s *Service) init() {
	s.mux = http.NewServeMux()
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.mux,
	}

	s.mux.HandleFunc("/health", health)

	s.mux.HandleFunc("POST /users", s.usersCfg.Register)
	s.mux.HandleFunc("GET /users", s.authService.Middleware(s.usersCfg.Info))
	s.mux.HandleFunc("DELETE /users", s.authService.Middleware(s.usersCfg.Delete))

	s.mux.HandleFunc("POST /login", s.authService.Login)
	s.mux.HandleFunc("POST /refresh", s.authService.Refresh)
	s.mux.HandleFunc("DELETE /logout", s.authService.Middleware(s.authService.Logout))

	s.mux.HandleFunc("POST /images", s.authService.Middleware(s.imagesCfg.Upload))
	s.mux.HandleFunc("GET /images", s.authService.Middleware(s.imagesCfg.Download))
	s.mux.HandleFunc("DELETE /images", s.authService.Middleware(s.imagesCfg.Delete))
}
