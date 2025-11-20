package strategy

import (
	"fmt"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"time"
)

type AvailableClassStrategy struct {
	timeout time.Duration
}

func NewAvailableClassStrategy(timeout time.Duration) *AvailableClassStrategy {
	return &AvailableClassStrategy{
		timeout: timeout,
	}
}

func (s *AvailableClassStrategy) Aggregate(targets []config.Target, origReq *http.Request) (*responcetypes.ApiResponse, error) {
	client := &http.Client{Timeout: s.timeout}

	if len(targets) < 2 {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "need at least 2 targets for chained aggregation",
		}, fmt.Errorf("bad config: expected 2 targets, got %d", len(targets))
	}

	parentEventId, err := extractIDFromPath(origReq, "/TODO!/")
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Error:      "user id not provided in path",
		}, err
	}

	_, _ = client, parentEventId

	return &responcetypes.ApiResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
		Data:       "TODO!!",
	}, nil
}
