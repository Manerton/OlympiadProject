package juryAssignmentsDto

import "github.com/google/uuid"

type JuryAssignmentsDTO struct {
	ID      uuid.UUID `json:"ID"`
	JuryID  uuid.UUID `json:"juryID" validate:"required"`
	EventID uuid.UUID `json:"eventID" validate:"required"`
}

type OneJuryManyAssignments struct {
	JuryID   uuid.UUID   `json:"juryID" validate:"required"`
	EventsID []uuid.UUID `json:"eventsID" validate:"required"`
}
