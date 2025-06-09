package auth_service

import (
	"context"
	"fmt"
	"log/slog"
	login_dto "main/internal/dto/auth/login"
	register_dto "main/internal/dto/auth/register"
	"main/internal/lib/crypt"
	"main/internal/lib/jwttoken"
	"main/internal/lib/liblogger"
	paricipant_mapper "main/internal/lib/mapper/participant_mapper"
	"main/internal/lib/mapper/user_mapper"
	"main/internal/models/participant"
	"main/internal/models/user"
	"main/internal/services/another_service"
	"main/internal/storage/orm"
	redisdb "main/internal/storage/redis"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, orm orm.ORM, email string) (*user.User, error)
	Create(ctx context.Context, orm orm.ORM, user *user.User) (uuid.UUID, error)
	Update(ctx context.Context, orm orm.ORM, user *user.User) error
}

type ParticipantRepository interface {
	Create(ctx context.Context, orm orm.ORM, participant *participant.Participant) (uuid.UUID, error)
}

type AuthService struct {
	log                   *slog.Logger
	jwtManager            *jwttoken.JWTManager
	db                    orm.ORM
	userRepository        UserRepository
	participantRepository ParticipantRepository
}

func New(log *slog.Logger, orm orm.ORM, jwtManager *jwttoken.JWTManager,
	userRepository UserRepository, participantRepository ParticipantRepository) *AuthService {
	return &AuthService{
		log:                   log,
		db:                    orm,
		jwtManager:            jwtManager,
		userRepository:        userRepository,
		participantRepository: participantRepository,
	}
}

func (s *AuthService) Login(ctx context.Context, loginRequest *login_dto.LoginRequestDTO) (*login_dto.AuthResultDTO, error) {
	const op = "services.auth_services.Login"
	const errMsg = "failed login"

	log := s.log.With(
		slog.String("op", op),
	)

	userResult, err := s.userRepository.GetByEmail(ctx, s.db, loginRequest.Email)
	if err != nil {
		log.Error("failed to get user by email: %s", loginRequest.Email, liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	if userResult == nil {
		log.Error("failed user not fount")
		return nil, fmt.Errorf("%s: %s", errMsg, "user is not exist")
	}

	// check password hash
	if !crypt.CheckPasswordHash(loginRequest.Password, userResult.PasswordHash) {
		log.Error("uncorrect password")
		return nil, fmt.Errorf("%s: uncorrent login or password", errMsg)
	}

	// create token
	token, err := s.jwtManager.CreateToken(*userResult)
	if err != nil {
		log.Error("failed when create token:", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	// create refresh token
	refreshToken, err := s.jwtManager.CreateRefreshToken(*userResult)
	if err != nil {
		log.Error("failed when create refresh token", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	return &login_dto.AuthResultDTO{
		AccessToken:      token,
		RefreshToken:     refreshToken,
		ExpiresInAccess:  int64(s.jwtManager.GetAccessDuration().Seconds()),
		ExpiresInRefresh: int64(s.jwtManager.GetRefreshDuration().Seconds()),
	}, err
}

func (s *AuthService) RegisterUser(ctx context.Context, registerUser *register_dto.RegusterUserRequestDTO) error {
	const op = "services.auth_service.RegisterUser"
	const errMsg = "failed register user"

	log := s.log.With(
		slog.String("op", op),
	)

	userFind, err := s.userRepository.GetByEmail(ctx, s.db, registerUser.Email)
	if err != nil {
		log.Error("failed when check user exist", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	if userFind != nil {
		log.Error("failed user exist")
		return fmt.Errorf("%s: %s", errMsg, "user already exist")
	}

	userModel := user_mapper.FromRegisterUserToModel(registerUser)
	_, err = s.userRepository.Create(ctx, s.db, &userModel)
	if err != nil {
		log.Error("failed when create user", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

func (s *AuthService) RegisterParticipant(ctx context.Context, registerRequst *register_dto.RegisterParticipantRequestDTO) error {
	const op = "services.auth_service.RegisterParticipant"
	const errMsg = "failed register participant"

	log := s.log.With(
		slog.String("op", op),
	)

	userFind, err := s.userRepository.GetByEmail(ctx, s.db, registerRequst.Email)
	if err != nil {
		log.Error("failed when check user exist", liblogger.Err(err))
		return fmt.Errorf("%s", err)
	}

	if userFind != nil {
		log.Error("failed user exist")
		return fmt.Errorf("%s: %s", errMsg, "user already exist")
	}

	userModel := user_mapper.FromRegisterToModel(registerRequst)
	transaction, err := s.db.TransactionBegin()
	if err != nil {
		log.Error("failed when begin transaction", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	userModel.PasswordHash, err = crypt.HashPassword(registerRequst.Password)
	if err != nil {
		log.Error("failed when hash password", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	userId, err := s.userRepository.Create(ctx, transaction, &userModel)
	if err != nil {
		transaction.TransactionRollback()
		log.Error("failed when create user", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}
	participantModel := paricipant_mapper.FromRegisterToModel(registerRequst, userId)
	_, err = s.participantRepository.Create(ctx, transaction, &participantModel)
	if err != nil {
		transaction.TransactionRollback()
		log.Error("failed when create participant", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	transaction.TransactionCommit()

	err = another_service.SendNotifyAcceptAccount(userModel.Email)
	if err != nil {
		log.Error("failed send notify on email: %w", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

func (s *AuthService) ActivateAccount(ctx context.Context, email string, userCode string) error {
	const op = "services.auth_service.ActivateAccount"
	const errMsg = "failed to activate account"

	log := s.log.With(
		slog.String("op", op),
	)

	trustCode, err := redisdb.GetActivationCode(email)
	if err != nil {
		log.Error("failed get trust code: %w", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	if trustCode != userCode {
		log.Error("userCode incorrent")
		return fmt.Errorf("%s", errMsg)
	}

	userModel, err := s.userRepository.GetByEmail(ctx, s.db, email)
	if err != nil {
		log.Error("failed to find user by email: %s", email, liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}
	// activate account
	userModel.Activated = true
	err = s.userRepository.Update(ctx, s.db, userModel)
	if err != nil {
		log.Error("failed update user activate status: %w", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}
