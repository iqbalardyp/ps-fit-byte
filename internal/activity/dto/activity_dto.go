package dto

import (
	"fit-byte/internal/activity/model"
	"time"
)

type ActivityResponse struct {
	ActivityId        string                 `json:"activityId"`
	ActivityType      model.ActivityTypeEnum `json:"activityType"`
	DoneAt            string                 `json:"doneAt"`
	DurationInMinutes int                    `json:"durationInMinutes"`
	CaloriesBurned    int                    `json:"caloriesBurned"`
	CreatedAt         string                 `json:"createdAt"`
	UpdatedAt         string                 `json:"updatedAt"`
}

type CreateAndUpdateActivityRequest struct {
	ActivityType      model.ActivityTypeEnum `json:"activityType" validate:"required,activity_type"`
	DoneAt            time.Time              `json:"doneAt" validate:"required,time_validator"`
	DurationInMinutes int                    `json:"durationInMinutes" validate:"required,min=1"`
}

type GetActivityRequest struct {
	Limit             int                     `query:"limit" validate:"omitempty,min=0"`
	Offset            int                     `query:"offset" validate:"omitempty,min=0"`
	ActivityType      *model.ActivityTypeEnum `query:"activityType" validate:"omitempty,activity_type"`
	DoneAtFrom        *time.Time              `query:"doneAtFrom" validate:"omitempty,time_validator"`
	DoneAtTo          *time.Time              `query:"doneAtTo" validate:"omitempty,time_validator"`
	CaloriesBurnedMin *int                    `query:"caloriesBurnedMin" validate:"omitempty,min=0"`
	CaloriesBurnedMax *int                    `query:"caloriesBurnedMax" validate:"omitempty,min=0"`
}
