package event_service

import (
	"fmt"
	"main/internal/dto/event_dto"
	"main/internal/models/event"
	"main/internal/repositories/event_repository"

	"gorm.io/gorm"
)

type EventService struct {
	db         *gorm.DB
	repository *event_repository.EventRepository
}

func NewEventService(db *gorm.DB, er *event_repository.EventRepository) *EventService {
	return &EventService{
		db:         db,
		repository: er,
	}
}

func (e *EventService) GetEventByID(id uint) (event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventByID"
	event, err := e.repository.GetEventByID(e.db, id)
	if err != nil {
		return event_dto.EventDTO{}, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertEventToDTO(event), nil
}

func (e *EventService) GetEventsByType(event_type event.EventType) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByType"
	events, err := e.repository.GetEventsByType(e.db, event_type)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConverManyEventsToDTO(events), nil
}

func (e *EventService) GetAllEvents() ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetAllEvents"
	events, err := e.repository.GetAllEvents(e.db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConverManyEventsToDTO(events), nil
}

func (e *EventService) CreateEvent(event_dto event_dto.EventDTO) (uint, error) {
	const op = "services.event_service.CreateEvent"
	event := ConvertDTOtoEvent(event_dto)
	id, err := e.repository.CreateEvent(e.db, event)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (e *EventService) UpdateEvent(event_dto event_dto.EventDTO) (uint, error) {
	const op = "services.event_service.UpdateEvent"
	event := ConvertDTOtoEvent(event_dto)
	id, err := e.repository.UpdateEvent(e.db, event)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (e *EventService) DeleteEvent(id uint) error {
	const op = "services.event_service.DeleteEvent"
	err := e.repository.DeleteEvent(e.db, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func ConvertDTOtoEvent(event_dto event_dto.EventDTO) event.Event {
	return event.Event{
		Model:           gorm.Model{ID: uint(event_dto.ID)},
		Name:            event_dto.Name,
		StartDate:       event_dto.StartDate,
		EndDate:         event_dto.EndDate,
		PreviousEventID: &event_dto.PreviousEventID,
		SubjectID:       &event_dto.SubjectID,
		AdditionalInfo:  event_dto.AdditionalInfo,
		EventType:       event_dto.EventType,
	}
}

func ConvertEventToDTO(event event.Event) event_dto.EventDTO {
	return event_dto.EventDTO{
		ID:              event.ID,
		Name:            event.Name,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		PreviousEventID: *event.PreviousEventID,
		SubjectID:       *event.SubjectID,
		AdditionalInfo:  event.AdditionalInfo,
		EventType:       event.EventType,
	}
}

func ConverManyEventsToDTO(events []event.Event) []event_dto.EventDTO {
	var events_dto []event_dto.EventDTO
	for _, event := range events {
		events_dto = append(events_dto, ConvertEventToDTO(event))
	}
	return events_dto
}
