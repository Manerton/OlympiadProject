package strategy

import (
	"errors"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"strings"
)

// AggregationStrategy интерфейс для всех стратегий агрегации
type AggregationStrategy interface {
	Aggregate(services []config.Target, origReq *http.Request) (*responcetypes.ApiResponse, error)
}

func buildTargetURL(targetURL string, id string) string {
	if strings.Contains(targetURL, "{id}") {
		return strings.Replace(targetURL, "{id}", id, 1)
	}
	// Склеиваем без двойных слэшей
	return strings.TrimSuffix(targetURL, "/") + "/" + strings.TrimPrefix(id, "/")
}

// Извлечь первый сегмент пути после префикса.
// prefix: может быть "/ApplicationEvent/" или "/ApplicationEvent"
func extractIDFromPath(r *http.Request, prefix string) (string, error) {
	path := r.URL.Path
	// Нормализуем префикс: убираем завершающий слэш
	normPrefix := strings.TrimSuffix(prefix, "/")

	// Если path имеет префикс (когда StripPrefix НЕ использовался в main.go)
	if after, ok := strings.CutPrefix(path, normPrefix); ok {
		path = after
	}

	// Убираем ведущие "/"
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		return "", errors.New("id not found in path")
	}

	// id — первый сегмент до следующего "/"
	parts := strings.SplitN(path, "/", 2)
	return parts[0], nil
}
