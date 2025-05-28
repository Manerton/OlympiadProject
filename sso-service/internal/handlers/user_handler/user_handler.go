package user_handler

import (
	user_dto "main/internal/dto/user"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserService interface {
	GetById(id string) (user_dto.UserResponseDTO, error)
	GetByListId(ids []string) ([]user_dto.UserResponseDTO, error)
}

type UserHandler struct {
	UserService UserService
}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {

	var id string = chi.URLParam(r, "id")

	userResponse, err := uh.UserService.GetById(id)
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

func (uh *UserHandler) GetUsersByListId(w http.ResponseWriter, r *http.Request) {
	type IdsReq struct {
		Ids []string `json:"ids"`
	}
	var ids IdsReq

	err := render.DecodeJSON(r.Body, &ids)
	if err != nil {
		render.JSON(w, r, response.ErrorResponse("failed to decode json"))
		return
	}

	usersResponse, err := uh.UserService.GetByListId(ids.Ids)
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
