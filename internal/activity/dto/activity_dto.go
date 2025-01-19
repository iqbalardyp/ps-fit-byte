package dto

import (
	"fit-byte/internal/activity/model"
	"time"
)

type ActivityResponse struct {
	ActivityId        int                    `json:"activityId"`
	ActivityType      model.ActivityTypeEnum `json:"activityType"`
	DoneAt            time.Time              `json:"doneAt"`
	DurationInMinutes int                    `json:"durationInMinutes"`
	CaloriesBurned    int                    `json:"caloriesBurned"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

type CreateAndUpdateActivityRequest struct {
	ActivityType      model.ActivityTypeEnum `json:"activityType" validate:"required,activity_type"`
	DoneAt            time.Time              `json:"doneAt" validate:"required,time_validator"`
	DurationInMinutes int                    `json:"durationInMinutes" validate:"required,min=1"`
}

type GetActivityRequest struct {
	Limit             int                     `json:"limit" validate:"omitempty,min=0"`
	Offset            int                     `json:"offset" validate:"omitempty,min=0"`
	ActivityType      *model.ActivityTypeEnum `json:"activityType" validate:"omitempty,activity_type"`
	DoneAtFrom        *time.Time              `json:"doneAtFrom" validate:"omitempty,time_validator"`
	DoneAtTo          *time.Time              `json:"doneAtTo" validate:"omitempty,time_validator"`
	CaloriesBurnedMin *int                    `json:"caloriesBurnedMin" validate:"omitempty,min=0"`
	CaloriesBurnedMax *int                    `json:"caloriesBurnedMax" validate:"omitempty,min=0"`
}
