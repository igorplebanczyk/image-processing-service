package server

import (
	"database/sql"
	"fmt"
	"image-processing-service/internal/auth"
	"image-processing-service/internal/database"
	"image-processing-service/internal/users"
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	Port               int
	DB                 *sql.DB
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func (cfg *Config) StartServer() error {
	userCfg := users.Config{
		UserRepo:         database.NewUserRepository(cfg.DB),
		RefreshTokenRepo: database.NewRefreshTokenRepository(cfg.DB),
	}

	authService := auth.NewService(userCfg.UserRepo, userCfg.RefreshTokenRepo, cfg.JWTSecret, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)

	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	mux.HandleFunc("/health", health)
	mux.HandleFunc("POST /register", userCfg.RegisterUser)
	mux.HandleFunc("POST /login", authService.Login)
	mux.HandleFunc("POST /refresh", authService.Refresh)
	mux.HandleFunc("DELETE /logout", authService.Middleware(userCfg.Logout))

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	slog.Info("Server starting", "port", srv.Addr)
	return nil
}
