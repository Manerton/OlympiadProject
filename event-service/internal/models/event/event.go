package event

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type EventType string

const (
	RegionalStage EventType = "REGIONAL_STAGE"
	Olympiad      EventType = "OLYMPIAD"
	Class         EventType = "CLASS"
	Stage         EventType = "STAGE"
	ViewWorks     EventType = "VIEW_WORKS"
	Appeal        EventType = "APPEAL"
)

type ClassCategoryType string

const (
	Class9     ClassCategoryType = "9"
	Class9_10  ClassCategoryType = "9-10"
	Class9_11  ClassCategoryType = "9-11"
	Class10    ClassCategoryType = "10"
	Class10_11 ClassCategoryType = "10-11"
	Class11    ClassCategoryType = "11"
)

const (
	Register = 1
	Approval = 2
	Finished = 3
)

type Event struct {
	ID            uuid.UUID         `gorm:"type:uuid;default:get_random_uuid();primaryKey"`
	Name          string            `gorm:"type:varchar(128);not null"`
	StartDate     time.Time         `gorm:"not null"`
	EndDate       time.Time         `gorm:"not null"`
	ClassCategory ClassCategoryType `gorm:"type:class_category"`
	// Temp for quic start
	Profiles pq.StringArray `gorm:"type:text[]"`
	Dates    pq.StringArray `gorm:"type:text[]"`

	EventType       EventType  `gorm:"type:event_type;not null"`
	PreviousEventID *uuid.UUID `gorm:"index"`
	PreviousEvent   *Event     `gorm:"foreignKey:PreviousEventID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Subject         int        `gorm:"type:int"`
	Status          int        `gorm:"type:int"`
	AdditionalInfo  string     `gorm:"type:text"`
}
