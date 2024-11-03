package main

import (
	"github.com/joho/godotenv"
	"image-processing-service/internal/server"
	"log/slog"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	cfg := server.Config{
		Addr: "8080",
	}
	err = cfg.StartServer()
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
