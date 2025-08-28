package strategy

import (
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
)

// AggregationStrategy интерфейс для всех стратегий агрегации
type AggregationStrategy interface {
	Aggregate(services []config.Target, origReq *http.Request) (*responcetypes.ApiResponse, error)
}
