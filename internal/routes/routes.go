package routes

import (
	"fit-byte/pkg/response"
	"net/http"

	activityHandler "fit-byte/internal/activity/handler"
	fileHandler "fit-byte/internal/file/handler"
	userHandler "fit-byte/internal/user/handler"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App             *echo.Echo
	ActivityHandler *activityHandler.ActivityHandler
	FileHandler     *fileHandler.FileHandler
	UserHandler     *userHandler.UserHandler
	AuthMiddleware  echo.MiddlewareFunc
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
	r.App.POST("/v1/register", r.UserHandler.RegisterUser)
	r.App.POST("/v1/login", r.UserHandler.LoginUser)
}
func (r *RouteConfig) setupAuthRoutes() {
	v1 := r.App.Group("/v1")

	r.setupActivityRoute(v1)
	r.setupFileRoutes(v1)
	r.setupUserRoute(v1)
}

func (r *RouteConfig) setupActivityRoute(api *echo.Group) {
	activity := api.Group("/activity", r.AuthMiddleware)
	activity.GET("", r.ActivityHandler.GetActivity)
	activity.POST("", r.ActivityHandler.CreateActivity)
	activity.PATCH("/:activityId", r.ActivityHandler.UpdateActivity)
	activity.DELETE("/:activityId", r.ActivityHandler.DeleteActivity)
}

func (r *RouteConfig) setupFileRoutes(api *echo.Group) {
	file := api.Group("/file", r.AuthMiddleware)
	file.POST("", r.FileHandler.UploadFile)
}

func (r *RouteConfig) setupUserRoute(api *echo.Group) {
	user := api.Group("/user", r.AuthMiddleware)
	user.GET("", r.UserHandler.GetUser)
	user.PATCH("", r.UserHandler.UpdateUser)
}
