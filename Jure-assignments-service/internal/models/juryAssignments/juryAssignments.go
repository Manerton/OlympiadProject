package juryassignments

type JuryAssignments struct {
	ID      uint `gorm:"primaryKey"`
	JuryID  uint `gorm:"not null;uniqueIndex:idx_jury_event"`
	EventID uint `gorm:"not null;uniqueIndex:idx_jury_event"`
}
