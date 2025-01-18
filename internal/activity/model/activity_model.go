package model

import (
	"fmt"
	"time"
)

type ActivityTypeEnum string

const (
	ActivityTypeEnumWalking    ActivityTypeEnum = "Walking"
	ActivityTypeEnumYoga       ActivityTypeEnum = "Yoga"
	ActivityTypeEnumStretching ActivityTypeEnum = "Stretching"
	ActivityTypeEnumCycling    ActivityTypeEnum = "Cycling"
	ActivityTypeEnumSwimming   ActivityTypeEnum = "Swimming"
	ActivityTypeEnumDancing    ActivityTypeEnum = "Dancing"
	ActivityTypeEnumHiking     ActivityTypeEnum = "Hiking"
	ActivityTypeEnumRunning    ActivityTypeEnum = "Running"
	ActivityTypeEnumHIIT       ActivityTypeEnum = "HIIT"
	ActivityTypeEnumJumpRope   ActivityTypeEnum = "JumpRope"
)

func (e *ActivityTypeEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ActivityTypeEnum(s)
	case string:
		*e = ActivityTypeEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for ActivityTypeEnum: %T", src)
	}
	return nil
}

type Activity struct {
	ID                int
	UserId            int
	ActivityType      ActivityTypeEnum
	DoneAt            time.Time
	DurationInMinutes int
	CaloriesBurned    int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
