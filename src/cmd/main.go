package main

import (
	"context"
	"fmt"
	authApp "image-processing-service/src/internal/auth/application"
	authInfra "image-processing-service/src/internal/auth/infrastructure"
	authInterface "image-processing-service/src/internal/auth/interfaces"
	"image-processing-service/src/internal/common/cache"
	"image-processing-service/src/internal/common/database"
	"image-processing-service/src/internal/common/database/transactions"
	"image-processing-service/src/internal/common/database/worker"
	"image-processing-service/src/internal/common/emails"
	_ "image-processing-service/src/internal/common/logs"
	"image-processing-service/src/internal/common/server"
	"image-processing-service/src/internal/common/storage"
	imagesApp "image-processing-service/src/internal/images/application"
	imagesInfra "image-processing-service/src/internal/images/infrastructure"
	imagesInterface "image-processing-service/src/internal/images/interfaces"
	usersApp "image-processing-service/src/internal/users/application"
	usersInfra "image-processing-service/src/internal/users/infrastructure"
	usersInterface "image-processing-service/src/internal/users/interfaces"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type application struct {
	serverService *server.Service
	dbService     *database.Service
	dbWorker      *worker.Worker
}

func main() {
	app := &application{}
	err := app.assemble()
	if err != nil {
		slog.Error("Init error: error assembling application", "error", err)
		panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go app.dbWorker.Start()
	go app.serverService.Start()

	slog.Info("Application started")

	<-ctx.Done()
	slog.Info("Received shutdown signal")

	app.dbWorker.Stop()
	app.dbService.Stop()
	app.serverService.Stop()

	slog.Info("Application shutdown")
}

func (a *application) assemble() error {
	// Get environment variables

	appPort := os.Getenv("APP_PORT")
	issuer := os.Getenv("APP_JWT_ISSUER")
	jwtSecret := os.Getenv("APP_JWT_SECRET")
	accessTokenExpiration := os.Getenv("APP_JWT_ACCESS_TOKEN_EXPIRATION")
	refreshTokenExpiration := os.Getenv("APP_JWT_REFRESH_TOKEN_EXPIRATION")

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

	mailHost := os.Getenv("MAIL_HOST")
	mailSenderEmail := os.Getenv("MAIL_SENDER_EMAIL")
	mailSenderPassword := os.Getenv("MAIL_SENDER_PASSWORD")

	// Convert environment variables to appropriate types

	appPortInt, err := strconv.Atoi(appPort)
	if err != nil {
		return fmt.Errorf("error converting port to integer: %w", err)
	}

	accessTokenExpirationInt, err := strconv.Atoi(accessTokenExpiration)
	if err != nil {
		return fmt.Errorf("error converting access token expiration to integer: %w", err)
	}
	accessTokenExpirationTime := time.Duration(accessTokenExpirationInt) * time.Minute

	refreshTokenExpirationInt, err := strconv.Atoi(refreshTokenExpiration)
	if err != nil {
		return fmt.Errorf("error converting refresh token expiration to integer: %w", err)
	}
	refreshTokenExpirationTime := time.Duration(refreshTokenExpirationInt) * time.Hour

	redisDBInt, err := strconv.Atoi(redisDB)
	if err != nil {
		return fmt.Errorf("error converting redis db to integer: %w", err)
	}

	// Setup common services

	dbService := database.NewService()
	db, err := dbService.Connect(postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	txProvider := transactions.NewTransactionProvider(db)
	dbWorker := worker.New(db, txProvider)

	cacheService, err := cache.NewService(redisHost, redisPort, redisPassword, redisDBInt)
	if err != nil {
		return fmt.Errorf("error creating cache: %w", err)
	}

	storageService, err := storage.NewService(
		azureStorageAccountName,
		azureStorageAccountKey,
		azureStorageAccountURL,
		azureStorageContainerName,
	)
	if err != nil {
		return fmt.Errorf("error creating storage: %w", err)
	}

	mailService, err := emails.NewService(mailHost, mailSenderEmail, mailSenderPassword)
	if err != nil {
		return fmt.Errorf("error creating email service: %w", err)
	}

	slog.Info("External services configured")

	// Assemble the application

	authUserRepo := authInfra.NewUserRepository(db)
	authRefreshTokenRepo := authInfra.NewRefreshTokenRepository(db, txProvider)
	authService := authApp.NewService(
		authUserRepo,
		authRefreshTokenRepo,
		jwtSecret,
		issuer,
		accessTokenExpirationTime,
		refreshTokenExpirationTime,
	)
	authAPI := authInterface.NewServer(authService)

	userRepo := usersInfra.NewUserRepository(db, txProvider)
	userService := usersApp.NewService(userRepo, mailService)
	userAPI := usersInterface.NewServer(userService)

	imageRepo := imagesInfra.NewImageRepository(db, txProvider)
	imageStorageRepo := imagesInfra.NewImageStorageRepository(storageService)
	imageCacheRepo := imagesInfra.NewImageCacheRepository(cacheService)
	imageService := imagesApp.NewService(imageRepo, imageStorageRepo, imageCacheRepo)
	imageAPI := imagesInterface.NewServer(imageService)

	serverService := server.NewService(appPortInt, authAPI, userAPI, imageAPI)

	a.serverService = serverService
	a.dbService = dbService
	a.dbWorker = dbWorker

	slog.Info("Application assembled")

	return nil
}
