package config

import (
	"github.com/go-playground/validator/v10"

	"fit-byte/internal/activity/model"
	"time"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("activity_type", activityTypeValidator)
	validate.RegisterValidation("time_validator", timeValidator)
	return validate
}

var validActivityTypes = map[model.ActivityTypeEnum]struct{}{
	model.ActivityTypeEnumWalking:    {},
	model.ActivityTypeEnumYoga:       {},
	model.ActivityTypeEnumStretching: {},
	model.ActivityTypeEnumCycling:    {},
	model.ActivityTypeEnumSwimming:   {},
	model.ActivityTypeEnumDancing:    {},
	model.ActivityTypeEnumHiking:     {},
	model.ActivityTypeEnumRunning:    {},
	model.ActivityTypeEnumHIIT:       {},
	model.ActivityTypeEnumJumpRope:   {},
}

func activityTypeValidator(fl validator.FieldLevel) bool {
	activity, ok := fl.Field().Interface().(model.ActivityTypeEnum)
	if !ok {
		return false
	}
	_, exists := validActivityTypes[activity]
	return exists
}

func timeValidator(fl validator.FieldLevel) bool {
	t, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// Example: Check if time is after 1970-01-01T00:00:00Z
	return !t.IsZero() && t.After(time.Unix(0, 0))
}
