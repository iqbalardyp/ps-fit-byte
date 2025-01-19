package repository

import (
	"context"
	"fit-byte/internal/activity/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
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

const queryUpdateActivity = `
	WITH
	payload AS (
		SELECT
			t.type,
			t.duration,
			t.calories_burned,
			t.done_at,
			t.updated_at
		FROM (
			VALUES (
				@type::enum_activity_types,
				@duration::INT,
				@calories_burned::INT,
				@done_at::TIMESTAMPTZ,
				@updated_at::TIMESTAMPTZ
			)
		) AS t(
			type,
			duration,
			calories_burned,
			done_at,
			updated_at
		)
	)
	UPDATE activities
	SET
		activity_type = COALESCE(payload.type, activities.activity_type),
		duration_in_minutes = COALESCE(payload.duration, activities.duration_in_minutes),
		calories_burned = COALESCE(payload.calories_burned, activities.calories_burned),
		done_at = COALESCE(payload.done_at, activities.done_at),
		updated_at = COALESCE(payload.updated_at, activities.updated_at)
	FROM payload
	WHERE
		activities.id = @activitiesId
	RETURNING
		activities.id,
		activities.activity_type,
		activities.done_at,
		activities.duration_in_minutes,
		activities.calories_burned,
		activities.created_at,
		activities.updated_at
		;
`

type PatchActivitiesParams struct {
	ActivityType      model.ActivityTypeEnum
	DoneAt            time.Time
	UpdatedAt         time.Time
	DurationInMinutes int
	CaloriesBurned    int
	ActivityId        int
	UserId            int
}

func (r *ActivityRepository) UpdateActivityRepo(ctx context.Context, arg PatchActivitiesParams) (*model.Activity, error) {
	var activity model.Activity

	args := pgx.NamedArgs{
		"type":            arg.ActivityType,
		"duration":        arg.DurationInMinutes,
		"calories_burned": arg.CaloriesBurned,
		"done_at":         arg.DoneAt,
		"updated_at":      arg.UpdatedAt,
		"activitiesId":    arg.ActivityId,
	}

	err := r.pool.QueryRow(ctx, queryUpdateActivity, args).Scan(
		&activity.ID,
		&activity.ActivityType,
		&activity.DoneAt,
		&activity.DurationInMinutes,
		&activity.CaloriesBurned,
		&activity.CreatedAt,
		&activity.UpdatedAt,
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to execute Update statements")
	}

	return &activity, nil
}

const queryDeleteActivity = `
DELETE FROM Activities WHERE user_id = @user_id AND id = @activityId
`

type DeleteActivitiesParams struct {
	ActivityId int
	UserId     int
}

func (r *ActivityRepository) DeleteActivity(ctx context.Context, arg DeleteActivitiesParams) error {
	args := pgx.NamedArgs{
		"user_id": arg.UserId,
		"id":      arg.ActivityId,
	}

	_, err := r.pool.Exec(ctx, queryDeleteActivity, args)
	if err != nil {
		return errors.Wrap(err, "failed to execute delete statements")
	}

	return nil
}
