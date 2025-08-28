package strategy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"strings"
	"time"
)

type ApplicationEvent struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	PreviousEvent  string `json:"previous_event_id"`
	Subject        string `json:"subject"`
	ClassNumber    int    `json:"class_number"`
	AdditionalInfo string `json:"additional_info"`
	Status         int    `json:"status"`
}

// Извлечь первый сегмент пути после префикса.
// prefix: может быть "/ApplicationEvent/" или "/ApplicationEvent"
func extractIDFromPath(r *http.Request, prefix string) (string, error) {
	path := r.URL.Path
	// Нормализуем префикс: убираем завершающий слэш
	normPrefix := strings.TrimSuffix(prefix, "/")

	// Если path имеет префикс (когда StripPrefix НЕ использовался в main.go)
	if strings.HasPrefix(path, normPrefix) {
		path = strings.TrimPrefix(path, normPrefix)
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

// Построить итоговый URL для target (подставить id).
// Если в targetURL есть {id} — заменить. Иначе аккуратно присоединить id.
func buildTargetURL(targetURL string, id string) string {
	if strings.Contains(targetURL, "{id}") {
		return strings.Replace(targetURL, "{id}", id, 1)
	}
	// Склеиваем без двойных слэшей
	return strings.TrimSuffix(targetURL, "/") + "/" + strings.TrimPrefix(id, "/")
}

type ApplicationEventStrategy struct {
	timeout time.Duration
}

func NewApplicationEventStrategy(timeout time.Duration) *ApplicationEventStrategy {
	return &ApplicationEventStrategy{timeout: timeout}
}

func (s *ApplicationEventStrategy) Aggregate(
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

	// === 1. первый сервис: /applications/user/{id} ===
	userID, err := extractIDFromPath(origReq, "/ApplicationEvent/")
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Error:      "user id not provided in path",
		}, err
	}

	appTarget := targets[0].URL
	appURL := buildTargetURL(appTarget, userID)

	req, _ := http.NewRequest("GET", appURL, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      fmt.Sprintf("call to applications failed: %v", err),
		}, err
	}
	defer resp.Body.Close()

	var appResp responcetypes.ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&appResp); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "invalid JSON from applications",
		}, err
	}

	applications, ok := appResp.Data.([]interface{})
	if !ok {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "unexpected applications format",
		}, fmt.Errorf("unexpected format in applications")
	}

	eventStatuses := make(map[string]int)
	var eventIDs []string
	for _, raw := range applications {
		app := raw.(map[string]interface{})
		eid := app["eventId"].(string)
		eventIDs = append(eventIDs, eid)
		eventStatuses[eid] = int(app["status"].(float64))
	}

	// === 2. второй сервис: /api/events/list (POST) ===
	body, _ := json.Marshal(map[string][]string{"ids": eventIDs})
	eventsURL := targets[1].URL

	req2, _ := http.NewRequest("POST", eventsURL, bytes.NewBuffer(body))
	req2.Header = origReq.Header.Clone()
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req2)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      fmt.Sprintf("call to events failed: %v", err),
		}, err
	}
	defer resp2.Body.Close()

	var evResp responcetypes.ApiResponse
	if err := json.NewDecoder(resp2.Body).Decode(&evResp); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "invalid JSON from events",
		}, err
	}

	events, ok := evResp.Data.([]interface{})
	if !ok {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "unexpected events format",
		}, fmt.Errorf("unexpected format in events")
	}

	// === 3. собираем ApplicationEvent ===
	var result []ApplicationEvent
	for _, raw := range events {
		ev := raw.(map[string]interface{})

		appEv := ApplicationEvent{
			ID:             getString(ev["id"]),
			Name:           getString(ev["name"]),
			StartDate:      getString(ev["start_date"]),
			EndDate:        getString(ev["end_date"]),
			PreviousEvent:  getString(ev["previous_event_id"]),
			Subject:        getString(ev["subject"]),
			ClassNumber:    int(ev["class_number"].(float64)),
			AdditionalInfo: getString(ev["additional_info"]),
			Status:         eventStatuses[getString(ev["id"])],
		}
		result = append(result, appEv)
	}

	return &responcetypes.ApiResponse{
		Status:     "success",
		StatusCode: 200,
		Data:       result,
	}, nil
}

func getString(v interface{}) string {
	if v == nil {
		return ""
	}
	return v.(string)
}
