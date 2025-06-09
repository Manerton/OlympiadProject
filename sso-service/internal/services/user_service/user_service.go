package user_service

import (
	"context"
	"fmt"
	"log/slog"
	user_dto "main/internal/dto/user"
	"main/internal/lib/liblogger"
	paricipant_mapper "main/internal/lib/mapper/participant_mapper"
	"main/internal/lib/mapper/user_mapper"
	"main/internal/models/participant"
	"main/internal/models/user"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (user.User, error)
	GetByListId(ctx context.Context, orm orm.ORM, ids []*uuid.UUID) ([]user.User, error)
	Update(ctx context.Context, orm orm.ORM, userModel *user.User) error
}

type ParticipantRepository interface {
	GetByUserId(ctx context.Context, orm orm.ORM, id uuid.UUID) (participant.Participant, error)
}

type UserService struct {
	db                    orm.ORM
	log                   *slog.Logger
	userRepository        UserRepository
	participantRepository ParticipantRepository
}

func New(log *slog.Logger, orm orm.ORM, userRepository UserRepository, participantRepository ParticipantRepository) *UserService {
	return &UserService{
		participantRepository: participantRepository,
		userRepository:        userRepository,
		db:                    orm,
		log:                   log,
	}
}

func (s *UserService) GetById(ctx context.Context, id string) (user_dto.UserResponseDTO, error) {
	const op = "services.user_services.GetById"
	const errMsg = "failed to find user"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse id from string to uuid", liblogger.Err(err))
		return user_dto.UserResponseDTO{}, fmt.Errorf("%s", errMsg)
	}

	userResult, err := s.userRepository.GetById(ctx, s.db, uid)
	if err != nil {
		log.Error("failed to get user", liblogger.Err(err))
		return user_dto.UserResponseDTO{}, fmt.Errorf("%s", errMsg)
	}

	return user_mapper.ToDTO(userResult), nil
}

func (s *UserService) GetParticipantUserById(ctx context.Context, id string) (user_dto.ParticipantUserResponseDTO, error) {
	const op = "services.user_services.GetParticipantUserById"
	const errMsg = "failed to find all info user"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse id from string to uuid", liblogger.Err(err))
		return user_dto.ParticipantUserResponseDTO{}, fmt.Errorf("%s", errMsg)
	}

	userResult, err := s.userRepository.GetById(ctx, s.db, uid)
	if err != nil {
		log.Error("failed to get user", liblogger.Err(err))
		return user_dto.ParticipantUserResponseDTO{}, fmt.Errorf("%s", errMsg)
	}

	participantResult, err := s.participantRepository.GetByUserId(ctx, s.db, userResult.ID)
	if err != nil {
		log.Error("failed to get participant", liblogger.Err(err))
		return user_dto.ParticipantUserResponseDTO{}, fmt.Errorf("%s", errMsg)
	}

	return user_dto.ParticipantUserResponseDTO{
		UserResponseDTO:        user_mapper.ToDTO(userResult),
		ParticipantResponseDTO: paricipant_mapper.ToDTO(participantResult),
	}, nil

}

func (s *UserService) GetByListId(ctx context.Context, ids []*string) ([]user_dto.UserResponseDTO, error) {
	const op = "services.user_handler.GetByListId"
	const errMsg = "failed to find users by list id"

	log := s.log.With(
		slog.String("op", op),
	)

	uids := make([]*uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uid, err := uuid.Parse(*id)
		if err != nil {
			log.Error("failed to parse id:", liblogger.Err(err))
			return nil, fmt.Errorf("%s", errMsg)
		}
		uids = append(uids, &uid)
	}

	usersResult, err := s.userRepository.GetByListId(ctx, s.db, uids)
	if err != nil {
		log.Error("failed to get users by list id", slog.Any("ids", ids), liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	usersDTO := make([]user_dto.UserResponseDTO, 0, len(usersResult))
	for _, userResult := range usersResult {
		usersDTO = append(usersDTO, user_mapper.ToDTO(userResult))
	}
	return usersDTO, nil
}

func (s *UserService) Update(ctx context.Context, userDto *user_dto.UpdateUserRequestDTO) error {
	const op = "services.user_sevice.Update"
	const errMsg = "failed update user"

	log := s.log.With(
		slog.String("op", op),
	)

	userModel := user_mapper.FromUpdateToModel(userDto)
	err := s.userRepository.Update(ctx, s.db, &userModel)
	if err != nil {
		log.Error("failed update user", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}
