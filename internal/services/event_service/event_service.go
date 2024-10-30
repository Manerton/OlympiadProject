package event_service

import (
	"errors"
	"fmt"
	"main/internal/dto/event_dto"
	"main/internal/models/event"
	"main/internal/repositories/event_repository"
	"sync"

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
func (s *EventService) GetAllEvents() ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetAllEvents"
	events, err := s.repository.GetAllEvents(s.db)
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

// Get list events by type
func (s *EventService) GetEventsByType(event_type event.EventType) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByType"
	events, err := s.repository.GetEventsByType(s.db, event_type)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

// Get list events by PreviousID
func (s *EventService) GetEventsByPreviousID(id uint) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByPreviousID"
	events, err := s.repository.GetEventsByPreviousID(s.db, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

func (s *EventService) GetEventsByListID(ids []uint) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByListID"
	var wg sync.WaitGroup
	events := make([]event.Event, len(ids))
	errors := make(chan error, len(ids))

	for i, id := range ids {
		wg.Add(1)
		go func(i int, id uint) {
			defer wg.Done()
			event, err := s.repository.GetEventByID(s.db, id)
			if err != nil {
				errors <- err
				return
			}
			events[i] = event
		}(i, id)
	}
	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return nil, fmt.Errorf("%s: %w", op, <-errors)
	}
	return ConvertManyEventsToDTO(events), nil
}

func (s *EventService) checkCorrectEventDTO(event_dto *event_dto.EventDTO) error {
	// Check date
	if event_dto.StartDate.IsZero() || event_dto.EndDate.IsZero() {
		return errors.New("date does not exist")
	}
	if !event_dto.StartDate.Before(event_dto.EndDate) {
		return errors.New("start date should be before end date")
	}
	// check have previousEvent
	previousEventID := event_dto.PreviousEventID
	if previousEventID != nil {
		previousEvent, err := s.repository.GetEventByID(s.db, *previousEventID)
		if err != nil {
			return err
		}

		// check correct date border
		if previousEvent.StartDate.After(event_dto.StartDate) || previousEvent.EndDate.Before(event_dto.EndDate) {
			return errors.New("incorrect date limits")
		}

		// set correct type
		switch previousEvent.EventType {
		case event.RegionalStage:
			if event_dto.SubjectID == nil {
				return errors.New("subject id does not exist")
			}
			event_dto.EventType = event.Olympiad
		case event.Olympiad:
			event_dto.EventType = event.Stage
		case event.Stage:
			// check stage cannot have more than one appeal
			appeals, err := s.repository.GetEventsByPreviousID(s.db, *previousEventID)
			if err != nil {
				return err
			}
			if len(appeals) > 0 {
				return errors.New("stage cannot have more than one appeal")
			}
			event_dto.EventType = event.Appeal
		}
	} else {
		event_dto.EventType = event.RegionalStage
	}
	return nil
}

// Create event
func (s *EventService) CreateEvent(event_dto event_dto.EventDTO) (uint, error) {
	const op = "services.event_service.CreateEvent"

	// // Check date
	// if event_dto.StartDate.IsZero() || event_dto.EndDate.IsZero() {
	// 	return 0, fmt.Errorf("%s: date does not exist", op)
	// }
	// if !event_dto.StartDate.Before(event_dto.EndDate) {
	// 	return 0, fmt.Errorf("%s: start date should be before end date", op)
	// }
	// // check have previousEvent
	// previousEventID := event_dto.PreviousEventID
	// if previousEventID != 0 {
	// 	previousEvent, err := s.repository.GetEventByID(s.db, previousEventID)
	// 	if err != nil {
	// 		return 0, fmt.Errorf("%s: %w", op, err)
	// 	}

	// 	// check correct date border
	// 	if previousEvent.StartDate.After(event_dto.StartDate) || previousEvent.EndDate.Before(event_dto.EndDate) {
	// 		return 0, fmt.Errorf("%s: incorrect date limits", op)
	// 	}

	// 	// set correct type
	// 	switch previousEvent.EventType {
	// 	case event.RegionalStage:
	// 		event_dto.EventType = event.Olympiad
	// 	case event.Olympiad:
	// 		event_dto.EventType = event.Stage
	// 	case event.Stage:
	// 		// check stage cannot have more than one appeal
	// 		appeals, err := s.repository.GetEventsByPreviousID(s.db, previousEventID)
	// 		if err != nil {
	// 			return 0, fmt.Errorf("%s: %w", op, err)
	// 		}
	// 		if len(appeals) > 0 {
	// 			return 0, fmt.Errorf("%s: stage cannot have more than one appeal", op)
	// 		}
	// 		event_dto.EventType = event.Appeal
	// 	}
	// } else {
	// 	event_dto.EventType = event.RegionalStage
	// }

	err := s.checkCorrectEventDTO(&event_dto)
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

// Update event
func (s *EventService) UpdateEvent(event_dto event_dto.EventDTO) (uint, error) {
	const op = "services.event_service.UpdateEvent"
	// Check correct data
	err := s.checkCorrectEventDTO(&event_dto)
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

func ConvertDTOtoEvent(event_dto event_dto.EventDTO) event.Event {
	return event.Event{
		Model:           gorm.Model{ID: uint(event_dto.ID)},
		Name:            event_dto.Name,
		StartDate:       event_dto.StartDate,
		EndDate:         event_dto.EndDate,
		PreviousEventID: event_dto.PreviousEventID,
		SubjectID:       event_dto.SubjectID,
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
		PreviousEventID: event.PreviousEventID,
		SubjectID:       event.SubjectID,
		AdditionalInfo:  event.AdditionalInfo,
		EventType:       event.EventType,
	}
}

func ConvertManyEventsToDTO(events []event.Event) []event_dto.EventDTO {
	var events_dto []event_dto.EventDTO
	for _, event := range events {
		events_dto = append(events_dto, ConvertEventToDTO(event))
	}
	return events_dto
}
