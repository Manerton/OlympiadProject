package jureAssignmentHandler

import (
	"fmt"
	"log/slog"
	"main/internal/dto/juryAssignmentsDto"
	"main/internal/lib/liblogger"
	"main/internal/lib/response"
	"main/internal/services/juryAssignmentsService"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type JureAssignmentHandler struct {
	service *juryAssignmentsService.JuryAssignmentsService
	log     *slog.Logger
}

func NewJureAssignmentHandler(js *juryAssignmentsService.JuryAssignmentsService, log *slog.Logger) *JureAssignmentHandler {
	return &JureAssignmentHandler{service: js, log: log}
}

func (h *JureAssignmentHandler) GetAllJuryAssignments(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.jureAssignmentHandler.GetAllJuryAssignments"
	log := h.log.With(
		slog.String("op", op),
	)

	dtos, err := h.service.GetAllJuryAssignments()
	if err != nil {
		log.Error("failed to get all jury assignments", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get all jury assignments"))
		return
	}
	log.Info("jury assignments getted", slog.Any("juryAssignments", dtos))

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   dtos,
	})
}

func (h *JureAssignmentHandler) GetJuryAssignmentsByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.jureAssignmentHandler.GetJuryAssignmentsByID"
	log := h.log.With(
		slog.String("op", op),
	)
	receivedID := chi.URLParam(r, "id")
	searchedID, err := strconv.ParseUint(receivedID, 10, 32)
	if err != nil {
		log.Error("failed to parse id to uint", slog.String("received id", receivedID), liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to parse id"))
		return
	}
	log.Info("id on request body decoded", slog.Any("id", searchedID))

	filter := juryAssignmentsDto.JuryAssignmentsDTO{ID: uint(searchedID)}
	dto, err := h.service.GetJuryAssignmentsByFilter(filter)
	if err != nil {
		log.Error("failed to get jury assignments", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get jury assignments"))
		return
	}
	log.Info("jury assignments getted", slog.Any("juryAssignments", dto))

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   dto,
	})
}

func (h *JureAssignmentHandler) GetAllJuryIDByEventID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.jureAssignmentHandler.GetAllJuryIDByEventID"
	log := h.log.With(
		slog.String("op", op),
	)
	receivedID := chi.URLParam(r, "event_id")
	searchedID, err := strconv.ParseUint(receivedID, 10, 32)
	if err != nil {
		log.Error("failed to parse id to uint", slog.String("received id", receivedID), liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to parse id"))
		return
	}
	log.Info("id on request body decoded", slog.Any("event_id", searchedID))
	filter := juryAssignmentsDto.JuryAssignmentsDTO{EventID: uint(searchedID)}
	fields := []string{"jury_id"}
	dto, err := h.service.GetPartOfAllJuryAssignmentsByFilter(fields, filter)
	if err != nil {
		log.Error("failed to get jury_id by event_id", slog.Any("jury_id", fields), slog.Any("event_id", searchedID), liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get jury_id"))
		return
	}
	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   dto,
	})
}

func (h *JureAssignmentHandler) CreateJuryAssignments(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.jureAssignmentHandler.CreateJuryAssignments"
	log := h.log.With(
		slog.String("op", op),
	)
	dto := juryAssignmentsDto.JuryAssignmentsDTO{}
	err := render.DecodeJSON(r.Body, &dto)
	if err != nil {
		log.Error("failed to decode request body", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to decode request"))
		return
	}
	log.Info("request body decoded", slog.Any("request", dto))
	err = validator.New().Struct(dto)
	if err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", liblogger.Err(err))
		render.JSON(w, r, response.Error(fmt.Sprintf("err %v", validateErr)))
		return
	}
	id, err := h.service.CreateJuryAssignments(dto)
	if err != nil {
		log.Error("failed to create JuryAssignments body", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to create JuryAssignments"))
		return
	}
	render.JSON(w, r, response.Success(fmt.Sprintf("id = %d", id)))
}
