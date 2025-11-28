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

type Event struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Subject int      `json:"subject"`
	Dates   []string `json:"dates"`
}

type ApplicationResponse struct {
	ID                 string `json:"id"`
	UserID             string `json:"userId"`
	SchoolID           string `json:"schoolId"`
	EventID            string `json:"eventId"`
	Profile            string `json:"profile"`
	ClassParticipation int    `json:"class_participation"`
	Status             int    `json:"status"`
}

type JuryNames struct {
	ID     string `json:"id"`
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}
