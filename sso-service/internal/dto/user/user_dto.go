package user_dto

type UserResponseDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birth_date"`
	Sex         string `json:"sex"`
	Role        string `json:"role"`
}

type UpdateUserRequestDTO struct {
	Email       *string `json:"email"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Surname     *string `json:"surname"`
	PhoneNumber *string `json:"phone_number"`
	BirthDate   *string `json:"birth_date"`
	Sex         *string `json:"sex"`
}
