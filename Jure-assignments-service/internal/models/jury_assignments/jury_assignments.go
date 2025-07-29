package jury_assignments

import "github.com/google/uuid"

type JuryAssignments struct {
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID  uuid.UUID `gorm:"not null;uniqueIndex:idx_jury_event"`
	EventID uuid.UUID `gorm:"not null;uniqueIndex:idx_jury_event"`
}
