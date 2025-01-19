package user_dto

type AuthRequestParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type AuthUser struct {
	ID             int
	HashedPassword string
}

type AuthResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type User struct {
	Name       *string `json:"name"`
	ImageURI   *string `json:"imageUri"`
	Height     *int    `json:"height"`
	HeightUnit *string `json:"heightUnit"`
	Weight     *int    `json:"weight"`
	WeightUnit *string `json:"weightUnit"`
	Preference *string `json:"preference"`
}

type GetUserResponse struct {
	Name       *string `json:"name"`
	ImageURI   *string `json:"imageUri"`
	Height     *int    `json:"height"`
	HeightUnit *string `json:"heightUnit"`
	Weight     *int    `json:"weight"`
	WeightUnit *string `json:"weightUnit"`
	Preference *string `json:"preference"`
	Email      string  `json:"email"`
}

type UpdateUserParams struct {
	Name       *string `json:"name" validate:"required,min=2,max=60"`
	ImageURI   *string `json:"imageUri" validate:"required,is_uri"`
	Height     *int    `json:"height"  validate:"required,min=3,max=250"`
	HeightUnit *string `json:"heightUnit" validate:"required,oneof=CM INCH"`
	Weight     *int    `json:"weight" validate:"required,min=10,max=1000"`
	WeightUnit *string `json:"weightUnit" validate:"required,oneof=KG LBS"`
	Preference *string `json:"preference" validate:"required,oneof=CARDIO WEIGHT"`
}
