package services

import (
	"encoding/json"
	"fmt"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"time"
)

// AggregateService отвечает за обращение к бэкендам и объединение ответов.
type AggregateService struct{}

// NewAggregateService конструктор
func NewAggregateService() *AggregateService {
	return &AggregateService{}
}

// Aggregate делает GET‑запросы к каждому URL в route.Services,
// передаёт туда токен от original request, собирает все ответы в слайс.
func (s *AggregateService) Aggregate(
	route config.Route,
	origReq *http.Request,
) ([]interface{}, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	var aggregated []interface{}

	for _, svcURL := range route.Services {
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
