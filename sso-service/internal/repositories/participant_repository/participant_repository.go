package participant_repository

import (
	"context"
	"fmt"
	"main/internal/models/participant"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type ParticipantRepository struct{}

func (r *ParticipantRepository) GetByUserId(ctx context.Context, orm orm.ORM, userId uuid.UUID) (participant.Participant, error) {
	const op = "repositories.participant_repository.GetByUserId"

	participantRes := participant.Participant{UserId: userId}
	err := orm.First(ctx, participant.Participant{}, nil, &participantRes)
	if err != nil {
		return participant.Participant{}, fmt.Errorf("%s: %w", op, err)
	}
	return participantRes, nil
}

func (r *ParticipantRepository) Create(ctx context.Context, orm orm.ORM, participant participant.Participant) (uuid.UUID, error) {
	const op = "repositories.participant_repository.Create"

	err := orm.Create(ctx, &participant)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return participant.ID, nil
}
func (r *ParticipantRepository) Update(ctx context.Context, orm orm.ORM, participant participant.Participant) error {
	const op = "repositories.participant_repository.Update"

	err := orm.Updates(ctx, &participant)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *ParticipantRepository) Delete(ctx context.Context, orm orm.ORM, id uuid.UUID) error {
	const op = "repositories.participant_repository.Delete"

	err := orm.Delete(ctx, &participant.Participant{ID: id})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
