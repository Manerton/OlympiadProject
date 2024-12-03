package juryAssignmentsDto

type JuryAssignmentsDTO struct {
	ID      uint `json:"ID"`
	JuryID  uint `json:"juryID" validate:"required"`
	EventID uint `json:"eventID" validate:"required"`
}

type OneJuryManyAssignments struct {
	JuryID   uint   `json:"juryID" validate:"required"`
	EventsID []uint `json:"eventsID" validate:"required"`
}
