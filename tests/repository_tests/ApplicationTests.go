package ApplicationHandlerTest

/*
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/exp/slog"

	"OlimpiadPortal/ApplicationService/internal/repositories/mocks" // Импортируйте mock-репозитории, если используете их
	ApplicationService "OlimpiadPortal/ApplicationService/internal/services/ApplicationService"
)

// Вспомогательная функция для создания тестового сервера с обработчиками
func setupRouter(service *ApplicationService.ApplicationService, logger *slog.Logger) *chi.Mux {
	handler := application_handler.NewApplicationHandler(service, logger)
	router := chi.NewRouter()
	router.Get("/applications", handler.GetAllApplications)
	router.Get("/applications/{id}", handler.GetApplicationByID)
	router.Post("/applications", handler.CreateApplication)
	router.Put("/applications/{id}", handler.UpdateApplicationStatus)
	router.Delete("/applications/{id}", handler.DeleteApplication)
	return router
}

// Тест для получения всех заявок
func TestGetAllApplications(t *testing.T) {
	mockService := new(mocks.ApplicationService)
	mockService.On("GetAllApplications").Return([]dto.ApplicationResponseDTO{}, nil)

	logger := slog.New(slog.HandlerOptions{})
	router := setupRouter(mockService, logger)

	req, _ := http.NewRequest("GET", "/applications", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var applications []dto.ApplicationResponseDTO
	json.NewDecoder(resp.Body).Decode(&applications)
	assert.NotNil(t, applications)
}

// Тест для получения заявки по ID
func TestGetApplicationByID(t *testing.T) {
	mockService := new(mocks.ApplicationService)
	appID := uint(1)
	mockService.On("GetApplicationByID", appID).Return(dto.ApplicationResponseDTO{ApplicationID: appID}, nil)

	logger := slog.New(slog.HandlerOptions{})
	router := setupRouter(mockService, logger)

	req, _ := http.NewRequest("GET", "/applications/"+strconv.Itoa(int(appID)), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var application dto.ApplicationResponseDTO
	json.NewDecoder(resp.Body).Decode(&application)
	assert.Equal(t, appID, application.ApplicationID)
}

// Тест для создания новой заявки
func TestCreateApplication(t *testing.T) {
	mockService := new(mocks.ApplicationService)
	mockService.On("CreateApplication", mock.Anything).Return(uint(1), nil)

	logger := slog.New(slog.HandlerOptions{})
	router := setupRouter(mockService, logger)

	applicationData := dto.CreateApplicationDTO{UserID: 1, EventID: 1}
	requestBody, _ := json.Marshal(applicationData)
	req, _ := http.NewRequest("POST", "/applications", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, uint(1), response["application_id"])
}

// Тест для обновления статуса заявки
func TestUpdateApplicationStatus(t *testing.T) {
	mockService := new(mocks.ApplicationService)
	appID := uint(1)
	status := true
	mockService.On("UpdateApplicationStatus", appID, &status).Return(nil)

	logger := slog.New(slog.HandlerOptions{})
	router := setupRouter(mockService, logger)

	statusData := dto.UpdateApplicationStatusDTO{Status: &status}
	requestBody, _ := json.Marshal(statusData)
	req, _ := http.NewRequest("PUT", "/applications/"+strconv.Itoa(int(appID)), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Статус заявки обновлен", response["message"])
}

// Тест для удаления заявки
func TestDeleteApplication(t *testing.T) {
	mockService := new(mocks.ApplicationService)
	appID := uint(1)
	mockService.On("DeleteApplication", appID).Return(nil)

	logger := slog.New(slog.HandlerOptions{})
	router := setupRouter(mockService, logger)

	req, _ := http.NewRequest("DELETE", "/applications/"+strconv.Itoa(int(appID)), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Заявка удалена", response["message"])
}
*/
