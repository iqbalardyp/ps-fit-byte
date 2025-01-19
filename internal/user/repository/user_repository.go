package repository

import (
	"context"
	"fit-byte/internal/user/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  email,
  hashed_password
) VALUES (
  $1, $2
) RETURNING id, email, hashed_password, username, user_image_uri, weight, height, height_unit, weight_unit, preference
`

type CreateUserParams struct {
	Email          string
	HashedPassword string
}

func (r *UserRepository) CreateUser(ctx context.Context, arg CreateUserParams) (model.User, error) {
	row := r.pool.QueryRow(ctx, createUser, arg.Email, arg.HashedPassword)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Username,
		&i.UserImageUri,
		&i.Weight,
		&i.Height,
		&i.HeightUnit,
		&i.WeightUnit,
		&i.Preference,
	)
	return i, err
}

const getUserFromEmail = `-- name: GetUserFromEmail :one
SELECT id, email, hashed_password, username, user_image_uri, weight, height, height_unit, weight_unit, preference FROM users
WHERE email = $1 LIMIT 1
`

func (r *UserRepository) GetUserFromEmail(ctx context.Context, email string) (model.User, error) {
	row := r.pool.QueryRow(ctx, getUserFromEmail, email)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Username,
		&i.UserImageUri,
		&i.Weight,
		&i.Height,
		&i.HeightUnit,
		&i.WeightUnit,
		&i.Preference,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET
  username = $1,
  user_image_uri = $2,
  weight = $3::int,
  height = $4::int,
  weight_unit = $5::enum_weight_units,
  height_unit = $6::enum_height_units,
  preference = $7::enum_preferences
WHERE
  id = $8
RETURNING id, email, hashed_password, username, user_image_uri, weight, height, height_unit, weight_unit, preference
`

type UpdateUserParams struct {
	Username     string
	UserImageUri string
	Weight       int
	Height       int
	WeightUnit   model.EnumWeightUnits
	HeightUnit   model.EnumHeightUnits
	Preference   model.EnumPreferences
	ID           int
}

func (r *UserRepository) UpdateUser(ctx context.Context, arg UpdateUserParams) (model.User, error) {
	row := r.pool.QueryRow(ctx, updateUser,
		arg.Username,
		arg.UserImageUri,
		arg.Weight,
		arg.Height,
		arg.WeightUnit,
		arg.HeightUnit,
		arg.Preference,
		arg.ID,
	)
	var i model.User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.HashedPassword,
		&i.Username,
		&i.UserImageUri,
		&i.Weight,
		&i.Height,
		&i.HeightUnit,
		&i.WeightUnit,
		&i.Preference,
	)
	return i, err
}
