package strategy

import (
	"encoding/json"
	"fmt"
	"main/internal/config"
	"main/internal/lib/errs"
	"main/internal/middleware/auth"
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
			Status:     errs.StatusError,
			StatusCode: http.StatusInternalServerError,
			Error:      "need at least 2 targets for chained aggregation",
		}, fmt.Errorf("bad config: expected 2 targets, got %d", len(targets))
	}

	parentEventId, err := extractIDFromPath(origReq, "/available-event/")
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     errs.StatusError,
			StatusCode: http.StatusBadRequest,
			Error:      "user id not provided in path",
		}, err
	}

	// Обращение к sso для получениея participant
	userInfo := origReq.Context().Value(auth.UserInfoKey{}).(auth.UserInfo)
	ssoTarget := targets[0].URL
	ssoURL := buildTargetURL(ssoTarget, userInfo.ID.String())

	req, err := http.NewRequest(http.MethodGet, ssoURL, nil)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     errs.StatusError,
			StatusCode: http.StatusInternalServerError,
			Error:      "failed create request",
		}, err
	}
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("call to applications failed: %v", err),
		}, err
	}
	defer resp.Body.Close()

	appResponse := responcetypes.ApiResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&appResponse); err != nil {
		return &responcetypes.ApiResponse{
			Status:     errs.StatusError,
			StatusCode: http.StatusBadRequest,
			Error:      "failed decode body",
		}, err
	}

	raw, _ := json.Marshal(appResponse.Data)

	var participant responcetypes.Participant
	if err := json.Unmarshal(raw, &participant); err != nil {
		return nil, fmt.Errorf("failed to parse participant: %w", err)
	}

	eventTarget := targets[1].URL
	eventURL := buildTargetURL(eventTarget, parentEventId) + fmt.Sprintf("?class=%s", participant.ClassNumber)
	req, err = http.NewRequest(http.MethodGet, eventURL, nil)
	req.Header = origReq.Header.Clone()

	resp, err = client.Do(req)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     errs.StatusError,
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Sprintf("call to applications failed: %v", err),
		}, nil
	}

	appResponse = responcetypes.ApiResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&appResponse); err != nil {
		return &responcetypes.ApiResponse{
			Status:     errs.StatusError,
			StatusCode: http.StatusBadRequest,
			Error:      "failed decode body",
		}, err
	}

	return &responcetypes.ApiResponse{
		Status:     errs.StatusSuccess,
		StatusCode: http.StatusOK,
		Data:       appResponse.Data,
	}, nil
}
