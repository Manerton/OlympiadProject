package register_dto

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
	City        string `json:"city"`
	Reason      string `json:"reason"`
	Citizenship string `json:"citizenship"`
	ClassNumber int    `json:"class_number"`
}

type RegusterUserRequestDTO struct {
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
