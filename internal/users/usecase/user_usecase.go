package user_usecase

import (
	"context"
	user_dto "fit-byte/internal/users/dto"
	user_repository "fit-byte/internal/users/repository"
	"fit-byte/pkg/bycript"
	"fit-byte/pkg/dotenv"
	"fit-byte/pkg/jwt"
)

type UserUsecase struct {
	UserRepo *user_repository.UserRepo
	Env      *dotenv.Env
}

func NewUserUsecase(repo *user_repository.UserRepo, env *dotenv.Env) *UserUsecase {
	return &UserUsecase{
		UserRepo: repo,
		Env:      env,
	}
}

func (u *UserUsecase) Register(ctx context.Context, payload *user_dto.AuthRequestParams) (*user_dto.AuthResponse, error) {
	// Hash Password
	hashedPassword, err := bycript.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	id, err := u.UserRepo.RegisterUser(ctx, &payload.Email, &hashedPassword)
	if err != nil {
		return nil, err
	}

	// Generate Token
	token, err := jwt.CreateToken(*id, u.Env.JWT_SECRET)
	if err != nil {
		return nil, err
	}

	authResponse := user_dto.AuthResponse{
		Email: payload.Email,
		Token: token,
	}
	return &authResponse, nil
}

func (u *UserUsecase) Login(ctx context.Context, payload *user_dto.AuthRequestParams) (*user_dto.AuthResponse, error) {
	var authResponse user_dto.AuthResponse

	user, err := u.UserRepo.GetUserByEmail(ctx, &payload.Email)
	if err != nil {
		return nil, err
	}

	// Compare password
	err = bycript.ComparePassword(payload.Password, user.HashedPassword)
	if err != nil {
		return nil, err
	}

	// Generate Token
	token, err := jwt.CreateToken(user.ID, u.Env.JWT_SECRET)
	if err != nil {
		return nil, err
	}

	authResponse.Email = payload.Email
	authResponse.Token = token
	return &authResponse, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, id *int) (*user_dto.GetUserResponse, error) {
	user, err := u.UserRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, id *int, payload *user_dto.UpdateUserParams) (*user_dto.User, error) {
	user, err := u.UserRepo.UpdateUser(ctx, id, payload)
	if err != nil {
		return nil, err
	}
	return user, nil
}
