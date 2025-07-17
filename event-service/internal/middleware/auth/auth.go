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
	id   string
	role int
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

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// // Retrieve the token from the cookie
		// cookie, err := r.Cookie("token")
		// if err != nil {
		// 	log.Println("Token missing in cookie")
		// 	http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		// 	return
		// }

		// Verify the token
		token, err := verifyToken(tokenString, []byte(key))
		if err != nil {
			log.Printf("Token verification failed: %v\n", err)
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Print information about the verified token
		fmt.Printf("Token verified successfully. Claims: %+v\n", token.Claims)

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid token claims")
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Get role
		role, ok := claims["role"].(float64)
		if !ok {
			log.Println("Role missing in JWT")
			http.Error(w, "Role missing in JWT", http.StatusForbidden)
			return
		}

		// Get id
		id, ok := claims["id"].(string)
		if !ok {
			log.Println("Id missing in JWT")
			http.Error(w, "Id missing in JWT", http.StatusForbidden)
			return
		}

		// Init struct
		userInfo := UserInfo{id: id, role: int(role)}
		// Add struct in context
		r = r.WithContext(context.WithValue(r.Context(), UserInfoKey{}, userInfo))

		// Continue with the next middleware or route handler
		next.ServeHTTP(w, r)
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
