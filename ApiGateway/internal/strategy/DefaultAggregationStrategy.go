package strategy

import (
	"encoding/json"
	"fmt"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"time"
)

// DefaultAggregationStrategy — возвращает все данные без фильтрации
type DefaultAggregationStrategy struct {
	timeout time.Duration
}

func NewDefaultAggregationStrategy(timeout time.Duration) *DefaultAggregationStrategy {
	return &DefaultAggregationStrategy{timeout: timeout}
}

func (s *DefaultAggregationStrategy) Aggregate(
	targets []config.Target,
	origReq *http.Request,
) (*responcetypes.ApiResponse, error) {
	client := &http.Client{Timeout: s.timeout}
	var aggregated []interface{}
	var resp responcetypes.ApiResponse
	for _, target := range targets {
		req, _ := http.NewRequest(origReq.Method, target.URL, nil)
		req.Header = origReq.Header.Clone()

		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("call to %s failed: %w", target.URL, err)
		}
		defer resp.Body.Close()

		var svcResp responcetypes.ApiResponse
		if err := json.NewDecoder(resp.Body).Decode(&svcResp); err != nil {
			return nil, fmt.Errorf("invalid JSON from %s: %w", target.URL, err)
		}
		if svcResp.StatusCode != 200 {
			return nil, fmt.Errorf("service %s error: %s", target.URL, svcResp.Error)
		}

		aggregated = append(aggregated, svcResp.Data)
	}
	resp.Data = aggregated
	return &resp, nil
}
