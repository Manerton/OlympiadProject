package handlers

import (
	"encoding/json"
	"log"
	"main/internal/config"
	"main/internal/middleware/auth"
	"main/internal/responcetypes"
	"main/internal/services"
	"main/internal/strategy"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// NewAggregateHandler возвращает http.Handler, который:
// 1) валидирует JWT (если нужно),
// 2) проверяет роль (если указаны route.Roles),
// 3) вызывает AggregateService для сбора данных,
// 4) оборачивает в ApiResponse и отдаёт JSON.

func getStrategy(prefix string) strategy.AggregationStrategy {
	switch prefix {
	case "/aplicationevent":
		return strategy.NewDefaultAggregationStrategy(5 * time.Second)
	default:
		return nil
	}
}

func NewAggregateHandler(route config.Route, jwtKey string) http.Handler {
	// Инициализация сервиса с таймаутом (можно вынести в конфиг)
	service := services.NewAggregateService(5 * time.Second)

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Выбор стратегии по префиксу маршрута
		var aggrstrategy = getStrategy(route.Prefix)

		// 2. Выполнение агрегации
		dataSlice, err := service.Aggregate(route, r, aggrstrategy)
		if err != nil {
			log.Printf("Aggregation failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. Формирование ответа
		resp := responcetypes.ApiResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Data:       dataSlice,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	})

	// Добавляем middleware аутентификации если нужно
	if !route.SkipAuth {
		handler = auth.AuthenticateMiddleware(jwtKey, handler)
		if len(route.Roles) > 0 {
			handler = auth.RoleBasedAccess(route.Roles...)(handler)
		}
	}

	return handler
}
func singleJoiningSlash(a, b string) string {
	a = strings.TrimSuffix(a, "/")
	b = strings.TrimPrefix(b, "/")
	if b == "" {
		return a
	}
	return a + "/" + b
}

// NewProxyHandler проксирует все запросы на route.Target и оборачивает JWT‑миддлварью
func NewProxyHandler(route config.Route, jwtKey string) http.Handler {
	// разбираем URL таргета
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

	// базовый handler
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// убираем префикс из пути
		//r.URL.Path = strings.TrimPrefix(r.URL.Path, route.Prefix)
		log.Printf("Proxying to %s%s", target.Host, r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	// JWT‑проверка, если нужно
	if !route.SkipAuth {
		handler = auth.AuthenticateMiddleware(jwtKey, handler)
	}

	return handler
}
