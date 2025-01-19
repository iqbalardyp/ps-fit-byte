package routes

import (
	"fit-byte/pkg/response"
	"net/http"

	activityHandler "fit-byte/internal/activity/handler"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App             *echo.Echo
	S3Client        *s3.Client
	ActivityHandler *activityHandler.ActivityHandler
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
	v1 := r.App.Group("/v1")

	r.setupActivityRoute(v1)
}

func (r *RouteConfig) setupActivityRoute(api *echo.Group) {
	// user := api.Group("/activity", r.AuthMiddleware)
	user := api.Group("/activity")
	user.GET("", r.ActivityHandler.GetActivity)
	user.POST("", r.ActivityHandler.CreateActivity)
	user.PATCH("/:activityId",r.ActivityHandler.UpdateActivity)
	user.DELETE("/:activityId",r.ActivityHandler.UpdateActivity)
}
