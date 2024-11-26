package event_dto

import (
	"main/internal/models/event"
	"time"
)

type EventDTO struct {
	ID              uint
	Name            string    `validate:"required"`
	StartDate       time.Time `validate:"required"`
	EndDate         time.Time `validate:"required"`
	EventType       event.EventType
	PreviousEventID *uint
	Subject         string
	AdditionalInfo  string
	Events          []EventDTO
}

type DetailsEvent struct {
	ID              uint            `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	StartDate       *time.Time      `json:"start_date,omitempty"` // указатель для учета нулевых значений
	EndDate         *time.Time      `json:"end_date,omitempty"`   // указатель для учета нулевых значений
	EventType       event.EventType `json:"event_type,omitempty"`
	PreviousEventID *uint           `json:"previous_event_id,omitempty"`
	Subject         string          `json:"subject,omitempty"`
	AdditionalInfo  string          `json:"additional_info,omitempty"`
}
