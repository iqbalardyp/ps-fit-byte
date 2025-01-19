package handler

import (
	"net/http"
	"strconv"

	"fit-byte/internal/activity/dto"
	"fit-byte/internal/activity/model/converter"
	"fit-byte/internal/activity/usecase"
	customErrors "fit-byte/pkg/custom-errors"
	"fit-byte/pkg/jwt"
	"fit-byte/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const (
	DEFAULT_LIMIT = 5
)

type ActivityHandler struct {
	UseCase  usecase.ActivityUseCase
	Validate *validator.Validate
}

func NewActivityHandler(useCase usecase.ActivityUseCase, validate *validator.Validate) *ActivityHandler {
	return &ActivityHandler{
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *ActivityHandler) GetActivity(ctx echo.Context) error {
	userData := ctx.Get("user").(*jwt.JwtClaim)

	var request = new(dto.GetActivityRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if request.Limit == 0 {
		request.Limit = DEFAULT_LIMIT
	}

	if err := c.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	activities, err := c.UseCase.GetActivity(ctx.Request().Context(), request, userData.Id)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	response := converter.ToActivityResponseList(*activities)

	return ctx.JSON(http.StatusOK, response)
}

func (c *ActivityHandler) CreateActivity(ctx echo.Context) error {
	userData := ctx.Get("user").(*jwt.JwtClaim)

	var request = new(dto.CreateAndUpdateActivityRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := c.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	activity, err := c.UseCase.CreateActivity(ctx.Request().Context(), request, userData.Id)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	response := converter.ToActivityResponse(*activity)

	return ctx.JSON(http.StatusCreated, response)
}

func (c *ActivityHandler) UpdateActivity(ctx echo.Context) error {
	userData := ctx.Get("user").(*jwt.JwtClaim)
	activityIdStr := ctx.Param("activityId")

	activityId, err := strconv.Atoi(activityIdStr)
	if activityIdStr == "" || err != nil {
		err = errors.Wrap(customErrors.ErrNotFound, "activity ID is required and must be a valid integer")
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	var request = new(dto.CreateAndUpdateActivityRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := c.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	activity, err := c.UseCase.UpdateActivity(ctx.Request().Context(), request, userData.Id, activityId)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	response := converter.ToActivityResponse(*activity)

	return ctx.JSON(http.StatusOK, response)
}

func (c *ActivityHandler) DeleteActivity(ctx echo.Context) error {
	userData := ctx.Get("user").(*jwt.JwtClaim)
	activityIdStr := ctx.Param("activityId")

	activityId, err := strconv.Atoi(activityIdStr)
	if activityIdStr == "" || err != nil {
		err = errors.Wrap(customErrors.ErrNotFound, "activity ID is required and must be a valid integer")
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	err = c.UseCase.DeleteActivity(ctx.Request().Context(), userData.Id, activityId)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  http.StatusText(http.StatusOK),
		Message: "deleted",
	})
}
