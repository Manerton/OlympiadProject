package strategy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/internal/config"
	"main/internal/lib/errs"
	"main/internal/responcetypes"
	"net/http"
	"time"
)

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

	raw, _ := json.Marshal(appResp.Data)

	var applications []responcetypes.ApplicationResponse
	if err := json.Unmarshal(raw, &applications); err != nil {
		return &responcetypes.ApiResponse{
			Status:     errs.StatusError,
			StatusCode: http.StatusBadRequest,
			Error:      "failed decode body",
		}, err
	}

	appIds := make(map[string]string)
	eventStatuses := make(map[string]int)
	eventProfile := make(map[string]string)
	eventClassParticipation := make(map[string]int)
	var eventIDs []string
	for _, raw := range applications {
		eventIDs = append(eventIDs, raw.EventID)
		appIds[raw.EventID] = raw.ID
		eventStatuses[raw.EventID] = raw.Status
		eventProfile[raw.EventID] = raw.Profile
		eventClassParticipation[raw.EventID] = raw.ClassParticipation
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

	raw, _ = json.Marshal(evResp.Data)
	var events []responcetypes.Event
	if err := json.Unmarshal(raw, &events); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "invalid JSON from events",
		}, err
	}

	// === 3. собираем ApplicationEvent ===
	var result []responcetypes.ApplicationEvent
	for _, event := range events {

		appEv := responcetypes.ApplicationEvent{
			ApplicationID:      appIds[event.ID],
			ID:                 event.ID,
			Name:               event.Name,
			Subject:            event.Subject,
			Dates:              event.Dates,
			Profile:            eventProfile[event.ID],
			Status:             eventStatuses[event.ID],
			ClassParticipation: eventClassParticipation[event.ID],
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
