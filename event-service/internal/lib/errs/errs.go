package errs

import "fmt"

type ApiError struct {
	Code     string
	HttpCode int
	Message  string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *ApiError) Wrap(message string) *ApiError {
	return &ApiError{
		Code:     e.Code,
		HttpCode: e.HttpCode,
		Message:  message,
	}
}

func IsApiError(err error) (*ApiError, bool) {
	apiErr, ok := err.(*ApiError)
	return apiErr, ok
}

var ()
