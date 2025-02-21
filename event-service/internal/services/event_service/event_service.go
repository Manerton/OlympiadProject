package event_service

import (
	"context"
	"errors"
	"fmt"
	"main/internal/dto/event_dto"
	"main/internal/models/event"
	"main/internal/models/subject"
	"main/internal/storage/orm"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type EventRepositoryInterface interface {
	GetEventByFilterAndFields(ctx context.Context, orm orm.ORM, filter event.Event, fields *[]string) (event.Event, error)
	GetEventsByFilterAndFields(ctx context.Context, orm orm.ORM, filter event.Event, fields *[]string, offset, limit *int) ([]event.Event, error)
	GetEventByID(ctx context.Context, orm orm.ORM, id uint) (event.Event, error)
	GetEventsByListID(ctx context.Context, orm orm.ORM, ids []uint) ([]event.Event, error)
	GetEventsByType(ctx context.Context, orm orm.ORM, eventType event.EventType, offset, limit *int) ([]event.Event, error)
	GetEventsByPreviousID(ctx context.Context, orm orm.ORM, previousID uint, offset, limit *int) ([]event.Event, error)
	GetAllEvents(ctx context.Context, orm orm.ORM, offset, limit *int) ([]event.Event, error)

	GetCountEventsByType(ctx context.Context, orm orm.ORM, eventType event.EventType) (int64, error)
	GetCountEventsByPreviousID(ctx context.Context, orm orm.ORM, previousID uint) (int64, error)

	CreateEvent(ctx context.Context, orm orm.ORM, event event.Event) (uint, error)
	UpdateEvent(ctx context.Context, orm orm.ORM, event event.Event) (uint, error)
	DeleteEvent(ctx context.Context, orm orm.ORM, id uint) error
}

type EventService struct {
	db         orm.ORM
	repository EventRepositoryInterface
}

func NewEventService(orm orm.ORM, er EventRepositoryInterface) *EventService {
	return &EventService{
		db:         orm,
		repository: er,
	}
}

// Get all events
func (s *EventService) GetAllEvents(ctx context.Context, offset, limit *int) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetAllEvents"
	events, err := s.repository.GetAllEvents(ctx, s.db, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

// Get event by id
func (s *EventService) GetEventByID(ctx context.Context, id uint) (event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventByID"
	if id == 0 {
		return event_dto.EventDTO{}, fmt.Errorf("%s: invalid ID %d", op, id)
	}

	event, err := s.repository.GetEventByID(ctx, s.db, id)
	if err != nil {
		return event_dto.EventDTO{}, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertEventToDTO(event), nil
}

// Get one event by filter and fileds.
//
// Example:
//
// filter -> EventDTO{name: "test", subject: "math"}.
//
// fields -> ["name", subject, start_date].
func (s *EventService) GetEventByFilterAndFields(ctx context.Context, filter event_dto.EventDTO, fields *[]string) (event_dto.DetailsEvent, error) {
	const op = "services.eventservice.GetEventByFilterAndFields"
	modelFilter := ConvertDTOtoEvent(filter)
	event, err := s.repository.GetEventByFilterAndFields(ctx, s.db, modelFilter, fields)
	if err != nil {
		return event_dto.DetailsEvent{}, fmt.Errorf("%s: %v", op, err)
	}
	return ConvertEventToDetails(event), nil
}

// Get list events by filter and fileds.
//
// Example:
//
// filter -> EventDTO{name: "test", subject: "math"}.
//
// fields -> ["name", subject, start_date].
func (s *EventService) GetEventsByFilterAndFields(ctx context.Context, filter event_dto.EventDTO, fields *[]string, offset, limit *int) ([]event_dto.DetailsEvent, error) {
	const op = "services.event_service.GetEventsByFilterAndFields"
	modelFilter := ConvertDTOtoEvent(filter)
	events, err := s.repository.GetEventsByFilterAndFields(ctx, s.db, modelFilter, fields, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return ConvertManyEventsToDetails(events), nil
}

// Get count events by event type (for pagination)
func (s *EventService) GetCountEventsByType(ctx context.Context, eventType event.EventType) (int64, error) {
	const op = "services.event_service.GetCountEventsByType"
	count, err := s.repository.GetCountEventsByType(ctx, s.db, eventType)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return count, err
}

// Get count events by previous id (for pagination)
func (s *EventService) GetCountEventsByPreviousID(ctx context.Context, id uint) (int64, error) {
	const op = "services.event_service.GetCountEvents"
	count, err := s.repository.GetCountEventsByPreviousID(ctx, s.db, id)
	if err != nil {
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return count, err
}

// Get list events by type
func (s *EventService) GetEventsByType(ctx context.Context, eventType event.EventType, offset, limit *int) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByType"
	events, err := s.repository.GetEventsByType(ctx, s.db, eventType, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

// Get events where type=stage and his childs
func (s *EventService) GetEventsTypeStageAndHisChilds(ctx context.Context, id uint) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsTypeStageAndHisChilds"
	// Get all event stage by previousID
	// tx := s.db.Begin()
	events, err := s.repository.GetEventsByPreviousID(ctx, s.db, id, nil, nil)
	if err != nil {
		// tx.Rollback()
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
			childs, err := s.repository.GetEventsByPreviousID(ctx, s.db, id, nil, nil)
			if err != nil {
				errors <- err
			}
			mx.Lock()
			*eventsDto[i].Events = ConvertManyEventsToDTO(childs)
			mx.Unlock()
		}()
	}
	wg.Wait()
	if len(errors) > 0 {
		// tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, <-errors)
	}
	// tx.Commit()
	return eventsDto, nil
}

// Get list events by PreviousID
func (s *EventService) GetEventsByPreviousID(ctx context.Context, id uint, offset, limit *int) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByPreviousID"
	events, err := s.repository.GetEventsByPreviousID(ctx, s.db, id, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

// Get list events by list id
func (s *EventService) GetEventsByListID(ctx context.Context, ids []uint) ([]event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventsByListID"
	events, err := s.repository.GetEventsByListID(ctx, s.db, ids)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ConvertManyEventsToDTO(events), nil
}

// Find olympiad who is parent for stage
func (s *EventService) getOlympiad(ctx context.Context, id uint) (event.Event, error) {

	previousEvent, err := s.repository.GetEventByID(ctx, s.db, id)
	if err != nil {
		return event.Event{}, err
	}
	if previousEvent.EventType == event.Olympiad {
		return previousEvent, nil
	} else {
		return s.getOlympiad(ctx, *previousEvent.PreviousEventID)
	}
}

func (s *EventService) checkCorrectEventDTO(ctx context.Context, eventDTO *event_dto.EventDTO, isUpdate bool) error {
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
		previousEvent, err := s.repository.GetEventByID(ctx, s.db, *previousEventID)
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
				viewWorks, err := s.repository.GetEventsByPreviousID(ctx, s.db, *previousEventID, nil, nil)
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
			tempPreviousEvent, err := s.getOlympiad(ctx, *previousEventID)
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

func (s *EventService) CreateEventsByJSON(ctx context.Context, eventDTO event_dto.EventDTO) error {
	const op = "services.event_service.CreateEventsByJSON"

	err := s.checkCorrectEventDTO(ctx, &eventDTO, false)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	eventModel := ConvertDTOtoEvent(eventDTO)
	_, err = s.repository.CreateEvent(ctx, s.db, eventModel)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.createEventRecursion(ctx, eventDTO, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *EventService) createEventRecursion(ctx context.Context, eventDTO event_dto.EventDTO, id *uint) error {
	const op = "services.event_service.createEventRecursion"
	if id != nil {
		eventDTO.PreviousEventID = id
	}
	err := s.checkCorrectEventDTO(ctx, &eventDTO, false)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	eventModel := ConvertDTOtoEvent(eventDTO)
	newId, err := s.repository.CreateEvent(ctx, s.db, eventModel)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if eventDTO.Events == nil {
		return nil
	}

	errgroup := errgroup.Group{}

	for _, event := range *eventDTO.Events {
		errgroup.Go(func() error {
			err = s.createEventRecursion(ctx, event, &newId)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err = errgroup.Wait(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Create event
func (s *EventService) CreateEvent(ctx context.Context, eventDTO event_dto.EventDTO) (uint, error) {
	const op = "services.event_service.CreateEvent"
	err := s.checkCorrectEventDTO(ctx, &eventDTO, false)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	eventModel := ConvertDTOtoEvent(eventDTO)
	// Auto create events for all subject
	if eventModel.EventType == event.RegionalStage {
		id, err := s.createEventsBySubjects(ctx, eventModel)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return id, nil
	}
	id, err := s.repository.CreateEvent(ctx, s.db, eventModel)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *EventService) createEventsBySubjects(ctx context.Context, eventModel event.Event) (uint, error) {
	const op = "services.event_service.createEventsBySubjects"
	tx, err := s.db.TransactionBegin()
	if err != nil {
		tx.TransactionRollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := s.repository.CreateEvent(ctx, tx, eventModel)
	if err != nil {
		tx.TransactionRollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	for _, subject := range subject.NewSubjectsStorage().GetAllSubject() {
		eventBySubject := event.Event{
			Name:            fmt.Sprintf("Олимпиада по %s", subject),
			PreviousEventID: &id,
			StartDate:       eventModel.StartDate,
			EndDate:         eventModel.EndDate,
			EventType:       event.Olympiad,
			Subject:         subject,
		}
		_, err := s.repository.CreateEvent(ctx, tx, eventBySubject)
		if err != nil {
			tx.TransactionRollback()
			return 0, fmt.Errorf("%s - create auto event by subjects: %w", op, err)
		}
	}
	tx.TransactionCommit()
	return id, nil
}

func (s *EventService) updateEventDTO(ctx context.Context, updatedEventDTO event_dto.EventDTO) (event_dto.EventDTO, error) {
	serchedID := updatedEventDTO.ID

	event, err := s.repository.GetEventByID(ctx, s.db, serchedID)
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
func (s *EventService) UpdateEvent(ctx context.Context, event_dto event_dto.EventDTO) (uint, error) {
	const op = "services.event_service.UpdateEvent"

	event_dto, err := s.updateEventDTO(ctx, event_dto)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = s.checkCorrectEventDTO(ctx, &event_dto, true)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	event := ConvertDTOtoEvent(event_dto)
	id, err := s.repository.UpdateEvent(ctx, s.db, event)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

// Delete event
func (s *EventService) DeleteEvent(ctx context.Context, id uint) error {
	const op = "services.event_service.DeleteEvent"
	err := s.repository.DeleteEvent(ctx, s.db, id)
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
