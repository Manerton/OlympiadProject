package user

import (
	"time"

	"github.com/google/uuid"
)

type RoleType int

const (
	ADMIN       = 0
	PARTICIPANT = 1
	JUDGE       = 2
	ORGANIZER   = 3
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string
	PasswordHash string
	Firstname    string
	Surname      string
	Patronymic   string
	PhoneNumber  string
	BirthDate    time.Time `gorm:"column:birhdate"`
	Gender       string
	Role         RoleType
	Activated    bool `gorm:"default:false"`
	CreatedAt    time.Time
}
