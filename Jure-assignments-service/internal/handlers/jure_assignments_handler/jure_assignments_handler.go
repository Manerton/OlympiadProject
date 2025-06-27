package jure_assignments_handler

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type JuryAssignmentsServiceInterface interface {
	GetAllJuryAssignments(context.Context) ([]juryAssignmentsDto.JuryAssignmentsDTO, error)
	GetJuryAssignmentsByID(ctx context.Context, id string) (juryAssignmentsDto.JuryAssignmentsDTO, error)
	GetAllJuryAssignmentsByFilter(context.Context, juryAssignmentsDto.JuryAssignmentsDTO) ([]juryAssignmentsDto.JuryAssignmentsDTO, error)
	GetPartOfAllJuryAssignmentsByFilter(context.Context, []string, juryAssignmentsDto.JuryAssignmentsDTO) ([]juryAssignmentsDto.JuryAssignmentsDTO, error)
	GetJuryAssignmentsByFilter(context.Context, juryAssignmentsDto.JuryAssignmentsDTO) (juryAssignmentsDto.JuryAssignmentsDTO, error)
	CreateManyAssignmentsByOneJury(context.Context, juryAssignmentsDto.OneJuryManyAssignments) ([]uuid.UUID, error)
	CreateJuryAssignments(context.Context, juryAssignmentsDto.JuryAssignmentsDTO) (uuid.UUID, error)
	UpdateJuryAssignments(context.Context, string, juryAssignmentsDto.JuryAssignmentsDTO) error
	DeleteJuryAssignments(context.Context, string) error
}

type JureAssignmentHandler struct {
	service JuryAssignmentsServiceInterface
}

func NewJureAssignmentHandler(js JuryAssignmentsServiceInterface, log *slog.Logger) *JureAssignmentHandler {
	return &JureAssignmentHandler{service: js}
}

func (h *JureAssignmentHandler) GetAllJuryAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtos, err := h.service.GetAllJuryAssignments(ctx)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to get all jury assignments"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   dtos,
	})
}

func (h *JureAssignmentHandler) GetJuryAssignmentsByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	receivedID := chi.URLParam(r, "id")

	dto, err := h.service.GetJuryAssignmentsByID(ctx, receivedID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to get jury assignments"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   dto,
	})
}

func (h *JureAssignmentHandler) GetAllJuryIDByEventID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	receivedID := chi.URLParam(r, "event_id")
	searchedID, err := uuid.Parse(receivedID)
	if err != nil {
		render.JSON(w, r, response.Error("failed to parse id"))
		return
	}

	filter := juryAssignmentsDto.JuryAssignmentsDTO{EventID: searchedID}
	fields := []string{"jury_id"}
	dto, err := h.service.GetPartOfAllJuryAssignmentsByFilter(ctx, fields, filter)
	if err != nil {
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   dto,
	})
}

func (h *JureAssignmentHandler) CreateJuryAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := juryAssignmentsDto.JuryAssignmentsDTO{}
	err := render.DecodeJSON(r.Body, &dto)
	if err != nil {
		render.JSON(w, r, response.Error("failed to decode request"))
		return
	}

	err = validator.New().Struct(dto)
	if err != nil {
		validateErr := err.(validator.ValidationErrors)
		render.JSON(w, r, response.Error(fmt.Sprintf("err %v", validateErr)))
		return
	}

	_, err = h.service.CreateJuryAssignments(ctx, dto)
	if err != nil {
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	render.JSON(w, r, response.Success(""))
}

func (h *JureAssignmentHandler) CreateManyAssignmentsByOneJury(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dto := juryAssignmentsDto.OneJuryManyAssignments{}
	err := render.DecodeJSON(r.Body, &dto)
	if err != nil {
		render.JSON(w, r, response.Error("failed to decode request"))
		return
	}

	ids, err := h.service.CreateManyAssignmentsByOneJury(ctx, dto)
	if err != nil {
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   ids,
	})

}

func (h *JureAssignmentHandler) UpdateJuryAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	dto := juryAssignmentsDto.JuryAssignmentsDTO{}
	err := render.DecodeJSON(r.Body, &dto)
	if err != nil {
		render.JSON(w, r, response.Error("failed to decode request"))
		return
	}

	err = h.service.UpdateJuryAssignments(ctx, id, dto)
	if err != nil {
		render.JSON(w, r, response.Error("failed to update JuryAssignments"))
		return
	}
	render.JSON(w, r, response.Success("success updated"))
}

func (h *JureAssignmentHandler) DeleteJuryAssignments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	receivedID := chi.URLParam(r, "id")

	err := h.service.DeleteJuryAssignments(ctx, receivedID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}
	render.JSON(w, r, response.Success("object deleted"))
}
