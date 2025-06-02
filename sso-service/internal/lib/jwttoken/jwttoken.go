package jwttoken

import (
	"fmt"
	"main/internal/models/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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

func (m *JWTManager) VerifyToken(token string) (bool, error) {
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if !tokenParsed.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}
