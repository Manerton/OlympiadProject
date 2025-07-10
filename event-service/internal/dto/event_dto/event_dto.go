package event_dto

import (
	"main/internal/models/event"
	"time"

	"github.com/google/uuid"
)

type EventDTOResponse struct {
	ID              uuid.UUID
	Name            string    `validate:"required"`
	StartDate       time.Time `validate:"required"`
	EndDate         time.Time `validate:"required"`
	EventType       event.EventType
	PreviousEventID *uuid.UUID
	Subject         string
	AdditionalInfo  string
	Events          *[]EventDTOResponse
}

type EventDTO struct {
	ID              uuid.UUID
	Name            string    `validate:"required"`
	StartDate       time.Time `validate:"required"`
	EndDate         time.Time `validate:"required"`
	EventType       event.EventType
	PreviousEventID *uuid.UUID
	Subject         string
	AdditionalInfo  string
	Events          *[]EventDTO
}

type DetailsEvent struct {
	ID              uuid.UUID       `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	StartDate       *time.Time      `json:"start_date,omitempty"` // указатель для учета нулевых значений
	EndDate         *time.Time      `json:"end_date,omitempty"`   // указатель для учета нулевых значений
	EventType       event.EventType `json:"event_type,omitempty"`
	PreviousEventID *uuid.UUID      `json:"previous_event_id,omitempty"`
	Subject         string          `json:"subject,omitempty"`
	AdditionalInfo  string          `json:"additional_info,omitempty"`
}

type CreateEventDTORequest struct {
	Name            string     `json:"name" validate:"required"`
	StartDate       time.Time  `json:"start_date" validate:"required"`
	EndDate         time.Time  `json:"end_date" validate:"required"`
	PreviousEventID *uuid.UUID `json:"previous_event_id,omitempty"`
	Subject         string     `json:"subject,omitempty"`
	AdditionalInfo  string     `json:"additional_info,omitempty"`
	ClassNumber     int        `json:"class_number,omitempty"` // <- добавить
	EventType       string     `json:"event_type,omitempty"`   // <- добавить
}

type UpdateEventDTORequest struct {
	Name           *string
	StartDate      *time.Time
	EndDate        *time.Time
	Subject        *string
	AdditionalInfo *string
}
