package event_service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"main/internal/dto/event_dto"
	"main/internal/lib/errs"
	"main/internal/lib/helper"
	"main/internal/lib/liblogger"
	"main/internal/lib/mapper/event_mapper"
	"main/internal/models/event"
	"main/internal/models/outbox"
	"main/internal/models/subject"
	"main/internal/storage/orm"
	"sync"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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
	UpdateEvent(ctx context.Context, orm orm.ORM, event event.Event) error
	DeleteEvent(ctx context.Context, orm orm.ORM, id uuid.UUID) error
}

type OutboxRepositoryInteface interface {
	Create(ctx context.Context, orm orm.ORM, outbox outbox.Outbox) (uuid.UUID, error)
}

type EventService struct {
	db                orm.ORM
	eventRepository   EventRepositoryInterface
	outboxRepository  OutboxRepositoryInteface
	kiznaiverServices map[string]string
	log               *slog.Logger
}

func NewEventService(orm orm.ORM, services map[string]string, er EventRepositoryInterface, or OutboxRepositoryInteface, log *slog.Logger) *EventService {
	return &EventService{
		db:                orm,
		kiznaiverServices: services,
		eventRepository:   er,
		outboxRepository:  or,
		log:               log,
	}
}

// Get all events
func (s *EventService) GetAllEvents(ctx context.Context, offset, limit *int) ([]event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetAllEvents"

	log := s.log.With(
		slog.String("op", op),
	)

	events, err := s.eventRepository.GetAllEvents(ctx, s.db, offset, limit)
	if err != nil {
		log.Error("failed to get all events: %v", liblogger.Err(err))
		return nil, errs.ErrInternalError
	}

	return event_mapper.ManyToDTO(events), nil
}

// Get event by id
func (s *EventService) GetEventByID(ctx context.Context, id string) (event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetEventByID"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", slog.String("id", id), liblogger.Err(err))
		return event_dto.EventDTOResponse{}, errs.ErrBadRequest.Wrap("invalid id")
	}

	event, err := s.eventRepository.GetEventByID(ctx, s.db, uid)
	if err != nil {
		log.Error("failed to get event by ID", liblogger.Err(err))
		return event_dto.EventDTOResponse{}, errs.ErrInternalError
	}
	return event_mapper.ToDTO(event), nil
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

	modelFilter := event_mapper.ConvertDTOtoEvent(filter)
	event, err := s.eventRepository.GetEventByFilterAndFields(ctx, s.db, modelFilter, fields)
	if err != nil {
		log.Error("failed to get event",
			slog.Any("filter", filter),
			slog.Any("fields", fields),
			liblogger.Err(err),
		)
		return event_dto.DetailsEvent{}, errs.ErrInternalError.Wrap("failed get event by filter and fields")
	}
	return event_mapper.ConvertEventToDetails(event), nil
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

	modelFilter := event_mapper.ConvertDTOtoEvent(filter)
	events, err := s.eventRepository.GetEventsByFilterAndFields(ctx, s.db, modelFilter, fields, offset, limit, order)
	if err != nil {
		log.Error("failed to get events",
			slog.Any("filter", filter),
			slog.Any("fields", fields),
			liblogger.Err(err),
		)
		return nil, errs.ErrInternalError.Wrap("failed get events by filter and fields")
	}
	return event_mapper.ConvertManyEventsToDetails(events), nil
}

// Get count events by event type (for pagination)
func (s *EventService) GetCountEventsByType(ctx context.Context, eventType event.EventType) (int64, error) {
	const op = "services.event_service.GetCountEventsByType"

	log := s.log.With(
		slog.String("op", op),
	)

	count, err := s.eventRepository.GetCountEventsByType(ctx, s.db, eventType)
	if err != nil {
		log.Error("failed to get count events by type",
			slog.Any("eventType", eventType),
			liblogger.Err(err),
		)
		return 0, errs.ErrInternalError.Wrap("failed get count events by event-type")
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
		log.Error("failed parse id", slog.String("id", id), liblogger.Err(err))
		return 0, errs.ErrBadRequest.Wrap("failed parse id to uuid")
	}

	count, err := s.eventRepository.GetCountEventsByPreviousID(ctx, s.db, uid)
	if err != nil {
		log.Error("failed to get count events by PreviousID",
			liblogger.Err(err),
		)
		return 0, errs.ErrInternalError.Wrap("failed get count events by prev id")
	}
	return count, err
}

// Get list events by type
func (s *EventService) GetEventsByType(ctx context.Context, eventType event.EventType, offset, limit *int, order *string) ([]event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetEventsByType"

	log := s.log.With(
		slog.String("op", op),
	)

	events, err := s.eventRepository.GetEventsByType(ctx, s.db, eventType, offset, limit, order)
	if err != nil {
		log.Error("failed to get events by type",
			slog.Any("eventType", eventType),
			liblogger.Err(err),
		)
		return nil, errs.ErrInternalError.Wrap("failed get events by type")
	}

	return event_mapper.ManyToDTO(events), nil
}

// Get events where type=stage and his childs
func (s *EventService) GetEventsTypeStageAndHisChilds(ctx context.Context, id string) ([]event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetEventsTypeStageAndHisChilds"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uuid", slog.String("id", id), liblogger.Err(err))
		return nil, errs.ErrBadRequest.Wrap("failed parse id")
	}

	// Get all event stage by previousID
	// tx := s.db.Begin()
	events, err := s.eventRepository.GetEventsByPreviousID(ctx, s.db, uid, nil, nil, nil)
	if err != nil {
		// tx.Rollback()
		log.Error("failed to get events type stage by PreviousID",
			slog.Any("id", id),
			liblogger.Err(err),
		)
		return nil, errs.ErrInternalError.Wrap("failed get events by prev id")
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
			childs, err := s.eventRepository.GetEventsByPreviousID(ctx, s.db, id, nil, nil, nil)
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
			liblogger.Err(<-errors),
		)
		return nil, errs.ErrInternalError.Wrap("failed get child by stage id")
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
			slog.String("id", previousId),
			liblogger.Err(err),
		)
		return nil, errs.ErrBadRequest.Wrap("failed parse id")
	}

	events, err := s.eventRepository.GetEventsByPreviousID(ctx, s.db, uid, offset, limit, order)
	if err != nil {
		log.Error("failed to get events by PreviousID",
			liblogger.Err(err),
		)
		return nil, errs.ErrInternalError.Wrap("failed get event by prev id")
	}
	return event_mapper.ManyToDTO(events), nil
}

// Get list events by list id
func (s *EventService) GetEventsByListID(ctx context.Context, ids []string) ([]event_dto.EventDTOResponse, error) {
	const op = "services.event_service.GetEventsByListID"

	log := s.log.With(
		slog.String("op", op),
	)

	uids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			return nil, errs.ErrBadRequest.Wrap("failed parse id")
		}
		uids = append(uids, uid)
	}

	events, err := s.eventRepository.GetEventsByListID(ctx, s.db, uids)
	if err != nil {
		log.Error("failed to get events by listID",
			liblogger.Err(err),
		)
		return nil, errs.ErrInternalError.Wrap("failed get events by list id")
	}
	log.Info("events getted", slog.Any("events", events))

	return event_mapper.ManyToDTO(events), nil
}

// Find olympiad who is parent for stage
func (s *EventService) getOlympiad(ctx context.Context, id uuid.UUID) (event.Event, error) {

	previousEvent, err := s.eventRepository.GetEventByID(ctx, s.db, id)
	if err != nil {
		return event.Event{}, err
	}
	if previousEvent.EventType == event.Olympiad {
		return previousEvent, nil
	} else {
		return s.getOlympiad(ctx, *previousEvent.PreviousEventID)
	}
}

func (s *EventService) checkCorrectEventDTO(ctx context.Context, eventModel *event.Event, isUpdate bool) error {
	const op = "services.event_service.checkCorrectEventDTO"

	// Check date
	if eventModel.StartDate.IsZero() || eventModel.EndDate.IsZero() {
		return fmt.Errorf("%s: %v", op, errors.New("date does not exist"))
	}
	if !eventModel.StartDate.Before(eventModel.EndDate) {
		return fmt.Errorf("%s: %v", op, errors.New("start date should be before end date"))
	}
	// check have previousEvent
	previousEventID := eventModel.PreviousEventID
	if previousEventID != nil {
		previousEvent, err := s.eventRepository.GetEventByID(ctx, s.db, *previousEventID)
		if err != nil {
			return fmt.Errorf("failed to get previousEvent: %v", err)
		}
		if !isUpdate {
			// set correct type
			switch previousEvent.EventType {
			case event.RegionalStage:
				if eventModel.Subject == 0 {
					return fmt.Errorf("%s: %v", op, errors.New("subject does not exist"))
				}
				eventModel.EventType = event.Olympiad
			case event.Olympiad:
				eventModel.EventType = event.Class
				eventModel.Subject = previousEvent.Subject
			case event.Class:
				eventModel.EventType = event.Stage
			case event.Stage:
				// check stage cannot have more than one ViewWorks
				viewWorks, err := s.eventRepository.GetEventsByPreviousID(ctx, s.db, *previousEventID, nil, nil, nil)
				if err != nil {
					return fmt.Errorf("%s: failed to get events when check stage: %v", op, err)
				}
				if len(viewWorks) == 0 {
					eventModel.EventType = event.ViewWorks
				} else if len(viewWorks) == 1 {
					eventModel.EventType = event.Appeal
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
		if eventModel.EventType == event.Stage || eventModel.EventType == event.Olympiad || eventModel.EventType == event.Class {
			if previousEvent.StartDate.After(eventModel.StartDate) || previousEvent.EndDate.Before(eventModel.EndDate) {
				return fmt.Errorf("%s: %v", op, errors.New("incorrect date limits"))
			}
		} else {
			tempPreviousEvent, err := s.getOlympiad(ctx, *previousEventID)
			if err != nil {
				return fmt.Errorf("%s, failed to find first Olympiad: %v", op, err)
			}
			// events type view works and appeal should be after endDate parrent, but before endDate Olympiad
			if previousEvent.EndDate.After(eventModel.StartDate) || tempPreviousEvent.EndDate.Before(eventModel.EndDate) {
				return fmt.Errorf("%s: %v", op, errors.New("incorrect date limits"))
			}
		}

	} else {
		eventModel.EventType = event.RegionalStage
	}
	return nil
}

// Create event
func (s *EventService) CreateEvent(ctx context.Context, eventDTO event_dto.CreateEventDTORequest) (uuid.UUID, error) {
	const op = "services.event_service.CreateEvent"
	log := s.log.With(
		slog.String("op", op),
	)

	eventModel := event_mapper.FromCreateToModel(eventDTO)

	err := s.checkCorrectEventDTO(ctx, &eventModel, false)
	if err != nil {
		log.Error("failed check correct event", liblogger.Err(err))
		return uuid.Nil, errs.ErrBadRequest.Wrap("incorrect event dto")
	}
	// Auto create events for all subject
	if eventModel.EventType == event.RegionalStage {
		id, err := s.createEventsBySubjects(ctx, eventModel)
		if err != nil {
			log.Error("failed create event by subjects", liblogger.Err(err))
			return uuid.Nil, errs.ErrInternalError.Wrap("failed auto create events for all subjects")
		}
		log.Info("events success created (with all subjects)", slog.Any("eventID", id))
		return id, nil
	}
	id, err := s.eventRepository.CreateEvent(ctx, s.db, eventModel)
	if err != nil {
		log.Error("failed to create event", liblogger.Err(err))
		return uuid.Nil, errs.ErrInternalError.Wrap("failed create event")
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
	id, err := s.eventRepository.CreateEvent(ctx, tx, eventModel)
	if err != nil {
		tx.TransactionRollback()
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	for num, subject := range subject.NewSubjectsStorage().GetAllSubject() {
		eventBySubject := event.Event{
			Name:            fmt.Sprintf("Олимпиада по %s", subject),
			PreviousEventID: &id,
			StartDate:       eventModel.StartDate,
			EndDate:         eventModel.EndDate,
			EventType:       event.Olympiad,
			Subject:         num,
		}
		_, err := s.eventRepository.CreateEvent(ctx, tx, eventBySubject)
		if err != nil {
			tx.TransactionRollback()
			return uuid.Nil, fmt.Errorf("%s - create auto event by subjects: %w", op, err)
		}
	}
	tx.TransactionCommit()
	return id, nil
}

func (s *EventService) updateEventDTO(ctx context.Context, updatedEvent event.Event, id uuid.UUID) (event.Event, error) {
	oldEvent, err := s.eventRepository.GetEventByID(ctx, s.db, id)
	if err != nil {
		return event.Event{}, nil
	}
	if updatedEvent.Name != "" {
		oldEvent.Name = updatedEvent.Name
	}
	if !updatedEvent.StartDate.IsZero() {
		oldEvent.StartDate = updatedEvent.StartDate
	}
	if !updatedEvent.EndDate.IsZero() {
		oldEvent.EndDate = updatedEvent.EndDate
	}
	if updatedEvent.ClassNumber != 0 {
		oldEvent.ClassNumber = updatedEvent.ClassNumber
	}
	if updatedEvent.Subject != 0 {
		oldEvent.Subject = updatedEvent.Subject
	}
	if updatedEvent.AdditionalInfo != "" {
		oldEvent.AdditionalInfo = updatedEvent.AdditionalInfo
	}

	return oldEvent, nil
}

// Update event
func (s *EventService) UpdateEvent(ctx context.Context, id string, eventDTO event_dto.UpdateEventDTORequest) error {
	const op = "services.event_service.UpdateEvent"
	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uuid", slog.String("id", id), liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("failed parse id")
	}

	event := event_mapper.FromUpdateToModel(eventDTO, uid)
	event, err = s.updateEventDTO(ctx, event, uid)
	if err != nil {
		return errs.ErrInternalError.Wrap("failed get and update dto model")
	}

	err = s.checkCorrectEventDTO(ctx, &event, true)
	if err != nil {
		log.Error("failed check correct event", liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("incorrect event data")
	}

	err = s.eventRepository.UpdateEvent(ctx, s.db, event)
	if err != nil {
		log.Error("failed update event", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed update event")
	}
	log.Info("event success updated")
	return nil
}

// Delete event
func (s *EventService) DeleteEvent(ctx context.Context, id string) error {
	const op = "services.event_service.DeleteEvent"
	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse id", slog.String("id", id), liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("failed parse id")
	}

	tx, err := s.db.TransactionBegin()
	if err != nil {
		log.Error("failed begin transaction", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed begin transaction to delete event")
	}

	// откат в случае ошибки
	defer func() {
		if err != nil {
			if rbErr := tx.TransactionRollback(); rbErr != nil {
				log.Error("rollback failed", liblogger.Err(rbErr))
			}
		}
	}()

	err = s.eventRepository.DeleteEvent(ctx, tx, uid)
	if err != nil {
		log.Error("failed to delete event", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed delete event")
	}

	for queueName, tableName := range s.kiznaiverServices {
		payload, err := helper.PayloadDeleteConstructor(tableName, uid.String())
		if err != nil {
			log.Error("failed create payload", liblogger.Err(err))
			return errs.ErrInternalError.Wrap("failed create payload")
		}

		outboxModel := outbox.Outbox{
			Payload:   datatypes.JSON(payload),
			QueueName: queueName,
		}
		_, err = s.outboxRepository.Create(ctx, tx, outboxModel)
		if err != nil {
			log.Error("failed create outbox", liblogger.Err(err))
			return errs.ErrInternalError.Wrap("failed create outbox message")
		}
	}

	err = tx.TransactionCommit()
	if err != nil {
		log.Error("failed transaction commit", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed transaction commit")
	}
	log.Info("event success deleted")
	return nil
}
