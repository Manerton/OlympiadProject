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
	"main/internal/models/user"
	"main/internal/storage/orm"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, orm orm.ORM, email string) (user.User, error)
}

type AuthService struct {
	log            *slog.Logger
	jwtManager     *jwttoken.JWTManager
	db             orm.ORM
	userRepository UserRepository
}

func NewAuthService(log *slog.Logger, jwtManager *jwttoken.JWTManager, userRepository UserRepository) *AuthService {
	return &AuthService{
		log:            log,
		jwtManager:     jwtManager,
		userRepository: userRepository,
	}
}

func (s *AuthService) Login(ctx context.Context, loginRequest *login_dto.LoginRequestDTO) (*login_dto.LoginResponseDTO, error) {
	const op = "services.auth_services.Login"
	const errMsg = "failed login"

	log := s.log.With(
		slog.String("op", op),
	)

	userResult, err := s.userRepository.GetByEmail(ctx, s.db, loginRequest.Email)
	if err != nil {
		log.Error("failed to get user by email: %s", loginRequest.Email, liblogger.Err(err))
		return nil, fmt.Errorf(errMsg)
	}

	// check password hash
	if !crypt.CheckPasswordHash(loginRequest.Password, userResult.PasswordHash) {
		log.Error("uncorrect password")
		return nil, fmt.Errorf("%s: uncorrent login or password", errMsg)
	}

	// create token
	token, err := s.jwtManager.CreateToken(userResult)
	if err != nil {
		log.Error("failed when create token:", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	// create refresh token
	refreshToken, err := s.jwtManager.CreateRefreshToken(userResult)
	if err != nil {
		log.Error("failed when create refresh token", liblogger.Err(err))
		return nil, fmt.Errorf("%s", errMsg)
	}

	return &login_dto.LoginResponseDTO{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, err
}

func (s *AuthService) Register(ctx context.Context, registerRequst *register_dto.RegisterParticipantRequestDTO) error {
	const op = "services.auth_service.Register"

	return nil
}
