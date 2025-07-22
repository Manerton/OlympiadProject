package event_mapper

import (
	"main/internal/dto/event_dto"
	"main/internal/models/event"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func FromCreateToModel(createDTO event_dto.CreateEventDTORequest) event.Event {

	subjectInt, err := strconv.Atoi(createDTO.Subject)
	if err != nil {

	}

	return event.Event{
		Name:            createDTO.Name,
		StartDate:       createDTO.StartDate,
		EndDate:         createDTO.EndDate,
		Subject:         subjectInt,
		ClassNumber:     createDTO.ClassNumber,
		PreviousEventID: createDTO.PreviousEventID,
		AdditionalInfo:  createDTO.AdditionalInfo,
	}
}

func FromUpdateToModel(updateDTO event_dto.UpdateEventDTORequest, id uuid.UUID) event.Event {
	var model event.Event
	model.ID = id

	if updateDTO.Name != nil {
		model.Name = *updateDTO.Name
	}
	if updateDTO.StartDate != nil {
		model.StartDate = *updateDTO.StartDate
	}
	if updateDTO.EndDate != nil {
		model.EndDate = *updateDTO.EndDate
	}
	if updateDTO.Subject != nil {
		subjectInt, err := strconv.Atoi(*updateDTO.Subject)
		if err != nil {

		}
		model.Subject = subjectInt
	}
	if updateDTO.ClassNumber != nil {
		model.ClassNumber = *updateDTO.ClassNumber
	}
	if updateDTO.AdditionalInfo != nil {
		model.AdditionalInfo = *updateDTO.AdditionalInfo
	}

	return model
}

func ToDTO(eventModel event.Event) event_dto.EventDTOResponse {

	subjectStr := strconv.Itoa(eventModel.Subject)

	return event_dto.EventDTOResponse{
		ID:              eventModel.ID,
		Name:            eventModel.Name,
		StartDate:       eventModel.StartDate,
		EndDate:         eventModel.EndDate,
		Subject:         subjectStr,
		ClassNumber:     eventModel.ClassNumber,
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

func ConvertDTOtoEvent(eventDTO event_dto.EventDTO) event.Event {

	subjectInt, err := strconv.Atoi(eventDTO.Subject)
	if err != nil {

	}

	return event.Event{
		ID:              eventDTO.ID,
		Name:            eventDTO.Name,
		StartDate:       eventDTO.StartDate,
		EndDate:         eventDTO.EndDate,
		PreviousEventID: eventDTO.PreviousEventID,
		Subject:         subjectInt,
		AdditionalInfo:  eventDTO.AdditionalInfo,
		EventType:       eventDTO.EventType,
	}
}

func ConvertEventToDTO(event event.Event) event_dto.EventDTO {

	subjectStr := strconv.Itoa(event.Subject)

	return event_dto.EventDTO{
		ID:              event.ID,
		Name:            event.Name,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		PreviousEventID: event.PreviousEventID,
		Subject:         subjectStr,
		AdditionalInfo:  event.AdditionalInfo,
		EventType:       event.EventType,
		Events:          &[]event_dto.EventDTO{},
	}
}

func ConvertEventToDetails(event event.Event) event_dto.DetailsEvent {
	var startDate, endDate *time.Time

	if !event.StartDate.IsZero() {
		startDate = &event.StartDate
	}
	if !event.EndDate.IsZero() {
		endDate = &event.EndDate
	}

	subjectStr := strconv.Itoa(event.Subject)

	return event_dto.DetailsEvent{
		ID:              event.ID,
		Name:            event.Name,
		StartDate:       startDate,
		EndDate:         endDate,
		PreviousEventID: event.PreviousEventID,
		Subject:         subjectStr,
		AdditionalInfo:  event.AdditionalInfo,
		EventType:       event.EventType,
	}
}

func ConvertManyEventsToDetails(events []event.Event) []event_dto.DetailsEvent {
	var eventsDTO []event_dto.DetailsEvent
	for _, event := range events {
		eventsDTO = append(eventsDTO, ConvertEventToDetails(event))
	}
	return eventsDTO
}
