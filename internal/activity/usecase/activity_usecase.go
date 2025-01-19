package usecase

import (
	"context"
	"time"

	"fit-byte/internal/activity/dto"
	"fit-byte/internal/activity/model"
	"fit-byte/internal/activity/repository"

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

func (c *ActivityUseCase) CreateActivity(ctx context.Context, request *dto.CreateAndUpdateActivityRequest, userid int) (*model.Activity, error) {

	caloriesBurned := calculateCalories(request.ActivityType, request.DurationInMinutes)
	arg := repository.CreateActivityParams{
		ActivityType:      request.ActivityType,
		DoneAt:            request.DoneAt,
		DurationInMinutes: request.DurationInMinutes,
		CaloriesBurned:    caloriesBurned,
		UserId:            userid,
	}

	activity, err := c.activityRepo.CreateActivity(ctx, arg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create activity")
	}

	return &activity, nil
}

func (c *ActivityUseCase) UpdateActivity(ctx context.Context, request *dto.CreateAndUpdateActivityRequest,activityId int, userId int)(*model.Activity, error){
	caloriesBurned := calculateCalories(request.ActivityType, request.DurationInMinutes)
	timeNow := time.Now()
	arg := repository.PatchActivitiesParams{
		ActivityType: request.ActivityType,
		DoneAt: request.DoneAt,
		UpdatedAt: timeNow,
		DurationInMinutes: request.DurationInMinutes,
		CaloriesBurned: caloriesBurned,
		ActivityId: activityId,
		UserId: userId,
	}

	activity, err := c.activityRepo.UpdateActivityRepo(ctx,arg)
	if err != nil {
		return nil, errors.Wrap(err,"failed to update Activity")
	}

	return activity,nil
}

func (c *ActivityUseCase) DeleteActivity(ctx context.Context, activityId int , userId int) error{

	arg := repository.DeleteActivitiesParams{
		ActivityId: activityId,
		UserId: userId,
	}

	return c.activityRepo.DeleteActivity(ctx, arg)
}
