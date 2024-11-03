package event_repository

import (
	"fmt"
	"main/internal/models/event"

	"gorm.io/gorm"
)

type EventRepository struct{}

// Get event by ID
func (r *EventRepository) GetEventByID(db *gorm.DB, id uint) (event.Event, error) {
	const op = "repositories.event_repository.GetEventById"

	eventRes := event.Event{ID: id}
	if err := db.First(&eventRes).Error; err != nil {
		return event.Event{}, fmt.Errorf("%s: %w", op, err)
	}
	return eventRes, nil
}

// Get slice events by EventType
func (r *EventRepository) GetEventsByType(db *gorm.DB, eventType event.EventType, offset, limit *int) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByType"
	eventsRes := []event.Event{}

	query := db.Model(event.Event{})
	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.Find(&eventsRes, event.Event{EventType: eventType}).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eventsRes, nil
}

// Get all events by PreviousID
func (r *EventRepository) GetEventsByPreviousID(db *gorm.DB, previousID uint, offset, limit *int) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByPreviousID"
	events := []event.Event{}

	query := db.Model(event.Event{})
	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.Find(&events, event.Event{PreviousEventID: &previousID}).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return events, nil
}

// Get all events with offset and limit
//
// Offset, limit can be nil
func (r *EventRepository) GetAllEvents(db *gorm.DB, offset, limit *int) ([]event.Event, error) {
	const op = "repositories.event_repository.GetAllEvents"
	eventsRes := []event.Event{}

	query := db.Model(event.Event{})
	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.Find(&eventsRes).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return eventsRes, nil
}

// Add new event in DB
func (r *EventRepository) CreateEvent(db *gorm.DB, event event.Event) (uint, error) {
	const op = "repositories.event_repository.CreateEvent"
	if err := db.Create(&event).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return event.ID, nil
}

// Update event
func (r *EventRepository) UpdateEvent(db *gorm.DB, event event.Event) (uint, error) {
	const op = "repositories.event_repository.UpdateEvent"
	if err := db.Updates(&event).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return event.ID, nil
}

// Delete event
func (r *EventRepository) DeleteEvent(db *gorm.DB, id uint) error {
	const op = "repositories.event_repository.DeleteEvent"
	if err := db.Delete(&event.Event{ID: id}).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
