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

type ApprovedApplicationEventStrategy struct {
	timeout time.Duration
}

func NewApprovedApplicationEventStrategy(timeout time.Duration) *ApprovedApplicationEventStrategy {
	return &ApprovedApplicationEventStrategy{timeout: timeout}
}

func (s *ApprovedApplicationEventStrategy) Aggregate(
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
	/////////ЭТО ГОВНОКОД НУЖНО ИСПРАВИТЬ
	eventStatuses := make(map[string]int)
	var eventIDs []string
	for _, raw := range applications {
		app := raw.(map[string]interface{})
		eid := app["eventId"].(string)
		/////////ЭТО ГОВНОКОД НУЖНО ИСПРАВИТЬ
		status := app["Status"].(string)
		if status == "2" {
			eventIDs = append(eventIDs, eid)
			eventStatuses[eid] = int(app["status"].(float64))
		}
		/////////ЭТО ГОВНОКОД НУЖНО ИСПРАВИТЬ

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
	var result []responcetypes.ApplicationEvent
	for _, raw := range events {
		ev := raw.(map[string]interface{})

		appEv := responcetypes.ApplicationEvent{
			ID:             getString(ev["id"]),
			Name:           getString(ev["name"]),
			StartDate:      getString(ev["start_date"]),
			EndDate:        getString(ev["end_date"]),
			PreviousEvent:  getString(ev["previous_event_id"]),
			Subject:        int(ev["subject"].(float64)),
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
