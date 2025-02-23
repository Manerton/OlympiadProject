package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type UserInfoKey struct{}

type UserInfo struct {
	id   float64
	role string
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

		// Retrieve the token from the cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			log.Println("Token missing in cookie")
			return
		}

		// Verify the token
		token, err := verifyToken(cookie.Value, []byte(key))
		if err != nil {
			log.Printf("Token verification failed: %v\n", err)
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
		role, ok := claims["role"].(string)
		if !ok {
			log.Println("Role missing in JWT")
			http.Error(w, "Role missing in JWT", http.StatusForbidden)
			return
		}

		// Get id
		id, ok := claims["id"].(float64)
		if !ok {
			log.Println("Id missing in JWT")
			http.Error(w, "Id missing in JWT", http.StatusForbidden)
			return
		}

		// Init struct
		userInfo := UserInfo{id: id, role: role}
		// Add struct in context
		r = r.WithContext(context.WithValue(r.Context(), UserInfoKey{}, userInfo))

		// Continue with the next middleware or route handler
		next.ServeHTTP(w, r)
	})
}

// Check role for access
func RoleBasedAccess(requiredRole string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Start rolebaseAccess")
			user := r.Context().Value(UserInfoKey{}).(UserInfo)

			if user.role != requiredRole {
				http.Error(w, "Access denied", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
