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
	ID             string `json:"id"`
	Name           string `json:"name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	PreviousEvent  string `json:"previous_event_id"`
	Subject        string `json:"subject"`
	ClassNumber    int    `json:"class_number"`
	AdditionalInfo string `json:"additional_info"`
	Status         int    `json:"status"`
}
