package main

import (
	"github.com/joho/godotenv"
	"image-processing-service/internal/services/auth"
	database2 "image-processing-service/internal/services/database"
	"image-processing-service/internal/services/server"
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

	dbService := database2.NewService()
	err = dbService.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		return
	}

	userCfg := &users.Config{
		UserRepo:         database2.NewUserRepository(dbService.DB),
		RefreshTokenRepo: database2.NewRefreshTokenRepository(dbService.DB),
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
