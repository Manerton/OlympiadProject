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
	GetParticipantUserById(ctx context.Context, id string) (user_dto.ParticipantUserResponseDTO, error)
	GetByListId(ctx context.Context, ids []*string) ([]user_dto.UserResponseDTO, error)

	Update(ctx context.Context, userDto *user_dto.UpdateUserRequestDTO) error
	Delete(ctx context.Context, id string) error
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

func (h *UserHandler) GetParticipantUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var id string = chi.URLParam(r, "id")

	participantUserResponse, err := h.UserService.GetParticipantUserById(ctx, id)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse("failed to find user"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: 200,
		Data:       participantUserResponse,
	})
}

func (h *UserHandler) GetUsersByListId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type IdsReq struct {
		Ids []*string `json:"ids"`
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

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userDto := &user_dto.UpdateUserRequestDTO{}

	err := render.DecodeJSON(r.Body, userDto)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse("failed decode json"))
		return
	}

	err = h.UserService.Update(ctx, userDto)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse("failed update user"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.SuccessResponse("user update success"))
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	err := h.UserService.Delete(ctx, id)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed delete user"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.SuccessResponse("user delete success"))
}
