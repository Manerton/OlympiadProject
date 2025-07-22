package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type UserInfoKey struct{}

type UserInfo struct {
	ID   float64
	Role string
}

func verifyToken(tokenString string, secret []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}

func AuthenticateMiddleware(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := verifyToken(tokenStr, []byte(secret))
		if err != nil {
			log.Printf("token error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		id, _ := claims["id"].(float64)
		role, _ := claims["role"].(string)
		ctx := context.WithValue(r.Context(), UserInfoKey{}, UserInfo{ID: id, Role: role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Check role for access
func RoleBasedAccess(requiredRole ...string) func(next http.Handler) http.Handler {
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
