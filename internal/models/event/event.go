package event

import (
	"time"
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
	ID              uint      `gorm:"primaryKey"`
	Name            string    `gorm:"type:varchar(128);not null"`
	StartDate       time.Time `gorm:"not null"`
	EndDate         time.Time `gorm:"not null"`
	EventType       EventType `gorm:"type:event_type;not null"`
	PreviousEventID *uint     `gorm:"index"`
	PreviousEvent   *Event    `gorm:"foreignKey:PreviousEventID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Subject         string    `gorm:"type:varchar(128)"`
	// SubjectID       *uint            `gorm:"index,constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	// Subject         *subject.Subject `gorm:"foreignKey:SubjectID"`
	AdditionalInfo string `gorm:"type:text"`
}
