package errs

import (
	"fmt"
	"net/http"
)

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

// Common errors
var (
	ErrInternalError = &ApiError{Code: "INTERNAL_SERVER_ERROR", HttpCode: http.StatusInternalServerError}
	ErrBadRequest    = &ApiError{Code: "BAD_REQUEST", HttpCode: http.StatusBadRequest}
)

// Token errors
var (
	ErrAccessToken  = &ApiError{Code: "ACCESS_TOKEN_GENERATION_FAILED", HttpCode: http.StatusInternalServerError, Message: "Ошибка при создании access токена"}
	ErrRefreshToken = &ApiError{Code: "REFRESH_TOKEN_GENERATION_FAILED", HttpCode: http.StatusInternalServerError, Message: "Ошибка при создании refresh токена"}
)

// Auth errors
var (
	ErrAuthFailed             = &ApiError{Code: "AUTH_FAILED", HttpCode: http.StatusUnauthorized}
	ErrUserAlreadyExists      = &ApiError{Code: "USER_ALREADY_EXISTS", HttpCode: http.StatusConflict, Message: "Пользователь с таким email уже существует"}
	ErrInvalidActivationCode  = &ApiError{Code: "INVALID_ACTIVATION_CODE", HttpCode: http.StatusBadRequest, Message: "invalid activation code"}
	ErrActivationCodeNotFound = &ApiError{Code: "ACTIVATION_CODE_NOT_FOUND", HttpCode: http.StatusBadRequest, Message: "activation code not found"}
)

// User errors
var (
	ErrUserNotFound      = &ApiError{Code: "USER_NOT_FOUND", HttpCode: http.StatusNotFound}
	ErrUserNotFoundEmail = &ApiError{Code: "USER_NOT_FOUND_EMAIL", HttpCode: http.StatusNotFound}
)
