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
	"sort"
	"strings"
	"time"
)

const timeSecond = 5 * time.Second

func getStrategy(prefix string) strategy.AggregationStrategy {
	switch prefix {
	case "/aplicationevent":
		return strategy.NewDefaultAggregationStrategy(timeSecond)
	case "/ApplicationEvent/":
		return strategy.NewApplicationEventStrategy(timeSecond)
	case "/history-event/":
		return strategy.NewHistoryEventStrategy(timeSecond)
	case "/events-appeal/":
		return strategy.NewApprovedApplicationEventStrategy(timeSecond)
	case "/jury-names/":
		return strategy.NewJuryNamesStrategy(timeSecond)
	case "/available-event/":
		return strategy.NewAvailableClassStrategy(timeSecond)
	case "/applications/create/":
		return strategy.NewApplicationCreateStrategy(timeSecond)
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

// Вспомогательная функция: убираем пустые строки из slice
func filterEmpty(s []string) []string {
	var result []string
	for _, str := range s {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
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
				// Query mode: базовый путь без изменений, params из path
				req.URL.Path = target.Path

				// Парсим target query на ключи (в порядке появления)
				targetQ, _ := url.ParseQuery(target.RawQuery)
				paramKeys := make([]string, 0, len(targetQ))
				for k := range targetQ {
					paramKeys = append(paramKeys, k)
				}
				// Сортируем для стабильности (опционально; если порядок не важен, убери sort)
				sort.Strings(paramKeys)

				// Парсим path на сегменты (e.g., "/3/abc" → []{"3", "abc"})
				pathSegments := strings.Split(strings.TrimPrefix(path, "/"), "/")
				pathSegments = filterEmpty(pathSegments) // Функция ниже, чтобы убрать пустые

				log.Printf("Proxy Debug - Query Mode: Param Keys: %v, Path Segments: %v", paramKeys, pathSegments)

				// Мержим query из target
				q := req.URL.Query()
				for k, v := range targetQ {
					for _, vv := range v {
						q.Add(k, vv) // Добавляем дефолтные значения
					}
				}

				// Маппим сегменты на ключи по порядку
				for i, key := range paramKeys {
					if i < len(pathSegments) {
						paramValue := pathSegments[i]
						q.Set(key, paramValue) // Set перезаписывает
						log.Printf("Proxy Debug - Set Param: %q = %q", key, paramValue)
					}
				}

				req.URL.RawQuery = q.Encode()
				log.Printf("Proxy Debug - Target Path: %q, Final Query: %q", target.Path, req.URL.RawQuery)
			} else {
				// Path mode: без изменений
				targetPath := target.Path
				joinedPath := singleJoiningSlash(targetPath, path)
				req.URL.Path = joinedPath

				// Мерж query (редко в path mode)
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
