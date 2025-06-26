package jury_assignments

import "github.com/google/uuid"

type JuryAssignments struct {
	ID      uuid.UUID `gorm:"primaryKey"`
	JuryID  uuid.UUID `gorm:"not null;uniqueIndex:idx_jury_event"`
	EventID uuid.UUID `gorm:"not null;uniqueIndex:idx_jury_event"`
}
