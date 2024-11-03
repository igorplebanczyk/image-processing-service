package main

import (
	"github.com/joho/godotenv"
	"image-processing-service/internal/database"
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

	_, err = database.Connect(os.Getenv("DB_CONN"))
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
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
