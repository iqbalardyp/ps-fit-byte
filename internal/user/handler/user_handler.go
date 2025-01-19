package handler

import (
	"net/http"

	"fit-byte/internal/user/dto"
	"fit-byte/internal/user/model/converter"
	"fit-byte/internal/user/usecase"
	customErrors "fit-byte/pkg/custom-errors"
	"fit-byte/pkg/jwt"
	"fit-byte/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type UserHandler struct {
	UseCase  usecase.UserUseCase
	Validate *validator.Validate
}

func NewUserHandler(useCase usecase.UserUseCase, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		UseCase:  useCase,
		Validate: validate,
	}
}

func (c *UserHandler) RegisterUser(ctx echo.Context) error {
	var request = new(dto.AuthRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := c.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	token, err := c.UseCase.Create(ctx.Request().Context(), request)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusCreated, &dto.AuthResponse{
		Email:       request.Email,
		AccessToken: *token,
	})
}

func (c *UserHandler) LoginUser(ctx echo.Context) error {
	var request = new(dto.AuthRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := c.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	token, err := c.UseCase.Login(ctx.Request().Context(), request)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	return ctx.JSON(http.StatusOK, &dto.AuthResponse{
		Email:       request.Email,
		AccessToken: *token,
	})
}

func (c *UserHandler) GetUser(ctx echo.Context) error {

	userData := ctx.Get("user").(*jwt.JwtClaim)
	user, err := c.UseCase.GetUser(ctx.Request().Context(), userData.Email)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	response := converter.ToUserResponse(*user)

	return ctx.JSON(http.StatusOK, response)
}

func (c *UserHandler) UpdateUser(ctx echo.Context) error {
	var request = new(dto.UpdateUserRequest)

	if err := ctx.Bind(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	if err := c.Validate.Struct(request); err != nil {
		err = errors.Wrap(customErrors.ErrBadRequest, err.Error())
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	userData := ctx.Get("user").(*jwt.JwtClaim)
	user, err := c.UseCase.UpdateUser(ctx.Request().Context(), request, userData.Id)
	if err != nil {
		return ctx.JSON(response.WriteErrorResponse(err))
	}

	response := converter.ToUserResponse(*user)

	return ctx.JSON(http.StatusOK, response)
}
