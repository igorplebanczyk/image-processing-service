package main

import (
	"github.com/joho/godotenv"
	"image-processing-service/internal/server"
	"image-processing-service/internal/user"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("Error loading .env file")
	}

	db, err := user.ConnectToDB(os.Getenv("DB_CONN"))
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		slog.Error("Error parsing port", "error", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	cfg := server.Config{
		Port:               port,
		DB:                 db,
		JWTSecret:          jwtSecret,
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 24 * time.Hour,
	}
	err = cfg.StartServer()
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
