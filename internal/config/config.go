package config

import (
	"fit-byte/db"
	"fit-byte/internal/routes"
	"time"

	fileHandler "fit-byte/internal/file/handler"
	fileUsecase "fit-byte/internal/file/usecase"

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
	fileUsecase := fileUsecase.NewFileUseCase(config.S3Client)
	fileHandler := fileHandler.NewFileHandler(fileUsecase, config.Log)

	// * Middleware
	config.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Timeout",
		Timeout:      30 * time.Second,
	}))

	routes := routes.RouteConfig{
		App:         config.App,
		FileHandler: fileHandler,
	}

	routes.SetupRoutes()
}
