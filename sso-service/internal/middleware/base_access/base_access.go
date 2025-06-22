package base_access

import (
	"log"
	"main/internal/lib/jwttoken"
	"main/internal/lib/response"
	"net/http"

	"github.com/go-chi/render"
)

func BaseAccess(jwtManager *jwttoken.JWTManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				log.Println("token missing in cookie")
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, response.ErrorResponse("missing token in cookie"))
				return
			}

			_, err = jwtManager.VerifyToken(cookie.Value)
			if err != nil {
				log.Printf("token verification failed: %v\n", err)
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, response.ErrorResponse("invalid token claims"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
