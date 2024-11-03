package server

import (
	"database/sql"
	"fmt"
	"image-processing-service/internal/database"
	"log/slog"
	"net/http"
)

type Config struct {
	Port int
	DB   *sql.DB
}

type ApiConfig struct {
	Repo *database.UserRepository
}

func (cfg *Config) StartServer() error {
	apiCfg := ApiConfig{
		Repo: database.NewUserRepository(cfg.DB),
	}

	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/register", apiCfg.RegisterUser)

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	slog.Info("Server starting", "port", srv.Addr)
	return nil
}
