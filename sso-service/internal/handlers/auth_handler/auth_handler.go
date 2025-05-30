package auth_handler

import (
	"context"
	login_dto "main/internal/dto/auth/login"
	register_dto "main/internal/dto/auth/register"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, loginRequest *login_dto.LoginRequestDTO) (*login_dto.LoginResponseDTO, error)
	Register(ctx context.Context, registerRequest *register_dto.RegisterParticipantRequestDTO) (uuid.UUID, error)
}

type AuthHandler struct {
	authService AuthService
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var loginRequest login_dto.LoginRequestDTO

	err := render.DecodeJSON(r.Body, &loginRequest)
	if err != nil {
		// Create error response
		render.JSON(w, r, "failed decode request data")
		return
	}

	loginResponse, err := h.authService.Login(ctx, &loginRequest)
	if err != nil {
		// Create error response
		render.JSON(w, r, map[string]string{
			"Error": "failted login",
		})
		return
	}

	render.JSON(w, r, loginResponse)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var registerRequest register_dto.RegisterParticipantRequestDTO

	err := render.DecodeJSON(r.Body, &registerRequest)
	if err != nil {
		render.JSON(w, r, "failed")
		return
	}

	registerResponse, err := h.authService.Register(ctx, &registerRequest)
	if err != nil {
		render.JSON(w, r, "failed")
		return
	}

	render.JSON(w, r, registerResponse)
}
