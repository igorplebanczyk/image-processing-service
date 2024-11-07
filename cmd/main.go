package main

import (
	"fmt"
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

const envPath string = ".env"

func main() {
	serverService, err := configure()
	if err != nil {
		slog.Error("Error configuring services", "error", err)
		return
	}

	err = serverService.Start()
	if err != nil {
		slog.Error("Error starting server", "error", err)
		return
	}
}

func configure() (*server.Service, error) {
	err := godotenv.Load(envPath)
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_CONN")
	jwtSecret := os.Getenv("JWT_SECRET")
	azureStorageAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	azureStorageAccountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	azureStorageAccountURL := os.Getenv("AZURE_STORAGE_ACCOUNT_URL")
	azureStorageContainerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	dbService := database.NewService()

	err = dbService.Connect(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	userRepo := database.NewUserRepository(dbService.DB)
	refreshTokenRepo := database.NewRefreshTokenRepository(dbService.DB)
	imageRepo := database.NewImageRepository(dbService.DB)

	authService := auth.NewService(userRepo, refreshTokenRepo, jwtSecret, 15*time.Minute, 15*time.Hour)
	storageService, err := storage.NewService(
		azureStorageAccountName,
		azureStorageAccountKey,
		azureStorageAccountURL,
		azureStorageContainerName,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating storage service: %w", err)
	}

	usersCfg := users.NewConfig(userRepo, refreshTokenRepo)
	imagesCfg := images.NewConfig(imageRepo, storageService)

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("error converting port to integer: %w", err)
	}

	serverService := server.NewService(portInt, authService, usersCfg, imagesCfg)

	return serverService, nil
}
