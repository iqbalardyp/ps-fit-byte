package dto

import (
	"fit-byte/internal/user/model"
)

type AuthRequest struct {
	Password string `json:"password" validate:"required,min=8,max=32"`
	Email    string `json:"email" validate:"required,email,min=1,max=255"`
}

type AuthResponse struct {
	Email       string `json:"email"`
	AccessToken string `json:"token"`
}

type UserResponse struct {
	Email        string                `json:"email"`
	Username     string                `json:"name"`
	UserImageUri string                `json:"imageUri"`
	Weight       int                   `json:"weight"`
	Height       int                   `json:"height"`
	WeightUnit   model.EnumWeightUnits `json:"weightUnit"`
	HeightUnit   model.EnumHeightUnits `json:"heightUnit"`
	Preference   model.EnumPreferences `json:"preference"`
}

type UpdateUserRequest struct {
	Username     string                `json:"name" validate:"required,min=2,max=60"`
	UserImageUri string                `json:"imageUri" validate:"required,valid_uri"`
	Weight       int                   `json:"weight" validate:"required,min=10,max=1000"`
	Height       int                   `json:"height" validate:"required,min=3,max=250"`
	WeightUnit   model.EnumWeightUnits `json:"weightUnit" validate:"required,oneof=KG LBS"`
	HeightUnit   model.EnumHeightUnits `json:"heightUnit" validate:"required,oneof=CM INCH"`
	Preference   model.EnumPreferences `json:"preference" validate:"required,oneof=CARDIO WEIGHT"`
}
