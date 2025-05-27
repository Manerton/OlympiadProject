package user_handler

import (
	user_dto "main/internal/dto/user"
	"net/http"
)

type UserService interface {
	GetById(id string) (user_dto.UserResponseDTO, error)
	GetByListId(ids []string) ([]user_dto.UserResponseDTO, error)
	Update()
}

type UserHandler struct {
	UserService UserService
}

func (uh *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {

	var id string = r.URL.Query().Get("id")

	_ = id

}
