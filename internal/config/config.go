package config

import (
	"fit-byte/db"
	"fit-byte/internal/routes"
	"time"

	activityHandler "fit-byte/internal/activity/handler"
	activityRepository "fit-byte/internal/activity/repository"
	activityUsecase "fit-byte/internal/activity/usecase"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	S3Client  *s3.Client
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

	routes := routes.RouteConfig{
		App:             config.App,
		S3Client:        config.S3Client,
		ActivityHandler: activityHandler,
	}

	routes.SetupRoutes()
}
