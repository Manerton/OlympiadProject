package register_dto

import (
	"encoding/json"
)

type RegisterParticipantRequestDTO struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstname"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birthdate"`
	Gender      int    `json:"gender"`

	Disability  int    `json:"disability"`
	SchoolId    string `json:"school_id"`
	Citizenship int    `json:"citizenship"`
	ClassNumber int    `json:"class_number"`
}

type RegisterUserRequestDTO struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstname"`
	Surname     string `json:"surname"`
	Patronymic  string `json:"patronymic"`
	PhoneNumber string `json:"phone_number"`
	BirthDate   string `json:"birthdate"`
	Gender      int    `json:"gender"`
	Role        int    `json:"role"`
}

func (dto *RegisterUserRequestDTO) UnmarshalJSON(data []byte) error {
	// Временная структура с полями-строками
	type Alias struct {
		Email       string      `json:"email"`
		Password    string      `json:"password"`
		FirstName   string      `json:"firstname"`
		Surname     string      `json:"surname"`
		Patronymic  string      `json:"patronymic"`
		PhoneNumber string      `json:"phone_number"`
		BirthDate   string      `json:"birthdate"`
		Gender      json.Number `json:"gender"` // Принимаем как строку
		Role        json.Number `json:"role"`   // Принимаем как строку
	}

	temp := &Alias{}
	if err := json.Unmarshal(data, temp); err != nil {
		return err
	}

	println(temp)

	// Конвертируем строки в числа
	dto.Email = temp.Email
	dto.Password = temp.Password
	dto.FirstName = temp.FirstName
	dto.Surname = temp.Surname
	dto.Patronymic = temp.Patronymic
	dto.PhoneNumber = temp.PhoneNumber
	dto.BirthDate = temp.BirthDate
	gender, _ := temp.Gender.Int64()
	role, _ := temp.Gender.Int64()

	dto.Gender = int(gender)
	dto.Role = int(role)

	return nil
}
