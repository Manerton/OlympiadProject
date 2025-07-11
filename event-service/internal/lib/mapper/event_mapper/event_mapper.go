package event_mapper

import (
	"main/internal/dto/event_dto"
	"main/internal/models/event"

	"github.com/google/uuid"
)

func FromCreateToModel(createDTO event_dto.CreateEventDTORequest) event.Event {
	return event.Event{
		Name:            createDTO.Name,
		StartDate:       createDTO.StartDate,
		EndDate:         createDTO.EndDate,
		Subject:         createDTO.Subject,
		PreviousEventID: createDTO.PreviousEventID,
		AdditionalInfo:  createDTO.AdditionalInfo,
	}
}

func FromUpdateToModel(updateDTO event_dto.UpdateEventDTORequest, id uuid.UUID) event.Event {
	return event.Event{
		ID:             id,
		Name:           *updateDTO.Name,
		StartDate:      *updateDTO.StartDate,
		EndDate:        *updateDTO.EndDate,
		Subject:        *updateDTO.Subject,
		AdditionalInfo: *updateDTO.AdditionalInfo,
	}
}

func ToDTO(eventModel event.Event) event_dto.EventDTOResponse {
	return event_dto.EventDTOResponse{
		ID:              eventModel.ID,
		Name:            eventModel.Name,
		StartDate:       eventModel.StartDate,
		EndDate:         eventModel.EndDate,
		Subject:         eventModel.Subject,
		AdditionalInfo:  eventModel.AdditionalInfo,
		PreviousEventID: eventModel.PreviousEventID,
		EventType:       eventModel.EventType,
		Events:          new([]event_dto.EventDTOResponse),
	}
}

func ManyToDTO(events []event.Event) []event_dto.EventDTOResponse {
	eventsResult := make([]event_dto.EventDTOResponse, 0, len(events))
	for _, eventRes := range events {
		eventsResult = append(eventsResult, ToDTO(eventRes))
	}

	return eventsResult
}
