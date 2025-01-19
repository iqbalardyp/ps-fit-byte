package config

import (
	"fit-byte/db"
	file_handler "fit-byte/internal/file/handler"
	file_usecase "fit-byte/internal/file/usecase"
	custom_middleware "fit-byte/internal/middleware"
	"fit-byte/internal/routes"
	user_handler "fit-byte/internal/users/handler"
	user_repository "fit-byte/internal/users/repository"
	user_usecase "fit-byte/internal/users/usecase"
	"fit-byte/pkg/dotenv"
	"time"

	activityHandler "fit-byte/internal/activity/handler"
	activityRepository "fit-byte/internal/activity/repository"
	activityUsecase "fit-byte/internal/activity/usecase"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type BootstrapConfig struct {
	Env        *dotenv.Env
	App        *echo.Echo
	DB         *db.Postgres
	Log        *logrus.Logger
	Validator  *validator.Validate
	S3Uploader *manager.Uploader
}

func Bootstrap(config *BootstrapConfig) {
	//activity
	activityRepo := activityRepository.NewActivityRepository(config.DB.Pool)
	activityUsecase := activityUsecase.NewActivityUseCase(*activityRepo)
	activityHandler := activityHandler.NewActivityHandler(*activityUsecase, config.Validator)

	// * Middleware
	config.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Timeout",
		Timeout:      30 * time.Second,
	}))

	userRepo := user_repository.NewUserRepo(config.DB.Pool)
	userUsecase := user_usecase.NewUserUsecase(userRepo, config.Env)
	userHandler := user_handler.NewUserHandler(config.Validator, userUsecase)

	fileUsecase := file_usecase.NewFileUseCase(config.S3Uploader, config.Env)
	fileHandler := file_handler.NewFileHandler(fileUsecase, config.Log)

	authMiddleware := custom_middleware.NewAuthMiddleware(config.Env)
	routes := routes.RouteConfig{
		App:             config.App,
		S3Uploader:      config.S3Uploader,
		ActivityHandler: activityHandler,
		UserHandler:     userHandler,
		FileHandler:     fileHandler,
		Middleware:      authMiddleware,
	}

	routes.SetupRoutes()
}
