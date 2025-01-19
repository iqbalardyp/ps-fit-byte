package routes

import (
	file_handler "fit-byte/internal/file/handler"
	custom_middleware "fit-byte/internal/middleware"
	user_handler "fit-byte/internal/users/handler"
	"fit-byte/pkg/response"
	"net/http"

	activityHandler "fit-byte/internal/activity/handler"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/labstack/echo/v4"
)

type RouteConfig struct {
	App             *echo.Echo
	S3Uploader      *manager.Uploader
	ActivityHandler *activityHandler.ActivityHandler
	UserHandler     *user_handler.UserHandler
	FileHandler     *file_handler.FileHandler
	Middleware      *custom_middleware.AuthConfig
}

func (r *RouteConfig) SetupRoutes() {
	r.App.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.BaseResponse{
			Status:  "Ok",
			Message: "",
		})
	})

	v1 := r.App.Group("/v1")
	r.setupPublicRoutes(v1)
	r.setupAuthRoutes(v1, r.Middleware.Authenticate())
}

func (r *RouteConfig) setupPublicRoutes(group *echo.Group) {
	group.POST("/register", r.UserHandler.Register)
	group.POST("/login", r.UserHandler.Login)
}

func (r *RouteConfig) setupAuthRoutes(group *echo.Group, m echo.MiddlewareFunc) {
	group.POST("/file", r.FileHandler.UploadFile, m)
	r.setupActivityRoute(group, m)
	r.setupUserRoutes(group, m)
}

func (r *RouteConfig) setupActivityRoute(api *echo.Group, m echo.MiddlewareFunc) {
	// user := api.Group("/activity", r.AuthMiddleware)
	user := api.Group("/activity", m)
	user.GET("", r.ActivityHandler.GetActivity, m)
	user.POST("", r.ActivityHandler.CreateActivity, m)
}

func (r *RouteConfig) setupUserRoutes(group *echo.Group, m echo.MiddlewareFunc) {
	group.GET("/user", r.UserHandler.GetUser, m)
	group.PATCH("/user", r.UserHandler.UpdateUser, m)
	user.PATCH("/:activityId",r.ActivityHandler.UpdateActivity)
	user.DELETE("/:activityId",r.ActivityHandler.UpdateActivity)
}
