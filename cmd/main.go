package main

import (
	"github.com/joho/godotenv"
	"image-processing-service/internal/auth"
	"image-processing-service/internal/database"
	"image-processing-service/internal/server"
	"image-processing-service/internal/users"
	"log/slog"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("Error loading .env file")
		return
	}

	dbService := database.NewService()
	err = dbService.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		return
	}

	userCfg := &users.Config{
		UserRepo:         database.NewUserRepository(dbService.DB),
		RefreshTokenRepo: database.NewRefreshTokenRepository(dbService.DB),
	}

	authService := auth.NewService(userCfg.UserRepo, userCfg.RefreshTokenRepo, os.Getenv("JWT_SECRET"), 15*time.Minute, 24*time.Hour)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		slog.Error("Error parsing port", "error", err)
		return
	}

	serverService := server.NewService(port, dbService, authService, userCfg)
	err = serverService.StartServer()
	if err != nil {
		slog.Error("Error starting server", "error", err)
		return
	}
}
