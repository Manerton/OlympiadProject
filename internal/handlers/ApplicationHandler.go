package ApplicationHandler

import (
	ApplicationDto "OlimpiadPortal/ApplicationService/internal/dto"
	"OlimpiadPortal/ApplicationService/internal/lib/response"
	application_service "OlimpiadPortal/ApplicationService/internal/services"
	"fmt"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	//"github.com/go-playground/validator/v10"
)

type ApplicationHandler struct {
	service *application_service.ApplicationService
	logger  *slog.Logger
}

// Конструктор обработчика заявок
func NewApplicationHandler(service *application_service.ApplicationService, logger *slog.Logger) *ApplicationHandler {
	return &ApplicationHandler{service: service, logger: logger}
}

// Получение всех заявок
func (h *ApplicationHandler) GetAllApplications(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Получение всех заявок")
	applications, err := h.service.GetAllApplications()

	if err != nil {
		h.logger.Error("Ошибка получения всех заявок", slog.Any("error", err))
		http.Error(w, "Не удалось получить заявки", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   applications,
	})

}

// Получение заявки по ID
func (h *ApplicationHandler) GetApplicationByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		h.logger.Error("Некорректный ID", slog.Any("error", err))
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Получение заявки по ID", slog.Uint64("id", id))
	application, err := h.service.GetApplicationByID(uint(id))
	if err != nil {
		h.logger.Error("Ошибка получения заявки", slog.Any("error", err))
		http.Error(w, "Заявка не найдена", http.StatusNotFound)
		return
	}
	//render.JSON(w, r, application)
	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   application,
	})
}

// Получение заявок пользователя по ID
func (h *ApplicationHandler) GetApplicationsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 32)
	if err != nil {
		h.logger.Error("Некорректный userID", slog.Any("error", err))
		http.Error(w, "Некорректный userID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Получение заявок пользователя", slog.Uint64("userID", userID))
	applications, err := h.service.GetApplicationsByUserID(uint(userID))
	if err != nil {
		h.logger.Error("Ошибка получения заявок пользователя", slog.Any("error", err))
		http.Error(w, "Не удалось получить заявки", http.StatusInternalServerError)
		return
	}
	//render.JSON(w, r, applications)
	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   applications,
	})
}

// Получение заявок по ID события
func (h *ApplicationHandler) GetApplicationsByEventID(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.ParseUint(chi.URLParam(r, "eventID"), 10, 32)
	if err != nil {
		h.logger.Error("Некорректный eventID", slog.Any("error", err))
		http.Error(w, "Некорректный eventID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Получение заявок по ID события", slog.Uint64("eventID", eventID))
	applications, err := h.service.GetApplicationsByEventID(uint(eventID))
	if err != nil {
		h.logger.Error("Ошибка получения заявок события", slog.Any("error", err))
		http.Error(w, "Не удалось получить заявки", http.StatusInternalServerError)
		return
	}
	//render.JSON(w, r, applications)
	render.JSON(w, r, response.ApiResponse{
		Status: response.StatusOK,
		Data:   applications,
	})
}

// Создание новой заявки
func (h *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var input ApplicationDto.CreateApplicationDTO
	if err := render.DecodeJSON(r.Body, &input); err != nil {
		h.logger.Error("Ошибка декодирования данных", slog.Any("error", err))
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	h.logger.Info("Создание новой заявки", slog.Any("user_id", input.UserID))
	id, err := h.service.CreateApplication(input)
	if err != nil {
		h.logger.Error("Ошибка создания заявки", slog.Any("error", err))
		http.Error(w, "Не удалось создать заявку", http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, map[string]interface{}{"application_id": id})
}

// Обновление статуса заявки
func (h *ApplicationHandler) UpdateApplicationStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		h.logger.Error("Некорректный ID", slog.Any("error", err))
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	var input ApplicationDto.UpdateApplicationStatusDTO
	if err := render.DecodeJSON(r.Body, &input); err != nil {
		h.logger.Error("Ошибка декодирования данных", slog.Any("error", err))
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	h.logger.Info("Обновление статуса заявки", slog.Uint64("id", id), slog.Any("status", input.Status))
	if err := h.service.UpdateApplicationStatus(uint(id), input); err != nil {
		h.logger.Error("Ошибка обновления статуса заявки", slog.Any("error", err))
		http.Error(w, "Не удалось обновить статус", http.StatusInternalServerError)
		return
	}
	//render.JSON(w, r, map[string]interface{}{"message": "Статус заявки обновлен"})
	render.JSON(w, r, response.Success(fmt.Sprintf("id = %d", id)))
}

// Удаление заявки
func (h *ApplicationHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		h.logger.Error("Некорректный ID", slog.Any("error", err))
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	h.logger.Info("Удаление заявки", slog.Uint64("id", id))
	if err := h.service.DeleteApplication(uint(id)); err != nil {
		h.logger.Error("Ошибка удаления заявки", slog.Any("error", err))
		http.Error(w, "Не удалось удалить заявку", http.StatusInternalServerError)
		return
	}

	//render.JSON(w, r, map[string]interface{}{"message": "Заявка удалена"})
	render.JSON(w, r, response.Success("Заявка удалена"))
}

/* // Пример: Получение информации о событии
type EventDetails struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

// Запрос в сервис событий (Event Service)
func getEventDetails(eventID uint) (EventDetails, error) {
	// Сделайте HTTP-запрос к Event Service
	resp, err := http.Get(fmt.Sprintf("http://event-service/events/%d", eventID))
	if err != nil {
		return EventDetails{}, err
	}
	defer resp.Body.Close()

	var details EventDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return EventDetails{}, err
	}

	return details, nil
} */

/* func syncEventDetails(application *Application) error {
	details, err := getEventDetails(application.EventID)
	if err != nil {
		return err
	}

	application.EventName = details.Name
	application.EventLocation = details.Location
	return nil
} */
