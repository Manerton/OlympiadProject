package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserInfoKey struct{}

type UserInfo struct {
	id   uuid.UUID
	role int
}

type TokenAccessClaims struct {
	Email string `json:"email"`
	Role  int    `json:"role"`
	jwt.RegisteredClaims
}

func ParseAccessTokenWithClaims(secret string, tokenStr string) (*TokenAccessClaims, error) {
	const op = "jwtManager.ParseAccessTokenWithClaims"

	fmt.Printf("key: %s, token: %s", secret, tokenStr)

	token, err := jwt.ParseWithClaims(tokenStr, &TokenAccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil // <- здесь преобразуем в []byte
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*TokenAccessClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// Function to verify JWT tokens
func verifyToken(tokenString string, secretKey []byte) (*jwt.Token, error) {

	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, fmt.Errorf("FAILED!: %w", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}

func AuthenticateMiddleware(next http.Handler, key string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := ParseAccessTokenWithClaims(key, tokenStr)
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
		ctx := context.WithValue(r.Context(), UserInfoKey{}, UserInfo{id: id, role: role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Check role for access
func RoleBasedAccess(requiredRole int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(UserInfoKey{}).(UserInfo)

			if user.role != requiredRole {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
