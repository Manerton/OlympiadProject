package main

import (
	"bytes"
	"flag"
	"io"
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
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://172.16.0.196:5173",
		}

		for _, o := range allowedOrigins {
			if origin == o {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin") // важно при кэшировании
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		log.Printf("CORS Middleware: Passing %s %s to next handler", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware2(cfg *config.AdditionalAddressesConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			for _, o := range cfg.AllowOrigins {
				if origin == o {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin") // важно при кэшировании
					break
				}
			}

			log.Println("ORIGIN", origin)

			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			log.Printf("CORS Middleware: Passing %s %s to next handler", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	cfgPath := flag.String("config", getConfigPath(), "path to config file")
	flag.Parse()

	cfg := config.GetConfig(*cfgPath)
	mux := http.NewServeMux()

	for _, route := range cfg.HTTPServers.Routes {
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

	loggedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// читаем тело запроса
		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // восстанавливаем body для следующего хендлера

		// логируем куки
		cookies := r.Cookies()

		// оборачиваем writer, чтобы перехватить ответ
		lrw := newLoggingResponseWriter(w)

		// вызываем основной mux
		mux.ServeHTTP(lrw, r)

		// Логируем
		log.Printf("➡️ Request: %s %s, Cookies: %+v, Body: %s",
			r.Method, r.URL.Path, cookies, string(bodyBytes))

		log.Printf("⬅️ Response: %d, Body: %s",
			lrw.statusCode, lrw.body.String())
	})

	// Применяем CORS middleware
	corsHandler := CORSMiddleware2(&cfg.AdditionalAddressesConfig)(loggedMux)

	log.Printf("API Gateway listening on %s", cfg.HTTPServerMain.GetAddress())
	log.Fatal(http.ListenAndServe(cfg.HTTPServerMain.GetAddress(), corsHandler))
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code

	// Логируем Set-Cookie заголовки
	if cookies := lrw.Header().Values("Set-Cookie"); len(cookies) > 0 {
		for _, c := range cookies {
			log.Printf("⬅️ Set-Cookie: %s", c)
		}
	}

	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b) // копируем в буфер
	return lrw.ResponseWriter.Write(b)
}
