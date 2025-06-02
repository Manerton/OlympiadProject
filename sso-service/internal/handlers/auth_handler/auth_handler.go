package auth_handler

import (
	"context"
	login_dto "main/internal/dto/auth/login"
	register_dto "main/internal/dto/auth/register"
	"main/internal/lib/response"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type AuthService interface {
	Login(ctx context.Context, loginRequest *login_dto.LoginRequestDTO) (*login_dto.AuthResultDTO, error)
	RegisterParticipant(ctx context.Context, registerRequest *register_dto.RegisterParticipantRequestDTO) error
}

type AuthHandler struct {
	authService AuthService
}

func New(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var loginRequest login_dto.LoginRequestDTO

	err := render.DecodeJSON(r.Body, &loginRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to decode json"))
		return
	}

	authResult, err := h.authService.Login(ctx, &loginRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	// Set the JWT token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    authResult.RefreshToken,
		HttpOnly: true,
		Secure:   false,
		// Domain:   "172.16.1.39",
		SameSite: http.SameSiteStrictMode,
		Path:     "/users/refresh",
		Expires:  time.Now().Add(24 * time.Hour), // Match the token expiration
	})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.ApiResponse{
		Status:     "SUCCESS",
		StatusCode: http.StatusOK,
		Data: login_dto.LoginResponseDTO{
			AccessToken: authResult.AccessToken,
			ExpiresIn:   authResult.ExpiresIn,
		},
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var registerRequest register_dto.RegisterParticipantRequestDTO

	err := render.DecodeJSON(r.Body, &registerRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to decode json"))
		return
	}

	err = h.authService.RegisterParticipant(ctx, &registerRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.SuccessResponse("Register success"))
}
