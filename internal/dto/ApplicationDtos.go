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
	ApplicationID uint      `json:"application_id"`
	UserID        uint      `json:"user_id"`
	EventID       uint      `json:"event_id"`
	Status        *bool     `json:"status"`
	SubmittedAt   time.Time `json:"submitted_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
