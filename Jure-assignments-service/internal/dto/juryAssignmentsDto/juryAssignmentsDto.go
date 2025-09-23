package juryAssignmentsDto

import "github.com/google/uuid"

type JuryAssignmentsResponseDTO struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id" validate:"required"`
	EventID uuid.UUID `json:"event_id" validate:"required"`
}

type JuryAssignmentsRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	EventID string `json:"event_id" validate:"required"`
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
	UserID   uuid.UUID   `json:"user_id" validate:"required"`
	EventIDs []uuid.UUID `json:"event_ids" validate:"required"`
}

type CreateOneAssigmentsManyJury struct {
	EventID string   `json:"event_id"`
	UserIDs []string `json:"user_ids"`
}

type DeleteManyAssigmentsJury struct {
	EventID string   `json:"event_id"`
	UserIDs []string `json:"user_ids"`
}
