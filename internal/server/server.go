package server

import (
	"database/sql"
	"fmt"
	"image-processing-service/internal/user"
	"log/slog"
	"net/http"
)

type Config struct {
	Port int
	DB   *sql.DB
}

func (cfg *Config) StartServer() error {
	userCfg := user.Config{
		Repo: user.NewRepository(cfg.DB),
	}

	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/register", userCfg.RegisterUser)

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	slog.Info("Server starting", "port", srv.Addr)
	return nil
}
