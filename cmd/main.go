package main

import (
	"context"
	"fmt"
	"image-processing-service/internal/images"
	"image-processing-service/internal/services/auth"
	"image-processing-service/internal/services/cache"
	"image-processing-service/internal/services/database"
	"image-processing-service/internal/services/server"
	"image-processing-service/internal/services/storage"
	"image-processing-service/internal/services/transformation"
	"image-processing-service/internal/services/worker"
	"image-processing-service/internal/users"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	slog.Info("Starting application...")

	cfg := &Config{}
	err := cfg.configure()
	if err != nil {
		slog.Error("Error configuring services", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go cfg.workerService.Start()
	go cfg.serverService.Start()

	<-ctx.Done()
	slog.Info("Shutting down...")
	cfg.transformationService.Close()
	cfg.serverService.Stop()
	cfg.workerService.Stop()
}

type Config struct {
	serverService         *server.Service
	workerService         *worker.Service
	transformationService *transformation.Service
}

func (cfg *Config) configure() error {
	port := os.Getenv("PORT")
	jwtSecret := os.Getenv("JWT_SECRET")
	postgresURL := os.Getenv("POSTGRES_URL")

	azureStorageAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	azureStorageAccountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	azureStorageAccountURL := os.Getenv("AZURE_STORAGE_ACCOUNT_URL")
	azureStorageContainerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("error converting port to integer: %w", err)
	}

	redisDBInt, err := strconv.Atoi(redisDB)
	if err != nil {
		return fmt.Errorf("error converting redis db to integer: %w", err)
	}

	dbService := database.New()

	err = dbService.Connect(postgresURL)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	userRepo := database.NewUserRepository(dbService.DB)
	refreshTokenRepo := database.NewRefreshTokenRepository(dbService.DB)
	imageRepo := database.NewImageRepository(dbService.DB)

	authService := auth.New(userRepo, refreshTokenRepo, jwtSecret, 15*time.Minute, 15*time.Hour)
	storageService, err := storage.New(
		azureStorageAccountName,
		azureStorageAccountKey,
		azureStorageAccountURL,
		azureStorageContainerName,
	)
	if err != nil {
		return fmt.Errorf("error creating storage service: %w", err)
	}
	cacheService, err := cache.New(redisAddr, redisPassword, redisDBInt)
	if err != nil {
		return fmt.Errorf("error creating cache service: %w", err)
	}
	transformationService := transformation.New(10, 100)

	usersCfg := users.NewConfig(userRepo, refreshTokenRepo)
	imagesCfg := images.NewConfig(imageRepo, storageService, cacheService, transformationService)

	cfg.transformationService = transformationService
	cfg.workerService = worker.New(refreshTokenRepo, time.Hour)
	cfg.serverService = server.New(portInt, authService, usersCfg, imagesCfg)

	return nil
}
