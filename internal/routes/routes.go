package routes

import (
	"fit-byte/pkg/response"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App      *echo.Echo
	S3Client *s3.Client
}

func (r *RouteConfig) SetupRoutes() {
	r.setupPublicRoutes()
	r.setupAuthRoutes()
}

func (r *RouteConfig) setupPublicRoutes() {
	r.App.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.BaseResponse{
			Status:  "Ok",
			Message: "",
		})
	})
}
func (r *RouteConfig) setupAuthRoutes() {
}
