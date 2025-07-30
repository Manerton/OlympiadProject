package handlers

import (
	"encoding/json"
	"log"
	"main/internal/config"          // ← добавить
	"main/internal/middleware/auth" // ← убедиться, что путь корректен
	"main/internal/responcetypes"
	"main/internal/services"
	"net/http" // ← добавить
	"net/http/httputil"
	"net/url"
	"strings"
)

// NewAggregateHandler возвращает http.Handler, который:
// 1) валидирует JWT (если нужно),
// 2) проверяет роль (если указаны route.Roles),
// 3) вызывает AggregateService для сбора данных,
// 4) оборачивает в ApiResponse и отдаёт JSON.
func NewAggregateHandler(route config.Route, jwtKey string) http.Handler {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dataSlice, err := services.NewAggregateService().Aggregate(route, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Упаковываем только поле Data
		resp := responcetypes.ApiResponse{
			Status:     "ok",
			StatusCode: http.StatusOK,
			Data:       dataSlice,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	if !route.SkipAuth {
		handler = auth.AuthenticateMiddleware(jwtKey, handler)
		if len(route.Roles) > 0 {
			handler = auth.RoleBasedAccess(route.Roles...)(handler)
		}
	}

	return handler
}

// NewProxyHandler проксирует все запросы на route.Target и оборачивает JWT‑миддлварью
func NewProxyHandler(route config.Route, jwtKey string) http.Handler {
	// разбираем URL таргета
	target, err := url.Parse(route.Target)
	if err != nil {
		log.Fatalf("bad target URL %s: %v", route.Target, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// базовый handler
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// убираем префикс из пути
		r.URL.Path = strings.TrimPrefix(r.URL.Path, route.Prefix)
		log.Printf("Proxying to %s%s", target.Host, r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	// JWT‑проверка, если нужно
	if !route.SkipAuth {
		handler = auth.AuthenticateMiddleware(jwtKey, handler)
	}

	return handler
}
