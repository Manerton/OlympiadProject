package participant

import (
	"main/internal/models/school"
	"main/internal/models/user"

	"github.com/google/uuid"
)

type Participant struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserId      uuid.UUID
	User        user.User
	OVZ         int
	SchoolId    uuid.UUID
	School      school.School
	City        string
	Reason      string
	Citizenship string
	ClassNumber int
}
