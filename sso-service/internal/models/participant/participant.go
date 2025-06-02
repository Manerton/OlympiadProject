package participant

import (
	"main/internal/models/user"

	"github.com/google/uuid"
)

type Participant struct {
	ID          uuid.UUID `gorm:"type:uuid;default:get_random_uuid();primaryKey"`
	UserId      uuid.UUID
	User        user.User
	OVZ         int
	SchoolName  string
	City        string
	Reason      string
	Citizenship string
	ClassNumber int
}
