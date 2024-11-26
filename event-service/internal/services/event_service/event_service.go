package event_service

import (
	"errors"
	"fmt"
	"main/internal/dto/event_dto"
	"main/internal/models/event"
	"main/internal/repositories/event_repository"
	"sync"
	"time"

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

// Get all events
func (s *EventService) GetAllEvents(offset, limit *int) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetAllEvents"
	events, err := s.repository.GetAllEvents(s.db, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

// Get event by id
func (s *EventService) GetEventByID(id uint) (event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventByID"
	if id == 0 {
		return event_dto.EventDTO{}, fmt.Errorf("%s: invalid ID %d", op, id)
	}

	event, err := s.repository.GetEventByID(s.db, id)
	if err != nil {
		return event_dto.EventDTO{}, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertEventToDTO(event), nil
}

func (s *EventService) GetEventByFilterAndFields(filter event_dto.EventDTO, fields *[]string) (event_dto.DetailsEvent, error) {
	const op = "services.event_service.GetEventByFilterAndFields"
	modelFilter := ConvertDTOtoEvent(filter)
	event, err := s.repository.GetEventByFilterAndFields(s.db, modelFilter, fields)
	if err != nil {
		return event_dto.DetailsEvent{}, fmt.Errorf("%s: %v", op, err)
	}
	return ConvertEventToDetails(event), nil
}

func (s *EventService) GetEventsByFilterAndFields(filter event_dto.EventDTO, fields *[]string, offset, limit *int) ([]event_dto.DetailsEvent, error) {
	const op = "services.event_service.GetEventsByFilterAndFields"
	modelFilter := ConvertDTOtoEvent(filter)
	events, err := s.repository.GetEventsByFilterAndFields(s.db, modelFilter, fields, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return ConvertManyEventsToDetails(events), nil
}

// Get list events by type
func (s *EventService) GetEventsByType(event_type event.EventType, offset, limit *int) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByType"
	events, err := s.repository.GetEventsByType(s.db, event_type, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

func (s *EventService) GetEventsTypeStageAndHisChilds(id uint) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsTypeStageAndHisChilds"
	// Get all event stage by previousID
	tx := s.db.Begin()
	events, err := s.repository.GetEventsByPreviousID(tx, id, nil, nil)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	eventsDto := ConvertManyEventsToDTO(events)

	wg := sync.WaitGroup{}
	mx := sync.Mutex{}
	errors := make(chan error, len(eventsDto))
	defer close(errors)
	// Get all childs by id Stage's
	for i := 0; i < len(eventsDto); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mx.Lock()
			id := eventsDto[i].ID
			mx.Unlock()
			childs, err := s.repository.GetEventsByPreviousID(tx, id, nil, nil)
			if err != nil {
				errors <- err
			}
			mx.Lock()
			eventsDto[i].Events = ConvertManyEventsToDTO(childs)
			mx.Unlock()
		}()
	}
	wg.Wait()
	if len(errors) > 0 {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, <-errors)
	}
	tx.Commit()
	return eventsDto, nil
}

// Get list events by PreviousID
func (s *EventService) GetEventsByPreviousID(id uint, offset, limit *int) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByPreviousID"
	events, err := s.repository.GetEventsByPreviousID(s.db, id, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

// Get list events by list id
func (s *EventService) GetEventsByListID(ids []uint) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByListID"
	events, err := s.repository.GetEventsByListID(s.db, ids)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

func (s *EventService) getOlympiad(id uint) (event.Event, error) {

	previousEvent, err := s.repository.GetEventByID(s.db, id)
	if err != nil {
		return event.Event{}, err
	}
	if previousEvent.EventType == event.Olympiad {
		return previousEvent, nil
	} else {
		return s.getOlympiad(*previousEvent.PreviousEventID)
	}
}

func (s *EventService) checkCorrectEventDTO(eventDTO *event_dto.EventDTO, isUpdate bool) error {
	// Check date
	if eventDTO.StartDate.IsZero() || eventDTO.EndDate.IsZero() {
		return errors.New("date does not exist")
	}
	if !eventDTO.StartDate.Before(eventDTO.EndDate) {
		return errors.New("start date should be before end date")
	}
	// check have previousEvent
	previousEventID := eventDTO.PreviousEventID
	if previousEventID != nil {
		previousEvent, err := s.repository.GetEventByID(s.db, *previousEventID)
		if err != nil {
			return err
		}
		if !isUpdate {
			// set correct type
			switch previousEvent.EventType {
			case event.RegionalStage:
				if eventDTO.Subject == "" {
					return errors.New("subject does not exist")
				}
				eventDTO.EventType = event.Olympiad
			case event.Olympiad:
				eventDTO.EventType = event.Stage
			case event.Stage:
				// check stage cannot have more than one ViewWorks
				viewWorks, err := s.repository.GetEventsByPreviousID(s.db, *previousEventID, nil, nil)
				if err != nil {
					return err
				}
				if len(viewWorks) == 0 {
					eventDTO.EventType = event.ViewWorks
				} else if len(viewWorks) == 1 {
					eventDTO.EventType = event.Appeal
				} else if len(viewWorks) > 2 {
					return errors.New("stage cannot have more than one view works")
				}
			}
			// case event.Stage:
			// 	// check stage cannot have more than one ViewWorks
			// 	viewWorks, err := s.repository.GetEventsByPreviousID(s.db, *previousEventID, nil, nil)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	if len(viewWorks) > 0 {
			// 		return errors.New("stage cannot have more than one view works")
			// 	}
			// 	eventDTO.EventType = event.ViewWorks
			// case event.ViewWorks:
			// 	// check ViewWorks cannot have more than one appeal
			// 	appeal, err := s.repository.GetEventsByPreviousID(s.db, *previousEventID, nil, nil)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	if len(appeal) > 0 {
			// 		return errors.New("view works cannot have more than one appeal")
			// 	}
			// 	eventDTO.EventType = event.Appeal
			// }
		}

		// check correct date border
		if eventDTO.EventType == event.Stage || eventDTO.EventType == event.Olympiad {
			if previousEvent.StartDate.After(eventDTO.StartDate) || previousEvent.EndDate.Before(eventDTO.EndDate) {
				return errors.New("incorrect date limits")
			}
		} else {
			tempPreviousEvent, err := s.getOlympiad(*previousEventID)
			if err != nil {
				return err
			}
			// events type view works and appeal should be after endDate parrent, but before endDate Olympiad
			if previousEvent.EndDate.After(eventDTO.StartDate) || tempPreviousEvent.EndDate.Before(eventDTO.EndDate) {
				return errors.New("incorrect date limits")
			}
		}

	} else {
		eventDTO.EventType = event.RegionalStage
	}
	return nil
}

// Create event
func (s *EventService) CreateEvent(event_dto event_dto.EventDTO) (uint, error) {
	const op = "services.event_service.CreateEvent"
	err := s.checkCorrectEventDTO(&event_dto, false)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	event := ConvertDTOtoEvent(event_dto)
	id, err := s.repository.CreateEvent(s.db, event)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *EventService) updateEventDTO(updatedEventDTO event_dto.EventDTO) (event_dto.EventDTO, error) {
	serchedID := updatedEventDTO.ID

	event, err := s.repository.GetEventByID(s.db, serchedID)
	newEventDTO := ConvertEventToDTO(event)
	if err != nil {
		return event_dto.EventDTO{}, nil
	}
	if updatedEventDTO.Name != "" {
		newEventDTO.Name = updatedEventDTO.Name
	}
	if !updatedEventDTO.StartDate.IsZero() {
		newEventDTO.StartDate = updatedEventDTO.StartDate
	}
	if !updatedEventDTO.EndDate.IsZero() {
		newEventDTO.EndDate = updatedEventDTO.EndDate
	}
	if updatedEventDTO.Subject != "" {
		newEventDTO.Subject = updatedEventDTO.Subject
	}
	if updatedEventDTO.AdditionalInfo != "" {
		newEventDTO.AdditionalInfo = updatedEventDTO.AdditionalInfo
	}

	return newEventDTO, nil
}

// Update event
func (s *EventService) UpdateEvent(event_dto event_dto.EventDTO) (uint, error) {
	const op = "services.event_service.UpdateEvent"

	event_dto, err := s.updateEventDTO(event_dto)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = s.checkCorrectEventDTO(&event_dto, true)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	event := ConvertDTOtoEvent(event_dto)
	id, err := s.repository.UpdateEvent(s.db, event)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

// Delete event
func (s *EventService) DeleteEvent(id uint) error {
	const op = "services.event_service.DeleteEvent"
	err := s.repository.DeleteEvent(s.db, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func ConvertDTOtoEvent(eventDTO event_dto.EventDTO) event.Event {
	return event.Event{
		ID:              eventDTO.ID,
		Name:            eventDTO.Name,
		StartDate:       eventDTO.StartDate,
		EndDate:         eventDTO.EndDate,
		PreviousEventID: eventDTO.PreviousEventID,
		Subject:         eventDTO.Subject,
		AdditionalInfo:  eventDTO.AdditionalInfo,
		EventType:       eventDTO.EventType,
	}
}

func ConvertEventToDTO(event event.Event) event_dto.EventDTO {
	return event_dto.EventDTO{
		ID:              event.ID,
		Name:            event.Name,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
		PreviousEventID: event.PreviousEventID,
		Subject:         event.Subject,
		AdditionalInfo:  event.AdditionalInfo,
		EventType:       event.EventType,
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

	return event_dto.DetailsEvent{
		ID:              event.ID,
		Name:            event.Name,
		StartDate:       startDate,
		EndDate:         endDate,
		PreviousEventID: event.PreviousEventID,
		Subject:         event.Subject,
		AdditionalInfo:  event.AdditionalInfo,
		EventType:       event.EventType,
	}
}

func ConvertManyEventsToDTO(events []event.Event) []event_dto.EventDTO {
	var eventsDTO []event_dto.EventDTO
	for _, event := range events {
		eventsDTO = append(eventsDTO, ConvertEventToDTO(event))
	}
	return eventsDTO
}

func ConvertManyEventsToDetails(events []event.Event) []event_dto.DetailsEvent {
	var eventsDTO []event_dto.DetailsEvent
	for _, event := range events {
		eventsDTO = append(eventsDTO, ConvertEventToDetails(event))
	}
	return eventsDTO
}
