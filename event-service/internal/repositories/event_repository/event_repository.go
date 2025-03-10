package event_repository

import (
	"context"
	"fmt"
	"main/internal/models/event"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type EventRepository struct{}

func (r *EventRepository) GetEventByFilterAndFields(ctx context.Context, orm orm.ORM, filter event.Event, fields *[]string) (event.Event, error) {
	const op = "repositories.event_repository.GetEventByFilter"
	if err := orm.First(ctx, event.Event{}, fields, &filter); err != nil {
		return event.Event{}, fmt.Errorf("%s: %w", op, err)
	}
	// query := db.Model(event.Event{})
	// // query := db.Debug().Model(event.Event{})
	// if fields != nil {
	// 	query.Select(*fields)
	// }

	// if err := query.First(&filter).Error; err != nil {
	// 	return event.Event{}, fmt.Errorf("%s: %w", op, err)
	// }
	return filter, nil
}

func (r *EventRepository) GetEventsByFilterAndFields(ctx context.Context, orm orm.ORM, filter event.Event, fields *[]string, offset, limit *int, order *string) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByFilterAndFields"
	eventRes := []event.Event{}

	if err := orm.Find(ctx, event.Event{}, fields, offset, limit, order, &eventRes, filter); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// query := db.Model(event.Event{})
	// if fields != nil {
	// 	query.Select(*fields)
	// }
	// if offset != nil {
	// 	query = query.Offset(*offset)
	// }
	// if limit != nil {
	// 	query = query.Limit(*limit)
	// }
	// eventRes := []event.Event{}

	// if err := query.Find(&eventRes, filter).Error; err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }
	return eventRes, nil
}

// Get event by ID
func (r *EventRepository) GetEventByID(ctx context.Context, orm orm.ORM, id uuid.UUID) (event.Event, error) {
	const op = "repositories.event_repository.GetEventByID"

	eventRes := event.Event{ID: id}
	if err := orm.First(ctx, event.Event{}, nil, &eventRes); err != nil {
		return event.Event{}, fmt.Errorf("%s: %w", op, err)
	}
	// if err := db.First(&eventRes).Error; err != nil {
	// 	return event.Event{}, fmt.Errorf("%s: %w", op, err)
	// }
	return eventRes, nil
}

// Get list events by list id
func (r *EventRepository) GetEventsByListID(ctx context.Context, orm orm.ORM, ids []uuid.UUID) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByListID"

	eventRes := []event.Event{}
	if err := orm.Find(ctx, event.Event{}, nil, nil, nil, nil, &eventRes, ids); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// if err := db.Find(&eventRes, ids).Error; err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }
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

	// query := db.Model(event.Event{})
	// if offset != nil {
	// 	query = query.Offset(*offset)
	// }
	// if limit != nil {
	// 	query = query.Limit(*limit)
	// }

	// if err := query.Find(&eventsRes, event.Event{EventType: eventType}).Error; err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }
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
	// query := db.Model(event.Event{})
	// if offset != nil {
	// 	query = query.Offset(*offset)
	// }
	// if limit != nil {
	// 	query = query.Limit(*limit)
	// }

	// if err := query.Find(&events, event.Event{PreviousEventID: &previousID}).Error; err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }
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

	// query := db.Model(event.Event{})
	// if offset != nil {
	// 	query = query.Offset(*offset)
	// }
	// if limit != nil {
	// 	query = query.Limit(*limit)
	// }

	// if err := query.Find(&eventsRes).Error; err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }
	return eventsRes, nil
}

func (r *EventRepository) GetCountEventsByType(ctx context.Context, orm orm.ORM, eventType event.EventType) (int64, error) {
	const op = "repositories.event_repository.GetCountEventsByType"
	var resultCount int64 = 0
	err := orm.Count(ctx, &event.Event{}, &resultCount, "event_type = ?", eventType)

	// err := db.Model(&event.Event{}).Where("event_type = ?", eventType).Count(&resultCount).Error
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return resultCount, nil
}

func (r *EventRepository) GetCountEventsByPreviousID(ctx context.Context, orm orm.ORM, previousID uuid.UUID) (int64, error) {
	const op = "repositories.event_repository.GetCountEventsByPreviousID"
	var resultCount int64 = 0
	err := orm.Count(ctx, &event.Event{}, &resultCount, "previous_event_id = ?", previousID)

	// err := db.Model(&event.Event{}).Where("previous_event_id = ?", previousID).Count(&resultCount).Error
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
	// if err := db.Create(&event).Error; err != nil {
	// 	return 0, fmt.Errorf("%s: %w", op, err)
	// }
	return event.ID, nil
}

// Update event
func (r *EventRepository) UpdateEvent(ctx context.Context, orm orm.ORM, event event.Event) (uuid.UUID, error) {
	const op = "repositories.event_repository.UpdateEvent"
	if err := orm.Updates(ctx, &event); err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	// if err := db.Updates(&event).Error; err != nil {
	// 	return 0, fmt.Errorf("%s: %w", op, err)
	// }
	return event.ID, nil
}

// Delete event
func (r *EventRepository) DeleteEvent(ctx context.Context, orm orm.ORM, id uuid.UUID) error {
	const op = "repositories.event_repository.DeleteEvent"
	if err := orm.Delete(ctx, &event.Event{ID: id}); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	// if err := db.Delete(&event.Event{ID: id}).Error; err != nil {
	// 	return fmt.Errorf("%s: %w", op, err)
	// }
	return nil
}
