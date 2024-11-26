package models

import "time"

type Application struct {
	ApplicationID uint `gorm:"primaryKey"`
	UserID        uint `gorm:"not null"`
	EventID       uint `gorm:"not null"`
	//ВРЕМЕННО ПОДУМАТЬ НАСЧЕТ КЭШИРОВАНИЯ
	EventName     string    `gorm:"not null"` //ВРЕМЕННО
	EventLocation string    `gorm:"not null"` //ВРЕМЕННО
	EventDate     time.Time `gorm:"not null"` //ВРЕМЕННО
	//ВРЕМЕННО ПОДУМАТЬ НАСЧЕТ КЭШИРОВАНИЯ
	Status      *bool     `gorm:"default:null"` // true = одобрено, false = отклонено, nil = не обработано
	SubmittedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
