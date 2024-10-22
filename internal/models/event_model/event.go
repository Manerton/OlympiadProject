package event

import (
	subject "main/internal/models/subject_model"
	"time"

	"gorm.io/gorm"
)

type EventType string

const (
	RegionalStage EventType = "REGIONAL_STAGE"
	Olympiad      EventType = "OLYMPIAD"
	Stage         EventType = "STAGE"
	Appeal        EventType = "APPEAL"
)

type Event struct {
	gorm.Model
	Name            string          `gorm:"type:varchar(128);not null"`
	StartDate       time.Time       `gorm:"not null"`
	EndDate         time.Time       `gorm:"not null"`
	EventType       EventType       `gorm:"type:event_type;not null"`
	PreviousEventID *uint           `gorm:"index"`
	PreviousEvent   *Event          `gorm:"foreignKey:PreviousEventID"`
	SubjectID       uint            `gorm:"index"`
	Subject         subject.Subject `gorm:"foreignKey:SubjectID"`
	AdditionalInfo  string          `gorm:"type:text"`
}
