package main

import (
	"github.com/joho/godotenv"
	"image-processing-service/internal/server"
	"log/slog"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		slog.Error("Error parsing port", "error", err)
	}

	cfg := server.Config{
		Port: port,
	}
	err = cfg.StartServer()
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
