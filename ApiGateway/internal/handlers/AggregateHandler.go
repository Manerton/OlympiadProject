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
	case "/history-event/":
		return strategy.NewHistoryEventStrategy(5 * time.Second)
	case "/events-appeal/":
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

			originalPath := req.URL.Path
			path := strings.TrimPrefix(originalPath, route.Prefix)

			log.Printf("Proxy Debug - Original Path: %q, Prefix: %q, Trimmed Path: %q", originalPath, route.Prefix, path)

			// Автоматическое определение mode по target.RawQuery
			isQueryMode := target.RawQuery != ""
			if path != "" && isQueryMode {
				// Query mode: базовый путь без изменений, param из path
				req.URL.Path = target.Path

				paramValue := strings.TrimPrefix(path, "/") // Значение без слеша

				// Парсим и мержим query из target
				q := req.URL.Query()
				if target.RawQuery != "" {
					targetQ, _ := url.ParseQuery(target.RawQuery)
					// Добавляем все params из target (кроме первого, если нужно перезаписать; здесь мержим все)
					for k, v := range targetQ {
						for _, vv := range v {
							q.Add(k, vv) // Add, чтобы не потерять множественные значения
						}
					}
					// Берём первый ключ из targetQ и перезаписываем его значением из path
					if len(targetQ) > 0 {
						firstKey := "" // Или итерацией: for k := range targetQ { firstKey = k; break }
						for k := range targetQ {
							firstKey = k
							break
						}
						q.Set(firstKey, paramValue) // Set перезаписывает
						log.Printf("Proxy Debug - Query Mode: Target Path: %q, First Param: %q=%q, Final Query: %q",
							target.Path, firstKey, paramValue, q.Encode())
					}
				}
				req.URL.RawQuery = q.Encode()
			} else {
				// Path mode: appending в путь
				targetPath := target.Path
				joinedPath := singleJoiningSlash(targetPath, path)
				req.URL.Path = joinedPath

				// Мержим query из target (если есть, но в path mode редко)
				if target.RawQuery != "" {
					q := req.URL.Query()
					targetQ, _ := url.ParseQuery(target.RawQuery)
					for k, v := range targetQ {
						for _, vv := range v {
							q.Add(k, vv)
						}
					}
					req.URL.RawQuery = q.Encode()
				}

				log.Printf("Proxy Debug - Path Mode: Target Path: %q, Joined Path: %q", targetPath, joinedPath)
			}

			req.Host = target.Host
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.Header.Set("X-Forwarded-For", req.RemoteAddr)

			modeStr := "path"
			if isQueryMode {
				modeStr = "query"
			}
			log.Printf("Proxying: %s %s?%s (auto-mode=%q, from target=%q + trimmed=%q) to %s",
				req.Method, req.URL.Path, req.URL.RawQuery, modeStr, target.Path, path, target.String())
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
