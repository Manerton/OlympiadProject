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
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Surname      string
	PhoneNumber  string
	BirthDate    time.Time
	Sex          string
	Role         RoleType
}
