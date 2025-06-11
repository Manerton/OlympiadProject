package register_dto

import "time"

type RegisterParticipantRequestDTO struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"firstname"`
	Surname     string    `json:"surname"`
	Patronymic  string    `json:"patronymic"`
	PhoneNumber string    `json:"phone_number"`
	BirthDate   time.Time `json:"birth_date"`
	Gender      string    `json:"gender"`

	OVZ         int    `json:"ovz"`
	SchoolName  string `json:"school_name"`
	City        string `json:"city"`
	Reason      string `json:"reason"`
	Citizenship string `json:"citizenship"`
	ClassNumber int    `json:"class_number"`
}

type RegusterUserRequestDTO struct {
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"firstname"`
	Surname     string    `json:"surname"`
	Patronymic  string    `json:"patronymic"`
	PhoneNumber string    `json:"phone_number"`
	BirthDate   time.Time `json:"birth_date"`
	Gender      string    `json:"gender"`
	Role        string    `json:"role"`
}
