package config

import (
	"fit-byte/db"
	"fit-byte/internal/routes"
	"time"

	auth "fit-byte/internal/middleware"

	activityHandler "fit-byte/internal/activity/handler"
	activityRepository "fit-byte/internal/activity/repository"
	activityUsecase "fit-byte/internal/activity/usecase"

	fileHandler "fit-byte/internal/file/handler"
	fileUsecase "fit-byte/internal/file/usecase"

	userHandler "fit-byte/internal/user/handler"
	userRepository "fit-byte/internal/user/repository"
	userUsecase "fit-byte/internal/user/usecase"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	App       *echo.Echo
	DB        *db.Postgres
	Log       *logrus.Logger
	Validator *validator.Validate
	S3Client  *manager.Uploader
}

func Bootstrap(config *BootstrapConfig) {
	//activity
	activityRepo := activityRepository.NewActivityRepository(config.DB.Pool)
	activityUsecase := activityUsecase.NewActivityUseCase(*activityRepo)
	activityHandler := activityHandler.NewActivityHandler(*activityUsecase, config.Validator)

	//file
	fileUsecase := fileUsecase.NewFileUseCase(config.S3Client)
	fileHandler := fileHandler.NewFileHandler(fileUsecase, config.Log)

	//user
	userRepo := userRepository.NewUserRepository(config.DB.Pool)
	userUseCase := userUsecase.NewUserUseCase(*userRepo)
	userHandler := userHandler.NewUserHandler(*userUseCase, config.Validator)

	// * Middleware
	config.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Timeout",
		Timeout:      30 * time.Second,
	}))
	authMiddleware := auth.Auth()

	routes := routes.RouteConfig{
		App:             config.App,
		ActivityHandler: activityHandler,
		FileHandler:     fileHandler,
		UserHandler:     userHandler,
		AuthMiddleware:  authMiddleware,
	}

	routes.SetupRoutes()
}
