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

// AllApplicationsAggregationStrategy стратегия для агрегации всех заявок
type AllApplicationsAggregationStrategy struct {
	timeout time.Duration
}

// NewAllApplicationsAggregationStrategy создает новую стратегию
func NewAllApplicationsAggregationStrategy(timeout time.Duration) *AllApplicationsAggregationStrategy {
	return &AllApplicationsAggregationStrategy{timeout: timeout}
}

// AggregatedApplicationResponse конечная структура ответа
type AggregatedApplicationResponse struct {
	ID           string `json:"id"`
	Firstname    string `json:"firstname"`
	Surname      string `json:"surname"`
	Patronymic   string `json:"patronymic"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Birthdate    string `json:"birthdate"`
	Gender       int    `json:"gender"`
	ClassNumber  int    `json:"classNumber"`
	Citizenship  int    `json:"citizenship"`
	Disability   int    `json:"disability"`
	SchoolName   string `json:"schoolName"`
	DistrictName string `json:"districtName"`
	OlympiadName string `json:"olympiadName"`
	Profile      string `json:"profile,omitempty"`
	Category     int    `json:"category"`
	Status       int    `json:"status"`
	Code         string `json:"code"`
	SubmittedAt  string `json:"submittedAt"`
}

// Aggregate выполняет агрегацию данных из 5 микросервисов
func (s *AllApplicationsAggregationStrategy) Aggregate(
	targets []config.Target,
	origReq *http.Request,
) (*responcetypes.ApiResponse, error) {
	fmt.Println("=== START AGGREGATION ===")

	client := &http.Client{Timeout: s.timeout}

	// Проверяем что у нас достаточно таргетов
	if len(targets) < 5 {
		return &responcetypes.ApiResponse{
			Status:     "error",
			StatusCode: 500,
			Error:      fmt.Sprintf("need 5 targets for aggregation, got %d", len(targets)),
		}, fmt.Errorf("bad config: expected 5 targets, got %d", len(targets))
	}

	fmt.Printf("Targets count: %d\n", len(targets))
	for i, target := range targets {
		fmt.Printf("Target %d: %s\n", i, target.URL)
	}

	// === 1. Получаем муниципалитеты ===
	fmt.Println("\n[1] Fetching districts...")
	districts, err := s.fetchDistricts(client, targets[0].URL, origReq)
	if err != nil {
		fmt.Printf("❌ Error fetching districts: %v\n", err)
		return s.errorResponse("failed to fetch districts", err)
	}
	fmt.Printf("✅ Got %d districts\n", len(districts))

	// === 2. Получаем все школы ===
	fmt.Println("\n[2] Fetching schools...")
	schools, err := s.fetchSchools(client, targets[1].URL, origReq)
	if err != nil {
		fmt.Printf("❌ Error fetching schools: %v\n", err)
		return s.errorResponse("failed to fetch schools", err)
	}
	fmt.Printf("✅ Got %d schools\n", len(schools))

	// === 3. Получаем все события (олимпиады) ===
	fmt.Println("\n[3] Fetching events...")
	events, err := s.fetchEvents(client, targets[2].URL, origReq)
	if err != nil {
		fmt.Printf("❌ Error fetching events: %v\n", err)
		return s.errorResponse("failed to fetch events", err)
	}
	fmt.Printf("✅ Got %d events\n", len(events))

	// === 4. Получаем все заявки ===
	fmt.Println("\n[4] Fetching applications...")
	applications, err := s.fetchApplications(client, targets[3].URL, origReq)
	if err != nil {
		fmt.Printf("❌ Error fetching applications: %v\n", err)
		return s.errorResponse("failed to fetch applications", err)
	}
	fmt.Printf("✅ Got %d applications\n", len(applications))

	// Если нет заявок - возвращаем пустой массив
	if len(applications) == 0 {
		fmt.Println("⚠️ No applications found, returning empty array")
		return &responcetypes.ApiResponse{
			Status:     "success",
			StatusCode: 200,
			Data:       []AggregatedApplicationResponse{},
		}, nil
	}

	// === 5. Получаем данные пользователей по userId из заявок ===
	fmt.Println("\n[5] Fetching users...")
	userIds := s.extractUserIds(applications)
	fmt.Printf("Extracted %d unique user IDs\n", len(userIds))

	users, err := s.fetchUsers(client, targets[4].URL, userIds, origReq)
	if err != nil {
		fmt.Printf("❌ Error fetching users: %v\n", err)
		return s.errorResponse("failed to fetch users", err)
	}
	fmt.Printf("✅ Got %d users\n", len(users))

	// === 6. Создаем мапы для быстрого доступа ===
	fmt.Println("\n[6] Creating lookup maps...")
	districtMap := s.createDistrictMap(districts)
	schoolMap := s.createSchoolMap(schools, districtMap)
	eventMap := s.createEventMap(events)
	userMap := s.createUserMap(users)

	// === 7. Агрегируем данные ===
	fmt.Println("\n[7] Aggregating final data...")
	aggregatedData := s.aggregateData(applications, schoolMap, eventMap, userMap)
	fmt.Printf("✅ Aggregated %d applications\n", len(aggregatedData))

	fmt.Println("=== AGGREGATION COMPLETE ===")
	return &responcetypes.ApiResponse{
		Status:     "success",
		StatusCode: 200,
		Data:       aggregatedData,
	}, nil
}

// fetchDistricts получает данные муниципалитетов
func (s *AllApplicationsAggregationStrategy) fetchDistricts(client *http.Client, url string, origReq *http.Request) ([]responcetypes.DistrictResponse, error) {
	fmt.Printf("Fetching districts from: %s\n", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request error: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем и логируем тело ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	fmt.Printf("Districts response status: %d\n", resp.StatusCode)
	fmt.Printf("Districts response body: %s\n", string(bodyBytes))

	// Парсим ответ
	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		fmt.Printf("Error unmarshaling districts response: %v\n", err)
		fmt.Printf("Raw response: %s\n", string(bodyBytes))
		return nil, fmt.Errorf("failed to unmarshal districts response: %v", err)
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("districts service returned status: %d, error: %s", apiResp.StatusCode, apiResp.Error)
	}

	// Преобразуем data в массив муниципалитетов
	raw, _ := json.Marshal(apiResp.Data)
	var districts []responcetypes.DistrictResponse
	if err := json.Unmarshal(raw, &districts); err != nil {
		fmt.Printf("Error unmarshaling districts data: %v\n", err)
		fmt.Printf("Raw data: %s\n", string(raw))
		return nil, err
	}

	return districts, nil
}

// fetchSchools получает данные школ
func (s *AllApplicationsAggregationStrategy) fetchSchools(client *http.Client, url string, origReq *http.Request) ([]responcetypes.SchoolResponse, error) {
	fmt.Printf("Fetching schools from: %s\n", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request error: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем и логируем тело ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	fmt.Printf("Schools response status: %d\n", resp.StatusCode)
	fmt.Printf("Schools response body (first 500 chars): %s\n", truncateString(string(bodyBytes), 500))

	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		fmt.Printf("Error unmarshaling schools response: %v\n", err)
		fmt.Printf("Raw response (first 300 chars): %s\n", truncateString(string(bodyBytes), 300))
		return nil, fmt.Errorf("failed to unmarshal schools response: %v", err)
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("schools service returned status: %d, error: %s", apiResp.StatusCode, apiResp.Error)
	}

	raw, _ := json.Marshal(apiResp.Data)
	var schools []responcetypes.SchoolResponse
	if err := json.Unmarshal(raw, &schools); err != nil {
		fmt.Printf("Error unmarshaling schools data: %v\n", err)
		fmt.Printf("Raw data (first 300 chars): %s\n", truncateString(string(raw), 300))
		return nil, err
	}

	return schools, nil
}

// fetchEvents получает данные событий (олимпиад)
func (s *AllApplicationsAggregationStrategy) fetchEvents(client *http.Client, url string, origReq *http.Request) ([]responcetypes.Event, error) {
	fmt.Printf("Fetching events from: %s\n", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request error: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем и логируем тело ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	fmt.Printf("Events response status: %d\n", resp.StatusCode)
	fmt.Printf("Events response body (first 500 chars): %s\n", truncateString(string(bodyBytes), 500))

	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		fmt.Printf("Error unmarshaling events response: %v\n", err)
		fmt.Printf("Raw response (first 300 chars): %s\n", truncateString(string(bodyBytes), 300))

		// Попробуем понять что это за структура
		var rawResponse interface{}
		if err2 := json.Unmarshal(bodyBytes, &rawResponse); err2 == nil {
			fmt.Printf("Response structure: %T\n", rawResponse)
		}

		return nil, fmt.Errorf("failed to unmarshal events response: %v", err)
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("events service returned status: %d, error: %s", apiResp.StatusCode, apiResp.Error)
	}

	raw, _ := json.Marshal(apiResp.Data)
	var events []responcetypes.Event
	if err := json.Unmarshal(raw, &events); err != nil {
		fmt.Printf("Error unmarshaling events data: %v\n", err)
		fmt.Printf("Raw data (first 300 chars): %s\n", truncateString(string(raw), 300))
		return nil, err
	}

	return events, nil
}

// fetchApplications получает данные заявок
func (s *AllApplicationsAggregationStrategy) fetchApplications(client *http.Client, url string, origReq *http.Request) ([]responcetypes.ApplicationResponse, error) {
	fmt.Printf("Fetching applications from: %s\n", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request error: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем и логируем тело ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	fmt.Printf("Applications response status: %d\n", resp.StatusCode)
	fmt.Printf("Applications response body (first 500 chars): %s\n", truncateString(string(bodyBytes), 500))

	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		fmt.Printf("Error unmarshaling applications response: %v\n", err)
		fmt.Printf("Raw response (first 300 chars): %s\n", truncateString(string(bodyBytes), 300))
		return nil, fmt.Errorf("failed to unmarshal applications response: %v", err)
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("applications service returned status: %d, error: %s", apiResp.StatusCode, apiResp.Error)
	}

	raw, _ := json.Marshal(apiResp.Data)
	var applications []responcetypes.ApplicationResponse
	if err := json.Unmarshal(raw, &applications); err != nil {
		fmt.Printf("Error unmarshaling applications data: %v\n", err)
		fmt.Printf("Raw data (first 300 chars): %s\n", truncateString(string(raw), 300))
		return nil, err
	}

	return applications, nil
}

// extractUserIds извлекает уникальные userId из заявок
func (s *AllApplicationsAggregationStrategy) extractUserIds(applications []responcetypes.ApplicationResponse) []string {
	userIdsMap := make(map[string]bool)
	var userIds []string

	for _, app := range applications {
		if _, exists := userIdsMap[app.UserID]; !exists && app.UserID != "" {
			userIdsMap[app.UserID] = true
			userIds = append(userIds, app.UserID)
		}
	}

	return userIds
}

// fetchUsers получает данные пользователей по массиву userId
func (s *AllApplicationsAggregationStrategy) fetchUsers(client *http.Client, url string, userIds []string, origReq *http.Request) ([]responcetypes.UserInfo, error) {
	fmt.Printf("Fetching users from: %s\n", url)
	fmt.Printf("User IDs to fetch: %v\n", userIds)

	if len(userIds) == 0 {
		return []responcetypes.UserInfo{}, nil
	}

	// Подготавливаем тело запроса с массивом ID
	requestBody := map[string][]string{"ids": userIds}
	bodyBytes, _ := json.Marshal(requestBody)
	fmt.Printf("Request body: %s\n", string(bodyBytes))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header = origReq.Header.Clone()
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request error: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Читаем и логируем тело ответа
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	fmt.Printf("Users response status: %d\n", resp.StatusCode)
	fmt.Printf("Users response body (first 1000 chars): %s\n", truncateString(string(respBodyBytes), 1000))

	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(respBodyBytes, &apiResp); err != nil {
		fmt.Printf("Error unmarshaling users response: %v\n", err)
		fmt.Printf("Raw response (first 500 chars): %s\n", truncateString(string(respBodyBytes), 500))

		// Попробуем распарсить как raw JSON
		var rawJSON interface{}
		if err2 := json.Unmarshal(respBodyBytes, &rawJSON); err2 == nil {
			fmt.Printf("Response type: %T\n", rawJSON)
		}

		return nil, fmt.Errorf("failed to unmarshal users response: %v", err)
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("users service returned status: %d, error: %s", apiResp.StatusCode, apiResp.Error)
	}

	raw, _ := json.Marshal(apiResp.Data)
	var users []responcetypes.UserInfo
	if err := json.Unmarshal(raw, &users); err != nil {
		fmt.Printf("Error unmarshaling users data: %v\n", err)
		fmt.Printf("Raw data (first 500 chars): %s\n", truncateString(string(raw), 500))

		// Попробуем понять структуру
		var rawData interface{}
		if err2 := json.Unmarshal(raw, &rawData); err2 == nil {
			fmt.Printf("Data type: %T\n", rawData)
			if arr, ok := rawData.([]interface{}); ok {
				fmt.Printf("Array length: %d\n", len(arr))
				if len(arr) > 0 {
					fmt.Printf("First element type: %T\n", arr[0])
				}
			}
		}

		return nil, err
	}

	return users, nil
}

// createDistrictMap создает map для быстрого доступа к муниципалитетам по ID
func (s *AllApplicationsAggregationStrategy) createDistrictMap(districts []responcetypes.DistrictResponse) map[string]responcetypes.DistrictResponse {
	districtMap := make(map[string]responcetypes.DistrictResponse)
	for _, district := range districts {
		districtMap[district.ID] = district
	}
	return districtMap
}

// createSchoolMap создает map для быстрого доступа к школам по ID
func (s *AllApplicationsAggregationStrategy) createSchoolMap(schools []responcetypes.SchoolResponse, districtMap map[string]responcetypes.DistrictResponse) map[string]responcetypes.SchoolWithDistrict {
	schoolMap := make(map[string]responcetypes.SchoolWithDistrict)

	for _, school := range schools {
		schoolWithDistrict := responcetypes.SchoolWithDistrict{
			School:   school,
			District: responcetypes.DistrictResponse{},
		}

		// Добавляем информацию о муниципалитете если есть
		if district, exists := districtMap[school.DistrictID]; exists {
			schoolWithDistrict.District = district
		}

		schoolMap[school.ID] = schoolWithDistrict
	}

	return schoolMap
}

// createEventMap создает map для быстрого доступа к событиям по ID
func (s *AllApplicationsAggregationStrategy) createEventMap(events []responcetypes.Event) map[string]responcetypes.Event {
	eventMap := make(map[string]responcetypes.Event)
	for _, event := range events {
		eventMap[event.ID] = event
	}
	return eventMap
}

// createUserMap создает map для быстрого доступа к пользователям по ID
func (s *AllApplicationsAggregationStrategy) createUserMap(users []responcetypes.UserInfo) map[string]responcetypes.UserInfo {
	userMap := make(map[string]responcetypes.UserInfo)
	for _, user := range users {
		userMap[user.ID] = user
	}
	return userMap
}

// aggregateData агрегирует данные в финальную структуру
func (s *AllApplicationsAggregationStrategy) aggregateData(
	applications []responcetypes.ApplicationResponse,
	schoolMap map[string]responcetypes.SchoolWithDistrict,
	eventMap map[string]responcetypes.Event,
	userMap map[string]responcetypes.UserInfo,
) []AggregatedApplicationResponse {
	var result []AggregatedApplicationResponse

	for i, app := range applications {
		fmt.Printf("Processing application %d/%d: ID=%s\n", i+1, len(applications), app.ID)

		// Получаем данные пользователя
		user, userExists := userMap[app.UserID]
		if !userExists {
			fmt.Printf("  ⚠️ User not found for ID: %s\n", app.UserID)
			continue
		}

		// Получаем данные школы
		schoolWithDistrict, schoolExists := schoolMap[app.SchoolID]
		if !schoolExists {
			fmt.Printf("  ⚠️ School not found for ID: %s\n", app.SchoolID)
			continue
		}

		// Получаем данные события
		event, eventExists := eventMap[app.EventID]
		if !eventExists {
			fmt.Printf("  ⚠️ Event not found for ID: %s\n", app.EventID)
			continue
		}

		// Форматируем дату рождения (если нужно)
		birthdate := user.Birthdate
		// Убираем временную зону если она есть в формате
		if len(birthdate) > 10 {
			birthdate = birthdate[:10]
		}

		// Создаем агрегированную запись
		aggregated := AggregatedApplicationResponse{
			ID:           app.ID,
			Firstname:    user.Firstname,
			Surname:      user.Surname,
			Patronymic:   user.Patronymic,
			Email:        user.Email,
			Phone:        user.PhoneNumber,
			Birthdate:    birthdate,
			Gender:       user.Gender,
			ClassNumber:  user.ClassNumber,
			Citizenship:  user.Citizenship,
			Disability:   user.Disability,
			SchoolName:   schoolWithDistrict.School.FullName,
			DistrictName: schoolWithDistrict.District.Name,
			OlympiadName: event.Name,
			Profile:      app.Profile,
			Category:     app.ClassParticipation,
			Status:       app.Status,
			Code:         app.Code,
			SubmittedAt:  app.SubmittedAt,
		}

		result = append(result, aggregated)
	}

	return result
}

// errorResponse создает стандартный ответ с ошибкой
func (s *AllApplicationsAggregationStrategy) errorResponse(message string, err error) (*responcetypes.ApiResponse, error) {
	return &responcetypes.ApiResponse{
		Status:     "error",
		StatusCode: 500,
		Error:      fmt.Sprintf("%s: %v", message, err),
	}, err
}

// truncateString обрезает строку до указанной длины
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}
