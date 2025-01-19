package usecase

import (
	"context"

	"fit-byte/internal/activity/dto"
	"fit-byte/internal/activity/model"
	"fit-byte/internal/activity/repository"
	customErrors "fit-byte/pkg/custom-errors"

	"github.com/pkg/errors"
)

type ActivityUseCase struct {
	activityRepo repository.ActivityRepository
}

func NewActivityUseCase(activityRepo repository.ActivityRepository) *ActivityUseCase {
	return &ActivityUseCase{
		activityRepo,
	}
}

func (c *ActivityUseCase) GetActivity(ctx context.Context, request *dto.GetActivityRequest, userid int) (*[]model.Activity, error) {

	arg := repository.ListActivitiesParams{
		Limit:             request.Limit,
		Offset:            request.Offset,
		ActivityType:      request.ActivityType,
		DoneAtFrom:        request.DoneAtFrom,
		DoneAtTo:          request.DoneAtTo,
		CaloriesBurnedMin: request.CaloriesBurnedMin,
		CaloriesBurnedMax: request.CaloriesBurnedMax,
		UserId:            userid,
	}

	activities, err := c.activityRepo.ListActivities(ctx, arg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get activities")
	}

	return &activities, nil
}

var activityCalories = map[model.ActivityTypeEnum]int{
	model.ActivityTypeEnumWalking:    4,
	model.ActivityTypeEnumYoga:       4,
	model.ActivityTypeEnumStretching: 4,
	model.ActivityTypeEnumCycling:    8,
	model.ActivityTypeEnumSwimming:   8,
	model.ActivityTypeEnumDancing:    8,
	model.ActivityTypeEnumHiking:     10,
	model.ActivityTypeEnumRunning:    10,
	model.ActivityTypeEnumHIIT:       10,
	model.ActivityTypeEnumJumpRope:   10,
}

func calculateCalories(activityType model.ActivityTypeEnum, duration int) int {
	caloriesPerMinute := activityCalories[activityType]
	return caloriesPerMinute * duration
}

func (c *ActivityUseCase) CreateActivity(ctx context.Context, request *dto.CreateAndUpdateActivityRequest, userId int) (*model.Activity, error) {

	caloriesBurned := calculateCalories(request.ActivityType, request.DurationInMinutes)
	arg := repository.CreateActivityParams{
		ActivityType:      request.ActivityType,
		DoneAt:            request.DoneAt,
		DurationInMinutes: request.DurationInMinutes,
		CaloriesBurned:    caloriesBurned,
		UserId:            userId,
	}

	activity, err := c.activityRepo.CreateActivity(ctx, arg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create activity")
	}

	return &activity, nil
}

func (c *ActivityUseCase) UpdateActivity(ctx context.Context, request *dto.CreateAndUpdateActivityRequest, userId int, activityId int) (*model.Activity, error) {

	caloriesBurned := calculateCalories(request.ActivityType, request.DurationInMinutes)
	arg := repository.UpdateActivityParams{
		ActivityType:      request.ActivityType,
		DoneAt:            request.DoneAt,
		DurationInMinutes: request.DurationInMinutes,
		CaloriesBurned:    caloriesBurned,
		Id:                activityId,
		UserId:            userId,
	}

	activity, err := c.activityRepo.UpdateActivity(ctx, arg)
	if err != nil {
		if errors.Is(err, customErrors.ErrNotFound) {
			return nil, errors.Wrap(customErrors.ErrNotFound, "Activity not found")
		}
		return nil, errors.Wrap(err, "failed to update activity")
	}

	return &activity, nil
}

func (c *ActivityUseCase) DeleteActivity(ctx context.Context, userId int, activityId int) error {

	arg := repository.GetAndDeleteActivityParams{
		Id:     activityId,
		UserId: userId,
	}

	_, err := c.activityRepo.GetActivity(ctx, arg)
	if err != nil {
		if errors.Is(err, customErrors.ErrNotFound) {
			return errors.Wrap(customErrors.ErrNotFound, "Activity not found")
		}
		return errors.Wrap(err, "failed to get activity")
	}

	err = c.activityRepo.DeleteActivity(ctx, arg)
	if err != nil {
		return errors.Wrap(err, "failed to delete activity")
	}

	return nil
}
