package user_handler

import (
	"context"
	user_dto "main/internal/dto/user"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserService interface {
	GetById(ctx context.Context, id string) (user_dto.UserResponseDTO, error)
	GetByListId(ctx context.Context, ids []string) ([]user_dto.UserResponseDTO, error)
}

type UserHandler struct {
	UserService UserService
}

func New(userService UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var id string = chi.URLParam(r, "id")

	userResponse, err := h.UserService.GetById(ctx, id)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse("failed to find user"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: 200,
		Data:       userResponse,
	})
}

func (h *UserHandler) GetUsersByListId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type IdsReq struct {
		Ids []string `json:"ids"`
	}
	var ids IdsReq

	err := render.DecodeJSON(r.Body, &ids)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse("failed to decode json"))
		return
	}

	usersResponse, err := h.UserService.GetByListId(ctx, ids.Ids)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse("failed to find users"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       usersResponse,
	})
}
