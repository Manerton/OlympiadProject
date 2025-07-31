package strategy

import (
	"encoding/json"
	"fmt"
	"main/internal/responcetypes"
	"net/http"
	"time"
)

// AggregationStrategy интерфейс для всех стратегий агрегации
type AggregationStrategy interface {
	Aggregate(services []string, origReq *http.Request) ([]interface{}, error)
	
}

// DefaultAggregationStrategy текущая реализация агрегации
type DefaultAggregationStrategy struct {
	timeout time.Duration
}

func NewDefaultAggregationStrategy(timeout time.Duration) *DefaultAggregationStrategy {
	return &DefaultAggregationStrategy{timeout: timeout}
}

func (s *DefaultAggregationStrategy) Aggregate(
	services []string,
	origReq *http.Request,
) ([]interface{}, error) {
	client := &http.Client{Timeout: s.timeout}
	var aggregated []interface{}

	for _, svcURL := range services {
		req, _ := http.NewRequest(origReq.Method, svcURL, nil)
		req.Header = origReq.Header.Clone()

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("call to %s failed: %w", svcURL, err)
		}
		defer resp.Body.Close()

		var svcResp responcetypes.ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&svcResp); err != nil {
			return nil, fmt.Errorf("invalid JSON from %s: %w", svcURL, err)
		}
		if svcResp.StatusCode != 200 {
			return nil, fmt.Errorf("service %s error: %s", svcURL, svcResp.Error)
		}

		aggregated = append(aggregated, svcResp.Data)
	}
	return aggregated, nil
}
