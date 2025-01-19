package usecase

import (
	"context"

	"fit-byte/pkg/bcrypt"
	customErrors "fit-byte/pkg/custom-errors"
	jwt "fit-byte/pkg/jwt"

	"fit-byte/internal/user/dto"
	"fit-byte/internal/user/model"
	"fit-byte/internal/user/repository"

	"github.com/pkg/errors"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (c *UserUseCase) Create(ctx context.Context, request *dto.AuthRequest) (*string, error) {
	hashedPassword, err := bcrypt.HashPassword(request.Password)
	if err != nil {
		return nil, err
	}

	arg := repository.CreateUserParams{
		Email:          request.Email,
		HashedPassword: hashedPassword,
	}

	user, err := c.userRepo.CreateUser(ctx, arg)

	if err != nil {
		if customErrors.ErrorCode(err) == customErrors.UniqueViolation {
			return nil, errors.Wrap(customErrors.ErrConflict, "email is exist")
		}
		return nil, errors.Wrap(err, "failed to create user")
	}

	token, err := jwt.CreateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (c *UserUseCase) Login(ctx context.Context, request *dto.AuthRequest) (*string, error) {

	user, err := c.userRepo.GetUserFromEmail(ctx, request.Email)

	if err != nil {
		if errors.Is(err, customErrors.ErrNotFound) {
			return nil, errors.Wrap(customErrors.ErrNotFound, "User not found")
		}
		return nil, errors.Wrap(err, "failed to get user")
	}

	err = bcrypt.ComparePassword(request.Password, user.HashedPassword)
	if err != nil {
		return nil, errors.Wrap(customErrors.ErrBadRequest, "password is wrong")
	}

	token, err := jwt.CreateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (c *UserUseCase) GetUser(ctx context.Context, userEmail string) (*model.User, error) {

	user, err := c.userRepo.GetUserFromEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, customErrors.ErrNotFound) {
			return nil, errors.Wrap(customErrors.ErrNotFound, "User not found")
		}
		return nil, errors.Wrap(err, "failed to get user")
	}

	return &user, nil
}

func (c *UserUseCase) UpdateUser(ctx context.Context, request *dto.UpdateUserRequest, userid int) (*model.User, error) {

	arg := repository.UpdateUserParams{
		ID:           userid,
		Username:     request.Username,
		UserImageUri: request.UserImageUri,
		Weight:       request.Weight,
		Height:       request.Height,
		WeightUnit:   request.WeightUnit,
		HeightUnit:   request.HeightUnit,
		Preference:   request.Preference,
	}

	user, err := c.userRepo.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, customErrors.ErrNotFound) {
			return nil, errors.Wrap(customErrors.ErrNotFound, "User not found")
		}
		return nil, errors.Wrap(err, "failed to update user")
	}

	return &user, nil
}
