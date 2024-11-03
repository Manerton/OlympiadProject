package event_dto

import (
	"main/internal/models/event"
	"time"
)

type EventDTO struct {
	ID              uint
	Name            string
	StartDate       time.Time
	EndDate         time.Time
	EventType       event.EventType
	PreviousEventID *uint
	Subject         string
	AdditionalInfo  string
}
