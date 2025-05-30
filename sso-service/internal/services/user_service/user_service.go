package user_service

import (
	"context"
	"fmt"
	"log/slog"
	user_dto "main/internal/dto/user"
	"main/internal/lib/liblogger"
	"main/internal/models/user"
	"main/internal/storage/orm"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (user.User, error)
	GetByListId(ctx context.Context, orm orm.ORM, ids []uuid.UUID) ([]user.User, error)
}

type UserService struct {
	userRepository UserRepository
	db             orm.ORM
	log            *slog.Logger
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
		return user_dto.UserResponseDTO{}, fmt.Errorf(errMsg)
	}

	userResult, err := s.userRepository.GetById(ctx, s.db, uid)
	if err != nil {
		log.Error("failed to get user", liblogger.Err(err))
		return user_dto.UserResponseDTO{}, fmt.Errorf(errMsg)
	}

	return user.ToDTO(userResult), nil
}

func (s *UserService) GetByListId(ctx context.Context, ids []string) ([]user_dto.UserResponseDTO, error) {
	const op = "service.user_handler.GetByListId"
	const errMsg = "failed to find users by list id"

	log := s.log.With(
		slog.String("op", op),
	)

	uids := make([]uuid.UUID, len(ids))
	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			log.Error("failed to parse id: %s", id, liblogger.Err(err))
			return nil, fmt.Errorf(errMsg)
		}
		uids = append(uids, uid)
	}

	usersResult, err := s.userRepository.GetByListId(ctx, s.db, uids)
	if err != nil {
		log.Error("failed to get users by list id: %v", ids, liblogger.Err(err))
		return nil, fmt.Errorf(errMsg)
	}

	usersDTO := make([]user_dto.UserResponseDTO, len(usersResult))
	for _, userResult := range usersResult {
		usersDTO = append(usersDTO, user.ToDTO(userResult))
	}
	return usersDTO, nil
}
