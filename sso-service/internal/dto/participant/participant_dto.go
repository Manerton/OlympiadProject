package participant_dto

type ParticipantResponseDTO struct {
	// ID          string `json:"id"`
	UserId      string `json:"user_id"`
	OVZ         int    `json:"ovz"`
	SchoolId    string `json:"school_id"`
	City        string `json:"city"`
	Reason      string `json:"reason"`
	Citizenship string `json:"citizenship"`
	ClassNumber int    `json:"class_number"`
}

type UpdateParticipantRequestDTO struct {
	OVZ         *string `json:"ovz"`
	SchoolId    *string `json:"school_id"`
	City        *string `json:"city"`
	Reason      *string `json:"reason"`
	Citizenship *string `json:"citizenship"`
	ClassNumber *int    `json:"class_number"`
}
