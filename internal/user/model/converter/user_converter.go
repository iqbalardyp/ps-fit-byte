package converter

import (
	"fit-byte/internal/user/dto"
	"fit-byte/internal/user/model"
	"fit-byte/pkg/helper"
)

func ToUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		Email:        user.Email,
		Username:     helper.DerefString(user.Username, ""),
		UserImageUri: helper.DerefString(user.UserImageUri, ""),
		Weight:       helper.DerefInt(user.Weight, 0),
		Height:       helper.DerefInt(user.Height, 0),
		WeightUnit:   helper.DerefGeneric[model.EnumWeightUnits](user.WeightUnit, ""),
		HeightUnit:   helper.DerefGeneric[model.EnumHeightUnits](user.HeightUnit, ""),
		Preference:   helper.DerefGeneric[model.EnumPreferences](user.Preference, ""),
	}
}
