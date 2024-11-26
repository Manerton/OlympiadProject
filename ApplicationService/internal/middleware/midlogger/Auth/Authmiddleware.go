package Auth

/*
import (
	"context"
	"net/http"
	"github.com/gorilla/sessions"
)

// AuthMiddleware checks the JWT token
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 	   tokenString := r.Header.Get("Authorization")
 	   if tokenString == "" {
 		   http.Error(w, "Unauthorized", http.StatusUnauthorized)
 		   return
 	   }

 	   // Parse token
 	   token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
 		   // Validate the algorithm
 		   if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
 			   return nil, http.ErrNoCookie
 		   }
 		   return jwtSecret, nil
 	   })
 	   if err != nil || !token.Valid {
 		   http.Error(w, "Unauthorized", http.StatusUnauthorized)
 		   return
 	   }

 	   // Extract user information
 	   claims, ok := token.Claims.(jwt.MapClaims)
 	   if !ok {
 		   http.Error(w, "Unauthorized", http.StatusUnauthorized)
 		   return
 	   }

 	   // Set the user information in context
 	   ctx := context.WithValue(r.Context(), "username", claims["username"])
 	   ctx = context.WithValue(ctx, "role", claims["role"])
 	   next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func RoleMiddleware(requiredRole int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Context().Value("role").(int)
			if role != requiredRole {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
*/
