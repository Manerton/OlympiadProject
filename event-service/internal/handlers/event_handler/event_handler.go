package event_handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"main/internal/dto/event_dto"
	"main/internal/lib/parsing"
	"main/internal/lib/request"
	"main/internal/lib/response"
	"main/internal/models/event"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type EventServiceInterface interface {
	GetAllEvents(ctx context.Context, offset, limit *int) ([]event_dto.EventDTOResponse, error)
	GetEventByID(ctx context.Context, id string) (event_dto.EventDTOResponse, error)

	GetEventByFilterAndFields(ctx context.Context, filter event_dto.EventDTO, fields *[]string) (event_dto.DetailsEvent, error)
	GetEventsByFilterAndFields(ctx context.Context, filter event_dto.EventDTO, fields *[]string, offset, limit *int, order *string) ([]event_dto.DetailsEvent, error)
	GetCountEventsByType(ctx context.Context, eventType event.EventType) (int64, error)
	GetCountEventsByPreviousID(ctx context.Context, previousID string) (int64, error)
	GetEventsByType(ctx context.Context, eventType event.EventType, offset, limit *int, order *string) ([]event_dto.EventDTOResponse, error)
	GetEventsTypeStageAndHisChilds(ctx context.Context, id uuid.UUID) ([]event_dto.EventDTOResponse, error)
	GetEventsByPreviousID(ctx context.Context, previousID string, offset, limit *int, order *string) ([]event_dto.EventDTOResponse, error)
	GetEventsByListID(ctx context.Context, ids []uuid.UUID) ([]event_dto.EventDTOResponse, error)

	CreateEvent(ctx context.Context, eventDTO event_dto.CreateEventDTORequest) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, id string, eventDTO event_dto.UpdateEventDTORequest) error
	DeleteEvent(ctx context.Context, id string) error
}

type EventHandler struct {
	service EventServiceInterface
}

func NewEventHandler(service EventServiceInterface) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) GetEventByFilterAndFields(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	detailRequest := request.DetailRequest{}

	err := render.DecodeJSON(r.Body, &detailRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to decode details"))
		return
	}

	eventDetails, err := h.service.GetEventByFilterAndFields(ctx, detailRequest.EventDTO, detailRequest.Fields)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get event"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventDetails,
	})
}

func (h *EventHandler) GetEventsByFilterAndFields(w http.ResponseWriter, r *http.Request) {
	detailRequest := request.DetailRequest{}

	ctx := r.Context()

	err := render.DecodeJSON(r.Body, &detailRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to decode details"))
		return
	}

	eventDetails, err := h.service.GetEventsByFilterAndFields(ctx,
		detailRequest.EventDTO, detailRequest.Fields, detailRequest.Offset, detailRequest.Limit, detailRequest.Order)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get events"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventDetails,
	})
}

func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	offset, limit, err := parsing.ParsePageLimitToOffsetLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get events"))
		return
	}

	eventsDTO, err := h.service.GetAllEvents(ctx, offset, limit)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get events"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventsDTO,
	})
}

func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	receivedID := chi.URLParam(r, "id")
	eventDTO, err := h.service.GetEventByID(ctx, receivedID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get event"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventDTO,
	})
}

func (h *EventHandler) GetEventsTypeRegionalStage(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	orderStr := r.URL.Query().Get("order")

	offset, limit, err := parsing.ParsePageLimitToOffsetLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get events"))
		return
	}

	eventsDTO, err := h.service.GetEventsByType(ctx, event.RegionalStage, offset, limit, &orderStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get events by type"))
		return
	}

	if offset == nil && limit == nil {
		render.JSON(w, r, response.ApiResponse{
			Status: response.StatusOK,
			Data:   eventsDTO,
		})
		return
	}

	count, err := h.service.GetCountEventsByType(ctx, event.RegionalStage)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get count events by type"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:   response.StatusOK,
		Data:     eventsDTO,
		Metadata: count,
	})

}

func (h *EventHandler) GetEventsByClassType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	orderStr := r.URL.Query().Get("order")

	offset, limit, err := parsing.ParsePageLimitToOffsetLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get events"))
		return
	}

	eventsDTO, err := h.service.GetEventsByType(ctx, event.Class, offset, limit, &orderStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get events by type"))
		return
	}

	if offset == nil && limit == nil {
		render.JSON(w, r, response.ApiResponse{
			Status: response.StatusOK,
			Data:   eventsDTO,
		})
		return
	}

	count, err := h.service.GetCountEventsByType(ctx, event.Class)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get count events by type"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventsDTO,
		Metadata: map[string]int{
			"amount": int(count),
		},
	})

}

func (h *EventHandler) GetCountEventTypeClass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	count, err := h.service.GetCountEventsByType(ctx, event.Class)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get count events by type"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:     response.StatusOK,
		StatusCode: http.StatusOK,
		Data:       count,
	})
}

func (h *EventHandler) GetEventsTypeStageAndHisChilds(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	receivedID := chi.URLParam(r, "id")
	searchedID, err := uuid.Parse(receivedID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to parse id"))
		return
	}

	eventsDto, err := h.service.GetEventsTypeStageAndHisChilds(ctx, searchedID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get events"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventsDto,
	})
}

func (h *EventHandler) GetEventsByPreviousID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	orderStr := r.URL.Query().Get("order")

	offset, limit, err := parsing.ParsePageLimitToOffsetLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to pares offset/limit"))
		return
	}

	receivedID := chi.URLParam(r, "id")

	eventsDTO, err := h.service.GetEventsByPreviousID(ctx, receivedID, offset, limit, &orderStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get events"))
		return
	}

	if offset == nil && limit == nil {
		render.JSON(w, r, response.ApiResponse{
			StatusCode: http.StatusOK,
			Status:     response.StatusOK,
			Data:       eventsDTO,
		})
		return
	}

	count, err := h.service.GetCountEventsByPreviousID(ctx, receivedID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get count events by type"))
		return
	}

	render.JSON(w, r, response.ApiResponse{
		Status:   response.StatusOK,
		Data:     eventsDTO,
		Metadata: int(count),
	})
}

func (h *EventHandler) GetEventsByListID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	type ReqIds struct {
		IDs []uuid.UUID `json:"ids"`
	}
	var ids ReqIds
	err := render.DecodeJSON(r.Body, &ids)
	if errors.Is(err, io.EOF) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("empty request"))
		return
	}
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failde to decode request"))
		return
	}

	eventsDTO, err := h.service.GetEventsByListID(ctx, ids.IDs)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to get event"))
		return
	}
	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   eventsDTO,
	})
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var eventDTO event_dto.CreateEventDTORequest

	err := render.DecodeJSON(r.Body, &eventDTO)
	if errors.Is(err, io.EOF) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("empty request"))
		return
	}
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failde to decode request"))
		return
	}

	err = validator.New().Struct(eventDTO)
	if err != nil {
		validateErr := err.(validator.ValidationErrors)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse(fmt.Sprintf("err %v", validateErr)))
		return
	}

	id, err := h.service.CreateEvent(ctx, eventDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to create event"))
		return
	}
	render.JSON(w, r, response.SuccessResponse(fmt.Sprintf("id = %v", id)))
}

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var eventDTO event_dto.UpdateEventDTORequest
	receivedID := chi.URLParam(r, "id")

	err := render.DecodeJSON(r.Body, &eventDTO)
	if errors.Is(err, io.EOF) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("empty request"))
		return
	}
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to decode event"))
		return
	}

	err = h.service.UpdateEvent(ctx, receivedID, eventDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to update event"))
		return
	}

	render.JSON(w, r, response.SuccessResponse("success update"))
}

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	receivedID := chi.URLParam(r, "id")

	err := h.service.DeleteEvent(ctx, receivedID)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorResponse("failed to delete event"))
		return
	}

	render.JSON(w, r, response.SuccessResponse("Object deleted"))
}
