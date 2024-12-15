package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Структура пользователя
type User struct {
	Email    string
	Name     string
	Password string
}

// Моковые пользователи
var users = map[uint]User{
	1: {
		Email:    "makar@mail.com",
		Name:     "Макаров Макар Макарович",
		Password: "makar",
	},
	2: {
		Email:    "jury@mail.com",
		Name:     "Иван Иванович Иванов",
		Password: "jury",
	},
	3: {
		Email:    "Organizer@mail.com",
		Name:     "Организатор организаторович",
		Password: "Organizer",
	},
	4: {
		Email:    "admin@mail.com",
		Name:     "Администратор администраторович",
		Password: "admin",
	},
}

func init() {
	// Изначальные пользователи с их данными
	plainUsers := map[uint]struct {
		Email    string
		Name     string
		Password string
	}{
		1: {
			Email:    "makar@mail.com",
			Name:     "Макаров Макар Макарович",
			Password: "makar",
		},
		2: {
			Email:    "jury@mail.com",
			Name:     "Иван Иванович Иванов",
			Password: "jury",
		},
		3: {
			Email:    "Organizer@mail.com",
			Name:     "Организатор организаторович",
			Password: "Organizer",
		},
		4: {
			Email:    "admin@mail.com",
			Name:     "Администратор администраторович",
			Password: "admin",
		},
	}

	// Хешируем пароли и заполняем `users`
	for id, data := range plainUsers {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(fmt.Sprintf("Не удалось захешировать пароль для пользователя %d: %v", id, err))
		}
		users[id] = User{
			Email:    data.Email,
			Name:     data.Name,
			Password: string(hashedPassword),
		}
	}
}

// Структура пользователя для jwt
type UserJwt struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

//1 - student
//2 - jury
//3 - organizer
//4 - admin

// Моковые пользователи для jwt
var usersjwt = map[uint]UserJwt{
	1: {
		ID:   1,
		Name: "Макаров Макар Макарович",
		Role: "1",
	},
	2: {
		ID:   2,
		Name: "Иван Иванович Иванов",
		Role: "2",
	},
	3: {
		ID:   3,
		Name: "Организатор организаторович",
		Role: "3",
	},
	4: {
		ID:   4,
		Name: "Администратор администраторович",
		Role: "4",
	},
}

// Add a new global variable for the secret key
var secretKey = []byte("lololololololololololololololololololololololololololololol")

// Function to create JWT tokens with claims
func createToken(id uint) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "localhost:8081", // Issuer
		"id":   usersjwt[id].ID,
		"name": usersjwt[id].Name,
		"role": usersjwt[id].Role,
		"exp":  time.Now().Add(time.Minute * 120).Unix(), // Expiration time
		"iat":  time.Now().Unix(),                        // Issued at
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Print information about the created token
	fmt.Printf("Token claims added: %+v\n", claims)
	return tokenString, nil
}

// Функция для создания Refresh Token(обычная 32 символьная уникальная строчка)
func createRefreshToken() (string, error) {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "localhost:8081",
		"token": fmt.Sprintf("%x", b),
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(), // Срок действия 7 дней

	})
	return claims.SignedString(secretKey)
}

// Function to verify JWT tokens
func verifyToken(tokenString string) (*jwt.Token, error) {
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

func authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token from the cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			fmt.Println("Token missing in cookie")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Verify the token
		token, err := verifyToken(cookie.Value)
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

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest

	// Parse the incoming JSON request
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Search for the user in the mock users map
	var user User
	var userID uint
	found := false
	for id, u := range users {
		if u.Email == loginReq.Email {
			user = u
			userID = id
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Verify the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create a JWT token for the user
	token, err := createToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := createRefreshToken()
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	// Set the JWT token in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Domain:   "localhost",
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour), // Match the token expiration
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: false,
		Secure:   false, // только по HTTPS
		Domain:   "localhost",
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Login successful")
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем refresh token из cookie
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token missing", http.StatusUnauthorized)
		return
	}

	// Проверяем refresh token
	token, err := verifyToken(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Извлекаем ID из токена
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["id"].(uint))

	// Генерация нового access token
	accessToken, err := createToken(userID)
	if err != nil {
		http.Error(w, "Error generating new access token", http.StatusInternalServerError)
		return
	}

	// Отправляем новый access token в cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Access token refreshed"))
}

func GetMyRoleAndID(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println("Token missing in cookie")
		return
	}

	// Verify the token
	token, err := verifyToken(cookie.Value)
	if err != nil {
		fmt.Printf("Token verification failed: %v\n", err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Invalid token claims")
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}
	fmt.Println(claims["role"], reflect.TypeOf(claims["role"]))

	// Get role
	role, ok := claims["role"].(string)
	if !ok {
		fmt.Println("Role missing in JWT")
		http.Error(w, "Role missing in JWT", http.StatusForbidden)
		return
	}
	fmt.Println(claims["id"], reflect.TypeOf(claims["id"]))
	// Get id
	id, ok := claims["id"].(float64)
	if !ok {
		fmt.Println("Id missing in JWT")
		http.Error(w, "Id missing in JWT", http.StatusForbidden)
		return
	}

	fmt.Println(claims["id"], reflect.TypeOf(claims["name"]))
	// Get id
	name, ok := claims["name"].(string)
	if !ok {
		fmt.Println("name missing in JWT")
		http.Error(w, "name missing in JWT", http.StatusForbidden)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"id":   uint(id),
		"role": role,
		"name": name,
	})
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Unix(0, 0), // Expire immediately
	})

	// Optionally clear the refresh token cookie if used
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}

func main() {

	// Инициализация маршрутизатора Chi
	r := chi.NewRouter()

	// Middleware для логирования запросов
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)

	// init cors
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // React URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // В секундах
	}
	r.Use(cors.Handler(corsOptions))
	// Определяем эндпоинт для получения текущего пользователя
	r.Post("/login", loginHandler)
	r.Post("/refresh", refreshHandler)
	r.Post("/logout", logoutHandler)
	r.Get("/my-info", GetMyRoleAndID)

	// Запускаем сервер
	fmt.Println("Server running on http://localhost:8081")
	http.ListenAndServe(":8081", r)
}
