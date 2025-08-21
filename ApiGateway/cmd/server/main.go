package main

import (
	"flag"
	"log"
	"main/internal/config"
	handlers "main/internal/handlers"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

func getConfigPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "../../config-yaml/config.yaml")
}

// corsMiddleware добавляет заголовки CORS к ответам
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Логируем запрос для отладки
		log.Printf("CORS Middleware: Handling %s %s from origin %s", r.Method, r.URL.Path, r.Header.Get("Origin"))

		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "http://172.16.0.196:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Expose-Headers", "Link")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "300")

		// Обрабатываем предварительные запросы OPTIONS
		if r.Method == http.MethodOptions {
			log.Printf("CORS Middleware: Handling OPTIONS request for %s", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		// Передаем запрос следующему обработчику
		log.Printf("CORS Middleware: Passing %s %s to next handler", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfgPath := flag.String("config", getConfigPath(), "path to config file")
	addr := "0.0.0.0:6611"
	flag.Parse()

	cfg := config.GetConfig(*cfgPath)
	mux := http.NewServeMux()

	for _, route := range cfg.HTTPServer.Routes {
		var handler http.Handler

		if route.Aggregate {
			handler = handlers.NewAggregateHandler(route, cfg.JwtTemp.Key)
		} else {
			handler = handlers.NewProxyHandler(route, cfg.JwtTemp.Key)
		}

		// Логируем регистрацию маршрута
		log.Printf("Registering route: %s (Aggregate: %v, Target: %s)", route.Prefix, route.Aggregate, route.Target)
		mux.Handle(route.Prefix, http.StripPrefix(strings.TrimSuffix(route.Prefix, "/"), handler))
	}

	// Логируем все запросы
	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		mux.ServeHTTP(w, r)
	})

	// Применяем CORS middleware
	corsHandler := corsMiddleware(loggedMux)

	log.Printf("API Gateway listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, corsHandler))
}
