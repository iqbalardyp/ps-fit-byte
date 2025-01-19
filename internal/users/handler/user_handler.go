package user_handler

import (
	user_dto "fit-byte/internal/users/dto"
	user_usecase "fit-byte/internal/users/usecase"
	customErrors "fit-byte/pkg/custom-errors"
	"fit-byte/pkg/jwt"
	"fit-byte/pkg/response"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type UserHandler struct {
	Validate    *validator.Validate
	UserUsecase *user_usecase.UserUsecase
}

func NewUserHandler(validator *validator.Validate, usecase *user_usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		Validate:    validator,
		UserUsecase: usecase,
	}
}

func (h *UserHandler) Register(ctx echo.Context) error {
	var payload user_dto.AuthRequestParams

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(response.WriteErrorResponse(customErrors.ErrBadRequest))
	}

	if err := h.Validate.Struct(&payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	user, err := h.UserUsecase.Register(ctx.Request().Context(), &payload)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(ctx echo.Context) error {
	var payload user_dto.AuthRequestParams

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(response.WriteErrorResponse(customErrors.ErrBadRequest))
	}

	if err := h.Validate.Struct(&payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	user, err := h.UserUsecase.Login(ctx.Request().Context(), &payload)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUser(ctx echo.Context) error {
	authUser := ctx.Get("user").(*jwt.JWTClaim)
	user, err := h.UserUsecase.GetUser(ctx.Request().Context(), &authUser.ID)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, &user)
}

func (h *UserHandler) UpdateUser(ctx echo.Context) error {
	var payload user_dto.UpdateUserParams

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(response.WriteErrorResponse(customErrors.ErrBadRequest))
	}

	if err := h.Validate.Struct(&payload); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	authUser := ctx.Get("user").(*jwt.JWTClaim)
	user, err := h.UserUsecase.UpdateUser(ctx.Request().Context(), &authUser.ID, &payload)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, &user)
}
