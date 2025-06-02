package user_dto

import participant_dto "main/internal/dto/participant"

type ParticipantUserResponseDTO struct {
	UserResponseDTO
	participant_dto.ParticipantResponseDTO
}

type UserResponseDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birth_date"`
	Gender      string `json:"gender"`
	Role        string `json:"role"`
	Status      bool   `json:"status"`
}

type UpdateUserRequestDTO struct {
	Email       *string `json:"email"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Surname     *string `json:"surname"`
	PhoneNumber *string `json:"phone_number"`
	BirthDate   *string `json:"birth_date"`
	Gender      *string `json:"gender"`
}
