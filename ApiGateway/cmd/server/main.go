/*
Простой API Gateway на Go с внешней конфигурацией маршрутов и JWT-валидацией.

Концепция:
1. Читаем YAML файл `config.yaml`, где описаны:
   - Секрет JWT
   - Маршруты: префикс пути → URL микросервиса
2. Строим middleware для проверки токена (JWT).
3. На основе маршрутов создаём ReverseProxy для каждого префикса.
4. Подключаем всё к HTTP-серверу.

*/

package main

import (
	"flag"
	"log"
	"main/internal/config"
	"main/internal/middleware/auth"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
)

func getConfigPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "../../config-yaml/config.yaml")
}

// singleJoiningSlash правильно объединяет пути со слэшами
func singleJoiningSlash(a, b string) string {
	a = strings.TrimSuffix(a, "/")
	b = strings.TrimPrefix(b, "/")
	if b == "" {
		return a
	}
	return a + "/" + b
}

func NewProxyHandler(route config.Route, jwtKey string) http.Handler {
	target, err := url.Parse(route.Target)
	if err != nil {
		log.Fatalf("bad target URL %s: %v", route.Target, err)
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = singleJoiningSlash(
				target.Path,
				strings.TrimPrefix(req.URL.Path, route.Prefix),
			)
			req.Host = target.Host
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		},
	}

	// Создаем базовый handler
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Proxying: %s -> %s%s", req.URL.Path, target.Host, req.URL.Path)
		proxy.ServeHTTP(w, req)
	})

	// Добавляем middleware аутентификации если нужно
	if !route.SkipAuth {
		handler = auth.AuthenticateMiddleware(jwtKey, handler)
	}

	return handler
}

func main() {
	cfgPath := flag.String("config", getConfigPath(), "path to config file")
	addr := "0.0.0.0:6611"
	flag.Parse()

	cfg := config.GetConfig(*cfgPath)
	mux := http.NewServeMux()

	// Регистрируем все маршруты из конфига
	for _, route := range cfg.HTTPServer.Routes {
		handler := NewProxyHandler(route, cfg.JwtTemp.Key)
		mux.Handle(
			route.Prefix,
			http.StripPrefix(
				strings.TrimSuffix(route.Prefix, "/"),
				handler,
			),
		)
	}

	// Добавляем логгирование всех запросов
	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})

	log.Printf("API Gateway listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, loggedMux))
}

// // jwtMiddleware проверяет токен в Authorization: Bearer <token>
// func jwtMiddleware(secret string, next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		auth := r.Header.Get("Authorization")
// 		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
// 			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
// 			return
// 		}
// 		tokenStr := strings.TrimPrefix(auth, "Bearer ")

// 		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 			// Только HMAC
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, jwt.ErrSignatureInvalid
// 			}
// 			return []byte(secret), nil
// 		})
// 		if err != nil || !token.Valid {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}
// 		// Всё ок, передаём дальше
// 		next.ServeHTTP(w, r)
// 	})
//}
