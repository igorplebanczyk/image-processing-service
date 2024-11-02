package main

import (
	"github.com/joho/godotenv"
	"image-processing-service/internal/server"
	"log/slog"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	mux.HandleFunc("/health", server.HealthHandler)

	err = srv.ListenAndServe()
	if err != nil {
		slog.Error("Error starting server: %v", err)
	}
	slog.Info("Server starting", "port", srv.Addr)
}
