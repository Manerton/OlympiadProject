package event_handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"main/internal/dto/event_dto"
	"main/internal/lib/liblogger"
	"main/internal/lib/response"
	"main/internal/models/event"
	"main/internal/services/event_service"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type EventHandler struct {
	service *event_service.EventService
	log     *slog.Logger
}

func NewEventHandler(service *event_service.EventService, log *slog.Logger) *EventHandler {
	return &EventHandler{service: service, log: log}
}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetAllEvents"
	log := h.log.With(
		slog.String("op", op),
	)

	events_dto, err := h.service.GetAllEvents()
	if err != nil {
		log.Error("failed to get events", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get events"))
		return
	}
	log.Info("subject getted", slog.Any("events", events_dto))

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   events_dto,
	})
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetEventByID"
	log := h.log.With(
		slog.String("op", op),
	)

	received_id := chi.URLParam(r, "id")
	searched_id, err := strconv.ParseUint(received_id, 10, 32)
	if err != nil {
		log.Error("failed to parse id to uint", slog.String("received id", received_id), liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to parse id"))
		return
	}
	log.Info("event_id on request body decoded", slog.Any("event_id", searched_id))

	event_dto, err := h.service.GetEventByID(uint(searched_id))
	if err != nil {
		log.Error("failed to get event", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get event"))
		return
	}
	log.Info("event getted", slog.Any("event", event_dto))

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   event_dto,
	})
}

func (h *EventHandler) GetAllEventsTypeRegionalStage(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetAllEventsTypeRegionalStage"
	log := h.log.With(
		slog.String("op", op),
	)

	events_dto, err := h.service.GetEventsByType(event.RegionalStage)
	if err != nil {
		log.Error("failed to get events by type", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get events by type"))
		return
	}
	log.Info("events getted", slog.Any("events", events_dto))

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   events_dto,
	})
}

func (h *EventHandler) GetEventsByPreviousID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.event_handler.GetEventsByPreviousID"
	log := h.log.With(
		slog.String("op", op),
	)

	received_id := chi.URLParam(r, "id")
	searched_id, err := strconv.ParseUint(received_id, 10, 32)
	if err != nil {
		log.Error("failed to parse id to uint", slog.String("received id", received_id), liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to parse id"))
		return
	}
	log.Info("event_id on request body decoded", slog.Any("event_id", searched_id))

	events_dto, err := h.service.GetEventsByPreviousID(uint(searched_id))
	if err != nil {
		log.Error("failed to get events by previous id")
		render.JSON(w, r, response.Error("failed to get events"))
		return
	}

	log.Info("events getted")
	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   events_dto,
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

	events_dto, err := h.service.GetEventsByListID(ids.IDs)
	if err != nil {
		log.Error("failed to get events", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to get event"))
		return
	}
	log.Info("events getted")
	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   events_dto,
	})
}

func (h *EventHandler) CreateEventTypeRegionalStage(w http.ResponseWriter, r *http.Request) {
	h.createEventWithType(w, r, event.RegionalStage)
}

func (h *EventHandler) CreateEventTypeOlympiad(w http.ResponseWriter, r *http.Request) {
	h.createEventWithType(w, r, event.Olympiad)
}

func (h *EventHandler) CreateEventTypeStage(w http.ResponseWriter, r *http.Request) {
	h.createEventWithType(w, r, event.Stage)
}

// TODO: Может ли быть более одной апеляции, нужна проверка
func (h *EventHandler) CreateEventTypeAppeal(w http.ResponseWriter, r *http.Request) {
	h.createEventWithType(w, r, event.Appeal)
}

func (h *EventHandler) createEventWithType(w http.ResponseWriter, r *http.Request, event_type event.EventType) {
	const op = "handlers.event_handlers.CreateEventWithType"
	log := h.log.With(
		slog.String("op", op),
	)

	var event_dto event_dto.EventDTO
	err := render.DecodeJSON(r.Body, &event_dto)
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
	event_dto.EventType = event_type
	log.Info("event on request body decoded", slog.Any("event", event_dto))

	id, err := h.service.CreateEvent(event_dto)
	if err != nil {
		log.Error("failed to create event", slog.Any("type", event_type), liblogger.Err(err))
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

	var event_dto event_dto.EventDTO
	err := render.DecodeJSON(r.Body, &event_dto)
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
	log.Info("event on request body decoded", slog.Any("event", event_dto))

	id, err := h.service.UpdateEvent(event_dto)
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

	received_id := chi.URLParam(r, "id")
	deleted_id, err := strconv.ParseUint(received_id, 10, 32)
	if err != nil {
		log.Error("failed to parse id to uint", slog.String("received id", received_id), liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to parse id"))
	}
	log.Info("event_id on request body decoded", slog.Any("event_id", deleted_id))

	err = h.service.DeleteEvent(uint(deleted_id))
	if err != nil {
		log.Error("failed to delete event", liblogger.Err(err))
		render.JSON(w, r, response.Error("failed to delete event"))
		return
	}
	log.Info("event deleted", slog.Any("id", deleted_id))

	render.JSON(w, r, response.Success("Object deleted"))
}
