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

func getStrategy(prefix string) strategy.AggregationStrategy {
	switch prefix {
	case "/aplicationevent":
		return strategy.NewDefaultAggregationStrategy(5 * time.Second)
	case "/ApplicationEvent/":
		return strategy.NewApplicationEventStrategy(5 * time.Second)
	case "/result/events-by-user/":
		return strategy.NewHistoryEventStrategy(5 * time.Second)
	default:
		return nil
	}
}

func NewAggregateHandler(route config.Route, jwtKey string) http.Handler {
	service := services.NewAggregateService(5 * time.Second)
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var aggrstrategy = getStrategy(route.Prefix)
		dataSlice, err := service.Aggregate(route, r, aggrstrategy)
		if err != nil {
			log.Printf("Aggregation failed: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

func NewProxyHandler(route config.Route, jwtKey string) http.Handler {
	target, err := url.Parse(route.Target)
	if err != nil {
		log.Fatalf("bad target URL %s: %v", route.Target, err)
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			// Извлекаем оставшуюся часть пути после префикса
			path := strings.TrimPrefix(req.URL.Path, route.Prefix)
			// Формируем новый путь, добавляя оставшуюся часть к target.Path
			req.URL.Path = singleJoiningSlash(target.Path, path)
			req.Host = target.Host
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.Header.Set("X-Forwarded-For", req.RemoteAddr)
			log.Printf("Proxying: %s %s to %s", req.Method, req.URL.Path, target.String())
		},
		ModifyResponse: func(resp *http.Response) error {
			// Удаляем любые CORS-заголовки, возвращаемые микросервисом
			log.Printf("Proxy Response: %d %s from %s", resp.StatusCode, resp.Request.URL.Path, target.Host)
			resp.Header.Del("Access-Control-Allow-Origin")
			resp.Header.Del("Access-Control-Allow-Methods")
			resp.Header.Del("Access-Control-Allow-Headers")
			resp.Header.Del("Access-Control-Expose-Headers")
			resp.Header.Del("Access-Control-Allow-Credentials")
			resp.Header.Del("Access-Control-Max-Age")
			// for i, c := range resp.Header["Set-Cookie"] {
			// 	// убираем Domain, чтобы браузер сохранял куку на gateway
			// 	c = strings.ReplaceAll(c, "Domain=172.16.1.39", "")
			// 	c = strings.ReplaceAll(c, "Domain=172.16.0.196", "")

			// 	// ставим SameSite=Lax вместо None (чтобы не требовал Secure при HTTP)
			// 	c = strings.ReplaceAll(c, "SameSite=None", "SameSite=Lax")

			// 	resp.Header["Set-Cookie"][i] = c
			// }
			return nil
		},
	}

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxying to %s%s", target.Host, r.URL.Path)
		proxy.ServeHTTP(w, r)
	})

	if !route.SkipAuth {
		handler = auth.AuthenticateMiddleware(jwtKey, handler)
	}

	return handler
}
