package routes

import (
	"fit-byte/pkg/response"
	"net/http"

	activityHandler "fit-byte/internal/activity/handler"
	fileHandler "fit-byte/internal/file/handler"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App             *echo.Echo
	ActivityHandler *activityHandler.ActivityHandler
	FileHandler     *fileHandler.FileHandler
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
	r.setupFileRoutes(v1)
}

func (r *RouteConfig) setupActivityRoute(api *echo.Group) {
	// user := api.Group("/activity", r.AuthMiddleware)
	user := api.Group("/activity")
	user.GET("", r.ActivityHandler.GetActivity)
	user.POST("", r.ActivityHandler.CreateActivity)
}

func (r *RouteConfig) setupFileRoutes(api *echo.Group) {
	// file := api.Group("/file", r.AuthMiddleware)
	file := api.Group("/file")
	file.POST("", r.FileHandler.UploadFile)
}
