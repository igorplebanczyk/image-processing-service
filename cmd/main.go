package main

import (
	"context"
	"fmt"
	authApp "image-processing-service/internal/auth/application"
	authInfra "image-processing-service/internal/auth/infrastructure"
	authInterface "image-processing-service/internal/auth/interfaces"
	"image-processing-service/internal/common/cache"
	"image-processing-service/internal/common/database"
	"image-processing-service/internal/common/database/transactions"
	"image-processing-service/internal/common/database/worker"
	"image-processing-service/internal/common/server"
	"image-processing-service/internal/common/storage"
	imagesApp "image-processing-service/internal/images/application"
	imagesInfra "image-processing-service/internal/images/infrastructure"
	imagesInterface "image-processing-service/internal/images/interfaces"
	usersApp "image-processing-service/internal/users/application"
	usersInfra "image-processing-service/internal/users/infrastructure"
	usersInterface "image-processing-service/internal/users/interfaces"
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
	err := cfg.assembleApplication()
	if err != nil {
		slog.Error("Error configuring services", "error", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go cfg.worker.Start()
	go cfg.serverService.Start()

	<-ctx.Done()
	slog.Info("Shutting down...")
	cfg.worker.Stop()
	cfg.dbService.Close()
	cfg.serverService.Stop()
}

type Config struct {
	dbService     *database.Service
	serverService *server.Service
	worker        *worker.Worker
}

func (cfg *Config) assembleApplication() error {
	port := os.Getenv("APP_PORT")
	jwtSecret := os.Getenv("APP_JWT_SECRET")

	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDB := os.Getenv("POSTGRES_DB")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	azureStorageAccountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	azureStorageAccountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	azureStorageAccountURL := os.Getenv("AZURE_STORAGE_ACCOUNT_URL")
	azureStorageContainerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("error converting port to integer: %w", err)
	}

	redisDBInt, err := strconv.Atoi(redisDB)
	if err != nil {
		return fmt.Errorf("error converting redis db to integer: %w", err)
	}

	dbService := database.NewService()
	db, err := dbService.Connect(postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	txProvider := transactions.NewTransactionProvider(db)
	workerService := worker.New(db, txProvider, time.Minute)

	cacheService, err := cache.NewService(redisHost, redisPort, redisPassword, redisDBInt)
	if err != nil {
		return fmt.Errorf("error creating cache: %w", err)
	}

	storageService, err := storage.NewService(azureStorageAccountName, azureStorageAccountKey, azureStorageAccountURL, azureStorageContainerName)
	if err != nil {
		return fmt.Errorf("error creating storage: %w", err)
	}

	authUserRepo := authInfra.NewUserRepository(db)
	authRefreshTokenRepo := authInfra.NewRefreshTokenRepository(db, txProvider)
	authService := authApp.NewService(authUserRepo, authRefreshTokenRepo, jwtSecret, "image-processing-service", 15*time.Minute, 15*time.Hour)
	authServer := authInterface.NewServer(authService)

	usersRepo := usersInfra.NewUserRepository(db, txProvider)
	userService := usersApp.NewService(usersRepo)
	usersServer := usersInterface.NewServer(userService)

	imagesRepo := imagesInfra.NewImageRepository(db, txProvider)
	imagesService := imagesApp.NewService(imagesRepo, storageService, cacheService, 10, 100)
	imagesServer := imagesInterface.NewServer(imagesService)

	serverService := server.NewService(portInt, authServer, usersServer, imagesServer)

	cfg.dbService = dbService
	cfg.serverService = serverService
	cfg.worker = workerService

	return nil
}
