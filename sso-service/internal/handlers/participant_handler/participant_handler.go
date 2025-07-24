package participant_handler

import (
	"context"
	participant_dto "main/internal/dto/participant"
	"main/internal/lib/parser"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ParticipantService interface {
	GetAll(ctx context.Context, page, limit *int) ([]participant_dto.ParticipantResponseDTO, error)
	GetById(ctx context.Context, id string) (participant_dto.ParticipantResponseDTO, error)
	GetByUserId(ctx context.Context, userId string) (participant_dto.ParticipantResponseDTO, error)

	GetCount(ctx context.Context) (int64, error)

	Update(ctx context.Context, id string, participantDTO participant_dto.UpdateParticipantRequestDTO) error
}

type ParticipantHandler struct {
	participantService ParticipantService
}

func New(participantSerivce ParticipantService) *ParticipantHandler {
	return &ParticipantHandler{
		participantService: participantSerivce,
	}
}

func (h *ParticipantHandler) GetCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	participantCount, err := h.participantService.GetCount(ctx)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed get participant users"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       participantCount,
	})

}

func (h *ParticipantHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	participantRes, err := h.participantService.GetById(ctx, id)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       participantRes,
	})
}

func (h *ParticipantHandler) GetByUserId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	participantRes, err := h.participantService.GetByUserId(ctx, id)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       participantRes,
	})

}

func (h *ParticipantHandler) GetAllParticipants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, limit, err := parser.ParsePageLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed parse page/limit"))
		return
	}

	participantResponse, err := h.participantService.GetAll(ctx, page, limit)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.SUCCESS,
		StatusCode: http.StatusOK,
		Data:       participantResponse,
	})
}

func (h *ParticipantHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	participantDTO := participant_dto.UpdateParticipantRequestDTO{}
	err := render.DecodeJSON(r.Body, &participantDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed decode json"))
		return
	}

	err = h.participantService.Update(ctx, id, participantDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(err.Error()))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.SuccessResponse("success update"))
}
