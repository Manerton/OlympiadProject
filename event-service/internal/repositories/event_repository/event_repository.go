package event_repository

import (
	"context"
	"fmt"
	"main/internal/models/event"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type EventRepository struct{}

// Get first event by parametrs filter and fields
//
// filter - Event fields used for the search
// Exemple: Event{name: "Test"}
//
// fileds - fields that you want to get, if all just use nil
func (r *EventRepository) GetEventByFilterAndFields(ctx context.Context, orm orm.ORM, filter event.Event, fields *[]string) (event.Event, error) {
	const op = "repositories.event_repository.GetEventByFilter"

	var eventRes event.Event
	if err := orm.First(ctx, event.Event{}, fields, &eventRes, filter); err != nil {
		return event.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	return eventRes, nil
}

func (r *EventRepository) GetEventsByFilterAndFields(ctx context.Context, orm orm.ORM, filter event.Event, fields *[]string, offset, limit *int, order *string) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByFilterAndFields"

	var eventRes []event.Event
	if err := orm.Find(ctx, event.Event{}, fields, offset, limit, order, &eventRes, filter); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return eventRes, nil
}

// Get event by ID
func (r *EventRepository) GetEventByID(ctx context.Context, orm orm.ORM, id uuid.UUID) (event.Event, error) {
	const op = "repositories.event_repository.GetEventByID"

	eventRes := event.Event{ID: id}
	if err := orm.First(ctx, event.Event{}, nil, &eventRes); err != nil {
		return event.Event{}, fmt.Errorf("%s: %w", op, err)
	}

	return eventRes, nil
}

// Get list events by list id
func (r *EventRepository) GetEventsByListID(ctx context.Context, orm orm.ORM, ids []uuid.UUID) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByListID"

	if len(ids) == 0 {
		return []event.Event{}, nil
	}

	eventRes := []event.Event{}

	// orderByName := "name"
	if err := orm.Find(ctx, event.Event{}, nil, nil, nil, nil, &eventRes, ids); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return eventRes, nil
}

// Get list events by EventType
// Offset, limit can be nil
func (r *EventRepository) GetEventsByType(ctx context.Context, orm orm.ORM, eventType event.EventType, offset, limit *int, order *string) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByType"
	eventsRes := []event.Event{}

	if err := orm.Find(ctx, event.Event{}, nil, offset, limit, order, &eventsRes, event.Event{EventType: eventType}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return eventsRes, nil
}

// Get all events by PreviousID
// Offset, limit can be nil
func (r *EventRepository) GetEventsByPreviousID(ctx context.Context, orm orm.ORM, previousID uuid.UUID, offset, limit *int, order *string) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByPreviousID"
	eventsRes := []event.Event{}

	if err := orm.Find(ctx, event.Event{}, nil, offset, limit, order, &eventsRes, event.Event{PreviousEventID: &previousID}); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return eventsRes, nil
}

// Get all events with offset and limit
// Offset, limit can be nil
func (r *EventRepository) GetAllEvents(ctx context.Context, orm orm.ORM, offset, limit *int) ([]event.Event, error) {
	const op = "repositories.event_repository.GetAllEvents"
	eventsRes := []event.Event{}

	if err := orm.Find(ctx, event.Event{}, nil, offset, limit, nil, &eventsRes); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return eventsRes, nil
}

func (r *EventRepository) GetAvailableEventsByClass(ctx context.Context, orm orm.ORM, id uuid.UUID, class int) ([]event.Event, error) {
	const op = "repositories.event_repository.GetAvailableEventsByClass"
	result := []event.Event{}
	err := orm.FindWithAdvancedQuery(ctx, event.Event{}, event.Event{PreviousEventID: &id}, &result, "class_number >= ?", class)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}

func (r *EventRepository) GetCountEventsByType(ctx context.Context, orm orm.ORM, eventType event.EventType) (int64, error) {
	const op = "repositories.event_repository.GetCountEventsByType"
	var resultCount int64 = 0
	err := orm.Count(ctx, &event.Event{}, &resultCount, "event_type = ?", eventType)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return resultCount, nil
}

func (r *EventRepository) GetCountEventsByPreviousID(ctx context.Context, orm orm.ORM, previousID uuid.UUID) (int64, error) {
	const op = "repositories.event_repository.GetCountEventsByPreviousID"
	var resultCount int64 = 0
	err := orm.Count(ctx, &event.Event{}, &resultCount, "previous_event_id = ?", previousID)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return resultCount, nil
}

// Add new event in DB
func (r *EventRepository) CreateEvent(ctx context.Context, orm orm.ORM, event event.Event) (uuid.UUID, error) {
	const op = "repositories.event_repository.CreateEvent"
	if err := orm.Create(ctx, &event); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return event.ID, nil
}

// Update event
func (r *EventRepository) UpdateEvent(ctx context.Context, orm orm.ORM, event event.Event) error {
	const op = "repositories.event_repository.UpdateEvent"
	if err := orm.Updates(ctx, &event); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Delete event
func (r *EventRepository) DeleteEvent(ctx context.Context, orm orm.ORM, id uuid.UUID) error {
	const op = "repositories.event_repository.DeleteEvent"
	if err := orm.Delete(ctx, &event.Event{ID: id}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
