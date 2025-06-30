package juryAssignmentsDto

import "github.com/google/uuid"

type JuryAssignmentsDTO struct {
	ID      uuid.UUID `json:"id"`
	JuryID  uuid.UUID `json:"jury_id" validate:"required"`
	EventID uuid.UUID `json:"evemt_id" validate:"required"`
}

type OneJuryManyAssignments struct {
	JuryID   uuid.UUID   `json:"jury_id" validate:"required"`
	EventIDs []uuid.UUID `json:"event_ids" validate:"required"`
}

type OneAssigmentsManyJury struct {
	EventID uuid.UUID   `json:"event_id"`
	JuryIDs []uuid.UUID `json:"jury_ids"`
}
