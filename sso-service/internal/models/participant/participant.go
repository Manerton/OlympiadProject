package participant

import (
	"main/internal/models/user"

	"github.com/google/uuid"
)

type Participant struct {
	ID          uuid.UUID
	UserId      uuid.UUID
	User        user.User
	OVZ         int
	SchoolName  string
	City        string
	Reason      string
	Citizenship string
	ClassNumber int
}
