package response

import (
	"main/internal/lib/errs"
	"net/http"
)

const (
	SUCCESS = "success"
	ERROR   = "error"
)

type ApiResponse struct {
	Status     string         `json:"status"`
	StatusCode int            `json:"status_code"`
	Data       interface{}    `json:"data,omitempty"`
	Error      *errs.ApiError `json:"error,omitempty"`
	Message    string         `json:"message,omitempty"`
	Metadata   interface{}    `json:"metadata,omitempty"`
}

func ErrorResponse(msg string, err error) ApiResponse {

	apiErr := err.(*errs.ApiError)

	return ApiResponse{
		Status:     ERROR,
		StatusCode: apiErr.HttpCode,
		Error:      apiErr,
		Message:    msg,
	}
}

func SuccessResponse(msg string) ApiResponse {
	return ApiResponse{
		Status:     SUCCESS,
		StatusCode: http.StatusOK,
		Message:    msg,
	}
}
