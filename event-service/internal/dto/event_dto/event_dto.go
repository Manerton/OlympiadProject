package event_dto

import (
	"main/internal/models/event"
	"time"

	"github.com/google/uuid"
)

type EventDTOResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`

	Dates    []string `json:"dates"`
	Profiles []string `json:"profiles"`

	EventType       event.EventType     `json:"event_type"`
	ClassCategory   string              `json:"class_category"`
	PreviousEventID *uuid.UUID          `json:"previous_event_id"`
	Subject         int                 `json:"subject"`
	AdditionalInfo  string              `json:"additional_info"`
	Status          int                 `json:"status"`
	Events          *[]EventDTOResponse `json:"events"`
}

type EventDTO struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `validate:"required"`
	StartDate       time.Time `validate:"required"`
	EndDate         time.Time `validate:"required"`
	EventType       event.EventType
	PreviousEventID *uuid.UUID
	ClassCategory   int `json:"class_category"`
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
	ClassCategory   string          `json:"class_category"`
	Subject         string          `json:"subject,omitempty"`
	AdditionalInfo  string          `json:"additional_info,omitempty"`
}

type CreateEventDTORequest struct {
	Name      string    `json:"name" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`

	Dates    []string `json:"dates"`
	Profiles []string `json:"profiles"`

	PreviousEventID *uuid.UUID `json:"previous_event_id,omitempty"`
	Subject         string     `json:"subject,omitempty"`
	AdditionalInfo  string     `json:"additional_info,omitempty"`
	ClassCategory   string     `json:"class_category,omitempty"` // <- добавить
	EventType       string     `json:"event_type,omitempty"`     // <- добавить
}

type UpdateEventDTORequest struct {
	Name      *string    `json:"name"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	Dates    *[]string `json:"dates"`
	Profiles *[]string `json:"profiles"`

	Subject        *string `json:"subject"`
	ClassCategory  *string `json:"class_category"`
	AdditionalInfo *string `json:"additional_info"`
	Status         *int    `json:"status"`
}
