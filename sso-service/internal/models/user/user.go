package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	Email       string
	Password    string
	FirstName   string
	LastName    string
	Surname     string
	PhoneNumber string
	BirthDate   time.Time
	Sex         string
}
