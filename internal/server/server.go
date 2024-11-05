package server

import (
	"fmt"
	"image-processing-service/internal/auth"
	"image-processing-service/internal/database"
	"image-processing-service/internal/users"
	"log/slog"
	"net/http"
)

type Service struct {
	port        int
	dbService   *database.Service
	authService *auth.Service
	userCfg     *users.Config
}

func NewService(port int, dbService *database.Service, authService *auth.Service, userCfg *users.Config) *Service {
	return &Service{
		port:        port,
		dbService:   dbService,
		authService: authService,
		userCfg:     userCfg,
	}
}

func (s *Service) StartServer() error {
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	mux.HandleFunc("/health", health)
	mux.HandleFunc("POST /register", s.userCfg.RegisterUser)
	mux.HandleFunc("POST /login", s.authService.Login)
	mux.HandleFunc("POST /refresh", s.authService.Refresh)
	mux.HandleFunc("DELETE /logout", s.authService.Middleware(s.userCfg.Logout))

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	slog.Info("Server starting", "port", srv.Addr)
	return nil
}
