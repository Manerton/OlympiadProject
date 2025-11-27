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

type ApplicationCreateStrategy struct {
	timeout time.Duration
}

func NewApplicationCreateStrategy(timeout time.Duration) *ApplicationCreateStrategy {
	return &ApplicationCreateStrategy{timeout: timeout}
}

// Входящий запрос: POST /applications/create/
// Body: { "userId": "...", "eventID": "..." }
func (s *ApplicationCreateStrategy) Aggregate(
	targets []config.Target,
	origReq *http.Request,
) (*responcetypes.ApiResponse, error) {

	client := &http.Client{Timeout: s.timeout}

	if len(targets) != 2 {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Error:      "expected exactly 2 targets",
		}, nil
	}

	// ========== 1. Читаем тело оригинального запроса ==========
	var incomingBody struct {
		UserID             string `json:"userId"`
		EventID            string `json:"eventID"`
		SchoolID           string `json:"schoolID"`
		Profile            string `json:"profile"`
		ClassParticipation int    `json:"class_participation"`
	}
	if err := json.NewDecoder(origReq.Body).Decode(&incomingBody); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Error:      "invalid request body",
		}, err
	}
	if incomingBody.UserID == "" || incomingBody.EventID == "" {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Error:      "userId and eventID are required",
		}, nil
	}

	userID := incomingBody.UserID
	eventID := incomingBody.EventID

	// ========== 2. Запрос к SSO: получаем participant и school_id ==========
	// ssoURL := targets[0].URL + userID // например http://sso-service:8181/api/participants/byuser/ + uuid

	// req1, err := http.NewRequest("GET", ssoURL, nil)
	// if err != nil {
	// 	return nil, err
	// }
	// req1.Header = origReq.Header.Clone() // пробрасываем Authorization и прочие заголовки

	// resp1, err := client.Do(req1)
	// if err != nil {
	// 	return &responcetypes.ApiResponse{
	// 		Status:     "error",
	// 		StatusCode: http.StatusBadGateway,
	// 		Error:      fmt.Sprintf("failed to call SSO service: %v", err),
	// 	}, err
	// }
	// defer resp1.Body.Close()

	// if resp1.StatusCode >= 400 {
	// 	return &responcetypes.ApiResponse{
	// 		Status:     "error",
	// 		StatusCode: resp1.StatusCode,
	// 		Error:      "SSO service returned error",
	// 	}, nil
	// }

	// var ssoResp struct {
	// 	Data map[string]interface{} `json:"data"`
	// }
	// if err := json.NewDecoder(resp1.Body).Decode(&ssoResp); err != nil {
	// 	return &responcetypes.ApiResponse{
	// 		Status:     "error",
	// 		StatusCode: http.StatusInternalServerError,
	// 		Error:      "failed to parse SSO response",
	// 	}, err
	// }

	// schoolID, ok := ssoResp.Data["school_id"].(string)
	// if !ok || schoolID == "" {
	// 	return &responcetypes.ApiResponse{
	// 		Status:     "error",
	// 		StatusCode: http.StatusUnprocessableEntity,
	// 		Error:      "school_id not found in participant data",
	// 	}, nil
	// }

	// ========== 3. Запрос к Application Service: создаём заявку ==========
	appPayload := map[string]any{
		"userId":              userID,
		"eventID":             eventID,
		"schoolId":            incomingBody.SchoolID,
		"profile":             incomingBody.Profile,
		"class_participation": incomingBody.ClassParticipation,
	}
	payloadBytes, _ := json.Marshal(appPayload)

	req2, err := http.NewRequest("POST", targets[1].URL, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	req2.Header = origReq.Header.Clone()
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := client.Do(req2)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadGateway,
			Error:      fmt.Sprintf("failed to call application-service: %v", err),
		}, err
	}
	defer resp2.Body.Close()

	// Пробрасываем ответ от application-service клиенту
	var finalBody bytes.Buffer
	_, err = finalBody.ReadFrom(resp2.Body)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Error:      "failed to read response body",
		}, err
	}

	return &responcetypes.ApiResponse{
		Status:     "success",
		StatusCode: resp2.StatusCode,
		Data:       json.RawMessage(finalBody.Bytes()),
	}, nil
}
