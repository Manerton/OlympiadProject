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
	GetAll(ctx context.Context, orm orm.ORM, offset, limit *int) ([]participant.Participant, error)
	GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (participant.Participant, error)
	GetByUserId(ctx context.Context, orm orm.ORM, id uuid.UUID) (participant.Participant, error)

	GetCount(ctx context.Context, orm orm.ORM) (int64, error)

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

func (s *ParticipantService) GetCount(ctx context.Context) (int64, error) {
	const op = "services.participant_service.GetCount"
	const errMsg = "failed count participants"

	log := s.log.With(
		slog.String("op", op),
	)

	participantCount, err := s.participantRepository.GetCount(ctx, s.db)
	if err != nil {
		log.Error("failed participant count users", liblogger.Err(err))
		return 0, fmt.Errorf("%s", errMsg)
	}

	return participantCount, nil
}

func (s *ParticipantService) GetById(ctx context.Context, id string) (participant_dto.ParticipantResponseDTO, error) {
	const op = "services.participant_service.GetById"
	const errMsg = "failed get participant"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uid", liblogger.Err(err))
		return participant_dto.ParticipantResponseDTO{}, fmt.Errorf("%s", errMsg)
	}

	participant, err := s.participantRepository.GetById(ctx, s.db, uid)
	if err != nil {
		log.Error("failed get participant", liblogger.Err(err))
		return participant_dto.ParticipantResponseDTO{}, fmt.Errorf("%s", errMsg)
	}

	return participant_mapper.ToDTO(participant), nil
}

func (s *ParticipantService) GetByUserId(ctx context.Context, id string) (participant_dto.ParticipantResponseDTO, error) {
	const op = "services.participant_service.GetByUserID"
	const errMsg = "failed get participant"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uid", liblogger.Err(err))
		return participant_dto.ParticipantResponseDTO{}, fmt.Errorf("%s", errMsg)
	}

	participantModel, err := s.participantRepository.GetByUserId(ctx, s.db, uid)
	if err != nil {
		log.Error("failed get participant by user id", liblogger.Err(err))
		return participant_dto.ParticipantResponseDTO{}, fmt.Errorf("%s", errMsg)
	}

	return participant_mapper.ToDTO(participantModel), nil

}

func (s *ParticipantService) GetAll(ctx context.Context, page, limit *int) ([]participant_dto.ParticipantResponseDTO, error) {
	const op = "services.participant_service.GetAll"
	const errMsg = "failed get participants"

	log := s.log.With(
		slog.String("op", op),
	)

	offset := new(int)
	if page != nil && limit != nil {
		*offset = (*page - 1) * (*limit)
	}

	participants, err := s.participantRepository.GetAll(ctx, s.db, offset, limit)
	if err != nil {
		log.Error("failed get participants", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	participantsDTO := make([]participant_dto.ParticipantResponseDTO, 0, len(participants))
	for _, participantModel := range participants {
		participantsDTO = append(participantsDTO, participant_mapper.ToDTO(participantModel))
	}
	return participantsDTO, nil

}

func (s *ParticipantService) Update(ctx context.Context, id string, participantDTO participant_dto.UpdateParticipantRequestDTO) error {
	const op = "services.participant_service.Update"
	const errMsg = "failed update participant"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", liblogger.Err(err))
		return fmt.Errorf("failed parse id")
	}

	participantModel, err := participant_mapper.FromUpdateToModel(participantDTO, uid)
	if err != nil {
		log.Error("failed convert update dto to model", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	err = s.participantRepository.Update(ctx, s.db, participantModel)
	if err != nil {
		log.Error("failed update participant", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}
	return nil
}
