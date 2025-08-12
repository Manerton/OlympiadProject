package user_service

import (
	"context"
	"log/slog"
	user_dto "main/internal/dto/user"
	"main/internal/lib/crypt"
	"main/internal/lib/errs"
	"main/internal/lib/liblogger"
	paricipant_mapper "main/internal/lib/mapper/participant_mapper"
	"main/internal/lib/mapper/user_mapper"
	"main/internal/models/participant"
	"main/internal/models/user"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx context.Context, orm orm.ORM, offset *int, limit *int) ([]user.User, error)
	GetByFilter(ctx context.Context, orm orm.ORM, user user.User) (user.User, error)
	GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (user.User, error)
	GetByListId(ctx context.Context, orm orm.ORM, ids []uuid.UUID) ([]user.User, error)

	GetCount(ctx context.Context, orm orm.ORM) (int64, error)

	Update(ctx context.Context, orm orm.ORM, userModel user.User) error
	Delete(ctx context.Context, orm orm.ORM, id uuid.UUID) error
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

func (s *UserService) GetCount(ctx context.Context) (int64, error) {
	const op = "services.user_services.GetCount"

	log := s.log.With(
		slog.String("op", op),
	)

	userCount, err := s.userRepository.GetCount(ctx, s.db)
	if err != nil {
		log.Error("failed get count users", liblogger.Err(err))
		return 0, errs.ErrInternalError.Wrap("failed get count users")
	}

	return userCount, nil
}

func (s *UserService) GetAll(ctx context.Context, page *int, limit *int) ([]user_dto.UserResponseDTO, error) {
	const op = "services.user_services.GetAll"

	log := s.log.With(
		slog.String("op", op),
	)

	offset := new(int)
	if page != nil && limit != nil {
		*offset = (*page - 1) * (*limit)
	}

	usersResult, err := s.userRepository.GetAll(ctx, s.db, offset, limit)
	if err != nil {
		log.Error("failed find users", liblogger.Err(err))
		return nil, errs.ErrInternalError.Wrap("failed get all users")
	}

	usersDTO := make([]user_dto.UserResponseDTO, 0, len(usersResult))
	for _, userResult := range usersResult {
		usersDTO = append(usersDTO, user_mapper.ToDTO(userResult))
	}

	return usersDTO, nil
}

func (s *UserService) GetByFilter(ctx context.Context, userDTO user_dto.SearchAttributesUserDTO) (user_dto.UserResponseDTO, error) {
	const op = "services.user_services.GetByFilter"

	log := s.log.With(
		slog.String("op", op),
	)

	userModel := user_mapper.FromSearchToModel(userDTO)
	userResult, err := s.userRepository.GetByFilter(ctx, s.db, userModel)
	if s.db.IsNotFound(err) {
		log.Warn("user not found", slog.Any("model", userModel), liblogger.Err(err))
		return user_dto.UserResponseDTO{}, errs.ErrUserNotFound
	}

	if err != nil {
		log.Error("failed get user", liblogger.Err(err))
		return user_dto.UserResponseDTO{}, errs.ErrInternalError.Wrap("failed find user")
	}

	return user_mapper.ToDTO(userResult), nil
}

func (s *UserService) GetById(ctx context.Context, id string) (user_dto.UserResponseDTO, error) {
	const op = "services.user_services.GetById"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse id from string to uuid", slog.String("id", id), liblogger.Err(err))
		return user_dto.UserResponseDTO{}, errs.ErrBadRequest.Wrap("failed parse uuid")
	}

	userResult, err := s.userRepository.GetById(ctx, s.db, uid)
	if s.db.IsNotFound(err) {
		log.Warn("user not found", slog.String("id", id), liblogger.Err(err))
		return user_dto.UserResponseDTO{}, errs.ErrUserNotFound
	}

	if err != nil {
		log.Error("failed to get user", liblogger.Err(err))
		return user_dto.UserResponseDTO{}, errs.ErrInternalError.Wrap("failed find user")
	}

	return user_mapper.ToDTO(userResult), nil
}

func (s *UserService) GetUserParticipantById(ctx context.Context, id string) (user_dto.UserParticipantResponseDTO, error) {
	const op = "services.user_services.GetUserParticipantById"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse id from string to uuid", slog.String("id", id), liblogger.Err(err))
		return user_dto.UserParticipantResponseDTO{}, errs.ErrBadRequest.Wrap("failed parse uuid")
	}

	userResult, err := s.userRepository.GetById(ctx, s.db, uid)
	if s.db.IsNotFound(err) {
		log.Warn("user not found", slog.String("id", id), liblogger.Err(err))
		return user_dto.UserParticipantResponseDTO{}, errs.ErrUserNotFound
	}

	if err != nil {
		log.Error("failed to get user", liblogger.Err(err))
		return user_dto.UserParticipantResponseDTO{}, errs.ErrInternalError.Wrap("failed find user")
	}

	participantResult, err := s.participantRepository.GetByUserId(ctx, s.db, userResult.ID)
	if s.db.IsNotFound(err) {
		log.Error("participant not found", slog.String("user id", userResult.ID.String()), liblogger.Err(err))
		return user_dto.UserParticipantResponseDTO{}, errs.ErrParticipantNotFound
	}

	if err != nil {
		log.Error("failed to get participant", liblogger.Err(err))
		return user_dto.UserParticipantResponseDTO{}, errs.ErrInternalError.Wrap("failed to get participant")
	}

	return user_dto.UserParticipantResponseDTO{
		UserResponseDTO:        user_mapper.ToDTO(userResult),
		ParticipantResponseDTO: paricipant_mapper.ToDTO(participantResult),
	}, nil

}

func (s *UserService) GetByListId(ctx context.Context, ids []string) ([]user_dto.UserResponseDTO, error) {
	const op = "services.user_handler.GetByListId"

	log := s.log.With(
		slog.String("op", op),
	)

	uids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			log.Error("failed to parse id", slog.String("id", id), liblogger.Err(err))
			return nil, errs.ErrBadRequest.Wrap("failed parse uuid")
		}
		uids = append(uids, uid)
	}

	usersResult, err := s.userRepository.GetByListId(ctx, s.db, uids)
	if err != nil {
		log.Error("failed to get users by list id", slog.Any("ids", ids), liblogger.Err(err))
		return nil, errs.ErrInternalError.Wrap("failed find users by list id")
	}

	usersDTO := make([]user_dto.UserResponseDTO, 0, len(usersResult))
	for _, userResult := range usersResult {
		usersDTO = append(usersDTO, user_mapper.ToDTO(userResult))
	}
	return usersDTO, nil
}

func (s *UserService) Update(ctx context.Context, id string, userDto user_dto.UpdateUserRequestDTO) error {
	const op = "services.user_sevice.Update"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", slog.String("id", id), liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("failed parse uuid")
	}

	userModel := user_mapper.FromUpdateToModel(userDto, uid)
	if userDto.Password != nil {
		userModel.PasswordHash, err = crypt.HashPassword(*userDto.Password)
		if err != nil {
			log.Error("failed hash password", liblogger.Err(err))
			return errs.ErrInternalError.Wrap("failed hash password")
		}
	}

	err = s.userRepository.Update(ctx, s.db, userModel)
	if err != nil {
		log.Error("failed update user", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed update user")
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	const op = "services.user_service.Delete"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id", slog.String("id", id), liblogger.Err(err))
		return errs.ErrBadRequest.Wrap("failed parse uuid")
	}

	err = s.userRepository.Delete(ctx, s.db, uid)
	if err != nil {
		log.Error("failed delete user", liblogger.Err(err))
		return errs.ErrInternalError.Wrap("failed delete user")
	}

	return nil
}
