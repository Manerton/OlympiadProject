package event_handler

import (
	"context"
	"fmt"
	"main/internal/dto/event_dto"
	"main/internal/lib/errs"
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
	GetEventsTypeStageAndHisChilds(ctx context.Context, id string) ([]event_dto.EventDTOResponse, error)
	GetEventsByPreviousID(ctx context.Context, previousID string, offset, limit *int, order *string) ([]event_dto.EventDTOResponse, error)
	GetEventsByListID(ctx context.Context, ids []string) ([]event_dto.EventDTOResponse, error)
	GetAvailableEventsByClass(ctx context.Context, parentEventId string, class string) ([]event_dto.EventDTOResponse, error)

	FinishedEvents(ctx context.Context, id string) error

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

// @Summery Search event by filter
// @Security BearerAuth
// @Description Получение события с нужными полями по фильтру полей
// @Tags events
// @Accept json
// @Produce json
// @Param credentials body request.DetailRequest true "Данные для поиска"
// @Success 200 {object} response.ApiResponse{data=event_dto.DetailsEvent}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/details/one [post]
func (h *EventHandler) GetEventByFilterAndFields(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	detailRequest := request.DetailRequest{}

	err := render.DecodeJSON(r.Body, &detailRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("")))
		return
	}

	eventDetails, err := h.service.GetEventByFilterAndFields(ctx, detailRequest.EventDTO, detailRequest.Fields)
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
		Status: response.StatusOK,
		Data:   eventDetails,
	})
}

// @Summery Search events by filter
// @Security BearerAuth
// @Description Получение событий с нужными полями по фильтру полей
// @Tags events
// @Accept json
// @Produce json
// @Param credentials body request.DetailRequest true "Данные для поиска"
// @Success 200 {object} response.ApiResponse{data=[]event_dto.DetailsEvent}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/details [post]
func (h *EventHandler) GetEventsByFilterAndFields(w http.ResponseWriter, r *http.Request) {
	detailRequest := request.DetailRequest{}

	ctx := r.Context()

	err := render.DecodeJSON(r.Body, &detailRequest)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode json")))
		return
	}

	eventDetails, err := h.service.GetEventsByFilterAndFields(ctx,
		detailRequest.EventDTO, detailRequest.Fields, detailRequest.Offset, detailRequest.Limit, detailRequest.Order)
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
		Status: response.StatusOK,
		Data:   eventDetails,
	})
}

// @Summery Get all
// @Security BearerAuth
// @Description Получение всех событий
// @Tags events
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Органичение на количество записей"
// @Success 200 {object} response.ApiResponse{data=[]event_dto.EventDTOResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events [get]
func (h *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	offset, limit, err := parsing.ParsePageLimitToOffsetLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed parse page/limit")))
		return
	}

	eventsDTO, err := h.service.GetAllEvents(ctx, offset, limit)
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
		Status:     response.StatusOK,
		StatusCode: http.StatusOK,
		Data:       eventsDTO,
	})
}

// @Summery Get event by id
// @Security BearerAuth
// @Description Получение события по id
// @Tags events
// @Produce json
// @Param id path string true "id события"
// @Success 200 {object} response.ApiResponse{data=event_dto.EventDTOResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/{id} [get]
func (h *EventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	receivedID := chi.URLParam(r, "id")
	eventDTO, err := h.service.GetEventByID(ctx, receivedID)
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
		Status: response.StatusOK,
		Data:   eventDTO,
	})
}

// @Summery Get events type regional stage
// @Security BearerAuth
// @Description Получение всех событий типа "Региональный этап"
// @Tags events
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Органичение на количество записей"
// @Param order query string false "Поле по которому необходимо сортировать и указания для сортировки Пример(name DESC)"
// @Success 200 {object} response.ApiResponse{data=[]event_dto.EventDTOResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/regional-stage [get]
func (h *EventHandler) GetEventsTypeRegionalStage(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	orderStr := r.URL.Query().Get("order")

	offset, limit, err := parsing.ParsePageLimitToOffsetLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed parse page/limit")))
		return
	}

	eventsDTO, err := h.service.GetEventsByType(ctx, event.RegionalStage, offset, limit, &orderStr)
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

	if offset == nil && limit == nil {
		render.JSON(w, r, response.ApiResponse{
			Status: response.StatusOK,
			Data:   eventsDTO,
		})
		return
	}

	count, err := h.service.GetCountEventsByType(ctx, event.RegionalStage)
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
		Status:   response.StatusOK,
		Data:     eventsDTO,
		Metadata: count,
	})

}

func (h *EventHandler) GetEventsTypeOlympiads(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := h.service.GetAllEvents(ctx, nil, nil)
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
		Status:     response.StatusOK,
		StatusCode: http.StatusOK,
		Data:       result,
	})
}

// @Summery Get events type class
// @Security BearerAuth
// @Description Получение всех событий типа "Класс"
// @Tags events
// @Produce json
// @Param page query int false "Номер страницы"
// @Param limit query int false "Органичение на количество записей"
// @Param order query string false "Поле по которому необходимо сортировать и указания для сортировки Пример(name DESC)"
// @Success 200 {object} response.ApiResponse{data=[]event_dto.EventDTOResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/class [get]
func (h *EventHandler) GetEventsClassType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	orderStr := r.URL.Query().Get("order")

	offset, limit, err := parsing.ParsePageLimitToOffsetLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed parse page/limit")))
		return
	}

	eventsDTO, err := h.service.GetEventsByType(ctx, event.Class, offset, limit, &orderStr)
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

	if offset == nil && limit == nil {
		render.JSON(w, r, response.ApiResponse{
			Status: response.StatusOK,
			Data:   eventsDTO,
		})
		return
	}

	count, err := h.service.GetCountEventsByType(ctx, event.Class)
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
		Status: response.StatusOK,
		Data:   eventsDTO,
		Metadata: map[string]int{
			"amount": int(count),
		},
	})
}

// @Summery Get count event type class
// @Description Получение количества событий типа "Класс"
// @Tags events
// @Produce json
// @Success 200 {object} response.ApiResponse{data=int}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/class/count [get]
func (h *EventHandler) GetCountEventTypeClass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	count, err := h.service.GetCountEventsByType(ctx, event.Class)
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
		Status:     response.StatusOK,
		StatusCode: http.StatusOK,
		Data:       count,
	})
}

// @Summery Get events type stage and his childs
// @Description Получение всех событий типа "Этап"
// @Tags events
// @Produce json
// @Param id path string true "id олимпиады типа Класс"
// @Success 200 {object} response.ApiResponse{data=[]event_dto.EventDTOResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/stages/{id} [get]
func (h *EventHandler) GetEventsTypeStageAndHisChilds(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	id := chi.URLParam(r, "id")

	eventsDto, err := h.service.GetEventsTypeStageAndHisChilds(ctx, id)
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
		Status:     response.StatusOK,
		StatusCode: http.StatusOK,
		Data:       eventsDto,
	})
}

func (h *EventHandler) GetAvailableEventsByClass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parentEventId := chi.URLParam(r, "id")
	classStr := r.URL.Query().Get("class")

	availableEvents, err := h.service.GetAvailableEventsByClass(ctx, parentEventId, classStr)
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
		Status:     response.StatusOK,
		StatusCode: http.StatusOK,
		Data:       availableEvents,
	})
}

// @Summery Get events by previous id
// @Description Получение всех дочерних событий первого уровня по id родителя
// @Tags events
// @Produce json
// @Param id path string true "id родителя (previous id)"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Органичение на количество записей"
// @Param order query string false "Поле по которому необходимо сортировать и указания для сортировки Пример(name DESC)"
// @Success 200 {object} response.ApiResponse{data=[]event_dto.EventDTOResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/child/{id} [get]
func (h *EventHandler) GetEventsByPreviousID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	orderStr := r.URL.Query().Get("order")

	offset, limit, err := parsing.ParsePageLimitToOffsetLimit(pageStr, limitStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed parse page/limit")))
		return
	}

	receivedID := chi.URLParam(r, "id")

	eventsDTO, err := h.service.GetEventsByPreviousID(ctx, receivedID, offset, limit, &orderStr)
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
		Status:   response.StatusOK,
		Data:     eventsDTO,
		Metadata: int(count),
	})
}

// @Summery Get events by list id
// @Security BearerAuth
// @Description Получение всех событий по списку id
// @Tags events
// @Accept json
// @Produce json
// @Param credentials body request.IdsRequest true "Список id событий"
// @Success 200 {object} response.ApiResponse{data=[]event_dto.EventDTOResponse}
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/list [post]
func (h *EventHandler) GetEventsByListID(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var ids request.IdsRequest
	err := render.DecodeJSON(r.Body, &ids)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode list id")))
		return
	}

	eventsDTO, err := h.service.GetEventsByListID(ctx, ids.Ids)
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
		Status: response.StatusOK,
		Data:   eventsDTO,
	})
}

// @Summery create event
// @Security BearerAuth
// @Description Создание события
// @Tags events
// @Accept json
// @Produce json
// @Param credentials body event_dto.CreateEventDTORequest true "Создание события"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/events [post]
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var eventDTO event_dto.CreateEventDTORequest

	err := render.DecodeJSON(r.Body, &eventDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode dto")))
		return
	}

	err = validator.New().Struct(eventDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("validate error")))
		return
	}

	id, err := h.service.CreateEvent(ctx, eventDTO)
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
	render.JSON(w, r, response.SuccessResponse(fmt.Sprintf("id = %v", id)))
}

// @Summery finished event
// @Security BearerAuth
// @Description Завершение событий
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "id события"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/finished/{id} [put]
func (h *EventHandler) FinishedEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	err := h.service.FinishedEvents(ctx, id)
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

	render.JSON(w, r, response.SuccessResponse("success finished"))
}

// @Summery update event
// @Security BearerAuth
// @Description Обновление события
// @Tags events
// @Accept json
// @Produce json
// @Param credentials body event_dto.UpdateEventDTORequest true "Данные для обновления события"
// @Param id path string true "id события"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/{id} [put]
func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var eventDTO event_dto.UpdateEventDTORequest
	receivedID := chi.URLParam(r, "id")

	err := render.DecodeJSON(r.Body, &eventDTO)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.ErrorApiResponse(errs.ErrBadRequest.Wrap("failed decode dto")))
		return
	}

	err = h.service.UpdateEvent(ctx, receivedID, eventDTO)
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

	render.JSON(w, r, response.SuccessResponse("success update"))
}

// @Summery Delete events
// @Security BearerAuth
// @Description Удаление события
// @Tags events
// @Produce json
// @Param id path string true "id события"
// @Success 200 {object} response.ApiResponse
// @Failure 400 {object} response.ApiResponse
// @Router /api/events/{id} [delete]
func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	receivedID := chi.URLParam(r, "id")

	err := h.service.DeleteEvent(ctx, receivedID)
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

	render.JSON(w, r, response.SuccessResponse("Object deleted"))
}
