package ApplicationHandler

import (
	ApplicationDto "OlimpiadPortal/ApplicationService/internal/dto"
	application_service "OlimpiadPortal/ApplicationService/internal/services"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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
	render.JSON(w, r, applications)
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
	render.JSON(w, r, application)
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
	render.JSON(w, r, map[string]interface{}{"message": "Статус заявки обновлен"})
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
	render.JSON(w, r, map[string]interface{}{"message": "Заявка удалена"})
}
