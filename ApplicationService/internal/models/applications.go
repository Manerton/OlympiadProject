package models

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID   uuid.UUID `gorm:"not null"`
	SchoolID uuid.UUID `gorm:"not null"`
	EventID  uuid.UUID `gorm:"not null"`
	//ВРЕМЕННО ПОДУМАТЬ НАСЧЕТ КЭШИРОВАНИЯ
	// EventName     string    `gorm:"not null"` //ВРЕМЕННО
	// EventLocation string    `gorm:"not null"` //ВРЕМЕННО
	// EventDate     time.Time `gorm:"not null"` //ВРЕМЕННО
	//ВРЕМЕННО ПОДУМАТЬ НАСЧЕТ КЭШИРОВАНИЯ
	Status      int       `gorm:"default:1"`    // 2 = одобрено, 3 = отклонено, 1 = не обработано
	Reason      int       `gorm:"default:null"` // 1 по результатам предудущего года, 2 по результатам текущего
	Code        string    `gorm:"dfeault:null"` // 09_111_25
	SubmittedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
