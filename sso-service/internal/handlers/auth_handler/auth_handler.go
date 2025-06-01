package auth_handler

import (
	"context"
	login_dto "main/internal/dto/auth/login"
	register_dto "main/internal/dto/auth/register"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/render"
)

type AuthService interface {
	Login(ctx context.Context, loginRequest *login_dto.LoginRequestDTO) (*login_dto.LoginResponseDTO, error)
	Register(ctx context.Context, registerRequest *register_dto.RegisterParticipantRequestDTO) error
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
		render.JSON(w, r, response.ErrorResponse("failed to decode json"))
		return
	}

	loginResponse, err := h.authService.Login(ctx, &loginRequest)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, loginResponse)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var registerRequest register_dto.RegisterParticipantRequestDTO

	err := render.DecodeJSON(r.Body, &registerRequest)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse("failed to decode json"))
		return
	}

	err = h.authService.Register(ctx, &registerRequest)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.SuccessResponse("Register success"))
}
