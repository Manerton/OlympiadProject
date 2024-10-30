package models

import "time"

type Application struct {
	ApplicationID uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null"`
	EventID       uint      `gorm:"not null"`
	Status        *bool     `gorm:"default:null"` // true = одобрено, false = отклонено, nil = не обработано
	SubmittedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}
