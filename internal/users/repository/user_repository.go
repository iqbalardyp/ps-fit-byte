package user_repository

import (
	"context"
	dto "fit-byte/internal/users/dto"
	customErrors "fit-byte/pkg/custom-errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

const (
	queryGetUserByEmail = "SELECT id, hashed_password FROM users WHERE email = @email;"
	queryCreateUser     = `
	INSERT INTO users(email, hashed_password)
	VALUES (@email, @hashedPassword)
	RETURNING id;`
	queryGetUserByID = `
	SELECT 
		name,
		email,
		image_uri,
		height,
		height_unit,
		weight,
		weight_unit,
		preference
	FROM users 
	WHERE id = @id;`
	queryUpdateUser = `
	WITH
	payload as (
		SELECT
			t.name,
			t.image_uri,
			t.height,
			t.height_unit,
			t.weight,
			t.weight_unit,
			t.preference
		FROM (VALUES
			(
				@name,
				@imageUri,
				@height::bigint,
				@heightUnit::enum_height_units,
				@weight::bigint,
				@weightUnit::enum_weight_units,
				@preference::enum_preferences
			)
		) as t(
			name,
			image_uri,
			height,
			height_unit,
			weight,
			weight_unit,
			preference
		)
	)
	UPDATE users
	SET 
		name = COALESCE(payload.name, users.name),
		image_uri = COALESCE(payload.image_uri, users.image_uri),
		height = payload.height,
		height_unit = payload.height_unit,
		weight = payload.weight,
		weight_unit = payload.weight_unit,
		preference = payload.preference
	FROM payload
	WHERE
		users.id = @id
	RETURNING
		users.name,
		users.image_uri,
		users.height,
		users.height_unit,
		users.weight,
		users.weight_unit,
		users.preference;`
)

func (r *UserRepo) GetUserByEmail(ctx context.Context, email *string) (*dto.AuthUser, error) {
	var user dto.AuthUser

	args := pgx.NamedArgs{
		"email": &email,
	}

	err := r.pool.QueryRow(ctx, queryGetUserByEmail, args).Scan(&user.ID, &user.HashedPassword)
	if err != nil {
		return nil, customErrors.HandlePgError(err, "failed to get user")
	}

	return &user, nil
}

func (r *UserRepo) RegisterUser(ctx context.Context, email, hashedPassword *string) (*int, error) {
	var id int
	args := pgx.NamedArgs{
		"email":          &email,
		"hashedPassword": &hashedPassword,
	}

	err := r.pool.QueryRow(ctx, queryCreateUser, args).Scan(&id)
	if err != nil {
		return nil, customErrors.HandlePgError(err, "failed to create user")
	}

	return &id, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id *int) (*dto.GetUserResponse, error) {
	var user dto.GetUserResponse

	args := pgx.NamedArgs{
		"id": &id,
	}
	err := r.pool.QueryRow(ctx, queryGetUserByID, args).Scan(
		&user.Name,
		&user.Email,
		&user.ImageURI,
		&user.Height,
		&user.HeightUnit,
		&user.Weight,
		&user.WeightUnit,
		&user.Preference,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, id *int, payload *dto.UpdateUserParams) (*dto.User, error) {
	var user dto.User
	args := pgx.NamedArgs{
		"id":         &id,
		"name":       &payload.Name,
		"imageUri":   &payload.ImageURI,
		"height":     &payload.Height,
		"heightUnit": &payload.HeightUnit,
		"weight":     &payload.Weight,
		"weightUnit": &payload.WeightUnit,
		"preference": &payload.Preference,
	}

	err := r.pool.QueryRow(ctx, queryUpdateUser, args).Scan(
		&user.Name,
		&user.ImageURI,
		&user.Height,
		&user.HeightUnit,
		&user.Weight,
		&user.WeightUnit,
		&user.Preference,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
