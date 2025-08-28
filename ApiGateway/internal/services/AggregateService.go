package services

import (
	"main/internal/config"
	"main/internal/responcetypes"
	"main/internal/strategy"
	"net/http"
	"time"
)

// AggregateService теперь использует стратегии
type AggregateService struct {
	defaultStrategy strategy.AggregationStrategy
}

func NewAggregateService(timeout time.Duration) *AggregateService {
	return &AggregateService{
		defaultStrategy: strategy.NewDefaultAggregationStrategy(timeout),
	}
}

// Aggregate использует стратегию по умолчанию (можно переопределить)
func (s *AggregateService) Aggregate(
	route config.Route,
	origReq *http.Request,
	strategy strategy.AggregationStrategy, // опциональная стратегия
) (*responcetypes.ApiResponse, error) {
	aggStrategy := s.defaultStrategy
	if strategy != nil {
		aggStrategy = strategy
	}

	return aggStrategy.Aggregate(route.Targets, origReq)
}
