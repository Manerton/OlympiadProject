package strategy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/internal/config"
	"main/internal/responcetypes"
	"net/http"
	"sync"
	"time"
)

type VerifyApplicationsStrategy struct {
	timeout time.Duration
}

func NewVerifyApplicationsStrategy(timeout time.Duration) *VerifyApplicationsStrategy {
	return &VerifyApplicationsStrategy{timeout: timeout}
}

func (s *VerifyApplicationsStrategy) Aggregate(
	targets []config.Target,
	origReq *http.Request,
) (*responcetypes.ApiResponse, error) {
	client := &http.Client{Timeout: s.timeout}

	// ========== 1. Читаем тело оригинального запроса ==========
	var incomingBody struct {
		ID   string `json:"id"`
		Role int    `json:"role"`
	}
	if err := json.NewDecoder(origReq.Body).Decode(&incomingBody); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Error:      "invalid request body",
		}, err
	}

	// ========== 2. Получаем список школ ==========
	var schoolIDs []string

	if incomingBody.Role == 0 {
		// ROLE = 0: делаем запрос к первому таргету для получения школ района
		districtURL := targets[0].URL + "?id=" + incomingBody.ID
		req1, err := http.NewRequest("GET", districtURL, nil)
		if err != nil {
			return nil, err
		}
		req1.Header = origReq.Header.Clone()

		resp1, err := client.Do(req1)
		if err != nil {
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusBadGateway,
				Error:      fmt.Sprintf("failed to call district service: %v", err),
			}, err
		}
		defer resp1.Body.Close()

		if resp1.StatusCode >= 400 {
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: resp1.StatusCode,
				Error:      "district service returned error",
			}, nil
		}

		var districtResp struct {
			Schools []struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		if err := json.NewDecoder(resp1.Body).Decode(&districtResp); err != nil {
			return &responcetypes.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusInternalServerError,
				Error:      "failed to parse district response",
			}, err
		}

		for _, school := range districtResp.Schools {
			schoolIDs = append(schoolIDs, school.ID)
		}
	} else {
		// ROLE != 0: используем переданный ID как школу
		schoolIDs = []string{incomingBody.ID}
	}

	if len(schoolIDs) == 0 {
		return &responcetypes.ApiResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Data:       json.RawMessage("[]"),
		}, nil
	}

	// ========== 3. Получаем заявки для школ ==========
	appPayload := map[string]any{
		"ids": schoolIDs,
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

	if resp2.StatusCode >= 400 {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: resp2.StatusCode,
			Error:      "application service returned error",
		}, nil
	}

	var applications []struct {
		ID                 string `json:"id"`
		UserID             string `json:"userId"`
		SchoolID           string `json:"schoolId"`
		EventID            string `json:"eventId"`
		Profile            string `json:"profile"`
		ClassParticipation int    `json:"class_participation"`
		Status             int    `json:"status"`
		Reason             int    `json:"reason"`
		Code               string `json:"code"`
		SubmittedAt        string `json:"submittedAt"`
		UpdatedAt          string `json:"updatedAt"`
	}

	if err := json.NewDecoder(resp2.Body).Decode(&applications); err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Error:      "failed to parse applications response",
		}, err
	}

	if len(applications) == 0 {
		return &responcetypes.ApiResponse{
			Status:     "success",
			StatusCode: http.StatusOK,
			Data:       json.RawMessage("[]"),
		}, nil
	}

	// ========== 4. Собираем уникальные UserID и EventID ==========
	userIDSet := make(map[string]bool)
	eventIDSet := make(map[string]bool)

	for _, app := range applications {
		userIDSet[app.UserID] = true
		eventIDSet[app.EventID] = true
	}

	userIDs := make([]string, 0, len(userIDSet))
	eventIDs := make([]string, 0, len(eventIDSet))

	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	for id := range eventIDSet {
		eventIDs = append(eventIDs, id)
	}

	// ========== 5. Параллельно запрашиваем пользователей и события ==========
	var users []struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		FirstName     string `json:"firstname"`
		Surname       string `json:"surname"`
		Patronymic    string `json:"patronymic"`
		PhoneNumber   string `json:"phone_number"`
		BirthDate     string `json:"birthdate"`
		Gender        int    `json:"gender"`
		Role          int    `json:"role"`
		Activated     bool   `json:"activated"`
		ParticipantID string `json:"participant_id"`
		Disability    string `json:"disability"`
		SchoolID      string `json:"school_id"`
		Citizenship   string `json:"citizenship"`
		ClassNumber   string `json:"class_number"`
	}

	var events []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var wg sync.WaitGroup
	var userErr, eventErr error

	// Запрос пользователей
	wg.Add(1)
	go func() {
		defer wg.Done()
		userPayload := map[string]any{"ids": userIDs}
		payloadBytes, _ := json.Marshal(userPayload)

		req3, err := http.NewRequest("POST", targets[2].URL, bytes.NewReader(payloadBytes))
		if err != nil {
			userErr = err
			return
		}
		req3.Header = origReq.Header.Clone()
		req3.Header.Set("Content-Type", "application/json")

		resp3, err := client.Do(req3)
		if err != nil {
			userErr = err
			return
		}
		defer resp3.Body.Close()

		if resp3.StatusCode >= 400 {
			userErr = fmt.Errorf("user service returned status %d", resp3.StatusCode)
			return
		}

		if err := json.NewDecoder(resp3.Body).Decode(&users); err != nil {
			userErr = err
		}
	}()

	// Запрос событий
	wg.Add(1)
	go func() {
		defer wg.Done()
		eventPayload := map[string]any{"ids": eventIDs}
		payloadBytes, _ := json.Marshal(eventPayload)

		req4, err := http.NewRequest("POST", targets[3].URL, bytes.NewReader(payloadBytes))
		if err != nil {
			eventErr = err
			return
		}
		req4.Header = origReq.Header.Clone()
		req4.Header.Set("Content-Type", "application/json")

		resp4, err := client.Do(req4)
		if err != nil {
			eventErr = err
			return
		}
		defer resp4.Body.Close()

		if resp4.StatusCode >= 400 {
			eventErr = fmt.Errorf("event service returned status %d", resp4.StatusCode)
			return
		}

		if err := json.NewDecoder(resp4.Body).Decode(&events); err != nil {
			eventErr = err
		}
	}()

	wg.Wait()

	if userErr != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadGateway,
			Error:      fmt.Sprintf("failed to get users: %v", userErr),
		}, nil
	}

	if eventErr != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusBadGateway,
			Error:      fmt.Sprintf("failed to get events: %v", eventErr),
		}, nil
	}

	// ========== 6. Создаем мапы для быстрого поиска ==========
	userMap := make(map[string]struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		FirstName     string `json:"firstname"`
		Surname       string `json:"surname"`
		Patronymic    string `json:"patronymic"`
		PhoneNumber   string `json:"phone_number"`
		BirthDate     string `json:"birthdate"`
		Gender        int    `json:"gender"`
		Role          int    `json:"role"`
		Activated     bool   `json:"activated"`
		ParticipantID string `json:"participant_id"`
		Disability    string `json:"disability"`
		SchoolID      string `json:"school_id"`
		Citizenship   string `json:"citizenship"`
		ClassNumber   string `json:"class_number"`
	})

	for _, user := range users {
		userMap[user.ID] = user
	}

	eventMap := make(map[string]string)
	for _, event := range events {
		eventMap[event.ID] = event.Name
	}

	// ========== 7. Формируем итоговый ответ ==========
	type ResultItem struct {
		ID           string `json:"id"`
		OlympiadName string `json:"olympiadName"`
		Category     string `json:"category"`
		Status       string `json:"status"`
		Surname      string `json:"surname"`
		FirstName    string `json:"firstName"`
		Patronymic   string `json:"patronymic"`
		Birthdate    string `json:"birthdate"`
		ClassNumber  string `json:"classNumber"`
		Citizenship  string `json:"citizenship"`
		Disability   string `json:"disability"`
	}

	results := make([]ResultItem, 0, len(applications))

	for _, app := range applications {
		user, userExists := userMap[app.UserID]
		eventName, eventExists := eventMap[app.EventID]

		if !userExists || !eventExists {
			continue
		}

		// Преобразуем статус в строку
		statusStr := "pending"
		switch app.Status {
		case 2:
			statusStr = "approved"
		case 3:
			statusStr = "rejected"
		}

		results = append(results, ResultItem{
			ID:           app.ID,
			OlympiadName: eventName,
			Category:     fmt.Sprintf("%d", app.ClassParticipation),
			Status:       statusStr,
			Surname:      user.Surname,
			FirstName:    user.FirstName,
			Patronymic:   user.Patronymic,
			Birthdate:    user.BirthDate,
			ClassNumber:  user.ClassNumber,
			Citizenship:  user.Citizenship,
			Disability:   user.Disability,
		})
	}

	responseData, err := json.Marshal(results)
	if err != nil {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Error:      "failed to marshal response",
		}, err
	}

	return &responcetypes.ApiResponse{
		Status:     "success",
		StatusCode: http.StatusOK,
		Data:       json.RawMessage(responseData),
	}, nil
}
