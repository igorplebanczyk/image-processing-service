package server

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Config struct {
	Port int
}

func (cfg *Config) StartServer() error {
	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	mux.HandleFunc("/health", healthHandler)

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %w", err)
	}

	slog.Info("Server starting", "port", srv.Addr)
	return nil
}
