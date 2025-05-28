package register_dto

type RegisterParticipantRequestDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Surname   string `json:"surname"`

	OVZ         string `json:"ovz"`
	SchoolName  string `json:"school_name"`
	City        string `json:"city"`
	Reason      string `json:"reason"`
	Citizenship string `json:"citizenship"`
	ClassNumber int    `json:"class_number"`
}
