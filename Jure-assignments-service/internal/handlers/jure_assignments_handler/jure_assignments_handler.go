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
	Update(context.Context, string, juryAssignmentsDto.UpdateJuryAssignmentsDTO) error
	Delete(context.Context, string) error
}

type JuryAssignmentHandler struct {
	service JuryAssignmentsServiceInterface
}

func NewJureAssignmentHandler(js JuryAssignmentsServiceInterface, log *slog.Logger) *JuryAssignmentHandler {
	return &JuryAssignmentHandler{service: js}
}

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

// func (h *JuryAssignmentHandler) CreateManyAssignmentsByOneJury(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	dto := juryAssignmentsDto.OneJuryManyAssignments{}
// 	err := render.DecodeJSON(r.Body, &dto)
// 	if err != nil {
// 		render.JSON(w, r, response.ErrorResponse("failed to decode request"))
// 		return
// 	}

// 	ids, err := h.service.CreateManyAssignmentsByOneJury(ctx, dto)
// 	if err != nil {
// 		render.JSON(w, r, response.ErrorResponse(err.Error()))
// 		return
// 	}

// 	render.JSON(w, r, response.ApiResponse{
// 		Status: response.SUCCESS,
// 		Data:   ids,
// 	})

// }

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
