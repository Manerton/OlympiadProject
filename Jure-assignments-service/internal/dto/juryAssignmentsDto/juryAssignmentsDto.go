package juryAssignmentsDto

import "github.com/google/uuid"

type JuryAssignmentsResponseDTO struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id" validate:"required"`
	EventID uuid.UUID `json:"event_id" validate:"required"`
}

type CreateJuryAssignmentsDTO struct {
	UserID  string `json:"user_id" validate:"required"`
	EventID string `json:"event_id" validate:"required"`
}

type UpdateJuryAssignmentsDTO struct {
	UserID  *string `json:"user_id" validate:"required"`
	EventID *string `json:"event_id" validate:"required"`
}

type OneJuryManyAssignments struct {
	JuryID   uuid.UUID   `json:"jury_id" validate:"required"`
	EventIDs []uuid.UUID `json:"event_ids" validate:"required"`
}

type CreateOneAssigmentsManyJury struct {
	EventID string   `json:"event_id"`
	JuryIDs []string `json:"jury_ids"`
}
