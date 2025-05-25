package event

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	RegionalStage EventType = "REGIONAL_STAGE"
	Olympiad      EventType = "OLYMPIAD"
	Stage         EventType = "STAGE"
	ViewWorks     EventType = "VIEW_WORKS"
	Appeal        EventType = "APPEAL"
)

type Event struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:get_random_uuid();primaryKey"`
	Name            string     `gorm:"type:varchar(128);not null"`
	StartDate       time.Time  `gorm:"not null"`
	EndDate         time.Time  `gorm:"not null"`
	EventType       EventType  `gorm:"type:event_type;not null"`
	PreviousEventID *uuid.UUID `gorm:"index"`
	PreviousEvent   *Event     `gorm:"foreignKey:PreviousEventID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Subject         string     `gorm:"type:varchar(128)"`
	AdditionalInfo  string     `gorm:"type:text"`
}
