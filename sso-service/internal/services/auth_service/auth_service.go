package auth_service

import (
	"context"
	"fmt"
	"log/slog"
	login_dto "main/internal/dto/auth/login"
	"main/internal/lib/liblogger"
	"main/internal/models/user"
)

type AuthRepository interface {
	GetByEmail(email string) (user.User, error)
}

type AuthService struct {
	log            *slog.Logger
	authRepository AuthRepository
}

func (s *AuthService) Login(ctx context.Context, loginRequest *login_dto.LoginRequestDTO) (*login_dto.LoginResponseDTO, error) {
	const op = "services.auth_services.Login"
	const errMsg = "failed login"

	log := s.log.With(
		slog.String("op", op),
	)

	userResult, err := s.authRepository.GetByEmail(loginRequest.Email)
	if err != nil {
		log.Error("failed to get user by email: %s", loginRequest.Email, liblogger.Err(err))
		return nil, fmt.Errorf(errMsg)
	}

	_ = userResult

	return nil, nil
}
