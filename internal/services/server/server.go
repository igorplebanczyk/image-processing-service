package server

import (
	"fmt"
	"image-processing-service/internal/images"
	"image-processing-service/internal/services/auth"
	"image-processing-service/internal/services/database"
	"image-processing-service/internal/users"
	"log/slog"
	"net/http"
)

type Service struct {
	port        int
	dbService   *database.Service
	authService *auth.Service
	usersCfg    *users.Config
	imagesCfg   *images.Config
}

func NewService(port int, dbService *database.Service, authService *auth.Service, usersCfg *users.Config, imagesCfg *images.Config) *Service {
	return &Service{
		port:        port,
		dbService:   dbService,
		authService: authService,
		usersCfg:    usersCfg,
		imagesCfg:   imagesCfg,
	}
}

func (s *Service) StartServer() error {
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	mux.HandleFunc("/health", health)

	mux.HandleFunc("GET /users", s.authService.Middleware(s.usersCfg.Info))
	mux.HandleFunc("POST /users", s.usersCfg.Register)
	mux.HandleFunc("DELETE /users", s.authService.Middleware(s.usersCfg.Delete))

	mux.HandleFunc("POST /login", s.authService.Login)
	mux.HandleFunc("POST /refresh", s.authService.Refresh)
	mux.HandleFunc("DELETE /logout", s.authService.Middleware(s.authService.Logout))

	mux.HandleFunc("POST /images", s.authService.Middleware(s.imagesCfg.Upload))
	mux.HandleFunc("GET /images", s.authService.Middleware(s.imagesCfg.Download))
	mux.HandleFunc("DELETE /images", s.authService.Middleware(s.imagesCfg.Delete))

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	slog.Info("Server starting", "port", srv.Addr)
	return nil
}
