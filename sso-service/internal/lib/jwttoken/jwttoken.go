package jwttoken

import (
	"fmt"
	"main/internal/models/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(user user.User, secretKey string, duration time.Duration) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(duration).Unix(),
	})

	tokenStr, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func CreateRefreshToken() (string, error) {
	return "", nil
}

func VerifyToken(token string, secretKey []byte) (bool, error) {
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return false, err
	}

	if !tokenParsed.Valid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}
