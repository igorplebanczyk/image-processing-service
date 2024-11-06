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

const envPath string = "../.env"

func main() {
	err := godotenv.Load(envPath)
	if err != nil {
		slog.Error("Error loading .env file")
		return
	}

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	azureStorageAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	azureStorageAccountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	azureStorageAccountURL := os.Getenv("AZURE_STORAGE_ACCOUNT_URL")
	azureStorageContainerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	dbService := database.NewService()

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
		slog.Error("Error creating storage service", "error", err)
		return
	}

	usersCfg := users.NewConfig(userRepo, refreshTokenRepo)
	imagesCfg := images.NewConfig(imageRepo, storageService)

	portInt, err := strconv.Atoi(port)
	if err != nil {
		slog.Error("Error parsing port", "error", err)
		return
	}

	serverService := server.NewService(portInt, dbService, authService, usersCfg, imagesCfg)

	err = dbService.Connect(dbURL)
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		return
	}

	err = serverService.StartServer()
	if err != nil {
		slog.Error("Error starting server", "error", err)
		return
	}
}
