package responcetypes

type ApiResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Message    string      `json:"message,omitempty"`
	Metadata   interface{} `json:"metadata,omitempty"`
}

type ApplicationEvent struct {
	ApplicationID      string   `json:"application_id"`
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	StartDate          string   `json:"start_date"`
	EndDate            string   `json:"end_date"`
	PreviousEvent      string   `json:"previous_event_id"`
	Dates              []string `json:"dates"`
	Profile            string   `json:"profile"`
	Subject            int      `json:"subject"`
	ClassNumber        int      `json:"class_number"`
	AdditionalInfo     string   `json:"additional_info"`
	Status             int      `json:"status"`
	ClassParticipation int      `json:"class_participation"`
}

type Participant struct {
	ID          string `json:"id"`
	UserId      string `json:"user_id"`
	Disability  string `json:"disability"`
	SchoolId    string `json:"school_id"`
	Citizenship string `json:"citizenship"`
	ClassNumber string `json:"class_number"`
}

//type Event struct {
//	ID      string   `json:"id"`
//	Name    string   `json:"name"`
//	Subject int      `json:"subject"`
//	Dates   []string `json:"dates"`
//}
//
//type ApplicationResponse struct {
//	ID                 string `json:"id"`
//	UserID             string `json:"userId"`
//	SchoolID           string `json:"schoolId"`
//	EventID            string `json:"eventId"`
//	Profile            string `json:"profile"`
//	ClassParticipation int    `json:"class_participation"`
//	Status             int    `json:"status"`
//}

type JuryNames struct {
	ID     string `json:"id"`
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}

type SchoolWithDistrict struct {
	School   SchoolResponse   `json:"school"`
	District DistrictResponse `json:"district"`
}

// DistrictResponse структура ответа для муниципалитетов
type DistrictResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Region int    `json:"region"`
}

// SchoolResponse структура ответа для школ
type SchoolResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	DistrictID string `json:"district_id"`
}

// UserInfo структура ответа для пользователей
type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Firstname     string `json:"firstname"`
	Surname       string `json:"surname"`
	Patronymic    string `json:"patronymic"`
	PhoneNumber   string `json:"phone_number"`
	Birthdate     string `json:"birthdate"`
	Gender        int    `json:"gender"`
	Role          int    `json:"role"`
	Activated     bool   `json:"activated"`
	ParticipantID string `json:"participant_id"`
	UserID        string `json:"user_id"`
	Disability    int    `json:"disability"`
	SchoolID      string `json:"school_id"`
	Citizenship   int    `json:"citizenship"`
	ClassNumber   int    `json:"class_number"`
}

// ApplicationResponse структура ответа для заявок
type ApplicationResponse struct {
	ID                 string `json:"id"`
	UserID             string `json:"userId"`
	SchoolID           string `json:"schoolId"`
	EventID            string `json:"eventId"`
	Profile            string `json:"profile"`
	ClassParticipation int    `json:"class_participation"`
	Status             int    `json:"status"`
	Reason             int    `json:"reason"`
	Code               string `json:"code"`
	SubmittedAt        string `json:"submittedAt"`
	UpdatedAt          string `json:"updatedAt"`
}

// Event структура ответа для событий
type Event struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	StartDate       string   `json:"start_date"`
	EndDate         string   `json:"end_date"`
	Dates           []string `json:"dates"`
	Profiles        []string `json:"profiles"`
	EventType       string   `json:"event_type"`
	ClassCategory   string   `json:"class_category"`
	PreviousEventID *string  `json:"previous_event_id"`
	Subject         int      `json:"subject"`
	AdditionalInfo  string   `json:"additional_info"`
	Status          int      `json:"status"`
	Events          []string `json:"events"`
}
