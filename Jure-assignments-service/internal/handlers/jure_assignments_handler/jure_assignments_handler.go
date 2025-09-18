package jure_assignments_handler

import (
	"context"
	"log/slog"
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/lib/errs"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type JuryAssignmentsServiceInterface interface {
	GetAllJuryAssignments(context.Context) ([]juryAssignmentsDto.JuryAssignmentsResponseDTO, error)
	GetJuryAssignmentsByID(ctx context.Context, id string) (juryAssignmentsDto.JuryAssignmentsResponseDTO, error)

	GetAllByEventId(ctx context.Context, id string) ([]juryAssignmentsDto.JuryAssignmentsResponseDTO, error)
	GetAllByJuryId(ctx context.Context, id string) ([]juryAssignmentsDto.JuryAssignmentsResponseDTO, error)

	// CreateManyAssignmentsByOneJury(context.Context, juryAssignmentsDto.OneJuryManyAssignments) ([]uuid.UUID, error)
	Create(context.Context, juryAssignmentsDto.CreateJuryAssignmentsDTO) (uuid.UUID, error)
	CreateMany(context.Context, juryAssignmentsDto.CreateOneAssigmentsManyJury) ([]uuid.UUID, error)

	Update(context.Context, string, juryAssignmentsDto.UpdateJuryAssignmentsDTO) error
	Delete(context.Context, string) error
}

type JuryAssignmentHandler struct {
	service JuryAssignmentsServiceInterface
}

func NewJureAssignmentHandler(js JuryAssignmentsServiceInterface, log *slog.Logger) *JuryAssignmentHandler {
	return &JuryAssignmentHandler{service: js}
}

// @Summery All jury-assignments
// @Description Получение всех связей жюри-событие
// @Tags jury-assignments
// @Produce json
// @Success 200 {object} response.ApiResponse{data=[]juryAssignmentsDto.JuryAssignmentsResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/jury-assignments [get]
func (h *JuryAssignmentHandler) GetAllJuryAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtos, err := h.service.GetAllJuryAssignments(ctx)
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
		Status: response.SUCCESS,
		Data:   dtos,
	})
}

// @Summery Get jury-assignments by id
// @Description Получения связи жюри-событие по id связи
// @Tags jury-assignments
// @Produce json
// @Param id path string true "id связи"
// @Success 200 {object} response.ApiResponse{data=juryAssignmentsDto.JuryAssignmentsResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/jury-assignments/{id} [get]
func (h *JuryAssignmentHandler) GetJuryAssignmentsByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	receivedID := chi.URLParam(r, "id")

	dto, err := h.service.GetJuryAssignmentsByID(ctx, receivedID)
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
		Status: response.SUCCESS,
		Data:   dto,
	})
}

// @Summery Get all jury-assignments by event-id
// @Description Получение всех связей жюри-событие по event-id
// @Tags jury-assignments
// @Produce json
// @Param id path string true "id события"
// @Success 200 {object} response.ApiResponse{data=[]juryAssignmentsDto.JuryAssignmentsResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/jury-assignments/event/{id} [get]
func (h *JuryAssignmentHandler) GetAllByEventId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventIdS := chi.URLParam(r, "id")

	result, err := h.service.GetAllByEventId(ctx, eventIdS)
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
		Data:       result,
	})
}

// @Summery Get all jury-assignments by user-id
// @Description Получение всех связей жюри-событие по user-id
// @Tags jury-assignments
// @Produce json
// @Param id path string true "id жюри"
// @Success 200 {object} response.ApiResponse{data=[]juryAssignmentsDto.JuryAssignmentsResponseDTO}
// @Failure 400 {object} response.ApiResponse
// @Router /api/jury-assignments/jury/{id} [get]
func (h *JuryAssignmentHandler) GetAllByJuryId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	juryId := chi.URLParam(r, "id")

	result, err := h.service.GetAllByJuryId(ctx, juryId)
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
		Data:       result,
	})
}

// @Summery Create jury-assignments
// @Description Создание связи записи жюри на событие
// @Tags jury-assignments
// @Accept json
// @Produce json
// @Param credentials body juryAssignmentsDto.CreateJuryAssignmentsDTO true "Данные для создания связи"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/jury-assignments [post]
func (h *JuryAssignmentHandler) CreateJuryAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := juryAssignmentsDto.CreateJuryAssignmentsDTO{}
	err := render.DecodeJSON(r.Body, &dto)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode body")))
		return
	}

	err = validator.New().Struct(dto)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed validate dto")))
		return
	}

	createdId, err := h.service.Create(ctx, dto)
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
	render.JSON(w, r, response.SuccessResponse(createdId.String()))
}

// @Summery Create many jury-assignments
// @Description Создание множественной связей нескольких жюри на событие
// @Tags jury-assignments
// @Accept json
// @Produce json
// @Param credentials body juryAssignmentsDto.CreateOneAssigmentsManyJury true "Данные для создания связи"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/jury-assignments/many [post]
func (h *JuryAssignmentHandler) CreateManyJuryAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := juryAssignmentsDto.CreateOneAssigmentsManyJury{}
	err := render.DecodeJSON(r.Body, &dto)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode body")))
		return
	}

	err = validator.New().Struct(dto)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed validate dto")))
		return
	}

	_, err = h.service.CreateMany(ctx, dto)
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
	render.JSON(w, r, response.SuccessResponse("Success create jury-assigments"))
}

// @Summery Update jury-assignments
// @Description Обновление связи записи жюри на событие
// @Tags jury-assignments
// @Accept json
// @Produce json
// @Param id path string true "id связи"
// @Param credentials body juryAssignmentsDto.UpdateJuryAssignmentsDTO true "Данные для обновления связи"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/jury-assignments/{id} [put]
func (h *JuryAssignmentHandler) UpdateJuryAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	dto := juryAssignmentsDto.UpdateJuryAssignmentsDTO{}
	err := render.DecodeJSON(r.Body, &dto)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode body")))
		return
	}

	err = h.service.Update(ctx, id, dto)
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
	render.JSON(w, r, response.SuccessResponse("success updated"))
}

// @Summery Delete jury-assignments
// @Description Удаление связи записи жюри на событие
// @Tags jury-assignments
// @Produce json
// @Param id path string true "id связи"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/jury-assignments/{id} [delete]
func (h *JuryAssignmentHandler) DeleteJuryAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	receivedID := chi.URLParam(r, "id")

	err := h.service.Delete(ctx, receivedID)
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
	render.JSON(w, r, response.SuccessResponse("object deleted"))
}
