package strategy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"time"
)

type JuryNamesStrategy struct {
	timeout time.Duration
}

func NewJuryNamesStrategy(timeout time.Duration) *JuryNamesStrategy {
	return &JuryNamesStrategy{timeout: timeout}
}

func (s *JuryNamesStrategy) Aggregate(
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

	// === 1. первый сервис: /jury-assignments/event/{id} ===
	eventID, err := extractIDFromPath(origReq, "/jury-names/")
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Error:      "event id not provided in path",
		}, err
	}

	appTarget := targets[0].URL
	appURL := buildTargetURL(appTarget, eventID)

	req, _ := http.NewRequest("GET", appURL, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      fmt.Sprintf("call to jury failed: %v", err),
		}, err
	}
	defer resp.Body.Close()

	var appResp responcetypes.ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&appResp); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "invalid JSON from jury",
		}, err
	}

	jury, ok := appResp.Data.([]interface{})
	if !ok {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "unexpected jury format",
		}, fmt.Errorf("unexpected format in jury")
	}

	var assignedIDs []string
	var userIDs []string
	for _, raw := range jury {
		jur := raw.(map[string]interface{})
		aid := jur["id"].(string)
		uid := jur["user_id"].(string)
		assignedIDs = append(assignedIDs, aid)
		userIDs = append(userIDs, uid)
	}
	/////////////////до этой части все ок!!!///////////////////////////////

	// === 2. второй сервис: /api/users/list (POST) ===
	body, _ := json.Marshal(map[string][]string{"ids": userIDs})
	ssoURL := targets[1].URL
	req2, _ := http.NewRequest("POST", ssoURL, bytes.NewBuffer(body))
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
	rawBody, _ := io.ReadAll(resp2.Body)
	resp2.Body = io.NopCloser(bytes.NewBuffer(rawBody))
	var userResp responcetypes.ApiResponse
	if err := json.NewDecoder(resp2.Body).Decode(&userResp); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "invalid JSON from events",
		}, err
	}
	users, ok := userResp.Data.([]interface{})
	if !ok {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      "unexpected events format",
		}, fmt.Errorf("unexpected format in events")
	}

	// === 3. собираем JuryNames ===
	userMap := make(map[string]map[string]interface{})
	for _, raw := range users {
		u := raw.(map[string]interface{})
		uid := getString(u["id"])
		userMap[uid] = u
	}

	var result []responcetypes.JuryNames
	for i, uid := range userIDs {
		assignID := assignedIDs[i]
		if u, ok := userMap[uid]; ok {
			// ФИО
			name := fmt.Sprintf("%s %s %s",
				getString(u["surname"]),
				getString(u["firstname"]),
				getString(u["patronymic"]),
			)

			result = append(result, responcetypes.JuryNames{
				ID:     assignID,
				UserId: uid,
				Name:   name,
			})
		}
	}

	return &responcetypes.ApiResponse{
		Status:     "success",
		StatusCode: 200,
		Data:       result,
	}, nil

}
