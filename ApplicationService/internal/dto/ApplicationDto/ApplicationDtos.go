package ApplicationDto

import (
	"time"

	"github.com/google/uuid"
)

// DTO для создания заявки
type CreateApplicationDTO struct {
	UserID  uuid.UUID `json:"userId" binding:"required"`
	EventID uuid.UUID `json:"eventId" binding:"required"`
}

// DTO для обновления статуса заявки
type UpdateApplicationDTO struct {
	Status int    `json:"status"`       // // 2 = одобрено, 3 = отклонено, 1 = не обработано
	Reason int    `gorm:"default:null"` // 1 по результатам предудущего года, 2 по результатам текущего
	Code   string `gorm:"default:null"` // 09_11_25
}

// DTO для возврата заявки
type ApplicationResponseDTO struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"userId"`
	EventID uuid.UUID `json:"eventId"`
	//ВРЕМЕННО ПОДУМАТЬ НАСЧЕТ КЭШИРОВАНИЯ
	// EventName     string    `json:"eventName"`     //ВРЕМЕННО
	// EventLocation string    `json:"eventLocation"` //ВРЕМЕННО
	// EventDate     time.Time `json:"eventDate"`     //ВРЕМЕННО
	//ВРЕМЕННО ПОДУМАТЬ НАСЧЕТ КЭШИРОВАНИЯ
	Status      int       `json:"status"` // // 2 = одобрено, 3 = отклонено, 1 = не обработано
	Reason      int       `json:"reason"` //
	Code        string    `json:"code"`   // 09_11_25
	SubmittedAt time.Time `json:"submittedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
