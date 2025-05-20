package user

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	Surname     string
	PhoneNumber string
}
