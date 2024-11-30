package main

import (
	"context"
	"fmt"
	authApplication "image-processing-service/src/internal/auth/application"
	authInfrastructure "image-processing-service/src/internal/auth/infrastructure"
	authInterfaces "image-processing-service/src/internal/auth/interfaces"
	"image-processing-service/src/internal/common/cache"
	"image-processing-service/src/internal/common/database"
	"image-processing-service/src/internal/common/database/tx"
	dbWorker "image-processing-service/src/internal/common/database/worker"
	"image-processing-service/src/internal/common/emails"
	_ "image-processing-service/src/internal/common/logs"
	_ "image-processing-service/src/internal/common/metrics"
	"image-processing-service/src/internal/common/server"
	"image-processing-service/src/internal/common/server/version"
	"image-processing-service/src/internal/common/storage"
	storageWorker "image-processing-service/src/internal/common/storage/worker"
	imagesApplication "image-processing-service/src/internal/images/application"
	"image-processing-service/src/internal/images/application/transformations"
	imagesInfrastructure "image-processing-service/src/internal/images/infrastructure"
	imagesInterfaces "image-processing-service/src/internal/images/interfaces"
	usersApplication "image-processing-service/src/internal/users/application"
	usersInfrastructure "image-processing-service/src/internal/users/infrastructure"
	usersInterfaces "image-processing-service/src/internal/users/interfaces"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type application struct {
	serverService          *server.Service
	dbService              *database.Service
	dbWorker               *dbWorker.Worker
	storageWorker          *storageWorker.Worker
	transformationsService *transformations.Service
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

	slog.Info("Init step 18: starting application")

	go app.dbWorker.Start()
	go app.storageWorker.Start()
	go app.serverService.Start()

	<-ctx.Done()
	slog.Info("Shutdown step 1: received signal to shutdown")

	app.transformationsService.Wait()
	app.dbWorker.Stop()
	app.storageWorker.Stop()
	app.dbService.Stop()
	app.serverService.Stop()

	slog.Info("Shutdown step 7: application shutdown")
}

// The assemble() function is quite a spaghetti code, but I feel like splitting it into smaller functions leads to
// unnecessary complexity. I figure it's not a big deal since it involves no complex logic. Also having the whole setup
// in a single function makes it easy to see how the application is assembled step by step.

func (a *application) assemble() error {
	// Get environment variables

	appPort := os.Getenv("APP_PORT")
	appVersion := os.Getenv("APP_VERSION")
	issuer := os.Getenv("APP_ISSUER")
	jwtSecret := os.Getenv("APP_JWT_SECRET")
	accessTokenExpiration := os.Getenv("APP_JWT_ACCESS_TOKEN_EXPIRATION")
	refreshTokenExpiration := os.Getenv("APP_JWT_REFRESH_TOKEN_EXPIRATION")
	otpExpiration := os.Getenv("APP_OTP_EXPIRATION")
	cacheExpiration := os.Getenv("APP_CACHE_EXPIRATION")

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

	slog.Info("Init step 3: environment variables loaded")

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

	otpExpirationInt, err := strconv.Atoi(otpExpiration)
	if err != nil {
		return fmt.Errorf("error converting otp expiration to integer: %w", err)
	}
	if otpExpirationInt < 0 {
		return fmt.Errorf("otp expiration must be greater than 0")
	}
	otpExpirationUint := uint(otpExpirationInt)

	cacheExpirationInt, err := strconv.Atoi(cacheExpiration)
	if err != nil {
		return fmt.Errorf("error converting cache expiration to integer: %w", err)
	}
	cacheExpirationTime := time.Duration(cacheExpirationInt) * time.Minute

	redisDBInt, err := strconv.Atoi(redisDB)
	if err != nil {
		return fmt.Errorf("error converting redis db to integer: %w", err)
	}

	slog.Info("Init step 4: environment variables converted")

	// Setup common services

	dbService := database.NewService()
	db, err := dbService.Connect(postgresUser, postgresPassword, postgresHost, postgresPort, postgresDB)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	txProvider := tx.NewProvider(db)

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

	transformationsService := transformations.NewService()

	slog.Info("Init step 12: all common services configured")

	// Assemble the application

	authUsersDBRepo := authInfrastructure.NewUsersDBRepository(db)
	authRefreshTokensDBRepo := authInfrastructure.NewRefreshTokensDBRepository(db, txProvider)
	authService := authApplication.NewService(
		authUsersDBRepo,
		authRefreshTokensDBRepo,
		mailService,
		jwtSecret,
		issuer,
		accessTokenExpirationTime,
		refreshTokenExpirationTime,
		otpExpirationUint,
	)
	authAPI := authInterfaces.NewAPI(authService, accessTokenExpirationTime, refreshTokenExpirationTime)

	slog.Info("Init step 13: auth module assembled")

	usersDBRepo := usersInfrastructure.NewUsersDBRepository(db, txProvider)
	usersService := usersApplication.NewService(usersDBRepo, mailService, issuer, otpExpirationUint)
	usersAPI := usersInterfaces.NewAPI(usersService)

	slog.Info("Init step 14: user module assembled")

	imagesDBRepo := imagesInfrastructure.NewImagesDBRepository(db, txProvider)
	imagesStorageRepo := imagesInfrastructure.NewImagesStorageRepository(storageService)
	imagesCacheRepo := imagesInfrastructure.NewImagesCacheRepository(cacheService)
	imagesService := imagesApplication.NewService(
		imagesDBRepo,
		imagesStorageRepo,
		imagesCacheRepo,
		transformationsService,
		cacheExpirationTime,
	)
	imagesAPI := imagesInterfaces.NewAPI(imagesService)

	slog.Info("Init step 15: image module assembled")

	version.Set(appVersion)
	serverService := server.NewService(appPortInt, authAPI, usersAPI, imagesAPI)

	a.serverService = serverService
	a.dbService = dbService
	a.dbWorker = dbWorker.New(db, txProvider)
	a.storageWorker = storageWorker.New(db, storageService)
	a.transformationsService = transformationsService

	slog.Info("Init step 17: application assembled")

	return nil
}
