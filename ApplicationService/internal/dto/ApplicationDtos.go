package ApplicationDto

import "time"

// DTO для создания заявки
type CreateApplicationDTO struct {
	UserID  uint `json:"user_id" binding:"required"`
	EventID uint `json:"event_id" binding:"required"`
}

// DTO для обновления статуса заявки
type UpdateApplicationStatusDTO struct {
	Status *bool `json:"status"` // Поддержка nil для необработанных заявок
}

// DTO для возврата заявки
type ApplicationResponseDTO struct {
	ApplicationID uint `json:"applicationID"`
	UserID        uint `json:"userID"`
	EventID       uint `json:"eventID"`
	//ВРЕМЕННО ПОДУМАТЬ НАСЧЕТ КЭШИРОВАНИЯ
	// EventName     string    `json:"eventName"`     //ВРЕМЕННО
	// EventLocation string    `json:"eventLocation"` //ВРЕМЕННО
	// EventDate     time.Time `json:"eventDate"`     //ВРЕМЕННО
	//ВРЕМЕННО ПОДУМАТЬ НАСЧЕТ КЭШИРОВАНИЯ
	Status      *bool     `json:"status"`
	SubmittedAt time.Time `json:"submittedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
