package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

// Function to verify JWT tokens
func verifyToken(tokenString string, secretKey string) (*jwt.Token, error) {

	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Return the verified token
	return token, nil
}

func authenticateMiddleware(next http.Handler, key string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token from the cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			fmt.Println("Token missing in cookie")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Verify the token
		token, err := verifyToken(cookie.Value, key)
		if err != nil {
			fmt.Printf("Token verification failed: %v\n", err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Print information about the verified token
		fmt.Printf("Token verified successfully. Claims: %+v\n", token.Claims)

		// Continue with the next middleware or route handler
		next.ServeHTTP(w, r)
	})
}
