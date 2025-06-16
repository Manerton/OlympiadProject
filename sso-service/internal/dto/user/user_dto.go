package user_dto

import participant_dto "main/internal/dto/participant"

type ParticipantUserResponseDTO struct {
	UserResponseDTO
	participant_dto.ParticipantResponseDTO
}

type UserResponseDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstname"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birthdate"`
	Gender      int    `json:"gender"`
	Role        int    `json:"role"`
	Activated   bool   `json:"activated"`
}

type UpdateUserRequestDTO struct {
	Email       *string `json:"email"`
	FirstName   *string `json:"firstname"`
	Surname     *string `json:"surname"`
	Patronymic  *string `json:"patronymic"`
	PhoneNumber *string `json:"phone_number"`
	BirthDate   *string `json:"birthdate"`
	Gender      *int    `json:"gender"`
}
