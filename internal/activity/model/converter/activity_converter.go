package converter

import (
	"fit-byte/internal/activity/dto"
	"fit-byte/internal/activity/model"
)

func ToActivityResponse(activity model.Activity) dto.ActivityResponse {
	return dto.ActivityResponse{
		ActivityId:        activity.ID,
		ActivityType:      activity.ActivityType,
		DoneAt:            activity.DoneAt,
		DurationInMinutes: activity.DurationInMinutes,
		CaloriesBurned:    activity.CaloriesBurned,
		CreatedAt:         activity.CreatedAt,
		UpdatedAt:         activity.UpdatedAt,
	}
}

func ToActivityResponseList(activities []model.Activity) []dto.ActivityResponse {
	// Ensure we always return an empty slice, not nil
	if activities == nil {
		return []dto.ActivityResponse{}
	}

	var responses []dto.ActivityResponse
	for _, activity := range activities {
		responses = append(responses, ToActivityResponse(activity))
	}
	return responses
}
