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
	event_res := event.Event{Model: gorm.Model{ID: id}}
	if err := db.First(&event_res).Error; err != nil {
		return event.Event{}, fmt.Errorf("%s: %w", op, err)
	}
	return event_res, nil
}

// Get slice events by EventType
func (r *EventRepository) GetEventsByType(db *gorm.DB, event_type event.EventType) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByType"
	events_res := []event.Event{}
	if err := db.Find(&events_res, event.Event{EventType: event_type}).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return events_res, nil
}

// Get all events
func (r *EventRepository) GetAllEvents(db *gorm.DB) ([]event.Event, error) {
	const op = "repositories.event_repository.GetAllEvents"
	events_res := []event.Event{}
	if err := db.Find(&events_res).Error; err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return events_res, nil
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
	if err := db.Delete(event.Event{}, id).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
