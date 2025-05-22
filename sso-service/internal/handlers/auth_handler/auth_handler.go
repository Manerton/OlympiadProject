package auth_handler

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) (uuid.UUID, error)
}

type AuthHandler struct {
	authService AuthService
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	token, err := ah.authService.Login(ctx, "", "")
	if err != nil {
		render.JSON(w, r, map[string]string{
			"Error": "failted login",
		})
	}

	_ = token

	render.JSON(w, r, map[string]any{
		"test": 1,
	})
}
