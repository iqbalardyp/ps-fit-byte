package repository

import (
	"context"
	"fit-byte/internal/activity/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ActivityRepository struct {
	pool *pgxpool.Pool
}

func NewActivityRepository(pool *pgxpool.Pool) *ActivityRepository {
	return &ActivityRepository{pool: pool}
}

const createActivity = `-- name: CreateActivity :one
INSERT INTO activities (
  activity_type,
  done_at,
  duration_in_minutes,
  calories_burned,
  user_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id, user_id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, updated_at
`

type CreateActivityParams struct {
	ActivityType      model.ActivityTypeEnum
	DoneAt            time.Time
	DurationInMinutes int
	CaloriesBurned    int
	UserId            int
}

func (r *ActivityRepository) CreateActivity(ctx context.Context, arg CreateActivityParams) (model.Activity, error) {
	row := r.pool.QueryRow(ctx, createActivity,
		arg.ActivityType,
		arg.DoneAt,
		arg.DurationInMinutes,
		arg.CaloriesBurned,
		arg.UserId,
	)
	var i model.Activity
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.ActivityType,
		&i.DoneAt,
		&i.DurationInMinutes,
		&i.CaloriesBurned,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listActivities = `-- name: ListActivities :many
SELECT id, user_id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, updated_at FROM activities
WHERE ($3::enum_activity_types IS NULL OR activity_type = $3::enum_activity_types)
  AND ($4::timestamptz IS NULL OR done_at >= $4::timestamptz)
  AND ($5::timestamptz IS NULL OR done_at <= $5::timestamptz)
  AND ($6::int IS NULL OR calories_burned >= $6::int)
  AND ($7::int IS NULL OR calories_burned <= $7::int)
  AND (user_id = $8::bigint)
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListActivitiesParams struct {
	Limit             int
	Offset            int
	ActivityType      *model.ActivityTypeEnum
	DoneAtFrom        *time.Time
	DoneAtTo          *time.Time
	CaloriesBurnedMin *int
	CaloriesBurnedMax *int
	UserId            int
}

func (r *ActivityRepository) ListActivities(ctx context.Context, arg ListActivitiesParams) ([]model.Activity, error) {
	rows, err := r.pool.Query(ctx, listActivities,
		arg.Limit,
		arg.Offset,
		arg.ActivityType,
		arg.DoneAtFrom,
		arg.DoneAtTo,
		arg.CaloriesBurnedMin,
		arg.CaloriesBurnedMax,
		arg.UserId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.Activity
	for rows.Next() {
		var i model.Activity
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.ActivityType,
			&i.DoneAt,
			&i.DurationInMinutes,
			&i.CaloriesBurned,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getActivity = `-- name: GetActivity :one
SELECT id, user_id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, updated_at FROM activities
WHERE (id = $1::bigint)
  AND (user_id = $2::bigint)
LIMIT 1
`

type GetAndDeleteActivityParams struct {
	Id     int
	UserId int
}

func (r *ActivityRepository) GetActivity(ctx context.Context, arg GetAndDeleteActivityParams) (model.Activity, error) {
	row := r.pool.QueryRow(ctx, getActivity, arg.Id, arg.UserId)
	var i model.Activity
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.ActivityType,
		&i.DoneAt,
		&i.DurationInMinutes,
		&i.CaloriesBurned,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateActivity = `-- name: UpdateActivity :one
UPDATE activities
SET
  activity_type = $1,
  done_at = $2,
  duration_in_minutes = $3,
  calories_burned = $4
WHERE
  (id = $5::bigint)
  AND (user_id = $6::bigint)
RETURNING id, user_id, activity_type, done_at, duration_in_minutes, calories_burned, created_at, updated_at
`

type UpdateActivityParams struct {
	ActivityType      model.ActivityTypeEnum
	DoneAt            time.Time
	DurationInMinutes int
	CaloriesBurned    int
	Id                int
	UserId            int
}

func (r *ActivityRepository) UpdateActivity(ctx context.Context, arg UpdateActivityParams) (model.Activity, error) {
	row := r.pool.QueryRow(ctx, updateActivity,
		arg.ActivityType,
		arg.DoneAt,
		arg.DurationInMinutes,
		arg.CaloriesBurned,
		arg.Id,
		arg.UserId,
	)
	var i model.Activity
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.ActivityType,
		&i.DoneAt,
		&i.DurationInMinutes,
		&i.CaloriesBurned,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteActivity = `-- name: DeleteActivity :exec
DELETE FROM activities 
WHERE
  (id = $1::bigint)
  AND (user_id = $2::bigint)
`

func (r *ActivityRepository) DeleteActivity(ctx context.Context, arg GetAndDeleteActivityParams) error {
	_, err := r.pool.Exec(ctx, deleteActivity, arg.Id, arg.UserId)
	return err
}
