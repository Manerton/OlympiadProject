package response

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type ApiResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Message    string      `json:"message,omitempty"`
	Metadata   interface{} `json:"metadata,omitempty"`
}

func ErrorResponse(msg string) ApiResponse {
	return ApiResponse{
		Status:     ERROR,
		StatusCode: 400,
		Message:    msg,
	}
}

func SuccessResponse(msg string) ApiResponse {
	return ApiResponse{
		Status:     SUCCESS,
		StatusCode: 200,
		Message:    msg,
	}
}
