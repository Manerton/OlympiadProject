package participant_service

import (
	"context"
	"fmt"
	"log/slog"
	participant_dto "main/internal/dto/participant"
	"main/internal/lib/liblogger"
	"main/internal/lib/mapper/participant_mapper"
	"main/internal/models/participant"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type ParticipantRepository interface {
	Update(ctx context.Context, orm orm.ORM, participant participant.Participant) error
}

type ParticipantService struct {
	log                   *slog.Logger
	db                    orm.ORM
	participantRepository ParticipantRepository
}

func New(log *slog.Logger, orm orm.ORM, participantRepository ParticipantRepository) *ParticipantService {
	return &ParticipantService{
		log:                   log,
		db:                    orm,
		participantRepository: participantRepository,
	}
}

func (s *ParticipantService) Update(ctx context.Context, id string, participantDTO participant_dto.UpdateParticipantRequestDTO) error {
	const op = "services.participant_service.Update"
	const errMsg = "failed update participant"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id: %w", liblogger.Err(err))
		return fmt.Errorf("failed parse id")
	}

	participantModel := participant_mapper.FromUpdateToModel(participantDTO, uid)
	err = s.participantRepository.Update(ctx, s.db, participantModel)
	if err != nil {
		log.Error("failed update participant: %w", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}
	return nil
}
