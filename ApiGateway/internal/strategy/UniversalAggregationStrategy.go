package strategy

import (
	"encoding/json"
	"fmt"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"time"
)

type UniversalAggregationStrategy struct {
	timeout time.Duration
}

func NewUniversalAggregationStrategy(timeout time.Duration) *UniversalAggregationStrategy {
	return &UniversalAggregationStrategy{timeout: timeout}
}

func (s *UniversalAggregationStrategy) Aggregate(
	targets []config.Target, // теперь не []string, а []Target
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

		// Если для этого target заданы поля — фильтруем
		if len(target.Fields) > 0 {
			if dataMap, ok := svcResp.Data.(map[string]interface{}); ok {
				filtered := make(map[string]interface{})
				for _, f := range target.Fields {
					if val, exists := dataMap[f]; exists {
						filtered[f] = val
					}
				}
				aggregated = append(aggregated, filtered)
				continue
			}
		}

		// иначе добавляем весь объект
		aggregated = append(aggregated, svcResp.Data)
	}
	resp.Data = aggregated
	return &resp, nil
}
