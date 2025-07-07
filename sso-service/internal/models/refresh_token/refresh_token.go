package refresh_token

import (
	"main/internal/models/user"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	User      user.User `gorm:"foreignKey:UserID;references:ID"`
	TokenHash string
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time
}
