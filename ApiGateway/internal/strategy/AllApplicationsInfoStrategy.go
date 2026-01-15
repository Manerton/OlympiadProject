package strategy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	client := &http.Client{Timeout: s.timeout}

	if len(targets) < 5 {
		err := fmt.Errorf("expected 5 targets, got %d", len(targets))
		log.Printf("aggregation config error: %v", err)

		return s.errorResponse("invalid aggregation configuration", err)
	}

	districts, err := s.fetchDistricts(client, targets[0].URL, origReq)
	if err != nil {
		log.Printf("failed to fetch districts: %v", err)
		return s.errorResponse("failed to fetch districts", err)
	}

	schools, err := s.fetchSchools(client, targets[1].URL, origReq)
	if err != nil {
		log.Printf("failed to fetch schools: %v", err)
		return s.errorResponse("failed to fetch schools", err)
	}

	events, err := s.fetchEvents(client, targets[2].URL, origReq)
	if err != nil {
		log.Printf("failed to fetch events: %v", err)
		return s.errorResponse("failed to fetch events", err)
	}

	applications, err := s.fetchApplications(client, targets[3].URL, origReq)
	if err != nil {
		log.Printf("failed to fetch applications: %v", err)
		return s.errorResponse("failed to fetch applications", err)
	}

	if len(applications) == 0 {
		return &responcetypes.ApiResponse{
			Status:     "success",
			StatusCode: 200,
			Data:       []AggregatedApplicationResponse{},
		}, nil
	}

	userIds := s.extractUserIds(applications)

	users, err := s.fetchUsers(client, targets[4].URL, userIds, origReq)
	if err != nil {
		log.Printf("failed to fetch users: %v", err)
		return s.errorResponse("failed to fetch users", err)
	}

	districtMap := s.createDistrictMap(districts)
	schoolMap := s.createSchoolMap(schools, districtMap)
	eventMap := s.createEventMap(events)
	userMap := s.createUserMap(users)

	aggregatedData := s.aggregateData(applications, schoolMap, eventMap, userMap)

	return &responcetypes.ApiResponse{
		Status:     "success",
		StatusCode: 200,
		Data:       aggregatedData,
	}, nil
}

// ===================== FETCH METHODS =====================

func (s *AllApplicationsAggregationStrategy) fetchDistricts(
	client *http.Client,
	url string,
	origReq *http.Request,
) ([]responcetypes.DistrictResponse, error) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("districts service returned status %d", apiResp.StatusCode)
	}

	raw, _ := json.Marshal(apiResp.Data)

	var districts []responcetypes.DistrictResponse
	if err := json.Unmarshal(raw, &districts); err != nil {
		return nil, err
	}

	return districts, nil
}

func (s *AllApplicationsAggregationStrategy) fetchSchools(
	client *http.Client,
	url string,
	origReq *http.Request,
) ([]responcetypes.SchoolResponse, error) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("schools service returned status %d", apiResp.StatusCode)
	}

	raw, _ := json.Marshal(apiResp.Data)

	var schools []responcetypes.SchoolResponse
	if err := json.Unmarshal(raw, &schools); err != nil {
		return nil, err
	}

	return schools, nil
}

func (s *AllApplicationsAggregationStrategy) fetchEvents(
	client *http.Client,
	url string,
	origReq *http.Request,
) ([]responcetypes.Event, error) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("events service returned status %d", apiResp.StatusCode)
	}

	raw, _ := json.Marshal(apiResp.Data)

	var events []responcetypes.Event
	if err := json.Unmarshal(raw, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *AllApplicationsAggregationStrategy) fetchApplications(
	client *http.Client,
	url string,
	origReq *http.Request,
) ([]responcetypes.ApplicationResponse, error) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header = origReq.Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("applications service returned status %d", apiResp.StatusCode)
	}

	raw, _ := json.Marshal(apiResp.Data)

	var applications []responcetypes.ApplicationResponse
	if err := json.Unmarshal(raw, &applications); err != nil {
		return nil, err
	}

	return applications, nil
}

// ===================== USERS =====================

func (s *AllApplicationsAggregationStrategy) extractUserIds(
	applications []responcetypes.ApplicationResponse,
) []string {

	idsMap := make(map[string]struct{})
	var ids []string

	for _, app := range applications {
		if app.UserID == "" {
			continue
		}
		if _, ok := idsMap[app.UserID]; !ok {
			idsMap[app.UserID] = struct{}{}
			ids = append(ids, app.UserID)
		}
	}

	return ids
}

func (s *AllApplicationsAggregationStrategy) fetchUsers(
	client *http.Client,
	url string,
	userIds []string,
	origReq *http.Request,
) ([]responcetypes.UserInfo, error) {

	if len(userIds) == 0 {
		return []responcetypes.UserInfo{}, nil
	}

	body, _ := json.Marshal(map[string][]string{"ids": userIds})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header = origReq.Header.Clone()
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp responcetypes.ApiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, err
	}

	if apiResp.StatusCode != 200 {
		return nil, fmt.Errorf("users service returned status %d", apiResp.StatusCode)
	}

	raw, _ := json.Marshal(apiResp.Data)

	var users []responcetypes.UserInfo
	if err := json.Unmarshal(raw, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// ===================== MAPS & AGGREGATION =====================

func (s *AllApplicationsAggregationStrategy) createDistrictMap(
	districts []responcetypes.DistrictResponse,
) map[string]responcetypes.DistrictResponse {

	m := make(map[string]responcetypes.DistrictResponse)
	for _, d := range districts {
		m[d.ID] = d
	}
	return m
}

func (s *AllApplicationsAggregationStrategy) createSchoolMap(
	schools []responcetypes.SchoolResponse,
	districtMap map[string]responcetypes.DistrictResponse,
) map[string]responcetypes.SchoolWithDistrict {

	m := make(map[string]responcetypes.SchoolWithDistrict)
	for _, school := range schools {
		item := responcetypes.SchoolWithDistrict{
			School:   school,
			District: districtMap[school.DistrictID],
		}
		m[school.ID] = item
	}
	return m
}

func (s *AllApplicationsAggregationStrategy) createEventMap(
	events []responcetypes.Event,
) map[string]responcetypes.Event {

	m := make(map[string]responcetypes.Event)
	for _, e := range events {
		m[e.ID] = e
	}
	return m
}

func (s *AllApplicationsAggregationStrategy) createUserMap(
	users []responcetypes.UserInfo,
) map[string]responcetypes.UserInfo {

	m := make(map[string]responcetypes.UserInfo)
	for _, u := range users {
		m[u.ID] = u
	}
	return m
}

func (s *AllApplicationsAggregationStrategy) aggregateData(
	applications []responcetypes.ApplicationResponse,
	schoolMap map[string]responcetypes.SchoolWithDistrict,
	eventMap map[string]responcetypes.Event,
	userMap map[string]responcetypes.UserInfo,
) []AggregatedApplicationResponse {

	var result []AggregatedApplicationResponse

	for _, app := range applications {
		user, okU := userMap[app.UserID]
		school, okS := schoolMap[app.SchoolID]
		event, okE := eventMap[app.EventID]

		if !okU || !okS || !okE {
			continue
		}

		birthdate := user.Birthdate
		if len(birthdate) > 10 {
			birthdate = birthdate[:10]
		}

		result = append(result, AggregatedApplicationResponse{
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
			SchoolName:   school.School.FullName,
			DistrictName: school.District.Name,
			OlympiadName: event.Name,
			Profile:      app.Profile,
			Category:     app.ClassParticipation,
			Status:       app.Status,
			Code:         app.Code,
			SubmittedAt:  app.SubmittedAt,
		})
	}

	return result
}

// ===================== ERROR =====================

func (s *AllApplicationsAggregationStrategy) errorResponse(
	message string,
	err error,
) (*responcetypes.ApiResponse, error) {

	return &responcetypes.ApiResponse{
		Status:     "error",
		StatusCode: 500,
		Error:      fmt.Sprintf("%s: %v", message, err),
	}, err
}
