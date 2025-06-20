package participant_handler

import (
	"context"
	participant_dto "main/internal/dto/participant"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ParticipantService interface {
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
