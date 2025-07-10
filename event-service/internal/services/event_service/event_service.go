package event_service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"main/internal/dto/event_dto"
	"main/internal/lib/liblogger"
	"main/internal/lib/mapper/event_mapper"
	"main/internal/models/event"
	"main/internal/models/subject"
	"main/internal/storage/orm"
	"sync"
	"time"

	"github.com/google/uuid"
)

type EventRepositoryInterface interface {
	GetEventByFilterAndFields(ctx context.Context, orm orm.ORM, filter event.Event, fields *[]string) (event.Event, error)
	GetEventsByFilterAndFields(ctx context.Context, orm orm.ORM, filter event.Event, fields *[]string, offset, limit *int, order *string) ([]event.Event, error)
	GetEventByID(ctx context.Context, orm orm.ORM, id uuid.UUID) (event.Event, error)
	GetEventsByListID(ctx context.Context, orm orm.ORM, ids []uuid.UUID) ([]event.Event, error)
	GetEventsByType(ctx context.Context, orm orm.ORM, eventType event.EventType, offset, limit *int, order *string) ([]event.Event, error)
	GetEventsByPreviousID(ctx context.Context, orm orm.ORM, previousID uuid.UUID, offset, limit *int, order *string) ([]event.Event, error)
	GetAllEvents(ctx context.Context, orm orm.ORM, offset, limit *int) ([]event.Event, error)

	GetCountEventsByType(ctx context.Context, orm orm.ORM, eventType event.EventType) (int64, error)
	GetCountEventsByPreviousID(ctx context.Context, orm orm.ORM, previousID uuid.UUID) (int64, error)

	CreateEvent(ctx context.Context, orm orm.ORM, event event.Event) (uuid.UUID, error)
	UpdateEvent(ctx context.Context, orm orm.ORM, event event.Event) (uuid.UUID, error)
	DeleteEvent(ctx context.Context, orm orm.ORM, id uuid.UUID) error
}

type EventService struct {
	db         orm.ORM
	repository EventRepositoryInterface
	log        *slog.Logger
}

func NewEventService(orm orm.ORM, er EventRepositoryInterface, log *slog.Logger) *EventService {
	return &EventService{
		db:         orm,
		repository: er,
		log:        log,
	}
}

// Get all events
func (s *EventService) GetAllEvents(ctx context.Context, offset, limit *int) ([]event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetAllEvents"

	log := s.log.With(
		slog.String("op", op),
	)

	events, err := s.repository.GetAllEvents(ctx, s.db, offset, limit)
	if err != nil {
		log.Error("failed to get all events: %v", liblogger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return event_mapper.ManyToDTO(events), nil
}

// Get event by id
func (s *EventService) GetEventByID(ctx context.Context, id string) (event_dto.EventDTO, error) {
	const op = "services.event_service.GetEventByID"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed with id", liblogger.Err(err))
		return event_dto.EventDTO{}, fmt.Errorf("failed to parse id")
	}

	if uid == uuid.Nil {
		log.Error("failed with id", slog.Any("invalid id:", uid))
		return event_dto.EventDTO{}, fmt.Errorf("%s: invalid ID %d", op, uid)
	}

	event, err := s.repository.GetEventByID(ctx, s.db, uid)
	if err != nil {
		log.Error("failed to get event by ID", liblogger.Err(err))
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

	log := s.log.With(
		slog.String("op", op),
	)

	modelFilter := ConvertDTOtoEvent(filter)
	event, err := s.repository.GetEventByFilterAndFields(ctx, s.db, modelFilter, fields)
	if err != nil {
		log.Error("failed to get event",
			slog.Any("filter", filter),
			slog.Any("fields", fields),
			liblogger.Err(err),
		)
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
func (s *EventService) GetEventsByFilterAndFields(ctx context.Context, filter event_dto.EventDTO, fields *[]string, offset, limit *int, order *string) ([]event_dto.DetailsEvent, error) {
	const op = "services.event_service.GetEventsByFilterAndFields"

	log := s.log.With(
		slog.String("op", op),
	)

	modelFilter := ConvertDTOtoEvent(filter)
	events, err := s.repository.GetEventsByFilterAndFields(ctx, s.db, modelFilter, fields, offset, limit, order)
	if err != nil {
		log.Error("failed to get events",
			slog.Any("filter", filter),
			slog.Any("fields", fields),
			liblogger.Err(err),
		)
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return ConvertManyEventsToDetails(events), nil
}

// Get count events by event type (for pagination)
func (s *EventService) GetCountEventsByType(ctx context.Context, eventType event.EventType) (int64, error) {
	const op = "services.event_service.GetCountEventsByType"

	log := s.log.With(
		slog.String("op", op),
	)

	count, err := s.repository.GetCountEventsByType(ctx, s.db, eventType)
	if err != nil {
		log.Error("failed to get count events by type",
			slog.Any("eventType", eventType),
			liblogger.Err(err),
		)
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return count, err
}

// Get count events by previous id (for pagination)
func (s *EventService) GetCountEventsByPreviousID(ctx context.Context, id string) (int64, error) {
	const op = "services.event_service.GetCountEvents"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", liblogger.Err(err))
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	count, err := s.repository.GetCountEventsByPreviousID(ctx, s.db, uid)
	if err != nil {
		log.Error("failed to get count events by PreviousID",
			slog.Any("id", id),
			liblogger.Err(err),
		)
		return 0, fmt.Errorf("%s: %v", op, err)
	}
	return count, err
}

// Get list events by type
func (s *EventService) GetEventsByType(ctx context.Context, eventType event.EventType, offset, limit *int, order *string) ([]event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetEventsByType"

	log := s.log.With(
		slog.String("op", op),
	)

	events, err := s.repository.GetEventsByType(ctx, s.db, eventType, offset, limit, order)
	if err != nil {
		log.Error("failed to get events by type",
			slog.Any("eventType", eventType),
			liblogger.Err(err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return event_mapper.ManyToDTO(events), nil
}

// Get events where type=stage and his childs
func (s *EventService) GetEventsTypeStageAndHisChilds(ctx context.Context, id uuid.UUID) ([]event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetEventsTypeStageAndHisChilds"

	log := s.log.With(
		slog.String("op", op),
	)
	// Get all event stage by previousID
	// tx := s.db.Begin()
	events, err := s.repository.GetEventsByPreviousID(ctx, s.db, id, nil, nil, nil)
	if err != nil {
		// tx.Rollback()
		log.Error("failed to get events type stage by PreviousID",
			slog.Any("id", id),
			liblogger.Err(err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	eventsDto := event_mapper.ManyToDTO(events)

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
			childs, err := s.repository.GetEventsByPreviousID(ctx, s.db, id, nil, nil, nil)
			if err != nil {
				errors <- err
			}
			*eventsDto[i].Events = event_mapper.ManyToDTO(childs)
			mx.Unlock()
		}()
	}
	wg.Wait()
	if len(errors) > 0 {
		// tx.Rollback()
		log.Error("failed to get event childs by id",
			slog.Any("id", id),
			liblogger.Err(err),
		)
		return nil, fmt.Errorf("%s: %w", op, <-errors)
	}
	// tx.Commit()
	return eventsDto, nil
}

// Get list events by PreviousID
func (s *EventService) GetEventsByPreviousID(ctx context.Context, previousId string, offset, limit *int, order *string) ([]event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetEventsByPreviousID"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(previousId)
	if err != nil {
		log.Error("failed to parse PreviousID",
			slog.Any("id", previousId),
			liblogger.Err(err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	events, err := s.repository.GetEventsByPreviousID(ctx, s.db, uid, offset, limit, order)
	if err != nil {
		log.Error("failed to get events by PreviousID",
			slog.Any("id", previousId),
			liblogger.Err(err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return event_mapper.ManyToDTO(events), nil
}

// Get list events by list id
func (s *EventService) GetEventsByListID(ctx context.Context, ids []uuid.UUID) ([]event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetEventsByListID"

	log := s.log.With(
		slog.String("op", op),
	)

	events, err := s.repository.GetEventsByListID(ctx, s.db, ids)
	if err != nil {
		log.Error("failed to get events by listID",
			slog.Any("ids", ids),
			liblogger.Err(err),
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("events getted", slog.Any("events", events))

	return event_mapper.ManyToDTO(events), nil
}

// Find olympiad who is parent for stage
func (s *EventService) getOlympiad(ctx context.Context, id uuid.UUID) (event.Event, error) {

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

func (s *EventService) checkCorrectEventDTO(ctx context.Context, eventDTO *event.Event, isUpdate bool) error {
	const op = "services.event_service.checkCorrectEventDTO"

	// Check date
	if eventDTO.StartDate.IsZero() || eventDTO.EndDate.IsZero() {
		return fmt.Errorf("%s: %v", op, errors.New("date does not exist"))
	}
	if !eventDTO.StartDate.Before(eventDTO.EndDate) {
		return fmt.Errorf("%s: %v", op, errors.New("start date should be before end date"))
	}
	// check have previousEvent
	previousEventID := eventDTO.PreviousEventID
	if previousEventID != nil {
		previousEvent, err := s.repository.GetEventByID(ctx, s.db, *previousEventID)
		if err != nil {
			return fmt.Errorf("failed to get previousEvent: %v", err)
		}
		if !isUpdate {
			// set correct type
			switch previousEvent.EventType {
			case event.RegionalStage:
				if eventDTO.Subject == "" {
					return fmt.Errorf("%s: %v", op, errors.New("subject does not exist"))
				}
				eventDTO.EventType = event.Olympiad
			case event.Olympiad:
				eventDTO.EventType = event.Class
			case event.Class:
				eventDTO.EventType = event.Stage
			case event.Stage:
				// check stage cannot have more than one ViewWorks
				viewWorks, err := s.repository.GetEventsByPreviousID(ctx, s.db, *previousEventID, nil, nil, nil)
				if err != nil {
					return fmt.Errorf("%s: failed to get events when check stage: %v", op, err)
				}
				if len(viewWorks) == 0 {
					eventDTO.EventType = event.ViewWorks
				} else if len(viewWorks) == 1 {
					eventDTO.EventType = event.Appeal
				} else if len(viewWorks) > 2 {
					return fmt.Errorf("%s: %v", op, errors.New("stage cannot have more than one view works"))
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
				return fmt.Errorf("%s: %v", op, errors.New("incorrect date limits"))
			}
		} else {
			tempPreviousEvent, err := s.getOlympiad(ctx, *previousEventID)
			if err != nil {
				return fmt.Errorf("%s, failed to find first Olympiad: %v", op, err)
			}
			// events type view works and appeal should be after endDate parrent, but before endDate Olympiad
			if previousEvent.EndDate.After(eventDTO.StartDate) || tempPreviousEvent.EndDate.Before(eventDTO.EndDate) {
				return fmt.Errorf("%s: %v", op, errors.New("incorrect date limits"))
			}
		}

	} else {
		eventDTO.EventType = event.RegionalStage
	}
	return nil
}

// func (s *EventService) multipleCreateEvents(ctx context.Context, eventDTO event_dto.EventDTO, id *uuid.UUID) error {
// 	const op = "service.event_service.multipleCreateEvents"

// 	errGroup := errgroup.Group{}

// 	// Additional support struct
// 	type StackEvent struct {
// 		event event_dto.EventDTO
// 		id    *uuid.UUID
// 	}

// 	// channel with events and ids
// 	eventChannel := make(chan StackEvent, 50)

// 	// count worker equal max procs
// 	workerCount := runtime.GOMAXPROCS(0)

// 	for i := 0; i < workerCount; i++ {
// 		errGroup.Go(func() error {
// 			for current := range eventChannel {

// 				current.event.PreviousEventID = current.id
// 				err := s.checkCorrectEventDTO(ctx, &current.event, false)
// 				if err != nil {
// 					return err
// 				}

// 				eventModel := ConvertDTOtoEvent(current.event)
// 				newId, err := s.repository.CreateEvent(ctx, s.db, eventModel)
// 				if err != nil {
// 					return fmt.Errorf("%s: %w", op, err)
// 				}

// 				for _, event := range *current.event.Events {
// 					eventChannel <- StackEvent{
// 						event: event,
// 						id:    &newId,
// 					}
// 				}
// 			}
// 			return nil
// 		})
// 	}

// 	return nil
// }

// Create event
func (s *EventService) CreateEvent(ctx context.Context, eventDTO event_dto.CreateEventDTORequest) (uuid.UUID, error) {
	const op = "services.event_service.CreateEvent"
	const errMsg = "failed create event"
	log := s.log.With(
		slog.String("op", op),
	)

	eventModel := event_mapper.FromCreateToModel(eventDTO)

	log.Debug("test", eventDTO)

	err := s.checkCorrectEventDTO(ctx, &eventModel, false)
	if err != nil {
		log.Error("failed check correct event", liblogger.Err(err))
		return uuid.Nil, fmt.Errorf("%s", errMsg)
	}
	// Auto create events for all subject
	if eventModel.EventType == event.RegionalStage {
		id, err := s.createEventsBySubjects(ctx, eventModel)
		if err != nil {
			log.Error("failed create event by subjects", liblogger.Err(err))
			return uuid.Nil, fmt.Errorf("%s", errMsg)
		}
		log.Info("events success created (with all subjects)", slog.Any("eventID", id))
		return id, nil
	}
	id, err := s.repository.CreateEvent(ctx, s.db, eventModel)
	if err != nil {
		log.Error("failed to create event", liblogger.Err(err))
		return uuid.Nil, fmt.Errorf("%s", errMsg)
	}
	log.Info("event success created (only one)", slog.Any("eventID", id))
	return id, nil
}

func (s *EventService) createEventsBySubjects(ctx context.Context, eventModel event.Event) (uuid.UUID, error) {
	const op = "services.event_service.createEventsBySubjects"
	tx, err := s.db.TransactionBegin()
	if err != nil {
		tx.TransactionRollback()
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	id, err := s.repository.CreateEvent(ctx, tx, eventModel)
	if err != nil {
		tx.TransactionRollback()
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
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
			return uuid.Nil, fmt.Errorf("%s - create auto event by subjects: %w", op, err)
		}
	}
	tx.TransactionCommit()
	return id, nil
}

// func (s *EventService) updateEventDTO(ctx context.Context, updatedEventDTO event_dto.UpdateEventDTORequest) (event_dto.EventDTO, error) {
// 	serchedID := updatedEventDTO.ID

// 	event, err := s.repository.GetEventByID(ctx, s.db, serchedID)
// 	newEventDTO := ConvertEventToDTO(event)
// 	if err != nil {
// 		return event_dto.EventDTO{}, nil
// 	}
// 	if updatedEventDTO.Name != nil {
// 		newEventDTO.Name = *updatedEventDTO.Name
// 	}
// 	if !updatedEventDTO.StartDate.IsZero() {
// 		newEventDTO.StartDate = updatedEventDTO.StartDate
// 	}
// 	if !updatedEventDTO.EndDate.IsZero() {
// 		newEventDTO.EndDate = updatedEventDTO.EndDate
// 	}
// 	if updatedEventDTO.Subject != "" {
// 		newEventDTO.Subject = updatedEventDTO.Subject
// 	}
// 	if updatedEventDTO.AdditionalInfo != "" {
// 		newEventDTO.AdditionalInfo = updatedEventDTO.AdditionalInfo
// 	}

// 	return newEventDTO, nil
// }

// Update event
func (s *EventService) UpdateEvent(ctx context.Context, id string, eventDTO event_dto.UpdateEventDTORequest) (uuid.UUID, error) {
	const op = "services.event_service.UpdateEvent"
	const errMsg = "failed update event"
	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uuid", liblogger.Err(err))
		return uuid.Nil, fmt.Errorf("%s", errMsg)
	}

	event := event_mapper.FromUpdateToModel(eventDTO, uid)

	err = s.checkCorrectEventDTO(ctx, &event, true)
	if err != nil {
		log.Error("failed check correct event", liblogger.Err(err))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	updatedId, err := s.repository.UpdateEvent(ctx, s.db, event)
	if err != nil {
		log.Error("failed update event", liblogger.Err(err))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("event success updated")
	return updatedId, nil
}

// Delete event
func (s *EventService) DeleteEvent(ctx context.Context, id string) error {
	const op = "services.event_service.DeleteEvent"
	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse id", liblogger.Err(err))
		return fmt.Errorf("failed parse id %s: %w", op, err)
	}

	err = s.repository.DeleteEvent(ctx, s.db, uid)
	if err != nil {
		log.Error("failed to delete event", liblogger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("event success deleted")
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

func ConvertManyEventsToDetails(events []event.Event) []event_dto.DetailsEvent {
	var eventsDTO []event_dto.DetailsEvent
	for _, event := range events {
		eventsDTO = append(eventsDTO, ConvertEventToDetails(event))
	}
	return eventsDTO
}
