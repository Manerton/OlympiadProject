package eventrepository

import (
	"fmt"
	"main/internal/dto/event_dto"
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
func (r *EventRepository) GetEventsByType(db *gorm.DB, eventType event.EventType) ([]event.Event, error) {
	const op = "repositories.event_repository.GetEventsByType"
	events_res := []event.Event{}
	if err := db.Find(&events_res, event.Event{EventType: eventType}).Error; err != nil {
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
func (r *EventRepository) CreateEvent(db *gorm.DB, event_dto event_dto.EventDTO) (uint, error) {
	const op = "repositories.event_repository.CreateEvent"
	new_event := event.Event{
		Name:            event_dto.Name,
		StartDate:       event_dto.StartDate,
		EndDate:         event_dto.EndDate,
		PreviousEventID: &event_dto.PreviousEventID,
		SubjectID:       &event_dto.SubjectID,
		AdditionalInfo:  event_dto.AdditionalInfo,
		EventType:       event_dto.EventType,
	}
	if err := db.Create(&new_event).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return new_event.ID, nil
}

// Update event
func (r *EventRepository) UpdateEvent(db *gorm.DB, event_dto event_dto.EventDTO) (uint, error) {
	const op = "repositories.event_repository.UpdateEvent"
	update_event := event.Event{
		Model:           gorm.Model{ID: uint(event_dto.ID)},
		Name:            event_dto.Name,
		StartDate:       event_dto.StartDate,
		EndDate:         event_dto.EndDate,
		PreviousEventID: &event_dto.PreviousEventID,
		SubjectID:       &event_dto.SubjectID,
		AdditionalInfo:  event_dto.AdditionalInfo,
		EventType:       event_dto.EventType,
	}
	if err := db.Updates(&update_event).Error; err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return update_event.ID, nil
}

// Delete event
func (r *EventRepository) DeleteEvent(db *gorm.DB, id uint) error {
	const op = "repositories.event_repository.DeleteEvent"
	if err := db.Delete(event.Event{}, id).Error; err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
