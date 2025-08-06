package user_handler

import (
	"context"
	user_dto "main/internal/dto/user"
	"main/internal/lib/errs"
	"main/internal/lib/parser"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserService interface {
	GetAll(ctx context.Context, page, limit *int) ([]user_dto.UserResponseDTO, error)
	GetByFilter(ctx context.Context, searchUser user_dto.SearchAttributesUserDTO) (user_dto.UserResponseDTO, error)
	GetById(ctx context.Context, id string) (user_dto.UserResponseDTO, error)
	GetUserParticipantById(ctx context.Context, id string) (user_dto.UserParticipantResponseDTO, error)
	GetByListId(ctx context.Context, ids []string) ([]user_dto.UserResponseDTO, error)

	GetCount(ctx context.Context) (int64, error)

	Update(ctx context.Context, id string, userDto user_dto.UpdateUserRequestDTO) error
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

// @Summery Users count
// @Security BearerAuth
// @Description Получение количества пользователей
// @Tags users
// @Produce json
// @Success 200 {object} response.ApiResponse{data=int}
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/count [get]
func (h *UserHandler) GetCountUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userCount, err := h.UserService.GetCount(ctx)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       userCount,
	})
}

// @Summery All users
// @Security BearerAuth
// @Description Получение всех пользователей
// @Tags users
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Ограничение на количество записей"
// @Success 200 {object} response.ApiResponse{data=[]user_dto.UserResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/users [get]
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit, err := parser.ParsePageLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("incorrect page/limit")))
		return
	}

	usersResponse, err := h.UserService.GetAll(ctx, page, limit)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       usersResponse,
	})
}

// @Summery Get user by filter
// @Security BearerAuth
// @Description Получение пользователя по фильтру из его полей
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body user_dto.SearchAttributesUserDTO true "Данные для поиска пользователя"
// @Success 200 {object} response.ApiResponse{data=user_dto.UserResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/filter [post]
func (h *UserHandler) GetUserByFilter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	searchDTO := user_dto.SearchAttributesUserDTO{}
	err := render.DecodeJSON(r.Body, &searchDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode body")))
		return
	}

	userResponse, err := h.UserService.GetByFilter(ctx, searchDTO)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       userResponse,
	})

}

// @Summery Get user by id
// @Security BearerAuth
// @Description Получение пользователя по id
// @Tags users
// @Produce json
// @Param id path string true "id пользователя"
// @Success 200 {object} response.ApiResponse{data=user_dto.UserResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var id string = chi.URLParam(r, "id")

	userResponse, err := h.UserService.GetById(ctx, id)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       userResponse,
	})
}

// @Summery Get full user participant info
// @Security BearerAuth
// @Description Получение всей информации о пользователе ученике
// @Tags users
// @Produce json
// @Param id path string true "id пользователя"
// @Success 200 {object} response.ApiResponse{data=user_dto.UserParticipantResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/all-info/{id} [get]
func (h *UserHandler) GetUserParticipantById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var id string = chi.URLParam(r, "id")

	participantUserResponse, err := h.UserService.GetUserParticipantById(ctx, id)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       participantUserResponse,
	})
}

// @Summery Get users by list id
// @Security BearerAuth
// @Description Получение пользователей по списку id
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body []string true "Список id пользователей"
// @Success 200 {object} response.ApiResponse{data=[]user_dto.UserResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/list [post]
func (h *UserHandler) GetUsersByListId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type IdsReq struct {
		Ids []string `json:"ids"`
	}
	var ids IdsReq

	err := render.DecodeJSON(r.Body, &ids)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode body")))
		return
	}

	usersResponse, err := h.UserService.GetByListId(ctx, ids.Ids)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       usersResponse,
	})
}

// @Summery Update user
// @Security BearerAuth
// @Description Обновление пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body user_dto.UpdateUserRequestDTO true "Данные для обновления пользователя"
// @Param id path string true "id пользователя"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	userDto := user_dto.UpdateUserRequestDTO{}

	err := render.DecodeJSON(r.Body, &userDto)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode body")))
		return
	}

	err = h.UserService.Update(ctx, id, userDto)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.SuccessResponse("user update success"))
}

// @Summery Delete user
// @Security BearerAuth
// @Description Удаление пользователя
// @Tags users
// @Produce json
// @Param id path string true "id пользователя"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	err := h.UserService.Delete(ctx, id)
	if err != nil {
		if apiErr, ok := errs.IsApiError(err); ok {
			render.Status(r, apiErr.HttpCode)
			render.JSON(w, r, response.ErrorApiResponse(apiErr))
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrInternalError))
		return
	}

	render.JSON(w, r, response.SuccessResponse("user delete success"))
}
