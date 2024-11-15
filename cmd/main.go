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
	"image-processing-service/internal/common/log"
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

const (
	issuer                 = "image-processing-service"
	accessTokenExpiration  = 15 * time.Minute
	refreshTokenExpiration = 15 * time.Hour
	dbWorkerInterval       = time.Hour
)

type application struct {
	serverService *server.Service
	dbService     *database.Service
	dbWorker      *worker.Worker
}

func main() {
	err := log.Setup("logs")
	if err != nil {
		panic(err)
	}

	slog.Info("Starting Application...")

	app := &application{}
	err = app.assemble()
	if err != nil {
		slog.Error("Error configuring services", "error", err)
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go app.dbWorker.Start()
	go app.serverService.Start()

	<-ctx.Done()
	slog.Info("Shutting down...")
	app.dbWorker.Stop()
	app.dbService.Stop()
	app.serverService.Stop()
}

func (a *application) assemble() error {
	// Get environment variables

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

	// Setup services

	dbService := database.NewService()
	db, err := dbService.Connect(postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	txProvider := transactions.NewTransactionProvider(db)
	dbWorker := worker.New(db, txProvider, dbWorkerInterval)

	cacheService, err := cache.NewService(redisHost, redisPort, redisPassword, redisDBInt)
	if err != nil {
		return fmt.Errorf("error creating cache: %w", err)
	}

	storageService, err := storage.NewService(azureStorageAccountName, azureStorageAccountKey, azureStorageAccountURL, azureStorageContainerName)
	if err != nil {
		return fmt.Errorf("error creating storage: %w", err)
	}

	// Assemble the application

	authUserRepo := authInfra.NewUserRepository(db)
	authRefreshTokenRepo := authInfra.NewRefreshTokenRepository(db, txProvider)
	authService := authApp.NewService(authUserRepo, authRefreshTokenRepo, jwtSecret, issuer, accessTokenExpiration, refreshTokenExpiration)
	authServer := authInterface.NewServer(authService)

	usersRepo := usersInfra.NewUserRepository(db, txProvider)
	userService := usersApp.NewService(usersRepo)
	usersServer := usersInterface.NewServer(userService)

	imagesRepo := imagesInfra.NewImageRepository(db, txProvider)
	imagesService := imagesApp.NewService(imagesRepo, storageService, cacheService, 10, 100)
	imagesServer := imagesInterface.NewServer(imagesService)

	serverService := server.NewService(portInt, authServer, usersServer, imagesServer)

	a.serverService = serverService
	a.dbService = dbService
	a.dbWorker = dbWorker

	return nil
}
