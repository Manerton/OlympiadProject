package strategy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"time"
)

type HistoryEventStrategy struct {
	timeout time.Duration
}

func NewHistoryEventStrategy(timeout time.Duration) *HistoryEventStrategy {
	return &HistoryEventStrategy{timeout: timeout}
}

func (s *HistoryEventStrategy) Aggregate(
	targets []config.Target,
	origReq *http.Request,
) (*responcetypes.ApiResponse, error) {
	client := &http.Client{Timeout: s.timeout}

	if len(targets) < 2 {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "need at least 2 targets for chained aggregation",
		}, fmt.Errorf("bad config: expected 2 targets, got %d", len(targets))
	}

	// === 1. первый сервис: /result/events-by-user/{id} ===
	userID, err := extractIDFromPath(origReq, "/result/events-by-user/")
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Error:      "user id not provided in path",
		}, err
	}

	firstURL := buildTargetURL(targets[0].URL, userID)

	req, _ := http.NewRequest("GET", firstURL, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      fmt.Sprintf("call to first service failed: %v", err),
		}, err
	}
	defer resp.Body.Close()

	var firstResp responcetypes.ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&firstResp); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "invalid JSON from first service",
		}, err
	}

	// предполагаем, что Data = []string (массив eventIDs)
	eventIDs, ok := firstResp.Data.([]interface{})
	if !ok {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "unexpected format from first service (expected []string)",
		}, fmt.Errorf("unexpected format: %#v", firstResp.Data)
	}

	// === 2. второй сервис: /api/events/list (POST) ===
	body, _ := json.Marshal(map[string]interface{}{"ids": eventIDs})
	secondURL := targets[1].URL

	req2, _ := http.NewRequest("POST", secondURL, bytes.NewBuffer(body))
	req2.Header = origReq.Header.Clone()
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req2)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      fmt.Sprintf("call to second service failed: %v", err),
		}, err
	}
	defer resp2.Body.Close()

	var secondResp responcetypes.ApiResponse
	if err := json.NewDecoder(resp2.Body).Decode(&secondResp); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "invalid JSON from second service",
		}, err
	}

	// возвращаем уже готовый ответ от второго сервиса
	return &responcetypes.ApiResponse{
		Status:     "success",
		StatusCode: 200,
		Data:       secondResp.Data,
	}, nil
}
