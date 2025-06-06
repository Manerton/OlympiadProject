package user

import (
	"time"

	"github.com/google/uuid"
)

type RoleType string

const (
	JUDGE       = "judge"
	PARTICIPANT = "participant"
	ADMIN       = "admin"
	ORGANIZER   = "organizer"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Surname      string
	PhoneNumber  string
	BirthDate    time.Time
	Gender       string
	Role         RoleType
	Activated    bool `gorm:"default:false"`
	CreatedAt    time.Time
}
