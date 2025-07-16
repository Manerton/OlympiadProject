package participant_dto

type ParticipantResponseDTO struct {
	ID          string `json:"id"`
	UserId      string `json:"user_id"`
	Disability  string `json:"disability"`
	SchoolId    string `json:"school_id"`
	Citizenship string `json:"citizenship"`
	ClassNumber string `json:"class_number"`
}

type UpdateParticipantRequestDTO struct {
	Disability  *string `json:"disability"`
	SchoolId    *string `json:"school_id"`
	Citizenship *string `json:"citizenship"`
	ClassNumber *string `json:"class_number"`
}
