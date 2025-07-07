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
	"main/internal/models/refresh_token"
	"main/internal/models/user"
	"main/internal/services/another_service"
	"main/internal/storage/orm"
	redisdb "main/internal/storage/redis"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (user.User, error)
	GetByEmail(ctx context.Context, orm orm.ORM, email string) (*user.User, error)
	Create(ctx context.Context, orm orm.ORM, user user.User) (uuid.UUID, error)
	Update(ctx context.Context, orm orm.ORM, user user.User) error
}

type ParticipantRepository interface {
	Create(ctx context.Context, orm orm.ORM, participant participant.Participant) (uuid.UUID, error)
}

type RefreshRepository interface {
	GetById(ctx context.Context, orm orm.ORM, id uuid.UUID) (refresh_token.RefreshToken, error)

	Create(ctx context.Context, orm orm.ORM, refreshToken refresh_token.RefreshToken) (uuid.UUID, error)
	Update(ctx context.Context, orm orm.ORM, conditions *refresh_token.RefreshToken, token refresh_token.RefreshToken) error
	Delete(ctx context.Context, orm orm.ORM, id uuid.UUID) error
}

type AuthService struct {
	log                   *slog.Logger
	jwtManager            *jwttoken.JWTManager
	db                    orm.ORM
	userRepository        UserRepository
	participantRepository ParticipantRepository
	refreshRepository     RefreshRepository
}

func New(log *slog.Logger, orm orm.ORM, jwtManager *jwttoken.JWTManager,
	userRepository UserRepository, participantRepository ParticipantRepository, refreshRepository RefreshRepository) *AuthService {
	return &AuthService{
		log:                   log,
		db:                    orm,
		jwtManager:            jwtManager,
		userRepository:        userRepository,
		participantRepository: participantRepository,
		refreshRepository:     refreshRepository,
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
		log.Error("failed to get user by email", loginRequest.Email, liblogger.Err(err))
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
		log.Error("failed when create token", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	// create refresh token
	refreshToken, err := s.preparationRefreshToken(ctx, *userResult)
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

func (s *AuthService) preparationRefreshToken(ctx context.Context, userResult user.User) (string, error) {
	// create refresh token
	refreshToken, err := s.jwtManager.CreateRefreshToken(userResult)
	if err != nil {
		return "", fmt.Errorf("failed when create refresh token: %w", err)
	}

	// hash token

	refreshTokenModel := refresh_token.RefreshToken{
		UserID:    userResult.ID,
		TokenHash: refreshToken,
		ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshDuration()),
	}

	// save refresh token
	_, err = s.refreshRepository.Create(ctx, s.db, refreshTokenModel)
	if err != nil {
		return "", fmt.Errorf("failed save refresh token: %w", err)
	}

	return refreshToken, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, registerUser *register_dto.RegisterUserRequestDTO) error {
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
		log.Error("failed user exist", slog.Any("user", userFind))
		return fmt.Errorf("%s: %s", errMsg, "user already exist")
	}

	userModel := user_mapper.FromRegisterUserToModel(registerUser)
	userModel.PasswordHash, err = crypt.HashPassword(registerUser.Password)
	if err != nil {
		log.Error("failed when hash password", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	_, err = s.userRepository.Create(ctx, s.db, userModel)
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

	userId, err := s.userRepository.Create(ctx, transaction, userModel)
	if err != nil {
		transaction.TransactionRollback()
		log.Error("failed when create user", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	participantModel := paricipant_mapper.FromRegisterToModel(registerRequst, userId)
	_, err = s.participantRepository.Create(ctx, transaction, participantModel)
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
		log.Error("failed get trust code", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	if trustCode != userCode {
		log.Error("userCode incorrent")
		return fmt.Errorf("%s", errMsg)
	}

	userModel, err := s.userRepository.GetByEmail(ctx, s.db, email)
	if err != nil {
		log.Error("failed to find user by email", email, liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}
	// activate account
	userModel.Activated = true
	err = s.userRepository.Update(ctx, s.db, *userModel)
	if err != nil {
		log.Error("failed update user activate status", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*login_dto.AuthResultDTO, error) {
	const op = "services.auth_service.Refresh"
	const errMsg = "failed to refresh tokens"

	log := s.log.With(
		slog.String("op", op),
	)

	// first check token
	nowToken, err := s.jwtManager.VerifyToken(refreshToken)
	if err != nil {
		log.Error("failed verify token", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	// get claims
	tokenClaims, err := s.jwtManager.GetRefreshClaims(nowToken)
	if err != nil {
		log.Error("failed get refresh token claims", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	// get token from db
	tokenDB, err := s.refreshRepository.GetById(ctx, s.db, tokenClaims.ID)
	if err != nil {
		log.Error("failed get refresh token", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	// check hash
	if !crypt.CheckPasswordHash(refreshToken, tokenDB.TokenHash) {
		log.Error("failed check hash token")
		return nil, fmt.Errorf("%s", errMsg)
	}

	// check token is revoked
	if tokenDB.Revoked {
		log.Error("Token was be revoked")
		return nil, fmt.Errorf("%s", errMsg)
	}

	// delete old token
	go func() {
		err := s.refreshRepository.Delete(ctx, s.db, tokenClaims.ID)
		if err != nil {
			log.Error("failed delete refresh token", liblogger.Err(err))
		}
	}()

	// find user
	userFind, err := s.userRepository.GetById(ctx, s.db, tokenClaims.UserId)
	if err != nil {
		log.Error("failed get user", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	// create token
	token, err := s.jwtManager.CreateToken(userFind)
	if err != nil {
		log.Error("failed when create token", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	// create refresh token
	newRefreshToken, err := s.preparationRefreshToken(ctx, userFind)
	if err != nil {
		log.Error("failed when create refresh token", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	return &login_dto.AuthResultDTO{
		AccessToken:      token,
		RefreshToken:     newRefreshToken,
		ExpiresInAccess:  int64(s.jwtManager.GetAccessDuration().Seconds()),
		ExpiresInRefresh: int64(s.jwtManager.GetRefreshDuration().Seconds()),
	}, err
}

func (s *AuthService) RevokeToken(ctx context.Context, id string) error {
	const op = "services.auth_service.RevokeToken"
	const errMsg = "failed to revoke token"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parse id to uuid", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	updates := refresh_token.RefreshToken{
		ID:      uid,
		Revoked: true,
	}

	err = s.refreshRepository.Update(ctx, s.db, nil, updates)
	if err != nil {
		log.Error("failed update refresh token", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}

func (s *AuthService) RevokeAllUserTokens(ctx context.Context, userId string) error {
	const op = "services.auth_service.RevokeAllUserTokens"
	const errMsg = "failed to revoke all user tokens"

	log := s.log.With(
		slog.String("op", op),
	)

	uid, err := uuid.Parse(userId)
	if err != nil {
		log.Error("failed parse user ID to uuid", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	conditions := refresh_token.RefreshToken{
		Revoked: true,
	}

	updates := refresh_token.RefreshToken{ID: uid}

	err = s.refreshRepository.Update(ctx, s.db, &conditions, updates)
	if err != nil {
		log.Error("failed updates all user tokens", liblogger.Err(err))
		return fmt.Errorf("%s", errMsg)
	}

	return nil
}
