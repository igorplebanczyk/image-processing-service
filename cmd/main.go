package main

import (
	"github.com/joho/godotenv"
	"image-processing-service/internal/images"
	"image-processing-service/internal/services/auth"
	"image-processing-service/internal/services/database"
	"image-processing-service/internal/services/server"
	"image-processing-service/internal/services/storage"
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

	userRepo := database.NewUserRepository(dbService.DB)
	refreshTokenRepo := database.NewRefreshTokenRepository(dbService.DB)
	imageRepo := database.NewImageRepository(dbService.DB)

	authService := auth.NewService(userRepo, refreshTokenRepo, os.Getenv("JWT_SECRET"), 15*time.Minute, 15*time.Hour)
	usersCfg := users.NewConfig(userRepo, refreshTokenRepo)

	storageService, err := storage.NewService(
		os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"),
		os.Getenv("AZURE_STORAGE_ACCOUNT_KEY"),
		os.Getenv("AZURE_STORAGE_ACCOUNT_URL"),
		os.Getenv("AZURE_STORAGE_CONTAINER_NAME"),
	)
	if err != nil {
		slog.Error("Error creating storage service", "error", err)
		return
	}

	imagesCfg := images.NewConfig(imageRepo, storageService)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		slog.Error("Error parsing port", "error", err)
		return
	}

	serverService := server.NewService(port, dbService, authService, usersCfg, imagesCfg)
	err = serverService.StartServer()
	if err != nil {
		slog.Error("Error starting server", "error", err)
		return
	}
}
