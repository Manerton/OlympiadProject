package jwttoken

import (
	"fmt"
	"main/internal/models/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	ID     uuid.UUID
	UserId uuid.UUID
}

type JWTManager struct {
	secretKey       []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewJWTManager(secretKey []byte, accessDuration time.Duration, refreshDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:       secretKey,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

func (m *JWTManager) GetAccessDuration() time.Duration {
	return m.accessDuration
}

func (m *JWTManager) GetRefreshDuration() time.Duration {
	return m.refreshDuration
}

func (m *JWTManager) GetRefreshClaims(token *jwt.Token) (*TokenClaims, error) {
	resultClaims, err := m.GetClaims(token)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed get claims")
	}

	// exist only refresh token
	id, ok := claims["id"].(string)
	if !ok {
		return nil, fmt.Errorf("failed get id on claims")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed parse id to uuid: %w", err)
	}
	resultClaims.ID = uid

	return resultClaims, nil
}

func (m *JWTManager) GetClaims(token *jwt.Token) (*TokenClaims, error) {
	resultClaims := &TokenClaims{}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed get claims")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("failed get user_id on claims")
	}

	uid, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("failed parse user_id to uuid: %w", err)
	}
	resultClaims.UserId = uid

	return resultClaims, nil
}

func (m *JWTManager) CreateToken(user user.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(m.accessDuration).Unix(),
	})

	tokenStr, err := claims.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (m *JWTManager) CreateRefreshToken(user user.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(m.refreshDuration).Unix(),
	})

	tokenStr, err := claims.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (m *JWTManager) VerifyToken(token string) (*jwt.Token, error) {
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !tokenParsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return tokenParsed, nil
}
