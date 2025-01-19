package converter

import (
	"fit-byte/internal/activity/dto"
	"fit-byte/internal/activity/model"
	"fit-byte/pkg/helper"
	"strconv"
)

func ToActivityResponse(activity model.Activity) dto.ActivityResponse {
	return dto.ActivityResponse{
		ActivityId:        strconv.Itoa(activity.ID),
		ActivityType:      activity.ActivityType,
		DoneAt:            helper.FormatTimeToUTC(activity.DoneAt),
		DurationInMinutes: activity.DurationInMinutes,
		CaloriesBurned:    activity.CaloriesBurned,
		CreatedAt:         helper.FormatTimeToUTC(activity.CreatedAt),
		UpdatedAt:         helper.FormatTimeToUTC(activity.UpdatedAt),
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
