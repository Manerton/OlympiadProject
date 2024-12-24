package eventhandler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"main/internal/dto/event_dto"
	"main/internal/lib/liblogger"
	"main/internal/lib/parsing"
	"main/internal/lib/request"
	"main/internal/lib/response"
	"main/internal/models/event"
	"main/internal/services/event_service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type EventHandler struct {
	service *event_service.EventService
	log     *slog.Logger
}

func NewEventHandler(service *event_service.EventService, log *slog.Logger) *EventHandler {
	return &EventHandler{service: service, log: log}
}

func (h *EventHandler) GetEventByFilterAndFields(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetEventByFilterAndFields"

	log := h.log.With(
		slog.String("op", op),
	)
	detailRequest := request.DetailRequest{}

	err := render.DecodeJSON(r.Body, &detailRequest)
	if err != nil {
		log.Error("failed to decode details", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to decode details"))
		return
	}

	eventDetails, err := h.service.GetEventByFilterAndFields(detailRequest.EventDTO, detailRequest.Fields)
	if err != nil {
		log.Error("failed to get event", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get event"))
		return
	}
	log.Info("event getted", slog.Any("event", eventDetails))

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventDetails,
	})
}

func (h *EventHandler) GetEventsByFilterAndFields(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetEventsByFilterAndFields"
	log := h.log.With(
		slog.String("op", op),
	)
	detailRequest := request.DetailRequest{}

	err := render.DecodeJSON(r.Body, &detailRequest)
	if err != nil {
		log.Error("failed to decode details", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to decode details"))
		return
	}

	eventDetails, err := h.service.GetEventsByFilterAndFields(detailRequest.EventDTO, detailRequest.Fields, detailRequest.Offset, detailRequest.Limit)
	if err != nil {
		log.Error("failed to get events", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get events"))
		return
	}
	log.Info("event getted", slog.Any("events", eventDetails))

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventDetails,
	})
}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetAllEvents"
	log := h.log.With(
		slog.String("op", op),
	)

	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset, limit, err := parsing.ParseOffsetLimit(offsetStr, limitStr)
	if err != nil {
		log.Error("failed to get events", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get events"))
		return
	}

	eventsDTO, err := h.service.GetAllEvents(offset, limit)
	if err != nil {
		log.Error("failed to get events", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get events"))
		return
	}
	log.Info("event getted", slog.Any("events", eventsDTO))

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventsDTO,
	})
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetEventByID"
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
	log.Info("event id on request body decoded", slog.Any("event id", searchedID))

	eventDTO, err := h.service.GetEventByID(uint(searchedID))
	if err != nil {
		log.Error("failed to get event", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get event"))
		return
	}
	log.Info("event getted", slog.Any("event", eventDTO))

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventDTO,
	})
}

func (h *EventHandler) GetEventsTypeRegionalStage(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetEventsTypeRegionalStage"
	log := h.log.With(
		slog.String("op", op),
	)

	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset, limit, err := parsing.ParseOffsetLimit(offsetStr, limitStr)
	if err != nil {
		log.Error("failed to get events", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get events"))
		return
	}

	eventsDTO, err := h.service.GetEventsByType(event.RegionalStage, offset, limit)
	if err != nil {
		log.Error("failed to get events by type", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get events by type"))
		return
	}
	log.Info("events getted", slog.Any("events", eventsDTO))

	if offset == nil && limit == nil {
		render.JSON(w, r, response.ApiResponse{
			Status: response.StatusOK,
			Data:   eventsDTO,
		})
		return
	}

	count, err := h.service.GetCountEventsByType(event.RegionalStage)
	if err != nil {
		log.Error("failed to get count events by type", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get count events by type"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data: response.PaginatedResponse{
			Events:     eventsDTO,
			TotalCount: int(count),
		},
	})

}

func (h *EventHandler) GetEventsTypeStageAndHisChilds(w http.ResponseWriter, r *http.Request) {
	const op = "handler.evend_handler.GetEventsTypeStageAndHisChilds"
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

	eventsDto, err := h.service.GetEventsTypeStageAndHisChilds(uint(searchedID))
	if err != nil {
		log.Error("failed to get events", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get events"))
		return
	}
	log.Info("events getted")

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventsDto,
	})
}

func (h *EventHandler) GetEventsByPreviousID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetEventsByPreviousID"
	log := h.log.With(
		slog.String("op", op),
	)

	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset, limit, err := parsing.ParseOffsetLimit(offsetStr, limitStr)
	if err != nil {
		log.Error("failed to parse offsetStr/limitStr", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to pares offset/limit"))
		return
	}

	receivedID := chi.URLParam(r, "id")
	searchedID, err := strconv.ParseUint(receivedID, 10, 32)
	if err != nil {
		log.Error("failed to parse id to uint", slog.String("received id", receivedID), liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to parse id"))
		return
	}
	log.Info("event id on request body decoded", slog.Any("event id", searchedID))

	eventsDTO, err := h.service.GetEventsByPreviousID(uint(searchedID), offset, limit)
	if err != nil {
		log.Error("failed to get events by previous id")
		render.JSON(w, r, response.Error("failed to get events"))
		return
	}

	log.Info("events getted")

	if offset == nil && limit == nil {
		render.JSON(w, r, response.ApiResponse{
			Status: response.StatusOK,
			Data:   eventsDTO,
		})
		return
	}

	count, err := h.service.GetCountEventsByPreviousID(uint(searchedID))
	if err != nil {
		log.Error("failed to get count events by type", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get count events by type"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data: response.PaginatedResponse{
			Events:     eventsDTO,
			TotalCount: int(count),
		},
	})
}

func (h *EventHandler) GetEventsByListID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetEventsByListID"
	log := h.log.With(
		slog.String("op", op),
	)

	type ReqIds struct {
		IDs []uint `json:"ids"`
	}
	var ids ReqIds
	err := render.DecodeJSON(r.Body, &ids)
	if errors.Is(err, io.EOF) {
		log.Error("render body is empty")
		render.JSON(w, r, response.Error("empty request"))
		return
	}
	if err != nil {
		log.Error("failed to decode request body", liblogger.Err(err))
		render.JSON(w, r, response.Error("failde to decode request"))
		return
	}
	log.Info("ids on request body decoded")

	eventsDTO, err := h.service.GetEventsByListID(ids.IDs)
	if err != nil {
		log.Error("failed to get events", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get event"))
		return
	}
	log.Info("events getted")
	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventsDTO,
	})
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handlers.CreateEvent"
	log := h.log.With(
		slog.String("op", op),
	)

	var eventDTO event_dto.EventDTO
	err := render.DecodeJSON(r.Body, &eventDTO)
	if errors.Is(err, io.EOF) {
		log.Error("render body is empty")
		render.JSON(w, r, response.Error("empty request"))
		return
	}
	if err != nil {
		log.Error("failed to decode request body", liblogger.Err(err))
		render.JSON(w, r, response.Error("failde to decode request"))
		return
	}
	log.Info("event on request body decoded", slog.Any("event", eventDTO))

	err = validator.New().Struct(eventDTO)
	if err != nil {
		validateErr := err.(validator.ValidationErrors)
		log.Error("invalid request", liblogger.Err(err))
		render.JSON(w, r, response.Error(fmt.Sprintf("err %v", validateErr)))
		return
	}

	id, err := h.service.CreateEvent(eventDTO)
	if err != nil {
		log.Error("failed to create event", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to create event"))
		return
	}
	log.Info("event created", slog.Any("id", id))
	render.JSON(w, r, response.Success(fmt.Sprintf("id = %d", id)))
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.UpdateEvent"
	log := h.log.With(
		slog.String("op", op),
	)

	var eventDTO event_dto.EventDTO
	receivedID := chi.URLParam(r, "id")
	updatedID, err := strconv.ParseUint(receivedID, 10, 32)
	if err != nil {
		log.Error("failed to parse id to uint", slog.String("received id", receivedID), liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to parse id"))
	}
	log.Info("event id on request body decoded", slog.Any("event id", updatedID))
	eventDTO.ID = uint(updatedID)

	err = render.DecodeJSON(r.Body, &eventDTO)
	if errors.Is(err, io.EOF) {
		log.Error("render body is empty")
		render.JSON(w, r, response.Error("empty request"))
		return
	}
	if err != nil {
		log.Error("failed to decode event", liblogger.Err(err))
		render.JSON(w, r, "failed to decode event")
		return
	}
	log.Info("event on request body decoded", slog.Any("event", eventDTO))

	id, err := h.service.UpdateEvent(eventDTO)
	if err != nil {
		log.Error("failed to update event", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to update event"))
		return
	}
	log.Info("event updated", slog.Any("id", id))

	render.JSON(w, r, response.Success(fmt.Sprintf("id = %d", id)))
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.DeleteEvent"
	log := h.log.With(
		slog.String("op", op),
	)

	receivedID := chi.URLParam(r, "id")
	deletedID, err := strconv.ParseUint(receivedID, 10, 32)
	if err != nil {
		log.Error("failed to parse id to uint", slog.String("received id", receivedID), liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to parse id"))
	}
	log.Info("event id on request body decoded", slog.Any("event id", deletedID))

	err = h.service.DeleteEvent(uint(deletedID))
	if err != nil {
		log.Error("failed to delete event", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to delete event"))
		return
	}
	log.Info("event deleted", slog.Any("id", deletedID))

	render.JSON(w, r, response.Success("Object deleted"))
}
