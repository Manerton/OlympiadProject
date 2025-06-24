package refresh_token

import (
	"os/user"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID
	User      user.User
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}
