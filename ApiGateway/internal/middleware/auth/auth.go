package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserInfoKey struct{}

type UserInfo struct {
	ID   uuid.UUID
	Role int
}

type TokenAccessClaims struct {
	Email string `json:"email"`
	Role  int    `json:"role"`
	jwt.RegisteredClaims
}

func ParseAccessTokenWithClaims(secret string, tokenStr string) (*TokenAccessClaims, error) {
	const op = "jwtManager.ParseAccessTokenWithClaims"

	token, err := jwt.ParseWithClaims(tokenStr, &TokenAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*TokenAccessClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Function to verify JWT tokens
// func verifyToken(tokenString string, secretKey []byte) (*jwt.Token, error) {

// 	// Parse the token with the secret key
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})

// 	// Check for verification errors
// 	if err != nil {
// 		return nil, fmt.Errorf("FAILED!: %w", err)
// 	}

// 	// Check if the token is valid
// 	if !token.Valid {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	// Return the verified token
// 	return token, nil
// }

func AuthenticateMiddleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := ParseAccessTokenWithClaims(secret, tokenStr)
		if err != nil {
			log.Printf("token error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		id, err := uuid.Parse(claims.Subject)
		if err != nil {
			log.Printf("token error: %v", err)
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}
		role := claims.Role
		ctx := context.WithValue(r.Context(), UserInfoKey{}, UserInfo{ID: id, Role: role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Check role for access
func RoleBasedAccess(requiredRole ...int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Start rolebaseAccess")
			user := r.Context().Value(UserInfoKey{}).(UserInfo)
			for _, role := range requiredRole {
				if user.Role == role {

					next.ServeHTTP(w, r)
					return

				}
			}
			http.Error(w, "Access denied", http.StatusForbidden)

		})
	}
}
