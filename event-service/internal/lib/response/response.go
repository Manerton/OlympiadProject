package response

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Events     interface{} `json:"events"`
	TotalCount int         `json:"totalCount"`
}

func Error(err string) ApiResponse {
	return ApiResponse{
		Status: StatusError,
		Error:  err,
	}
}

func Success(msg string) ApiResponse {
	return ApiResponse{
		Status:  StatusOK,
		Message: msg,
	}
}
